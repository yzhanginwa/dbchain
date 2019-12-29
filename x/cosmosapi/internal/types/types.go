package types

import (
	"fmt"
	"strings"
	"crypto/sha256"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	PollStatusNew = "new"
	PollStatusReady = "ready"
	PollStatusFinished = "finished"

	PollKeyPrefix = "poll:"
)


//////////
//      //
// Poll //
//      //
//////////

// the key would be like "poll:[name]"
type Poll struct {
        Id string                 `json:"id"`
	Title string              `json:"title"`
	Owner sdk.AccAddress      `json:"owner"`
        Status string             `json:"status"`
        Choices []string          `json:"choices"`
        Voters  []sdk.AccAddress  `json:"voters"`
        Winner string             `json:"winner"`
}

func NewPollId(title string, owner sdk.AccAddress) string {
	message := fmt.Sprintf("%s%s", title, owner)
	digest := sha256.Sum256([]byte(message))
        return fmt.Sprintf("%x", digest)
}

func PollKey(id string) string {
	return fmt.Sprintf("%s%s", PollKeyPrefix, id)
}

func NewPoll() Poll {
	return Poll {
		Status: "new",
	}
}

// implement fmt.Stringer
func (p Poll) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Title: %s`, p.Title))
}

////////////
//        //
// Ballot //
//        //
////////////

//the key would be like "voting:[poll_id]:[address]"
type Ballot struct {
	ChoiceIndex []uint16       `json:"choiceindex"`
	Voter       sdk.AccAddress `json:"voter"`
}

func BallotKey(pollId string, address sdk.AccAddress) string {
	return fmt.Sprintf("ballot:%s:%s", pollId, address)
}

func NewBallot(ballot []uint16) Ballot {
	return Ballot {
		ChoiceIndex: ballot,
	}
}

///////////////
//           //
// User-Poll //
//           //
///////////////

type UserPoll struct {
	Address sdk.AccAddress   `json:"address"`
	OwnedPolls []string      `json:"owned_polls"`
        InvitedPolls []string    `json:"invited_polls"`
}

func UserPollKey(address sdk.AccAddress) string {
	return fmt.Sprintf("user:%s", address)

}

func NewUserPoll(addr sdk.AccAddress, owned []string, invited []string) UserPoll {
	return UserPoll {
		Address: addr,
		OwnedPolls: owned,
		InvitedPolls: invited,
	}
}

