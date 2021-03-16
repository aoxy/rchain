package bitcoin

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/spf13/viper"
)

type Blockchain struct {
	Tip []byte
	Db  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

var blocksBucket = "bc-bucket"

func NewBlockchain(address string) *Blockchain {

	// tip有尾部，尖端的意思，在这里tip存储的是最后一个块的哈希
	var tip []byte

	// 1.打开一个数据库文件
	// 2.检查文件里面是否已经存储了一个区块链
	// 3.如果有：
	// 		创建一个新的Blockchain实例；
	// 		设置 Blockchain 实例的 tip 为数据库中存储的最后一个块的哈希
	// 4.如果没有：
	// 		创建创世块；
	// 		存储到数据库；
	// 		将创世块哈希保存为最后一个块的哈希；
	// 		创建一个新的 Blockchain 实例，初始时 tip 指向创世块
	db, err := bolt.Open(viper.GetString("bitcoin.dbFile"), 0600, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, "New Block.")
		genesis := NewGenesisBlock(cbtx)
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			fmt.Println("创建Bucket")
			b, err = tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				fmt.Println("CreateBucket error ", err)
				return err
			}
			fmt.Println("CreateBucket success")
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				fmt.Println("Put1 error ", err)
				return err
			}
			fmt.Println("Put1 success")
			err = b.Put([]byte("1"), genesis.Hash)
			if err != nil {
				fmt.Println("Put2 error ", err)
				return err
			}
			fmt.Println("Put2 success")
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("1"))
		}
		return nil
	})

	fmt.Println("ok...", tip)
	return &Blockchain{tip, db}
}

func (c *Blockchain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	// BoltDB 事务类型（只读）
	err := c.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("1"))
		return nil
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	newBlock := NewBlock(transactions, lastHash)

	// BoltDB 事务类型（读写）
	err = c.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}
		err = b.Put([]byte("1"), newBlock.Hash)
		if err != nil {
			return err
		}
		c.Tip = newBlock.Hash
		return nil
	})

	if err != nil {
		fmt.Println("AddBlock error ", err)
	}
}

// 找到可以花费的交易
func (c *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := c.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

// 找到包含未花费输出的交易
func (c *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := c.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if tx.IsCoinBase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTXs
}

// 迭代器的初始状态为链中的 tip，因此区块将从尾到头（创世块为头），也就是从最新的到最旧的进行获取。
func (c *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{c.Tip, c.Db}
}

// BlockchainIterator 只会做一件事情：返回链中的下一个块。
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	i.currentHash = block.PrevBlockHash
	return block
}
