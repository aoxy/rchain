package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const targetBits = 24

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}
	return pow
}

func (p *ProofOfWork) Run() (int, []byte) {
	// hashInt是hash的整形表示
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Println("Mining the block containing")
	maxNonce := math.MaxInt64
	for nonce < maxNonce {
		data := p.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(p.target) == -1 {
			fmt.Printf("\r%x \n\n", hash)
			break
		} else {
			nonce++
		}
	}
	return nonce, hash[:]
}

// 对工作量证明进行验证
func (p *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := p.prepareData(p.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(p.target) == -1
}

// 工作量证明算法必须要将存储在区块里面的交易考虑进去
// 从而保证区块链交易存储的一致性和可靠性
func (p *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			p.block.PrevBlockHash,
			p.block.HashTransactions(),
			IntToHex(p.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func IntToHex(num int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}
