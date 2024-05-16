package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

type Block struct {
	timestamp    int64
	nonce        int
	prevHash     [32]byte
	transactions []*Transaction
}

func (b *Block) Print() {
	fmt.Printf("Timestamp         %d\n", b.timestamp)
	fmt.Printf("Nonce             %d\n", b.nonce)
	fmt.Printf("Previous_hash     %x\n", b.prevHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int            `json:"nonce"`
		PrevHash     [32]byte       `json:"prevHash"`
		Timestamp    int64          `json:"timestamp"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Nonce:        b.nonce,
		PrevHash:     b.prevHash,
		Timestamp:    b.timestamp,
		Transactions: b.transactions,
	})
}

func newBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.prevHash = previousHash
	b.transactions = transactions
	return b
}

type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
}

func newBlockchain(blockchainAddress string) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.createBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) createBlock(nonce int, previousHash [32]byte) *Block {
	b := newBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d  %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func (bn *Blockchain) LastBlock() *Block {
	return bn.chain[len(bn.chain)-1]
}

func (bc *Blockchain) addTransaction(sender string, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions, NewTransaction(t.senderAddress, t.recipientAddress, t.value))
	}
	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.addTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.createBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
}

// User amount
func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalValue float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockchainAddress == t.recipientAddress {
				totalValue += value
			}
			if blockchainAddress == t.senderAddress {
				totalValue -= value
			}
		}
	}
	return totalValue
}

type Transaction struct {
	senderAddress    string
	recipientAddress string
	value            float32
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{senderAddress: sender, recipientAddress: recipient, value: value}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blockchain_address    %s\n", t.senderAddress)
	fmt.Printf(" recipient_blockchain_address %s\n", t.recipientAddress)
	fmt.Printf(" transaction_value            %.1f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"senderAddress"`
		Recipient string  `json:"recipientAddress"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderAddress,
		Recipient: t.recipientAddress,
		Value:     t.value,
	})
}

func init() {
	log.SetPrefix(("Blockchain: "))
}

func main() {
	myBlockchainAddress := "my_blockchain_address"
	blockchain := newBlockchain(myBlockchainAddress)
	blockchain.Print()

	blockchain.addTransaction("A", "B", 1.0)
	blockchain.Mining()
	blockchain.Print()

	blockchain.addTransaction("C", "D", 10.0)
	blockchain.addTransaction("F", "H", 5.7)
	blockchain.Mining()
	blockchain.Print()

	fmt.Printf("my %.1f\n", blockchain.CalculateTotalAmount("my_blockchain_address"))
	fmt.Printf("C %.1f\n", blockchain.CalculateTotalAmount("F"))
	fmt.Printf("D %.1f\n", blockchain.CalculateTotalAmount("D"))
}
