package keeper

import (
	"fmt"
	"sort"
	"time"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gogoproto/proto"
	"github.com/osmosis-labs/osmosis/v15/osmoutils/sumtree"

	"github.com/dymensionxyz/dymension/v3/x/lockup/types"
)

// WithdrawAllMaturedLocks withdraws every lock thats in the process of unlocking, and has finished unlocking by
// the current block time.
func (k Keeper) WithdrawAllMaturedLocks(ctx sdk.Context) {
	k.unlockFromIterator(ctx, k.LockIteratorBeforeTime(ctx, ctx.BlockTime()))
}

// GetModuleBalance returns full balance of the module.
func (k Keeper) GetModuleBalance(ctx sdk.Context) sdk.Coins {
	acc := k.ak.GetModuleAccount(ctx, types.ModuleName)
	return k.bk.GetAllBalances(ctx, acc.GetAddress())
}

// GetModuleLockedCoins Returns locked balance of the module.
func (k Keeper) GetModuleLockedCoins(ctx sdk.Context) sdk.Coins {
	// all not unlocking + not finished unlocking
	notUnlockingCoins := k.getCoinsFromIterator(ctx, k.LockIterator(ctx, false))
	unlockingCoins := k.getCoinsFromIterator(ctx, k.LockIteratorAfterTime(ctx, ctx.BlockTime()))
	return notUnlockingCoins.Add(unlockingCoins...)
}

// GetPeriodLocksAccumulation returns the total amount of query.Denom tokens locked for longer than
// query.Duration.
func (k Keeper) GetPeriodLocksAccumulation(ctx sdk.Context, query types.QueryCondition) math.Int {
	beginKey := accumulationKey(query.Duration)
	return k.accumulationStore(ctx, query.Denom).SubsetAccumulation(beginKey, nil)
}

// BeginUnlockAllNotUnlockings begins unlock for all not unlocking locks of the given account.
func (k Keeper) BeginUnlockAllNotUnlockings(ctx sdk.Context, account sdk.AccAddress) ([]types.PeriodLock, error) {
	locks, err := k.beginUnlockFromIterator(ctx, k.AccountLockIterator(ctx, false, account))
	return locks, err
}

// AddToExistingLock adds the given coin to the existing lock with the same owner and duration.
// Returns the updated lock ID if successfully added coin, returns 0 and error when a lock with
// given condition does not exist, or if fails to add to lock.
func (k Keeper) AddToExistingLock(ctx sdk.Context, owner sdk.AccAddress, coin sdk.Coin, duration time.Duration) (uint64, error) {
	locks := k.GetAccountLockedDurationNotUnlockingOnly(ctx, owner, coin.Denom, duration)

	// if no lock exists for the given owner + denom + duration, return an error
	if len(locks) < 1 {
		return 0, errorsmod.Wrapf(types.ErrLockupNotFound, "lock with denom %s before duration %s does not exist", coin.Denom, duration.String())
	}
	// there should only be a single lock with the same duration + token, thus we take the first lock
	lock := &locks[0]

	_, err := k.AddTokensToLockByID(ctx, lock.ID, owner, coin)
	if err != nil {
		return 0, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// update the timestamp of the lock
	lock, err = k.GetLockByID(ctx, lock.ID)
	if err != nil {
		return 0, err
	}
	lock.UpdatedAt = ctx.BlockTime()
	err = k.setLock(ctx, *lock)
	if err != nil {
		return 0, err
	}

	return lock.ID, nil
}

// HasLock returns true if lock with the given condition exists
func (k Keeper) HasLock(ctx sdk.Context, owner sdk.AccAddress, denom string, duration time.Duration) bool {
	locks := k.GetAccountLockedDurationNotUnlockingOnly(ctx, owner, denom, duration)
	return len(locks) > 0
}

// HasUnlockingLock returns true if unlocking lock with the given condition exists
func (k Keeper) HasUnlockingLock(ctx sdk.Context, owner sdk.AccAddress, denom string, duration time.Duration) bool {
	locks := k.GetAccountLockedDurationUnlockingOnly(ctx, owner, denom, duration)
	return len(locks) > 0
}

// AddTokensToLockByID locks additional tokens into an existing lock with the given ID.
// Tokens locked are sent and kept in the module account.
// This method alters the lock state in store, thus we do a sanity check to ensure
// lock owner matches the given owner.
func (k Keeper) AddTokensToLockByID(ctx sdk.Context, lockID uint64, owner sdk.AccAddress, tokensToAdd sdk.Coin) (*types.PeriodLock, error) {
	lock, err := k.GetLockByID(ctx, lockID)
	if err != nil {
		return nil, err
	}

	if lock.GetOwner() != owner.String() {
		return nil, types.ErrNotLockOwner
	}

	lock.Coins = lock.Coins.Add(tokensToAdd)
	err = k.lock(ctx, *lock, sdk.NewCoins(tokensToAdd))
	if err != nil {
		return nil, err
	}

	if k.hooks == nil {
		return lock, nil
	}

	k.hooks.AfterAddTokensToLock(ctx, lock.OwnerAddress(), lock.GetID(), sdk.NewCoins(tokensToAdd))

	return lock, nil
}

// CreateLock creates a new lock with the specified duration for the owner.
// Returns an error in the following conditions:
//   - account does not have enough balance
func (k Keeper) CreateLock(ctx sdk.Context, owner sdk.AccAddress, coins sdk.Coins, duration time.Duration) (types.PeriodLock, error) {
	ID := k.GetLastLockID(ctx) + 1
	// unlock time is initially set without a value, gets set as unlock start time + duration
	// when unlocking starts.
	lock := types.NewPeriodLock(ID, owner, duration, time.Time{}, coins, ctx.BlockTime())
	err := k.lock(ctx, lock, lock.Coins)
	if err != nil {
		return lock, err
	}

	// add lock refs into not unlocking queue
	err = k.addLockRefs(ctx, lock)
	if err != nil {
		return lock, err
	}

	k.SetLastLockID(ctx, lock.ID)
	return lock, nil
}

// lock is an internal utility to lock coins and set corresponding states.
// This is only called by either of the two possible entry points to lock tokens.
// 1. CreateLock
// 2. AddTokensToLockByID
func (k Keeper) lock(ctx sdk.Context, lock types.PeriodLock, tokensToLock sdk.Coins) error {
	owner, err := sdk.AccAddressFromBech32(lock.Owner)
	if err != nil {
		return err
	}
	if err := k.bk.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, tokensToLock); err != nil {
		return err
	}

	// store lock object into the store
	err = k.setLock(ctx, lock)
	if err != nil {
		return err
	}

	// add to accumulation store
	for _, coin := range tokensToLock {
		k.accumulationStore(ctx, coin.Denom).Increase(accumulationKey(lock.Duration), coin.Amount)
	}

	k.hooks.OnTokenLocked(ctx, owner, lock.ID, lock.Coins, lock.Duration, lock.EndTime)
	return nil
}

// BeginUnlock is a utility to start unlocking coins from NotUnlocking queue.
func (k Keeper) BeginUnlock(ctx sdk.Context, lockID uint64, coins sdk.Coins) (uint64, error) {
	lock, err := k.GetLockByID(ctx, lockID)
	if err != nil {
		return 0, err
	}

	unlockingLock, err := k.beginUnlock(ctx, *lock, coins)
	return unlockingLock, err
}

// beginUnlock unlocks specified tokens from the given lock. Existing lock refs
// of not unlocking queue are deleted and new lock refs are then added.
// EndTime of the lock is set within this method.
// Coins provided as the parameter does not require to have all the tokens in the lock,
// as we allow partial unlockings of a lock.
// Returns lock id, new lock id if the lock was split, else same lock id.
func (k Keeper) beginUnlock(ctx sdk.Context, lock types.PeriodLock, coins sdk.Coins) (uint64, error) {
	// sanity check
	if !coins.IsAllLTE(lock.Coins) {
		return 0, fmt.Errorf("requested amount to unlock exceeds locked tokens")
	}

	if lock.IsUnlocking() {
		return 0, fmt.Errorf("trying to unlock a lock that is already unlocking")
	}

	// If the amount were unlocking is empty, or the entire coins amount, unlock the entire lock.
	// Otherwise, split the lock into two locks, and fully unlock the newly created lock.
	// (By virtue, the newly created lock we split into should have the unlock amount)
	if len(coins) != 0 && !coins.Equal(lock.Coins) {
		splitLock, err := k.splitLock(ctx, lock, coins, false)
		if err != nil {
			return 0, err
		}
		lock = splitLock
	}

	// remove existing lock refs from not unlocking queue
	err := k.deleteLockRefs(ctx, types.KeyPrefixNotUnlocking, lock)
	if err != nil {
		return 0, err
	}

	// store lock with the end time set to current block time + duration
	lock.EndTime = ctx.BlockTime().Add(lock.Duration)
	err = k.setLock(ctx, lock)
	if err != nil {
		return 0, err
	}

	// add lock refs into unlocking queue
	err = k.addLockRefs(ctx, lock)
	if err != nil {
		return 0, err
	}

	if k.hooks != nil {
		k.hooks.OnStartUnlock(ctx, lock.OwnerAddress(), lock.ID, lock.Coins, lock.Duration, lock.EndTime)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		createBeginUnlockEvent(&lock),
	})

	return lock.ID, nil
}

func (k Keeper) clearKeysByPrefix(ctx sdk.Context, prefix []byte) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, prefix)
	defer iterator.Close() // nolint: errcheck

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

func (k Keeper) ClearAccumulationStores(ctx sdk.Context) {
	k.clearKeysByPrefix(ctx, types.KeyPrefixLockAccumulation)
}

func (k Keeper) BeginForceUnlockWithEndTime(ctx sdk.Context, lockID uint64, endTime time.Time) error {
	lock, err := k.GetLockByID(ctx, lockID)
	if err != nil {
		return err
	}
	return k.beginForceUnlockWithEndTime(ctx, *lock, endTime)
}

func (k Keeper) beginForceUnlockWithEndTime(ctx sdk.Context, lock types.PeriodLock, endTime time.Time) error {
	// remove lock refs from not unlocking queue if exists
	err := k.deleteLockRefs(ctx, types.KeyPrefixNotUnlocking, lock)
	if err != nil {
		return err
	}

	// store lock with end time set
	lock.EndTime = endTime
	err = k.setLock(ctx, lock)
	if err != nil {
		return err
	}

	// add lock refs into unlocking queue
	err = k.addLockRefs(ctx, lock)
	if err != nil {
		return err
	}

	if k.hooks != nil {
		k.hooks.OnStartUnlock(ctx, lock.OwnerAddress(), lock.ID, lock.Coins, lock.Duration, lock.EndTime)
	}

	return nil
}

// UnlockMaturedLock finishes unlocking by sending back the locked tokens from the module accounts
// to the owner. This method requires lock to be matured, having passed the endtime of the lock.
func (k Keeper) UnlockMaturedLock(ctx sdk.Context, lockID uint64) error {
	lock, err := k.GetLockByID(ctx, lockID)
	if err != nil {
		return err
	}

	// validation for current time and unlock time
	curTime := ctx.BlockTime()
	if !lock.IsUnlocking() {
		return fmt.Errorf("lock hasn't started unlocking yet")
	}
	if curTime.Before(lock.EndTime) {
		return fmt.Errorf("lock is not unlockable yet: %s >= %s", curTime.String(), lock.EndTime.String())
	}

	return k.unlockMaturedLockInternalLogic(ctx, *lock)
}

// PartialForceUnlock begins partial ForceUnlock of given lock for the given amount of coins.
// ForceUnlocks the lock as a whole when provided coins are empty, or coin provided equals amount of coins in the lock.
// This also supports the case of lock in an unbonding status.
func (k Keeper) PartialForceUnlock(ctx sdk.Context, lock types.PeriodLock, coins sdk.Coins) error {
	// sanity check
	if !coins.IsAllLTE(lock.Coins) {
		return fmt.Errorf("requested amount to unlock exceeds locked tokens")
	}

	// split lock to support partial force unlock.
	// (By virtue, the newly created lock we split into should have the unlock amount)
	if len(coins) != 0 && !coins.Equal(lock.Coins) {
		splitLock, err := k.splitLock(ctx, lock, coins, true)
		if err != nil {
			return err
		}
		lock = splitLock
	}

	return k.ForceUnlock(ctx, lock)
}

// ForceUnlock ignores unlock duration and immediately unlocks the lock and refunds tokens to lock owner.
func (k Keeper) ForceUnlock(ctx sdk.Context, lock types.PeriodLock) error {
	// Steps:

	// 2) If lock is bonded, move it to unlocking
	// 3) Run logic to delete unlocking metadata, and send tokens to owner.

	if !lock.IsUnlocking() {
		_, err := k.BeginUnlock(ctx, lock.ID, nil)
		if err != nil {
			return err
		}
	}
	// NOTE: This caused a bug! BeginUnlock changes the owner the lock.EndTime
	// This shows the bad API design of not using lock.ID in every public function.
	lockPtr, err := k.GetLockByID(ctx, lock.ID)
	if err != nil {
		return err
	}
	return k.unlockMaturedLockInternalLogic(ctx, *lockPtr)
}

// unlockMaturedLockInternalLogic handles internal logic for finishing unlocking matured locks.
func (k Keeper) unlockMaturedLockInternalLogic(ctx sdk.Context, lock types.PeriodLock) error {
	owner, err := sdk.AccAddressFromBech32(lock.Owner)
	if err != nil {
		return err
	}

	// send coins back to owner
	if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, lock.Coins); err != nil {
		return err
	}

	k.deleteLock(ctx, lock.ID)

	// delete lock refs from the unlocking queue
	err = k.deleteLockRefs(ctx, types.KeyPrefixUnlocking, lock)
	if err != nil {
		return err
	}

	// remove from accumulation store
	for _, coin := range lock.Coins {
		k.accumulationStore(ctx, coin.Denom).Decrease(accumulationKey(lock.Duration), coin.Amount)
	}

	k.hooks.OnTokenUnlocked(ctx, owner, lock.ID, lock.Coins, lock.Duration, lock.EndTime)
	return nil
}

// ExtendLockup changes the existing lock duration to the given lock duration.
// Updating lock duration would fail on either of the following conditions.
// 1. Only lock owner is able to change the duration of the lock.
// 2. Locks that are unlocking are not allowed to change duration.
// 3. Locks that have synthetic lockup are not allowed to change.
// 4. Provided duration should be greater than the original duration.
func (k Keeper) ExtendLockup(ctx sdk.Context, lockID uint64, owner sdk.AccAddress, newDuration time.Duration) error {
	lock, err := k.GetLockByID(ctx, lockID)
	if err != nil {
		return err
	}

	if lock.GetOwner() != owner.String() {
		return types.ErrNotLockOwner
	}

	if lock.IsUnlocking() {
		return fmt.Errorf("cannot edit unlocking lockup for lock %d", lock.ID)
	}

	oldDuration := lock.GetDuration()
	if newDuration <= oldDuration {
		return fmt.Errorf("new duration should be greater than the original")
	}

	// completely delete existing lock refs
	err = k.deleteLockRefs(ctx, unlockingPrefix(lock.IsUnlocking()), *lock)
	if err != nil {
		return err
	}

	// update accumulation store
	for _, coin := range lock.Coins {
		k.accumulationStore(ctx, coin.Denom).Decrease(accumulationKey(lock.Duration), coin.Amount)
		k.accumulationStore(ctx, coin.Denom).Increase(accumulationKey(newDuration), coin.Amount)
	}

	lock.Duration = newDuration

	// add lock refs with the new duration
	err = k.addLockRefs(ctx, *lock)
	if err != nil {
		return err
	}

	err = k.setLock(ctx, *lock)
	if err != nil {
		return err
	}

	k.hooks.OnLockupExtend(ctx,
		lock.GetID(),
		oldDuration,
		lock.GetDuration(),
	)

	return nil
}

// InitializeAllLocks takes a set of locks, and initializes state to be storing
// them all correctly. This utilizes batch optimizations to improve efficiency,
// as this becomes a bottleneck at chain initialization & upgrades.
func (k Keeper) InitializeAllLocks(ctx sdk.Context, locks []types.PeriodLock) error {
	// index by coin.Denom, them duration -> amt
	// We accumulate the accumulation store entries separately,
	// to avoid hitting the myriad of slowdowns in the SDK iterator creation process.
	// We then save these once to the accumulation store at the end.
	accumulationStoreEntries := make(map[string]map[time.Duration]math.Int)
	denoms := []string{}
	for i, lock := range locks {
		if i%25000 == 0 {
			msg := fmt.Sprintf("Reset %d lock refs, cur lock ID %d", i, lock.ID)
			ctx.Logger().Info(msg)
		}
		err := k.setLockAndAddLockRefs(ctx, lock)
		if err != nil {
			return err
		}

		// Add to the accumulation store cache
		for _, coin := range lock.Coins {
			// update or create the new map from duration -> Int for this denom.
			var curDurationMap map[time.Duration]math.Int
			if durationMap, ok := accumulationStoreEntries[coin.Denom]; ok {
				curDurationMap = durationMap
				// update or create new amount in the duration map
				newAmt := coin.Amount
				if curAmt, ok := durationMap[lock.Duration]; ok {
					newAmt = newAmt.Add(curAmt)
				}
				curDurationMap[lock.Duration] = newAmt
			} else {
				denoms = append(denoms, coin.Denom)
				curDurationMap = map[time.Duration]math.Int{lock.Duration: coin.Amount}
			}
			accumulationStoreEntries[coin.Denom] = curDurationMap
		}
	}

	// deterministically iterate over durationMap cache.
	sort.Strings(denoms)
	for _, denom := range denoms {
		curDurationMap := accumulationStoreEntries[denom]
		durations := make([]time.Duration, 0, len(curDurationMap))
		for duration := range curDurationMap {
			durations = append(durations, duration)
		}
		sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
		// now that we have a sorted list of durations for this denom,
		// add them all to accumulation store
		msg := fmt.Sprintf("Setting accumulation entries for locks for %s, there are %d distinct durations",
			denom, len(durations))
		ctx.Logger().Info(msg)
		for _, d := range durations {
			amt := curDurationMap[d]
			k.accumulationStore(ctx, denom).Increase(accumulationKey(d), amt)
		}
	}

	return nil
}

func (k Keeper) accumulationStore(ctx sdk.Context, denom string) sumtree.Tree {
	return sumtree.NewTree(prefix.NewStore(ctx.KVStore(k.storeKey), accumulationStorePrefix(denom)), 10)
}

// setLock is a utility to store lock object into the store.
func (k Keeper) setLock(ctx sdk.Context, lock types.PeriodLock) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := proto.Marshal(&lock)
	if err != nil {
		return err
	}
	store.Set(lockStoreKey(lock.ID), bz)
	return nil
}

// setLockAndAddLockRefs sets the lock, and resets all of its lock references
// This puts the lock into a 'clean' state, aside from the AccumulationStore.
func (k Keeper) setLockAndAddLockRefs(ctx sdk.Context, lock types.PeriodLock) error {
	err := k.setLock(ctx, lock)
	if err != nil {
		return err
	}

	return k.addLockRefs(ctx, lock)
}

// deleteLock removes the lock object from the state.
func (k Keeper) deleteLock(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(lockStoreKey(id))
}

// splitLock splits a lock with the given amount, and stores split new lock to the state.
// Returns the new lock after modifying the state of the old lock.
func (k Keeper) splitLock(ctx sdk.Context, lock types.PeriodLock, coins sdk.Coins, forceUnlock bool) (types.PeriodLock, error) {
	if !forceUnlock && lock.IsUnlocking() {
		return types.PeriodLock{}, fmt.Errorf("cannot split unlocking lock")
	}

	lock.Coins = lock.Coins.Sub(coins...)
	err := k.setLock(ctx, lock)
	if err != nil {
		return types.PeriodLock{}, err
	}

	// create a new lock
	splitLockID := k.GetLastLockID(ctx) + 1
	k.SetLastLockID(ctx, splitLockID)

	splitLock := types.NewPeriodLock(splitLockID, lock.OwnerAddress(), lock.Duration, lock.EndTime, coins, ctx.BlockTime())

	err = k.setLock(ctx, splitLock)
	return splitLock, err
}

func (k Keeper) getCoinsFromLocks(locks []types.PeriodLock) sdk.Coins {
	coins := sdk.Coins{}
	for _, lock := range locks {
		coins = coins.Add(lock.Coins...)
	}
	return coins
}

// SetLock sets a lock in the store (public for migrations).
func (k Keeper) SetLock(ctx sdk.Context, lock types.PeriodLock) error {
	return k.setLock(ctx, lock)
}
