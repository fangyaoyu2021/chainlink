package core

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
)

// UpkeepTriggerID returns the identifier using the given upkeepID and trigger.
// It follows the same logic as the contract, but performs it locally.
func UpkeepTriggerID(id *big.Int, trigger []byte) (string, error) {
	idBytes := id.Bytes()

	combined := append(idBytes, trigger...)

	triggerIDBytes := crypto.Keccak256(combined)
	triggerID := hex.EncodeToString(triggerIDBytes)

	return triggerID, nil
}

// UpkeepWorkID returns the identifier using the given upkeepID and trigger extension(tx hash and log index).
func UpkeepWorkID(id *big.Int, trigger ocr2keepers.Trigger) (string, error) {
	extensionBytes, err := json.Marshal(trigger.Extension)
	if err != nil {
		return "", err
	}

	combined := fmt.Sprintf("%s%s", id, extensionBytes)
	hash := md5.Sum([]byte(combined))
	return hex.EncodeToString(hash[:]), nil
}

func NewUpkeepPayload(uid *big.Int, tp int, trigger ocr2keepers.Trigger, checkData []byte) (ocr2keepers.UpkeepPayload, error) {
	// construct payload
	p := ocr2keepers.UpkeepPayload{
		Upkeep: ocr2keepers.ConfiguredUpkeep{
			ID:     ocr2keepers.UpkeepIdentifier(uid.Bytes()),
			Type:   tp,
			Config: struct{}{}, // empty struct by default
		},
		Trigger:   trigger,
		CheckData: checkData,
	}
	// set work id based on upkeep id and trigger
	p.WorkID, _ = UpkeepWorkID(uid, trigger)

	// manually convert trigger to triggerWrapper
	triggerW := triggerWrapper{
		BlockNum:  uint32(trigger.BlockNumber),
		BlockHash: common.HexToHash(trigger.BlockHash),
	}

	switch UpkeepType(tp) {
	case LogTrigger:
		trExt, ok := trigger.Extension.(LogTriggerExtension)
		if !ok {
			return ocr2keepers.UpkeepPayload{}, fmt.Errorf("unrecognized trigger extension data")
		}
		hex, err := common.ParseHexOrString(trExt.TxHash)
		if err != nil {
			return ocr2keepers.UpkeepPayload{}, fmt.Errorf("tx hash parse error: %w", err)
		}
		triggerW.TxHash = common.BytesToHash(hex[:])
		triggerW.LogIndex = uint32(trExt.LogIndex)
	default:
	}

	// get trigger in bytes
	triggerBytes, err := PackTrigger(uid, triggerW)
	if err != nil {
		return ocr2keepers.UpkeepPayload{}, fmt.Errorf("%w: failed to pack trigger", err)
	}

	// set trigger id based on upkeep id and trigger
	p.ID, _ = UpkeepTriggerID(uid, triggerBytes)

	// end
	return p, nil
}
