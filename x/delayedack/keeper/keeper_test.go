package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/dymensionxyz/dymension/v3/app/apptesting"
	commontypes "github.com/dymensionxyz/dymension/v3/x/common/types"
	"github.com/dymensionxyz/dymension/v3/x/delayedack/types"
)

const delayedAckEventType = "delayedack"

type DelayedAckTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(DelayedAckTestSuite))
}

func (suite *DelayedAckTestSuite) SetupTest() {
	app := apptesting.Setup(suite.T())
	ctx := app.NewContext(false)

	suite.App = app
	suite.Ctx = ctx
}

func (suite *DelayedAckTestSuite) FinalizePacket(ctx sdk.Context, p commontypes.RollappPacket) {
	suite.T().Helper()

	handler := suite.App.MsgServiceRouter().Handler(new(types.MsgFinalizePacket))
	resp, err := handler(ctx, &types.MsgFinalizePacket{
		Sender:            apptesting.CreateRandomAccounts(1)[0].String(),
		RollappId:         p.RollappId,
		PacketProofHeight: p.ProofHeight,
		PacketType:        p.Type,
		PacketSrcChannel:  p.Packet.SourceChannel,
		PacketSequence:    p.Packet.Sequence,
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
}
