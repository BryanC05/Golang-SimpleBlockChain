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

// ## 1. The Block Struct
// This defines the structure of a single block in our chain.
type Block struct {
	Timestamp     int64  // The time when the block was created
	Data          string // The data (e.g., transaction details)
	PrevBlockHash string // The hash of the previous block in the chain
	Hash          string // The hash of the current block
	Nonce         int    // The number used in the Proof-of-Work
}

// ## 2. The Blockchain Struct
// This holds the chain of blocks and the PoW difficulty.
type Blockchain struct {
	Blocks     []*Block
	Difficulty int
}

// calculateHash generates the SHA256 hash for a block.
// It combines all the key block fields into one string and hashes it.
func calculateHash(timestamp int64, data string, prevHash string, nonce int) string {
	// 1. Concatenate all the parts of the block
	record := strconv.FormatInt(timestamp, 10) + data + prevHash + strconv.Itoa(nonce)

	// 2. Create a new SHA256 hash object
	h := sha256.New()

	// 3. Write the data to the hash object
	h.Write([]byte(record))

	// 4. Get the final hash sum
	hashed := h.Sum(nil)

	// 5. Return the hash as a hex-encoded string
	return hex.EncodeToString(hashed)
}

// ## 3. The Proof-of-Work (PoW) Algorithm
// mineBlock repeatedly hashes the block's data until it finds a hash
// with the required number of leading zeros (the "difficulty").
func mineBlock(difficulty int, timestamp int64, data string, prevHash string) (int, string) {
	nonce := 0
	target := strings.Repeat("0", difficulty) // e.g., "0000" if difficulty is 4

	for {
		// Calculate the hash with the current nonce
		hash := calculateHash(timestamp, data, prevHash, nonce)

		// Check if the hash has the required prefix
		if strings.HasPrefix(hash, target) {
			// Found it!
			return nonce, hash
		}

		// Didn't find it, increment nonce and try again
		nonce++
	}
}

// NewBlock creates a new, mined block.
func NewBlock(data string, prevHash string, difficulty int) *Block {
	timestamp := time.Now().Unix()

	// Call the mining function to get the valid nonce and hash
	nonce, hash := mineBlock(difficulty, timestamp, data, prevHash)

	// Create the block with the found hash and nonce
	block := &Block{
		Timestamp:     timestamp,
		Data:          data,
		PrevBlockHash: prevHash,
		Hash:          hash,
		Nonce:         nonce,
	}

	return block
}

// NewGenesisBlock creates the very first block in the chain.
func NewGenesisBlock(difficulty int) *Block {
	// The Genesis Block has no previous hash.
	return NewBlock("Genesis Block", "", difficulty)
}

// NewBlockchain creates a new blockchain with a Genesis Block.
func NewBlockchain(difficulty int) *Blockchain {
	// Create the first block
	genesisBlock := NewGenesisBlock(difficulty)

	// Return a new blockchain with the genesis block
	return &Blockchain{
		Blocks:     []*Block{genesisBlock},
		Difficulty: difficulty,
	}
}

// ## 4. A Way to Add New Blocks
// AddBlock adds a new block to the blockchain.
func (bc *Blockchain) AddBlock(data string) {
	// Get the previous block from the chain
	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	// Create a new block, using the previous block's hash
	newBlock := NewBlock(data, prevBlock.Hash, bc.Difficulty)

	// Append the new block to the chain
	bc.Blocks = append(bc.Blocks, newBlock)
}

// ## 5. Main Function (Demonstration)
func main() {
	// Set the PoW difficulty.
	// 4 is good for a quick demo. Try 5 or 6 to see it slow down.
	const difficulty = 4

	// --- Start the blockchain ---
	log.Println("Creating a new blockchain...")
	bc := NewBlockchain(difficulty)
	log.Println("Blockchain created!")

	// --- Add blocks ---
	log.Println("Mining block 1 (Send 1 BTC to Alex)...")
	bc.AddBlock("Send 1 BTC to Alex")
	log.Println("Block 1 added.")

	log.Println("Mining block 2 (Send 2 BTC to Ben)...")
	bc.AddBlock("Send 2 BTC to Ben")
	log.Println("Block 2 added.")

	// --- Print the blockchain ---
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
