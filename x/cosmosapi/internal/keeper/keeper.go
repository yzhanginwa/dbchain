package keeper

import (
	"os"
	"fmt"
	"errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/yzhanginwa/rcv-chain/x/rcvchain/internal/types"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	logger = defaultLogger()
)

func defaultLogger() log.Logger {
        return log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("ethan1", "ethan2")
}

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	CoinKeeper bank.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}


// NewKeeper creates new instances of the rcvchain Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		CoinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// Check if the poll id is present in the store or not
func (k Keeper) IsPollPresent(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(types.PollKey(id)))
}


// Create a new poll
func (k Keeper) CreatePoll(ctx sdk.Context, title string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	var poll types.Poll = types.NewPoll()
        poll.Id = types.NewPollId(title, owner)
        poll.Title = title
	poll.Owner = owner
	//logger.Info(poll.Title)
	store.Set([]byte(types.PollKey(poll.Id)), k.cdc.MustMarshalBinaryBare(poll))
}

// Save or update UserPoll
func (k Keeper) SaveUserPolls(ctx sdk.Context, id string, addr sdk.AccAddress, isOwner bool) {
	store := ctx.KVStore(k.storeKey)
	var userPoll types.UserPoll
	key := types.UserPollKey(addr)
	bz := store.Get([]byte(key))
	if bz == nil {
		if isOwner {
			userPoll = types.NewUserPoll(addr, []string{id}, []string{})
		} else {
			userPoll = types.NewUserPoll(addr, []string{}, []string{id})
		}
	} else {
		k.cdc.MustUnmarshalBinaryBare(bz, &userPoll)
		if isOwner {
  			for _, item := range userPoll.OwnedPolls {
				if id == item {
					return
				}
			}
			userPoll.OwnedPolls = append(userPoll.OwnedPolls, id)

		} else {
  			for _, item := range userPoll.InvitedPolls {
				if id == item {
					return
				}
			}
			userPoll.InvitedPolls = append(userPoll.InvitedPolls, id)
		}
	}
	store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(userPoll))
}

// Gets a user poll for a address
func (k Keeper) GetUserPolls(ctx sdk.Context, addr sdk.AccAddress) (types.UserPoll, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.UserPollKey(addr)))
	if bz == nil {
		return types.NewUserPoll(addr, nil, nil), nil
	}
	var userPoll types.UserPoll
	k.cdc.MustUnmarshalBinaryBare(bz, &userPoll)
	return userPoll, nil
}

// Gets a poll for an id
func (k Keeper) GetPoll(ctx sdk.Context, id string) (types.Poll, error) {
	store := ctx.KVStore(k.storeKey)
        // TODO: check if bz is nil
	bz := store.Get([]byte(types.PollKey(id)))
	if bz == nil {
		return types.Poll{}, errors.New("not found poll")
	}
	var poll types.Poll
	k.cdc.MustUnmarshalBinaryBare(bz, &poll)
	return poll, nil
}

// Get poll status
func (k Keeper) GetPollStatus(ctx sdk.Context, id string) string {
    poll, err := k.GetPoll(ctx, id)
    if err != nil {
	return ""
     }
    return poll.Status
}

// Get poll title
func (k Keeper) GetPollTitle(ctx sdk.Context, id string) string {
    poll, err := k.GetPoll(ctx, id)
    if err != nil {
	return ""
     }
    return poll.Title
}

// Get an iterator over all polls in which the keys are the poll:title and the values are the poll
func (k Keeper) GetPollsIterator(ctx sdk.Context) sdk.Iterator {
        store := ctx.KVStore(k.storeKey)
        return sdk.KVStorePrefixIterator(store, nil)
}

// Add a choice
func (k Keeper) AddChoice(ctx sdk.Context, id string, choice string, owner sdk.AccAddress) {
	poll, err := k.GetPoll(ctx, id)
	if err != nil {
		return
	}
	store := ctx.KVStore(k.storeKey)

	choices := poll.Choices

	for _, item := range choices {
		if choice == item {
			return
		}
	}
        poll.Choices = append(poll.Choices, choice)
	store.Set([]byte(types.PollKey(id)), k.cdc.MustMarshalBinaryBare(poll))
}

// Invite a voter
func (k Keeper) InviteVoter(ctx sdk.Context, id string, voter sdk.AccAddress, owner sdk.AccAddress) {
	poll, err := k.GetPoll(ctx, id)
	if err != nil {
		return
	}
	voters := poll.Voters

	for _, item := range voters {
		if voter.Equals(item) {
			return
		}
	}
	poll.Voters = append(voters, voter)

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.PollKey(id)), k.cdc.MustMarshalBinaryBare(poll))
}

// Begin voting
func (k Keeper) BeginVoting(ctx sdk.Context, id string, owner sdk.AccAddress) {
	poll, err := k.GetPoll(ctx, id)
	if err != nil { return }
	if !owner.Equals(poll.Owner) { return }

	poll.Status = types.PollStatusReady

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.PollKey(id)), k.cdc.MustMarshalBinaryBare(poll))
}

// End voting
func (k Keeper) EndVoting(ctx sdk.Context, id string, owner sdk.AccAddress) {
	poll, err := k.GetPoll(ctx, id)
	if err != nil { return }
	if !owner.Equals(poll.Owner) { return }

	poll.Status = types.PollStatusFinished
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.PollKey(id)), k.cdc.MustMarshalBinaryBare(poll))

        k.calculateAndSaveWinner(ctx, poll)
}

// To vote
func (k Keeper) CreateBallot(ctx sdk.Context, id string, votes []uint16, voter sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.BallotKey(id, voter)

        if nil != store.Get([]byte(key)) {
                return
        }

	poll, err := k.GetPoll(ctx, id)
	if err != nil {
		return
	}

	choices := poll.Choices
	choicesLen := uint16(len(choices))
	for _, choice := range votes {
		if choice >= choicesLen {
			return
		}
	}

	value := types.NewBallot(votes)
	store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(value))

	k.tryToFinishPoll(ctx, poll)
}

// Get ballot
func (k Keeper) GetBallot(ctx sdk.Context, id string, voter sdk.AccAddress) ([]string, error) {
	poll, err := k.GetPoll(ctx, id)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("not valid key %s", id))
	}

	store := ctx.KVStore(k.storeKey)
	key := types.BallotKey(poll.Id, voter)
        bz := store.Get([]byte(key))
        if bz == nil {
                return nil, nil
        }
	var ballot types.Ballot
	k.cdc.MustUnmarshalBinaryBare(bz, &ballot)

	var result []string
	for _, choice := range ballot.ChoiceIndex {
		result = append(result, poll.Choices[choice])
	}

	return result, nil
}

func (k Keeper) tryToFinishPoll(ctx sdk.Context, poll types.Poll) {
	store := ctx.KVStore(k.storeKey)

	for _, voter := range poll.Voters {
		key := types.BallotKey(poll.Id, voter)
		bz := store.Get([]byte(key))
		if bz == nil {
			return
		}
	}

	k.EndVoting(ctx, poll.Id, poll.Owner)
}

// Calculate and save the winder
func (k Keeper) calculateAndSaveWinner(ctx sdk.Context, poll types.Poll) {
	voters := poll.Voters

        var ballots [][]uint16
	store := ctx.KVStore(k.storeKey)

	for _, voter := range voters {
		ballotKey := types.BallotKey(poll.Id, voter)
		bz := store.Get([]byte(ballotKey))
		if bz == nil { continue }
	        var ballot types.Ballot
		k.cdc.MustUnmarshalBinaryBare(bz, &ballot)
		ballots = append(ballots, ballot.ChoiceIndex)
	}

	if len(ballots) == 0 {
		// TODO: mark the poll status to "failed"
		return
	}

	winner, _ := rcvRound(ballots)

	poll.Winner = poll.Choices[winner]
	store.Set([]byte(types.PollKey(poll.Id)), k.cdc.MustMarshalBinaryBare(poll))
}

//////////////////////
//                  //
// helper functions //
//                  //
//////////////////////

func rcvRound(matrix [][]uint16) (uint16, bool) {
        winner, loser, ok := topWinnerOrLoser(matrix)
        if ok {
                return winner, true
        } else {
                popLoser(matrix, loser)
                return rcvRound(matrix)
        }
}

func topWinnerOrLoser(matrix [][]uint16) (uint16, uint16, bool) {
        var count = make(map[uint16]int)
        rows := len(matrix)

        for _, list := range matrix {
                if len(list) > 0 {
                        u16 := list[0]
                        _, e := count[u16]
                        if e {
                                count[u16] += 1
                        } else {
                                count[u16] = 1
                        }

                        if count[u16] * 2 > rows {
                                return u16, 0, true
                        }
                }
        }

        var minKey uint16
        var minValue int = rows       // initialize the lowest candidate to the largest possible number

        for k, v := range count {
                if v < int(minValue) {
                        minKey = k
                        minValue = v
                }
        }
        return 0, minKey, false
}

func popLoser(matrix[][]uint16, loser uint16) {
        for index, list := range matrix {
                if len(list) > 0 {
                        if list[0] == loser {
                                matrix[index] = list[1:]
                        }
                }
        }
}

