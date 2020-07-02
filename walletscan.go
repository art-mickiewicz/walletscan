package main

import (
    "os"
    "runtime"
	"fmt"
	"log"
    "bufio"
    "strings"
    "errors"
    "robpike.io/filter"

	"github.com/miguelmota/go-ethereum-hdwallet"
    "github.com/tyler-smith/go-bip39/wordlists"
    "github.com/ethereum/go-ethereum/accounts"
    "github.com/ethereum/go-ethereum/common"
    "gopkg.in/godo.v2/glob"
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

func checkWallet(mnemonic string, account accounts.Account, results chan Result) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
        results <- Result{Success: false, Mnemonic: mnemonic}
        return
	}

    deriveAccount(wallet)
    results <- Result{Success: wallet.Contains(account), Mnemonic: mnemonic}
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

func makeThreadedGuesses(g Guess, numThreads int) []Guess {
    out := make([]Guess, numThreads)
    vPerThread := g.NumVariants() / int64(numThreads)
    for i := 0; i < numThreads; i++ {
        out[i] = g
        out[i].Seek = vPerThread * int64(i)
        if i < numThreads - 1 {
            out[i].Limit = vPerThread * (int64(i) + 1) - 1
        }
    }
    return out
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

    guess := createGuessFromString(guessString)
    numVariants := guess.NumVariants()
    fmt.Printf("Number of variants: %d\n", numVariants)

    numThreads := runtime.NumCPU()
    threadedGuesses := makeThreadedGuesses(guess, numThreads)
    runtime.GOMAXPROCS(numThreads)
    results := make(chan Result, numThreads)
    var checkedNum int64
    var percent int
    Main:
    for {
        dispatched := 0
        for i := range threadedGuesses {
            mnemonic, err := threadedGuesses[i].Next()
            if err != nil {
                continue
            }
            go checkWallet(mnemonic, account, results)
            dispatched++
        }

        /* All variants checked */
        if dispatched == 0 {
            break Main
        }

        for i := 0; i < dispatched; i++ {
            res := <-results
            checkedNum++
            newPercent := int(checkedNum * 100 / numVariants)
            if newPercent > percent {
                percent = newPercent
                fmt.Printf("%d %%\n", percent)
            }
            if res.Success {
                fmt.Printf("SUCCESS! %s\n", res.Mnemonic)
                break Main
            }
        }
    }
    fmt.Printf("Mnemonics checked: %d\n", checkedNum)
}

func wordlistForGlob(s string) []string {
    regexp := glob.Globexp(s)
    matches := func (w string) bool {
        return regexp.MatchString(w)
    }
    return filter.Choose(wordlists.English, matches).([]string)
}

func createGuessFromString(guessString string) Guess {
    words := strings.Fields(guessString)
    var guess Guess
    wordCount := len(guess.Words)
    if len(words) != wordCount {
        log.Fatal("Guess is misformed")
    }
    for i := 0; i < wordCount; i++ {
        guess.Words[i] = wordlistForGlob(words[i])
    }

    return guess
}

type Guess struct {
    Words [12][]string
    Seek int64
    Limit int64
}

func (g Guess) NumVariants() int64 {
    var ret int64 = 1
    for i := range g.Words {
        ret *= int64(len(g.Words[i]))
    }

    return ret
}

func (g *Guess) Next() (string, error) {
    seek := g.Seek
    if g.Limit > 0 && seek > g.Limit {
        return "", errors.New("Seek reached the end")
    }

    wordCount := len(g.Words)
    out := make([]string, wordCount)

    for i := wordCount - 1; i >= 0; i-- {
        posCount := int64(len(g.Words[i]))
        wordIdx := seek % posCount
        seek /= posCount
        out[i] = g.Words[i][wordIdx]
    }
    if seek > 0 {
        return "", errors.New("Seek reached the end")
    }

    g.Seek++
    return strings.Join(out, " "), nil
}

type Result struct {
    Success bool
    Mnemonic string
}
