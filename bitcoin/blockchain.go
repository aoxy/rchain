package bitcoin

import (
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

func NewBlockchain() *Blockchain {

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
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			genesis := NewGenesisBlock()
			b, err = tx.CreateBucket([]byte(blocksBucket))
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("1"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("1"))
		}
		return nil
	})

	return &Blockchain{tip, db}
}

func (c *Blockchain) AddBlock(data string) {
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

	newBlock := NewBlock(data, lastHash)

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
