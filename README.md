# go-bip39

Go implementation of Bitcoin BIP39,containing full test compatible with bitcoinjs and other well-known mnemonic implementations of bip39.
The main features of the implementation are:
1) Multilingual support with multilingual test of code
2) Standard of utf and bytes based on NFKD encoding of bitcoinjs bip39
3) Full coverage testing
The implementation is based on https://github.com/tyler-smith/go-bip39
and resolved problem of  encoding utf bytes based on NFKD encoding of bitcoinjs bip39 that was not supported in the package and added multilingual support of bip39 to the project.
The project tried to complete the tests and resolve bottlenecks of the mentioned project

# Examples
```go
package main

// The example contains two functions of bip39 package. You can see other functions in bip39 package
import (
	"encoding/hex"
	"fmt"

	"github.com/decen-one/go-bip39"
)

func main() {
	//NewRandMnemonic(language string,mnemonicSize int) (mnemonic space separated wordlist string, error )
	//language should be in ["chinesesimplified", "chinesetraditional", "czech", "english", "french", "italian", "japanese", "korean", "portuguese", "spanish"]
	//mnemonicSize should be in [12, 15, 18, 21, 24]
	words, err := bip39.NewRandMnemonic("english", 12)
	fmt.Println(words)
	fmt.Println(err)

	//convert words to corresponding seed checking words and language and validating the mnemonic words first
	//language should be in ["chinesesimplified", "chinesetraditional", "czech", "english", "french", "italian", "japanese", "korean", "portuguese", "spanish"]
	//mnemonic should be list of words with length of  [12, 15, 18, 21, 24] words in space separated string
	//password is passphrase string
	seed, err := bip39.NewSeedWithErrorChecking("english", words, "password")
	fmt.Println(hex.EncodeToString(seed))
	fmt.Println(err)
}

```