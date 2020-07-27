
/* database */
create database newchaindb;
use newchaindb;

/* create user */
create user 'newchain'@'localhost' identified by 'password';
-- create user 'newchain'@'%' identified by 'password';

/* Grant */
grant all on newchaindb.* to 'newchain'@'localhost';
-- grant all on newchaindb.* to 'newchain'@'%';

/* https://github.com/ethereum/wiki/wiki/JSON-RPC#returns-26 */
/* blocks */
CREATE TABLE IF NOT EXISTS blocks (
    block_number BIGINT UNSIGNED NOT NULL, /* QUANTITY */
    block_hash CHAR(64), /* DATA(32) */
    block_parent_hash CHAR(64), /* DATA(32) */
    block_nonce BIGINT UNSIGNED, /* DATA(8), BIGINT UNSIGNED */
    block_sha3_uncles CHAR(64), /* DATA(32) */
    block_transactions_root CHAR(64), /* DATA(32) */
    block_state_root CHAR(64), /* DATA(32) */
    block_receipts_root CHAR(64), /* DATA(32) */ 
    block_miner CHAR(40), /* DATA(20) */
    block_difficulty BIGINT UNSIGNED, /* QUANTITY */
    block_total_difficulty BIGINT UNSIGNED, /* QUANTITY */
    block_size BIGINT UNSIGNED, /* QUANTITY */  
    block_gas_limit BIGINT UNSIGNED, /* QUANTITY */
    block_gas_used BIGINT UNSIGNED, /* QUANTITY */
    block_timestamp TIMESTAMP, /* QUANTITY from_unixtime */
    block_transaction_count BIGINT, /* len(transactions) */
    block_signer CHAR(40), /* DATA(20) */

    PRIMARY KEY  (`block_number`)
)engine=innodb row_format=compressed;


/* https://github.com/ethereum/wiki/wiki/JSON-RPC#returns-28 */
/* https://github.com/ethereum/wiki/wiki/JSON-RPC#parameters-22 */
/* transactions */
CREATE TABLE IF NOT EXISTS transactions (
    tx_hash CHAR(64) NOT NULL, /* DATA(32) + "0x" */
    tx_block_number BIGINT UNSIGNED, /* QUANTITY */
    tx_nonce BIGINT UNSIGNED, /* QUANTITY */
    tx_from CHAR(40), /* DATA(20) */
    tx_to  CHAR(40), /* DATA(20) */
    tx_value VARCHAR(32), /* QUANTITY */
    tx_gas BIGINT UNSIGNED, /* QUANTITY */
    tx_gas_price BIGINT UNSIGNED, /* QUANTITY */

    PRIMARY KEY  (`tx_hash`),
    FOREIGN key(tx_block_number) REFERENCES blocks(block_number)
)engine=innodb row_format=compressed;