package simulation

import (
	"math/rand"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	dymsimtypes "github.com/dymensionxyz/dymension/v3/simulation/types"
	"github.com/dymensionxyz/dymension/v3/x/incentives/keeper"
	"github.com/dymensionxyz/dymension/v3/x/incentives/types"
	lockuptypes "github.com/dymensionxyz/dymension/v3/x/lockup/types"
)

// Simulation operation weights constants.
const (
	DefaultWeightMsgCreateGauge int = 100
	DefaultWeightMsgAddToGauge  int = 100
	OpWeightMsgCreateGauge          = "op_weight_msg_create_gauge" //nolint:gosec
	OpWeightMsgAddToGauge           = "op_weight_msg_add_to_gauge" //nolint:gosec
)

// WeightedOperations returns all the operations from the module with their respective weights.
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONCodec,
	txCfg client.TxConfig,

	ak dymsimtypes.AccountKeeper,
	bk dymsimtypes.BankKeeper,
	ek types.EpochKeeper,
	k keeper.Keeper,
) simulation.WeightedOperations {
	var (
		weightMsgCreateGauge int
		weightMsgAddToGauge  int
	)

	appParams.GetOrGenerate(
		OpWeightMsgCreateGauge, &weightMsgCreateGauge, nil,
		func(*rand.Rand) { weightMsgCreateGauge = DefaultWeightMsgCreateGauge },
	)

	appParams.GetOrGenerate(
		OpWeightMsgAddToGauge, &weightMsgAddToGauge, nil,
		func(*rand.Rand) { weightMsgAddToGauge = DefaultWeightMsgAddToGauge },
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreateGauge,
			SimulateMsgCreateGauge(txCfg, ak, bk, ek, k),
		),
		simulation.NewWeightedOperation(
			weightMsgAddToGauge,
			SimulateMsgAddToGauge(txCfg, ak, bk, k),
		),
	}
}

// genRewardCoins generates a random number of coin denoms with a respective random value for each coin.
func genRewardCoins(r *rand.Rand, coins sdk.Coins, fee math.Int) (res sdk.Coins) {
	numCoins := 1 + r.Intn(min(coins.Len(), 1))
	denomIndices := r.Perm(numCoins)
	for i := 0; i < numCoins; i++ {
		var (
			amt math.Int
			err error
		)
		denom := coins[denomIndices[i]].Denom
		if denom == sdk.DefaultBondDenom {
			amt, err = simtypes.RandPositiveInt(r, coins[i].Amount.Sub(fee))
			if err != nil {
				panic(err)
			}
		} else {
			amt, err = simtypes.RandPositiveInt(r, coins[i].Amount)
			if err != nil {
				panic(err)
			}
		}
		res = append(res, sdk.Coin{Denom: denom, Amount: amt})
	}
	return
}

// genQueryCondition returns a single lockup QueryCondition, which is generated from a single coin randomly selected from the provided coin array
func genQueryCondition(r *rand.Rand, blocktime time.Time, coins sdk.Coins, durations []time.Duration) lockuptypes.QueryCondition {
	denom := coins[r.Intn(len(coins))].Denom
	durationIndex := r.Intn(len(durations))
	duration := durations[durationIndex]
	// lock_age is not used in simulation, set to 0
	return lockuptypes.QueryCondition{
		Denom:    denom,
		Duration: duration,
	}
}

// SimulateMsgCreateGauge generates and executes a MsgCreateGauge with random parameters
func SimulateMsgCreateGauge(
	txConfig client.TxConfig,
	ak dymsimtypes.AccountKeeper,
	bk dymsimtypes.BankKeeper,
	ek types.EpochKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		params := k.GetParams(ctx)
		// we always expect that we add no more than 1 denom to the gauge in simulation
		fee := params.CreateGaugeBaseFee.Add(params.AddDenomFee.MulRaw(1))
		feeCoin := sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: fee}

		simAccount, _ := simtypes.RandomAcc(r, accs)
		simCoins := bk.SpendableCoins(ctx, simAccount.Address)
		if simCoins.AmountOf(sdk.DefaultBondDenom).LT(fee) {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(&types.MsgCreateGauge{}), "Account have no coin"), nil, nil
		}

		distributeTo := genQueryCondition(r, ctx.BlockTime(), simCoins, types.DefaultGenesis().LockableDurations)
		rewards := genRewardCoins(r, simCoins, fee)
		startTimeSecs := r.Intn(1 * 60 * 60 * 24 * 7) // range of 1 week
		startTime := ctx.BlockTime().Add(time.Duration(startTimeSecs) * time.Second)
		numEpochsPaidOver := uint64(1) // == 1 since we only support perpetual gauges

		msg := &types.MsgCreateGauge{
			Owner:             simAccount.Address.String(),
			IsPerpetual:       true, // all gauges are perpetual
			GaugeType:         types.GaugeType_GAUGE_TYPE_ASSET,
			Asset:             &distributeTo,
			Coins:             rewards,
			StartTime:         startTime,
			NumEpochsPaidOver: numEpochsPaidOver,
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           txConfig,
			Cdc:             nil,
			Msg:             msg,
			CoinsSpentInMsg: rewards.Add(feeCoin),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgAddToGauge generates and executes a MsgAddToGauge with random parameters
func SimulateMsgAddToGauge(
	txConfig client.TxConfig,
	ak dymsimtypes.AccountKeeper,
	bk dymsimtypes.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		gauge := dymsimtypes.RandomGauge(ctx, r, k)
		if gauge == nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(&types.MsgAddToGauge{}), "No gauge exists"), nil, nil
		} else if gauge.IsFinishedGauge(ctx.BlockTime()) {
			// TODO: Ideally we'd still run this but expect failure.
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(&types.MsgAddToGauge{}), "Selected a gauge that is finished"), nil, nil
		}

		params := k.GetParams(ctx)
		// we always expect that we add no more than 1 denom to the gauge in simulation
		fee := params.AddToGaugeBaseFee.Add(params.AddDenomFee.MulRaw(int64(1 + len(gauge.Coins))))
		feeCoin := sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: fee}

		simAccount, _ := simtypes.RandomAcc(r, accs)
		simCoins := bk.SpendableCoins(ctx, simAccount.Address)
		if simCoins.AmountOf(sdk.DefaultBondDenom).LT(fee) {
			return simtypes.NoOpMsg(
				types.ModuleName, sdk.MsgTypeURL(&types.MsgAddToGauge{}), "Account have no coin"), nil, nil
		}

		rewards := genRewardCoins(r, simCoins, fee)
		msg := &types.MsgAddToGauge{
			Owner:   simAccount.Address.String(),
			GaugeId: gauge.Id,
			Rewards: rewards,
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           txConfig,
			Cdc:             nil,
			Msg:             msg,
			CoinsSpentInMsg: rewards.Add(feeCoin),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
