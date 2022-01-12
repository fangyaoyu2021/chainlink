package terratxm

import (
	"time"

	"github.com/smartcontractkit/terra.go/msg"
)

// State represents the state of a given terra msg
type State string

var (
	// Unstarted means queued but not processed
	Unstarted State = "unstarted"
	// Broadcasted means included in the mempool of a node
	Broadcasted State = "broadcasted"
	// Confirmed means we're able to retrieve the txhash of the tx which broadcasted the msg
	Confirmed State = "confirmed"
	// Errored means the msg reverted in simulation OR the tx containing the message timed out waiting to be confirmed
	// TODO: when we add gas bumping, we'll address that timeout case
	Errored State = "errored"
)

// TerraMsg a terra msg
type TerraMsg struct {
	ID         int64
	ContractID string
	State      State
	Msg        []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// In memory only
	ExecuteContract *msg.ExecuteContract
}

func getMsgs(tms []TerraMsg) []msg.Msg {
	var msgs []msg.Msg
	for i := range tms {
		msgs = append(msgs, tms[i].ExecuteContract)
	}
	return msgs
}

func getIDs(tms []TerraMsg) []int64 {
	var ids []int64
	for i := range tms {
		ids = append(ids, tms[i].ID)
	}
	return ids
}
