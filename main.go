package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          string
	PrevBlockHash string
	Hash          string
	Nonce         int
}

type Blockchain struct {
	Blocks     []*Block
	Difficulty int
}

func calculateHash(timestamp int64, data string, prevHash string, nonce int) string {
	record := strconv.FormatInt(timestamp, 10) + data + prevHash + strconv.Itoa(nonce)

	h := sha256.New()

	h.Write([]byte(record))

	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}
func mineBlock(difficulty int, timestamp int64, data string, prevHash string) (int, string) {
	nonce := 0
	target := strings.Repeat("0", difficulty)

	for {
		hash := calculateHash(timestamp, data, prevHash, nonce)

		if strings.HasPrefix(hash, target) {
			return nonce, hash
		}

		nonce++
	}
}

func NewBlock(data string, prevHash string, difficulty int) *Block {
	timestamp := time.Now().Unix()

	nonce, hash := mineBlock(difficulty, timestamp, data, prevHash)

	block := &Block{
		Timestamp:     timestamp,
		Data:          data,
		PrevBlockHash: prevHash,
		Hash:          hash,
		Nonce:         nonce,
	}

	return block
}

func NewGenesisBlock(difficulty int) *Block {
	return NewBlock("Genesis Block", "", difficulty)
}

func NewBlockchain(difficulty int) *Blockchain {
	genesisBlock := NewGenesisBlock(difficulty)

	return &Blockchain{
		Blocks:     []*Block{genesisBlock},
		Difficulty: difficulty,
	}
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	newBlock := NewBlock(data, prevBlock.Hash, bc.Difficulty)

	bc.Blocks = append(bc.Blocks, newBlock)
}

func main() {
	const difficulty = 4

	log.Println("Creating a new blockchain...")
	bc := NewBlockchain(difficulty)
	log.Println("Blockchain created!")

	log.Println("Mining block 1 (Send 1 BTC to Alex)...")
	bc.AddBlock("Send 1 BTC to Alex")
	log.Println("Block 1 added.")

	log.Println("Mining block 2 (Send 2 BTC to Ben)...")
	bc.AddBlock("Send 2 BTC to Ben")
	log.Println("Block 2 added.")

	log.Println("\n--- Printing Blockchain ---")
	for i, block := range bc.Blocks {
		fmt.Printf("======= Block %d =======\n", i)
		fmt.Printf("Data:          %s\n", block.Data)
		fmt.Printf("Timestamp:     %d\n", block.Timestamp)
		fmt.Printf("Prev. Hash:    %s\n", block.PrevBlockHash)
		fmt.Printf("Hash:          %s\n", block.Hash)
		fmt.Printf("Nonce:         %d\n", block.Nonce)
		fmt.Printf("PoW Difficulty:  %d\n", bc.Difficulty)
		fmt.Println()
	}
}

