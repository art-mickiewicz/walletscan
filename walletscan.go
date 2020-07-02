package main

import (
	"fmt"
	"log"

	"github.com/miguelmota/go-ethereum-hdwallet"
)

func main() {
	mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex())

	// path = hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/1")
	// account, err = wallet.Derive(path, false)
	// if err != nil {
	//	log.Fatal(err)
	// }

	// fmt.Println(account.Address.Hex())
}
