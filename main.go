package main

import (
	"fmt"
	"log"
	"goBlockchain/wallet"
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
}
