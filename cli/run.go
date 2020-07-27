package cli

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/sha3"
	db "upper.io/db.v3"
	"upper.io/db.v3/mysql"
)

func (cli *CLI) buildRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "run [--start 0] [--end 100] [--loop] [--delay 5]",
		Short:                 "Get NewChain blocks and store in database",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			start, err := cmd.Flags().GetInt64("start")
			if err != nil {
				fmt.Println("Get Start Error: ", start)
				fmt.Println(cmd.UsageString())
				return
			}
			end, err := cmd.Flags().GetInt64("end")
			if err != nil {
				fmt.Println("Get end Error: ", start)
				fmt.Println(cmd.UsageString())
				return
			}
			startBlock := big.NewInt(start)
			endBlock := big.NewInt(end)
			loop, err := cmd.Flags().GetBool("loop")
			if err != nil {
				fmt.Println(err)
				fmt.Println(cmd.UsageString())
				return
			}

			if start == 0 && end == 0 && loop == false {
				fmt.Println("Specify at least one parameter: start, end or loop")
				fmt.Println(cmd.UsageString())

				return
			}

			blockDelay, err := cmd.Flags().GetInt64("delay")
			if err != nil {
				fmt.Println(err)
				fmt.Println(cmd.UsageString())
				return
			}

			dbtype := "mysql"
			var settings = mysql.ConnectionURL{
				Host:     viper.GetString(dbtype + ".Host"),
				Database: viper.GetString(dbtype + ".Database"),
				User:     viper.GetString(dbtype + ".User"),
				Password: viper.GetString(dbtype + ".Password"),
			}

			cli.handleBlocks(settings, startBlock, endBlock, loop, blockDelay)

			return
		},
	}

	cmd.Flags().Bool("loop", false, "loop get the block")
	cmd.Flags().Int64("start", 0, "The `number` of start block")
	cmd.Flags().Int64("end", 0, "The `number` of end block, not include")
	cmd.Flags().Int64("delay", 0, "The `number` of block delay")

	return cmd
}

func (cli *CLI) handleBlocks(settings db.ConnectionURL, startBlock, endBlock *big.Int, loop bool, blockDelay int64) {

	// database
	newdb, err := mysql.Open(settings)
	if err != nil {
		fmt.Println("Open db error: ", err)
		return
	}
	defer newdb.Close()

	newdb.Exec("set time_zone='+00:00';")

	//blocksCollection := newdb.Collection("blocks")
	//transactionsCollection := newdb.Collection("transactions")

	// log
	logFile, err := os.OpenFile(cli.logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("open log file error")
	}
	defer logFile.Close()
	errorLog := log.New(logFile, "[Error]", log.Llongfile|log.Ldate|log.Ltime)

	// start
	if loop && startBlock.Cmp(big.NewInt(0)) == 0 && endBlock.Cmp(big.NewInt(0)) == 0 {
		// get max block number
		row, err := newdb.QueryRow(`select max(block_number) max from blocks`)
		if err != nil {
			fmt.Println("Error get max block number from database: ", err)
			errorLog.Println("Error get max block number from database: ", err)
		}
		var number int
		var numStr string
		row.Scan(&numStr)
		if numStr != "" {
			number, err = strconv.Atoi(numStr)
			if err == nil {
				startBlock.SetInt64(int64(number + 1))
			}
		}
	}

	// rpc
	client, err := ethclient.Dial(cli.rpcURL)
	if err != nil {
		fmt.Println("Dial Error: ", err)
		return
	}
	ctx := context.Background()

	// get block interval
	blockLatest, err := client.BlockByNumber(ctx, nil)
	if err != nil {
		fmt.Println("BlockByNumber latest Error: ", err)
		errorLog.Println("BlockByNumber latest Error: ", err)
		return
	}
	blockLatesParentHash := blockLatest.ParentHash()
	blockLatesParent, err := client.BlockByHash(ctx, blockLatesParentHash)
	if err != nil {
		fmt.Println("BlockByNumber latest parent Error: ", err)
		errorLog.Println("BlockByNumber latest parent Error: ", err)
		return
	}
	blockIntervalUint64 := blockLatest.Time() - blockLatesParent.Time()
	blockInterval := time.Duration(blockIntervalUint64)
	fmt.Printf("Block Interval is %s second\n", blockIntervalUint64)

	isSynced := false // set true when block near NewChain newwork
	blockLatestHeight := blockLatest.Number().Int64()
	for i := startBlock; loop || i.Cmp(endBlock) <= 0; i = i.Add(i, big.NewInt(1)) {
		block, err := client.BlockByNumber(ctx, i)
		if err != nil {
			if isSynced {
				i = i.Sub(i, big.NewInt(1))
				errorLog.Println(err, " isSynced ", i.String())
				time.Sleep(time.Second * blockInterval)
				continue
			} else {
				blocklatest, errl := client.BlockByNumber(ctx, nil)
				if errl != nil {
					i = i.Sub(i, big.NewInt(1))
					fmt.Println("BlockByNumber latest Error: ", errl)
					errorLog.Println("BlockByNumber latest Error: ", errl)
					time.Sleep(time.Second * blockInterval)
					continue
				}

				numberlatest := blocklatest.Number()
				if numberlatest.Sub(numberlatest, big.NewInt(8)).Cmp(i) > 0 {
					//  error: latest - 8 > i
					fmt.Printf("BlockByNumber: Error(%v) with number(%s)\n", err, i.String())
					errorLog.Printf("BlockByNumber: Error(%v) with number(%s)\n", err, i.String())
				} else if numberlatest.Add(numberlatest, big.NewInt(16)).Cmp(i) > 0 {
					// near: latest - 8 <= i < lates + 8
					fmt.Printf("BlockByNumber: %s near the latest block(%s)\n", i.String(), blocklatest.Number().String())
					errorLog.Printf("BlockByNumber: %s near the latest block(%s)\n", i.String(), blocklatest.Number().String())
					isSynced = true
				} else {
					// i is too big: latest + 8 < i
					fmt.Printf("BlockByNumber: %s is much larger than the latest block(%s)\n", i.String(), blocklatest.Number())
					errorLog.Printf("BlockByNumber: %s is much larger than the latest block(%s)\n", i.String(), blocklatest.Number())
					i = blocklatest.Number()
					isSynced = true
				}
			}
			i = i.Sub(i, big.NewInt(1))
			time.Sleep(time.Second * blockInterval)
			continue
		}
		if blockDelay > 0 {
			blockHeightDelay := blockLatestHeight - block.Number().Int64()
			if blockHeightDelay < blockDelay {
				blockLatest, err = client.BlockByNumber(ctx, nil)
				if err != nil {
					fmt.Println("BlockByNumber latest Error: ", err)
					errorLog.Println("BlockByNumber latest Error: ", err)
					i = i.Sub(i, big.NewInt(1))
					time.Sleep(time.Second * blockInterval)
					continue
				}
				blockLatestHeight = blockLatest.Number().Int64()
				blockHeightDelay = blockLatestHeight - block.Number().Int64()
				if blockHeightDelay < blockDelay {
					i = i.Sub(i, big.NewInt(1))
					fmt.Printf("Requires the current block height(%s) to be delayed by %d heights than the latest(%s)\n",
						block.Number().String(), blockDelay, blockLatest.Number().String())
					errorLog.Printf("Requires the current block height(%s) to be delayed by %d heights than the latest(%s)\n",
						block.Number().String(), blockDelay, blockLatest.Number().String())
					time.Sleep(time.Second * blockInterval)
					continue
				}
			}
		}
		fmt.Println("Block: ", block.Number(), block.Hash().String())

		var numberStr string
		number := block.Number()
		if number == nil {
			fmt.Println("Block number nil")
			errorLog.Println("Block number nil", block.Hash().String())
			continue
		}
		numberStr = number.String()

		var sizeStr string
		size := new(big.Int)
		if size != nil {
			size.SetString(block.Size().String(), 10)
			sizeStr = size.String() // 123.00B
		}

		var difficultyStr string
		if block.Difficulty() != nil {
			difficultyStr = block.Difficulty().String()
		}

		var timestamp time.Time
		if block.Time() == 0 {
			timestamp = time.Unix(1, 0).UTC() // not null for TIMESTAMP
		} else {
			timestamp = time.Unix(int64(block.Time()), 0).UTC()
		}

		transactions := block.Transactions()

		// Get signer for clique
		var signer common.Address
		var (
			extraVanity = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
			extraSeal   = 65 // Fixed number of extra-data suffix bytes reserved for signer seal
		)
		header := block.Header()
		if header.Number.Uint64() != 0 && len(header.Extra) >= extraVanity+extraSeal { // clique fixed number(32 + 65)
			signature := header.Extra[len(header.Extra)-extraSeal:]
			pubkey, err := crypto.Ecrecover(sigHash(header).Bytes(), signature)
			if err != nil {
				errorLog.Println(err)
			}
			copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

		} else {
			signer = header.Coinbase
		}

		if newdb == nil {
			newdb, err = mysql.Open(settings)
			if err != nil {
				i = i.Sub(i, big.NewInt(1))
				fmt.Println("Open db error: ", err)
				errorLog.Println("Open db Error: ", err)
				time.Sleep(time.Second)
				continue
			}
			defer newdb.Close()

			newdb.Exec("set time_zone='+00:00';")
		}
		dbtx, err := newdb.NewTx(newdb.Context())
		if err != nil {
			i = i.Sub(i, big.NewInt(1))
			fmt.Println("Database NewTX Error: ", err)
			errorLog.Println("Database NewTX Error: ", err)
			time.Sleep(time.Second)
			if newdb != nil {
				newdb.Close()
				newdb = nil
			}
			newdb = nil
			continue
		}

		_, err = dbtx.Collection("blocks").Insert(Blocks{
			Number:           numberStr,
			Hash:             removeHexPrefix(block.Hash().String()),
			ParentHash:       removeHexPrefix(block.ParentHash().String()), // tmp
			Nonce:            block.Nonce(),
			Sha3Uncles:       removeHexPrefix(block.UncleHash().String()),
			TransactionsRoot: removeHexPrefix(block.TxHash().String()),
			ReceiptsRoot:     removeHexPrefix(block.ReceiptHash().String()),
			StateRoot:        removeHexPrefix(block.Root().String()),
			Miner:            removeHexPrefix(block.Coinbase().String()),
			Difficulty:       difficultyStr,
			TotalDifficulty:  difficultyStr, // TODO:
			Size:             sizeStr,
			GasLimit:         block.GasLimit(),
			GasUsed:          block.GasUsed(),
			Timestamp:        timestamp,
			TransactionCount: transactions.Len(),
			Signer:           removeHexPrefix(signer.String()),
		})

		if err != nil {
			errorLog.Println(err)
		}

		for _, tx := range block.Transactions() {

			fmt.Println("Transactions: ", tx.Hash().String())

			var signer types.Signer = types.FrontierSigner{}
			if tx.Protected() {
				signer = types.NewEIP155Signer(tx.ChainId())
			}
			from, err := types.Sender(signer, tx)
			//from, err := types.Sender(types.MakeSigner(params.MainnetChainConfig, block.Number()), tx)
			if err != nil {
				errorLog.Println("TransactionSender ", err, tx.Hash().String(), block.Number().String())
			}
			to := tx.To()
			var toStr string
			if to != nil {
				toStr = to.String()
			}

			_, err = dbtx.Collection("transactions").Insert(Transactions{
				Hash:        removeHexPrefix(tx.Hash().String()),
				BlockNumber: numberStr,
				Nonce:       tx.Nonce(),
				From:        removeHexPrefix(from.String()),
				To:          removeHexPrefix(toStr),
				Value:       removeHexPrefix(tx.Value().String()),
				Gas:         tx.Gas(),
				GasPrice:    tx.GasPrice().String(),
			})

			if err != nil {
				// jsonHeader, _ := json.Marshal(block.Header())
				// jsonTxs, _ := json.Marshal(block.Transactions())
				// errorLog.Println(err, block.Number().String(), string(jsonHeader), string(jsonTxs))
				errorLog.Println(err, block.Number().String())
			}

		}
		err = dbtx.Commit()
		if err != nil {
			errorLog.Println(err, block.Number().String())
			if err.Error() == "invalid connection" {
				i = i.Sub(i, big.NewInt(1))
				if newdb != nil {
					newdb.Close()
					newdb = nil
				}
				continue
			}
		}
	}
}

func removeHexPrefix(str string) string {
	if len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X') {
		return str[2:]
	}
	return str
}

func sigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()

	rlp.Encode(hasher, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-65], // Yes, this will panic if extra is too short
		header.MixDigest,
		header.Nonce,
	})
	hasher.Sum(hash[:0])
	return hash
}
