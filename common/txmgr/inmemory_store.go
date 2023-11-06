package txmgr

import (
	"context"
	"fmt"
	"sync"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/label"
	"gopkg.in/guregu/null.v4"
)

var (
	// ErrInvalidChainID is returned when the chain ID is invalid
	ErrInvalidChainID = fmt.Errorf("invalid chain ID")
	// ErrTxnNotFound is returned when a transaction is not found
	ErrTxnNotFound = fmt.Errorf("transaction not found")
	// ErrExistingIdempotencyKey is returned when a transaction with the same idempotency key already exists
	ErrExistingIdempotencyKey = fmt.Errorf("transaction with idempotency key already exists")
	// ErrAddressNotFound is returned when an address is not found
	ErrAddressNotFound = fmt.Errorf("address not found")
)

// Store and update all transaction state as files
// Read from the files to restore state at startup
// Delete files when transactions are completed or reaped

// Life of a Transaction
// 1. Transaction Request is created
// 2. Transaction Request is submitted to the Transaction Manager
// 3. Transaction Manager creates and persists a new transaction (unstarted) from the transaction request (not persisted)
// 4. Transaction Manager sends the transaction (unstarted) to the Broadcaster Unstarted Queue
// 4. Transaction Manager prunes the Unstarted Queue based on the transaction prune strategy

// NOTE(jtw): Only one transaction per address can be in_progress at a time
// NOTE(jtw): Only one broadcasted attempt exists per transaction the rest are errored or abandoned
// 1. Broadcaster assigns a sequence number to the transaction
// 2. Broadcaster creates and persists a new transaction attempt (in_progress) from the transaction (in_progress)
// 3. Broadcaster asks the Checker to check if the transaction should not be sent
// 4. Broadcaster asks the Attempt builder to figure out gas fee for the transaction
// 5. Broadcaster attempts to send the Transaction to TransactionClient to be published on-chain
// 6. Broadcaster updates the transaction attempt (broadcast) and transaction (unconfirmed)
// 7. Broadcaster increments global sequence number for address for next transaction attempt

// NOTE(jtw): Only one receipt should exist per confirmed transaction
// 1. Confirmer listens and reads new Head events from the Chain
// 2. Confirmer sets the last known block number for the transaction attempts that have been broadcast
// 3. Confirmer checks for missing receipts for transactions that have been broadcast
// 4. Confirmer sets transactions that have failed to (unconfirmed) which will be retried by the resender
// 5. Confirmer sets transactions that have been confirmed to (confirmed) and creates a new receipt which is persisted

type InMemoryStore[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	chainID CHAIN_ID

	keyStore txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ]
	txStore  txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE]

	pendingLock sync.Mutex
	// NOTE(jtw): we might need to watch out for txns that finish and are removed from the pending map
	pendingIdempotencyKeys map[string]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]

	unstarted  map[ADDR]chan *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	inprogress map[ADDR]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
}

// NewInMemoryStore returns a new InMemoryStore
func NewInMemoryStore[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
](
	chainID CHAIN_ID,
	keyStore txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ],
	txStore txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
) (*InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], error) {
	tm := InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		chainID:  chainID,
		keyStore: keyStore,
		txStore:  txStore,

		pendingIdempotencyKeys: map[string]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{},

		unstarted:  map[ADDR]chan *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{},
		inprogress: map[ADDR]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{},
	}

	addresses, err := keyStore.EnabledAddressesForChain(chainID)
	if err != nil {
		return nil, fmt.Errorf("new_in_memory_store: %w", err)
	}
	for _, fromAddr := range addresses {
		// Channel Buffer is set to something high to prevent blocking and allow the pruning to happen
		tm.unstarted[fromAddr] = make(chan *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], 100)
	}

	return &tm, nil
}

// CreateTransaction creates a new transaction for a given txRequest.
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CreateTransaction(ctx context.Context, txRequest txmgrtypes.TxRequest[ADDR, TX_HASH], chainID CHAIN_ID) (txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	// TODO(jtw): do generic checks

	// Persist Transaction to persistent storage
	tx, err := ms.txStore.CreateTransaction(ctx, txRequest, chainID)
	if err != nil {
		return tx, fmt.Errorf("create_transaction: %w", err)
	}
	if err := ms.sendTxToBroadcaster(tx); err != nil {
		return tx, fmt.Errorf("create_transaction: %w", err)
	}

	return tx, nil
}

// FindTxWithIdempotencyKey returns a transaction with the given idempotency key
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxWithIdempotencyKey(ctx context.Context, idempotencyKey string, chainID CHAIN_ID) (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	// TODO(jtw): this is a change from current functionality... it returns nil, nil if nothing found in other implementations
	if ms.chainID.String() != chainID.String() {
		return nil, fmt.Errorf("find_tx_with_idempotency_key: %w", ErrInvalidChainID)
	}
	if idempotencyKey == "" {
		return nil, fmt.Errorf("find_tx_with_idempotency_key: idempotency key cannot be empty")
	}

	ms.pendingLock.Lock()
	defer ms.pendingLock.Unlock()

	tx, ok := ms.pendingIdempotencyKeys[idempotencyKey]
	if !ok {
		return nil, fmt.Errorf("find_tx_with_idempotency_key: %w", ErrTxnNotFound)
	}

	return tx, nil
}

// CheckTxQueueCapacity checks if the queue capacity has been reached for a given address
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CheckTxQueueCapacity(ctx context.Context, fromAddress ADDR, maxQueuedTransactions uint64, chainID CHAIN_ID) error {
	if maxQueuedTransactions == 0 {
		return nil
	}
	if ms.chainID.String() != chainID.String() {
		return fmt.Errorf("check_tx_queue_capacity: %w", ErrInvalidChainID)
	}
	if _, ok := ms.unstarted[fromAddress]; !ok {
		return fmt.Errorf("check_tx_queue_capacity: %w", ErrAddressNotFound)
	}

	count := uint64(len(ms.unstarted[fromAddress]))
	if count >= maxQueuedTransactions {
		return fmt.Errorf("check_tx_queue_capacity: cannot create transaction; too many unstarted transactions in the queue (%v/%v). %s", count, maxQueuedTransactions, label.MaxQueuedTransactionsWarning)
	}

	return nil
}

/////
// BROADCASTER FUNCTIONS
/////

// FindLatestSequence returns the latest sequence number for a given address
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindLatestSequence(ctx context.Context, fromAddress ADDR, chainID CHAIN_ID) (SEQ, error) {
	// query the persistent storage since this method only gets called when the broadcaster is starting up.
	// It is used to initialize the in-memory sequence map in the broadcaster
	// NOTE(jtw): should the nextSequenceMap be moved to the in-memory store?

	// TODO(jtw): do generic checks
	// TODO(jtw): maybe this should be handled now
	seq, err := ms.txStore.FindLatestSequence(ctx, fromAddress, chainID)
	if err != nil {
		return seq, fmt.Errorf("find_latest_sequence: %w", err)
	}

	return seq, nil
}

// CountUnconfirmedTransactions returns the number of unconfirmed transactions for a given address.
// Unconfirmed transactions are transactions that have been broadcast but not confirmed on-chain.
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CountUnconfirmedTransactions(ctx context.Context, fromAddress ADDR, chainID CHAIN_ID) (uint32, error) {
	// NOTE(jtw): used to calculate total inflight transactions
	if ms.chainID.String() != chainID.String() {
		return 0, fmt.Errorf("count_unstarted_transactions: %w", ErrInvalidChainID)
	}
	if _, ok := ms.unstarted[fromAddress]; !ok {
		return 0, fmt.Errorf("count_unstarted_transactions: %w", ErrAddressNotFound)
	}

	// TODO(jtw): NEEDS TO BE UPDATED TO USE IN MEMORY STORE
	return ms.txStore.CountUnconfirmedTransactions(ctx, fromAddress, chainID)
}

// CountUnstartedTransactions returns the number of unstarted transactions for a given address.
// Unstarted transactions are transactions that have not been broadcast yet.
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CountUnstartedTransactions(ctx context.Context, fromAddress ADDR, chainID CHAIN_ID) (uint32, error) {
	// NOTE(jtw): used to calculate total inflight transactions
	if ms.chainID.String() != chainID.String() {
		return 0, fmt.Errorf("count_unstarted_transactions: %w", ErrInvalidChainID)
	}
	if _, ok := ms.unstarted[fromAddress]; !ok {
		return 0, fmt.Errorf("count_unstarted_transactions: %w", ErrAddressNotFound)
	}

	return uint32(len(ms.unstarted[fromAddress])), nil
}

// UpdateTxUnstartedToInProgress updates a transaction from unstarted to in_progress.
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) UpdateTxUnstartedToInProgress(
	ctx context.Context,
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	attempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) error {
	if tx.Sequence == nil {
		return fmt.Errorf("update_tx_unstarted_to_in_progress: in_progress transaction must have a sequence number")
	}
	if tx.State != TxUnstarted {
		return fmt.Errorf("update_tx_unstarted_to_in_progress: can only transition to in_progress from unstarted, transaction is currently %s", tx.State)
	}
	if attempt.State != txmgrtypes.TxAttemptInProgress {
		return fmt.Errorf("update_tx_unstarted_to_in_progress: attempt state must be in_progress")
	}

	// Persist to persistent storage
	if err := ms.txStore.UpdateTxUnstartedToInProgress(ctx, tx, attempt); err != nil {
		return fmt.Errorf("update_tx_unstarted_to_in_progress: %w", err)
	}
	tx.TxAttempts = append(tx.TxAttempts, *attempt)

	// Update in memory store
	ms.inprogress[tx.FromAddress] = tx

	return nil
}

// GetTxInProgress returns the in_progress transaction for a given address.
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) GetTxInProgress(ctx context.Context, fromAddress ADDR) (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	tx, ok := ms.inprogress[fromAddress]
	if !ok {
		return nil, nil
	}
	if len(tx.TxAttempts) != 1 || tx.TxAttempts[0].State != txmgrtypes.TxAttemptInProgress {
		return nil, fmt.Errorf("get_tx_in_progress: expected in_progress transaction %v to have exactly one unsent attempt. "+
			"Your database is in an inconsistent state and this node will not function correctly until the problem is resolved", tx.ID)
	}

	return tx, nil
}

// UpdateTxAttemptInProgressToBroadcast updates a transaction attempt from in_progress to broadcast.
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) UpdateTxAttemptInProgressToBroadcast(
	ctx context.Context,
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	attempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	newAttemptState txmgrtypes.TxAttemptState,
) error {
	// TODO(jtw)
	return fmt.Errorf("update_tx_attempt_in_progress_to_broadcast: not implemented")
}

// FindNextUnstartedTransactionFromAddress returns the next unstarted transaction for a given address.
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindNextUnstartedTransactionFromAddress(ctx context.Context, tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], fromAddress ADDR, chainID CHAIN_ID) error {
	if ms.chainID.String() != chainID.String() {
		return fmt.Errorf("find_next_unstarted_transaction_from_address: %w", ErrInvalidChainID)
	}
	if _, ok := ms.unstarted[fromAddress]; !ok {
		return fmt.Errorf("find_next_unstarted_transaction_from_address: %w", ErrAddressNotFound)
	}

	return fmt.Errorf("find_next_unstarted_transaction_from_address: not implemented")
}

// SaveReplacementInProgressAttempt saves a replacement attempt for a transaction that is in_progress.
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SaveReplacementInProgressAttempt(
	ctx context.Context,
	oldAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	replacementAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) error {
	// TODO(jtw)
	return fmt.Errorf("save_replacement_in_progress_attempt: not implemented")
}

// UpdateTxFatalError updates a transaction to fatal_error.
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) UpdateTxFatalError(ctx context.Context, tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	// TODO(jtw)
	return fmt.Errorf("update_tx_fatal_error: not implemented")
}

// Close closes the InMemoryStore
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Close() {
	// Close the event recorder
	ms.txStore.Close()

	// Close all channels
	for _, ch := range ms.unstarted {
		close(ch)
	}

	// Clear all pending requests
	ms.pendingLock.Lock()
	clear(ms.pendingIdempotencyKeys)
	ms.pendingLock.Unlock()
}

// Abandon removes all transactions for a given address
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Abandon(ctx context.Context, chainID CHAIN_ID, addr ADDR) error {
	if ms.chainID.String() != chainID.String() {
		return fmt.Errorf("abandon: %w", ErrInvalidChainID)
	}

	// TODO(jtw): do generic checks
	// Mark all persisted transactions as abandoned
	if err := ms.txStore.Abandon(ctx, chainID, addr); err != nil {
		return err
	}

	// Mark all unstarted transactions as abandoned
	close(ms.unstarted[addr])
	for tx := range ms.unstarted[addr] {
		tx.State = TxFatalError
		tx.Sequence = nil
		tx.Error = null.NewString("abandoned", true)
	}
	// reset the unstarted channel
	ms.unstarted[addr] = make(chan *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], 100)

	// Mark all inprogress transactions as abandoned
	if tx, ok := ms.inprogress[addr]; ok {
		tx.State = TxFatalError
		tx.Sequence = nil
		tx.Error = null.NewString("abandoned", true)
	}
	ms.inprogress[addr] = nil

	// TODO(jtw): Mark all unconfirmed transactions as abandoned

	// Mark all pending transactions as abandoned
	for _, tx := range ms.pendingIdempotencyKeys {
		if tx.FromAddress == addr {
			tx.State = TxFatalError
			tx.Sequence = nil
			tx.Error = null.NewString("abandoned", true)
		}
	}
	// TODO(jtw): SHOULD THE REAPER BE RESPONSIBLE FOR CLEARING THE PENDING MAPS?

	return nil
}

// TODO(jtw): change naming to something more appropriate
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) sendTxToBroadcaster(tx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	// TODO(jtw); HANDLE PRUNING STEP

	select {
	// Add the request to the Unstarted channel to be processed by the Broadcaster
	case ms.unstarted[tx.FromAddress] <- &tx:
		// Persist to persistent storage

		ms.pendingLock.Lock()
		if tx.IdempotencyKey != nil {
			ms.pendingIdempotencyKeys[*tx.IdempotencyKey] = &tx
		}
		ms.pendingLock.Unlock()

		return nil
	default:
		// Return an error if the Manager Queue Capacity has been reached
		return fmt.Errorf("transaction manager queue capacity has been reached")
	}
}