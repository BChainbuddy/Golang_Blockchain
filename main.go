package main

import (
	"fmt"
	"goBlockchain/wallet"
	"log"
)

func init() {
	log.SetPrefix(("Blockchain: "))
}

func main() {
	w := wallet.NewWallet()
	fmt.Println(w.PrivateKey())
	fmt.Println(w.PublicKey())
	fmt.Println(w.PrivateKeyStr())
	fmt.Println(w.PublicKeyStr())
	fmt.Println(w.BlockChainAddress())
}
