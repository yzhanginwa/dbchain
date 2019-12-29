package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is the module name router key
const RouterKey = ModuleName // this was defined in your key.go file

///////////////////
//               //
// MsgCreatePoll //
//               //
///////////////////

// MsgCreatePoll defines a CreatePoll message
type MsgCreatePoll struct {
	Title string         `json:"title"`
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgCreatePoll is a constructor function for MsgCreatPoll
func NewMsgCreatePoll(title string, owner sdk.AccAddress) MsgCreatePoll {
	return MsgCreatePoll {
		Title: title,
		Owner: owner,
	}
}

// Route should return the name of the module
func (msg MsgCreatePoll) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreatePoll) Type() string { return "create_poll" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreatePoll) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Title) == 0 {
		return sdk.ErrUnknownRequest("Title cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreatePoll) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreatePoll) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

//////////////////
//              //
// MsgAddChoice //
//              //
//////////////////

// MsgAddChoice defines the AddChoice message
type MsgAddChoice struct {
	Id     string        `json:"id"`
	Choice string        `json:"choice"`
	Owner sdk.AccAddress `json:"owner"`
}

func NewMsgAddChoice(id string, choice string, owner sdk.AccAddress) MsgAddChoice {
	return MsgAddChoice {
		Id:     id,
		Choice: choice,
		Owner:  owner,
	}
}

// Route should return the name of the module
func (msg MsgAddChoice) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddChoice) Type() string { return "add_choice" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddChoice) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Id) == 0 {
		return sdk.ErrUnknownRequest("Poll cannot be empty")
	}
	if len(msg.Choice) == 0 {
		return sdk.ErrUnknownRequest("Choice cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAddChoice) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddChoice) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

////////////////////
//                //
// MsgInviteVoter //
//                //
////////////////////

// MsgInviteVoter defines the InviteVoter message
type MsgInviteVoter struct {
	Id    string         `json:"id"`
	Voter sdk.AccAddress `json:"voter"`
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgInviteVoter is the constructor function for MsgInviteVoter
func NewMsgInviteVoter(id string, voter sdk.AccAddress, owner sdk.AccAddress) MsgInviteVoter {
	return MsgInviteVoter{
		Id:    id,
		Voter: voter,
		Owner: owner,
	}
}

// Route should return the name of the module
func (msg MsgInviteVoter) Route() string { return RouterKey }

// Type should return the action
func (msg MsgInviteVoter) Type() string { return "invite_voter" }

// ValidateBasic runs stateless checks on the message
func (msg MsgInviteVoter) ValidateBasic() sdk.Error {
	if msg.Voter.Empty() {
		return sdk.ErrInvalidAddress(msg.Voter.String())
	}
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Id) == 0 {
		return sdk.ErrUnknownRequest("Id cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgInviteVoter) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgInviteVoter) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

////////////////////
//                //
// MsgBeginVoting //
//                //
////////////////////

// MsgBeginVoting defines a BeginVoting message
type MsgBeginVoting struct {
	Id    string         `json:"id"`
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgBeginVoting is a constructor function for MsgBeginVoting
func NewMsgBeginVoting(id string, owner sdk.AccAddress) MsgBeginVoting {
	return MsgBeginVoting {
		Id:    id,
		Owner: owner,
	}
}

// Route should return the name of the module
func (msg MsgBeginVoting) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBeginVoting) Type() string { return "begin_voting" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBeginVoting) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Id) == 0 {
		return sdk.ErrUnknownRequest("Id cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgBeginVoting) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgBeginVoting) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

/////////////////////
//                 //
// MsgCreateBallot //
//                 //
/////////////////////


// MsgEndVoting defines a EndVoting message
type MsgCreateBallot struct {
	Id    string         `json:"id"`
	Votes []uint16       `json:"votes"`
	Voter sdk.AccAddress `json:"voter"`
}

// NewMsgCreateBallot is a constructor function for MsgCreateBallot
func NewMsgCreateBallot(id string, votes []uint16, voter sdk.AccAddress) MsgCreateBallot {
	return MsgCreateBallot{
		Id:    id,
		Votes: votes,
		Voter: voter,
	}
}

// Route should return the name of the module
func (msg MsgCreateBallot) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateBallot) Type() string { return "create_ballot" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateBallot) ValidateBasic() sdk.Error {
	if msg.Voter.Empty() {
		return sdk.ErrInvalidAddress(msg.Voter.String())
	}
	if len(msg.Id) == 0 {
		return sdk.ErrUnknownRequest("Id cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateBallot) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateBallot) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Voter}
}

//////////////////
//              //
// MsgEndVoting //
//              //
//////////////////

// MsgEndVoting defines a EndVoting message
type MsgEndVoting struct {
	Id    string         `json:"id"`
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgEndVoting is a constructor function for MsgEndVoting
func NewMsgEndVoting(id string, owner sdk.AccAddress) MsgEndVoting {
	return MsgEndVoting{
		Id:    id,
		Owner: owner,
	}
}

// Route should return the name of the module
func (msg MsgEndVoting) Route() string { return RouterKey }

// Type should return the action
func (msg MsgEndVoting) Type() string { return "end_voting" }

// ValidateBasic runs stateless checks on the message
func (msg MsgEndVoting) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Id) == 0 {
		return sdk.ErrUnknownRequest("Id cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgEndVoting) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgEndVoting) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

