package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/dymensionxyz/dymension/v3/x/delayedack/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/osmosis-labs/osmosis/v15/x/epochs/types"

	commontypes "github.com/dymensionxyz/dymension/v3/x/common/types"
	eibctypes "github.com/dymensionxyz/dymension/v3/x/eibc/types"
)

/* -------------------------------------------------------------------------- */
/*                                 eIBC Hooks                                 */
/* -------------------------------------------------------------------------- */

var _ eibctypes.EIBCHooks = eibcHooks{}

const (
	deletePacketsBatchSize = 1000
)

type eibcHooks struct {
	eibctypes.BaseEIBCHook
	Keeper
}

func (k Keeper) GetEIBCHooks() eibctypes.EIBCHooks {
	return eibcHooks{
		BaseEIBCHook: eibctypes.BaseEIBCHook{},
		Keeper:       k,
	}
}

// AfterDemandOrderFulfilled is called every time a demand order is fulfilled.
// Once it is fulfilled the underlying packet recipient should be updated to the fulfiller.
func (k eibcHooks) AfterDemandOrderFulfilled(ctx sdk.Context, o *eibctypes.DemandOrder, receiverAddr string) error {
	if o.CompletionHook != nil {
		err := k.RunOrderCompletionHook(ctx, o, o.PriceAmount())
		if err != nil {
			return errorsmod.Wrap(err, "run completion hook")
		}
	}

	err := k.UpdateRollappPacketTransferAddress(ctx, o.TrackingPacketKey, receiverAddr)
	if err != nil {
		return err
	}
	return nil
}

/* -------------------------------------------------------------------------- */
/*                                 epoch hooks                                */
/* -------------------------------------------------------------------------- */
var _ epochstypes.EpochHooks = epochHooks{}

type epochHooks struct {
	Keeper
}

func (k Keeper) GetEpochHooks() epochstypes.EpochHooks {
	return epochHooks{
		Keeper: k,
	}
}

// BeforeEpochStart is the epoch start hook.
func (e epochHooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return nil
}

// AfterEpochEnd is the epoch end hook.
// We want to clean up the demand orders that are with underlying packet status which are finalized.
func (e epochHooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	params := e.GetParams(ctx)

	if epochIdentifier != params.EpochIdentifier {
		return nil
	}

	listFilter := types.ByStatus(commontypes.Status_FINALIZED).Take(int(deletePacketsBatchSize))
	count := 0

	// Get batch of rollapp packets with status != PENDING and delete them
	for toDeletePackets := e.ListRollappPackets(ctx, listFilter); len(toDeletePackets) > 0; toDeletePackets = e.ListRollappPackets(ctx, listFilter) {
		e.Logger(ctx).Debug("Deleting rollapp packets", "num_packets", len(toDeletePackets))

		count += len(toDeletePackets)

		for _, packet := range toDeletePackets {
			e.DeleteRollappPacket(ctx, &packet)
		}

		// if the total number of deleted packets reaches the hard limit for the epoch, stop deleting packets
		if count >= int(params.DeletePacketsEpochLimit) {
			break
		}
	}
	return nil
}
