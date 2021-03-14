package bitcoin

type Blockchain struct {
	Blocks []*Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: []*Block{NewGenesisBlock()},
	}
}

func (c *Blockchain) AddBlock(data string) {
	prevBlock := c.Blocks[len(c.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	c.Blocks = append(c.Blocks, newBlock)
}
