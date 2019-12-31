package types

import (
    "fmt"
    "strings"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

// QueryResPoll Queries Result Payload for a poll query
type QueryResPoll struct {
    Id string                 `json:"id"`
    Title string              `json:"title"`
    Owner sdk.AccAddress      `json:"owner"`
    Status string             `json:"status"`
    Choices []string          `json:"choices"`
    Voters  []sdk.AccAddress  `json:"voters"`
    Winner string             `json:"winner"`
}

// implement fmt.Stringer
func (r QueryResPoll) String() string {
    return strings.Join([]string{r.Title, r.Status}, "\n")
}

////////////////

// QueryResResolve Queries Result Payload for a resolve query
type QueryResStatus struct {
    Status string `json:"status"`
}

// implement fmt.Stringer
func (r QueryResStatus) String() string {
    return r.Status
}

////////////////

// QueryResTitles Queries Result Payload for a titles query
type QueryResTitles []string

// implement fmt.Stringer
func (n QueryResTitles) String() string {
    return strings.Join(n[:], "\n")
}

////////////////

type QueryResBallot []string

func (b QueryResBallot) String() string {
    return strings.Join(b, "\n")
}

////////////////

type QueryResUserPolls struct {
    Address sdk.AccAddress   `json:"address"`
    OwnedPolls []string      `json:"owned_polls"`
    InvitedPolls []string    `json:"invited_polls"`
}

func (up QueryResUserPolls) String() string {
    return fmt.Sprintf("UserPoll:%s", up.Address)
}

