package cli

import (
	"time"
)

// Blocks NewChain Blocks
type Blocks struct {
	Number           string    `db:"block_number"`
	Hash             string    `db:"block_hash"`
	ParentHash       string    `db:"block_parent_hash"`
	Nonce            uint64    `db:"block_nonce"`
	Sha3Uncles       string    `db:"block_sha3_uncles"`
	TransactionsRoot string    `db:"block_transactions_root"`
	StateRoot        string    `db:"block_state_root"`
	ReceiptsRoot     string    `db:"block_receipts_root"`
	Miner            string    `db:"block_miner"`
	Difficulty       string    `db:"block_difficulty"`
	TotalDifficulty  string    `db:"block_total_difficulty"`
	Size             string    `db:"block_size"`
	GasLimit         uint64    `db:"block_gas_limit"`
	GasUsed          uint64    `db:"block_gas_used"`
	Timestamp        time.Time `db:"block_timestamp"`
	TransactionCount int       `db:"block_transaction_count"`
	Signer           string    `db:"block_signer"`
}

// Transactions NewChain Transactions
type Transactions struct {
	Hash        string `db:"tx_hash"`
	BlockNumber string `db:"tx_block_number"`
	Nonce       uint64 `db:"tx_nonce"`
	From        string `db:"tx_from"`
	To          string `db:"tx_to"`
	Value       string `db:"tx_value"`
	Gas         uint64 `db:"tx_gas"`
	GasPrice    string `db:"tx_gas_price"`
}
