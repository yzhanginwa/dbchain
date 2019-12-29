package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreatePoll{}, "rcvchain/CreatePoll", nil)
	cdc.RegisterConcrete(MsgAddChoice{}, "rcvchain/AddChoice", nil)
	cdc.RegisterConcrete(MsgInviteVoter{}, "rcvchain/InviteVoter", nil)
	cdc.RegisterConcrete(MsgBeginVoting{}, "rcvchain/BeginVoting", nil)
	cdc.RegisterConcrete(MsgCreateBallot{}, "rcvchain/CreateBallot", nil)
	cdc.RegisterConcrete(MsgEndVoting{}, "rcvchain/EndVoting", nil)
}

