package main

import (
    "os"
	"fmt"
	"log"
    "bufio"

	"github.com/miguelmota/go-ethereum-hdwallet"
//    "github.com/tyler-smith/go-bip39/wordlists"
    "github.com/ethereum/go-ethereum/accounts"
    "github.com/ethereum/go-ethereum/common"
)

func guessFromFile(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)
    scanner.Scan()

    return scanner.Text()
}

func checkWallet(mnemonic string, account accounts.Account) bool {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}
    deriveAccount(wallet)

    return wallet.Contains(account)
}

func deriveAccount(wallet *hdwallet.Wallet) {
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	// path = hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/1")
	_, err := wallet.Derive(path, true)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(account.Address.Hex())
}

func main() {
    argv := os.Args
    if len(argv) < 3 {
        log.Fatal("Not enough arguments")
    }

    searchAddress := common.HexToAddress(argv[1])
    guessFilename := argv[2]
    account := accounts.Account{Address: searchAddress}
    guessString := guessFromFile(guessFilename)

    fmt.Printf("Looking who owns %s\n", account.Address.Hex())
    fmt.Printf("Guess: %s\n", guessString)

	mnemonic := guessString
    fmt.Println(checkWallet(mnemonic, account))

    //print(wordlists.English[2])
}

type Guess struct {
    w01 []string
    w02 []string
    w03 []string
    w04 []string
    w05 []string
    w06 []string
    w07 []string
    w08 []string
    w09 []string
    w10 []string
    w11 []string
    w12 []string
}
