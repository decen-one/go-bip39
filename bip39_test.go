package bip39

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/decen-one/go-bip39/assert"
	"github.com/decen-one/go-bip39/wordlist"
)

var languages = [10]string{"chinese-simplified", "chinese-traditional", "czech", "english", "french", "italian", "japanese", "korean", "portuguese", "spanish"}

func TestGetWordList(t *testing.T) {
	for _, lang := range languages {
		switch lang {
		case "chinese-simplified":
			assert.EqualStringsSlices(t, wordlist.ChineseSimplified, wordList[lang])
		case "chinese-traditional":
			assert.EqualStringsSlices(t, wordlist.ChineseTraditional, wordList[lang])
		case "czech":
			assert.EqualStringsSlices(t, wordlist.Czech, wordList[lang])
		case "english":
			assert.EqualStringsSlices(t, wordlist.English, wordList[lang])
		case "french":
			assert.EqualStringsSlices(t, wordlist.French, wordList[lang])
		case "italian":
			assert.EqualStringsSlices(t, wordlist.Italian, wordList[lang])
		case "japanese":
			assert.EqualStringsSlices(t, wordlist.Japanese, wordList[lang])
		case "korean":
			assert.EqualStringsSlices(t, wordlist.Korean, wordList[lang])
		case "portuguese":
			assert.EqualStringsSlices(t, wordlist.Portuguese, wordList[lang])
		case "spanish":
			assert.EqualStringsSlices(t, wordlist.Spanish, wordList[lang])
		}
	}

}

func TestGetWordIndex(t *testing.T) {
	for _, lang := range languages {
		words, _ := GetWordList(lang)
		for expectedIdx, word := range words {
			actualIdx, err := GetWordIndex(lang, word)
			if expectedIdx != actualIdx {
				fmt.Println(lang, expectedIdx, actualIdx, word)
				fmt.Println(wordMap[lang])
				fmt.Println(wordList[lang])

			}

			assert.Nil(t, err)
			assert.Equal(t, actualIdx, expectedIdx)
		}

		for _, word := range []string{"a", "set", "of", "invalid", "words"} {
			actualIdx, err := GetWordIndex(lang, word)
			assert.NotNil(t, err)
			assert.Equal(t, actualIdx, -1)
		}

	}

}

func TestNewMnemonic(t *testing.T) {
	for _, vector := range testVectors() {
		entropy, err := hex.DecodeString(vector.entropy)
		assert.Nil(t, err)

		mnemonic, err := NewMnemonic(vector.lang, entropy)
		assert.Nil(t, err)
		if vector.mnemonic != mnemonic {
			fmt.Println(vector, "x"+mnemonic+"x")
		}
		assert.EqualString(t, vector.mnemonic, mnemonic)

		_, err = NewSeedWithErrorChecking(vector.lang, mnemonic, vector.password)
		assert.Nil(t, err)

		seed := NewSeed(mnemonic, vector.password)
		assert.EqualString(t, vector.seed, hex.EncodeToString(seed))
	}
}

func TestNewMnemonicInvalidEntropy(t *testing.T) {
	for _, lang := range languages {
		_, err := NewMnemonic(lang, []byte{})
		assert.NotNil(t, err)

	}

}

func TestNewSeedWithErrorCheckingInvalidMnemonics(t *testing.T) {
	for _, vector := range badMnemonicSentences() {
		_, err := NewSeedWithErrorChecking(vector.lang, vector.mnemonic, "TREZOR")
		assert.NotNil(t, err)
	}
}

func TestIsMnemonicValid(t *testing.T) {
	for _, vector := range badMnemonicSentences() {
		assert.False(t, IsMnemonicValid(vector.lang, vector.mnemonic))
	}

	for _, vector := range testVectors() {
		assert.True(t, IsMnemonicValid(vector.lang, vector.mnemonic))
	}
}

func TestMnemonicToByteArrayWithRawIsEqualToEntropyFromMnemonic(t *testing.T) {
	for _, vector := range testVectors() {
		rawEntropy, err := MnemonicToByteArray(vector.lang, vector.mnemonic, true)
		assert.Nil(t, err)
		rawEntropy2, err := EntropyFromMnemonic(vector.lang, vector.mnemonic)
		assert.Nil(t, err)
		assert.True(t, bytes.Equal(rawEntropy, rawEntropy2))
	}
}

func TestMnemonicToByteArrayInvalidMnemonics(t *testing.T) {
	for _, vector := range badMnemonicSentences() {
		_, err := MnemonicToByteArray(vector.lang, vector.mnemonic)
		assert.NotNil(t, err)
	}

	_, err := MnemonicToByteArray("english", "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon yellow")
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrChecksumIncorrect)

	_, err = MnemonicToByteArray("english", "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon angry")
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrInvalidMnemonic)
}

func TestNewEntropy(t *testing.T) {
	// Good tests.
	for i := 128; i <= 256; i += 32 {
		_, err := NewEntropy(i)
		assert.Nil(t, err)
	}

	// Bad Values
	for i := 0; i <= 256; i++ {
		if i%32 != 0 {
			_, err := NewEntropy(i)
			assert.NotNil(t, err)
		}
	}
}
func TestNewEntropyWithMnemonicSize(t *testing.T) {
	// Good tests.
	for i := 12; i <= 24; i += 3 {
		_, err := NewEntropyWithMnemonicSize(i)
		assert.Nil(t, err)
	}

	// Bad Values
	for i := 0; i <= 24; i++ {
		if i%3 != 0 {
			_, err := NewEntropyWithMnemonicSize(i)
			assert.NotNil(t, err)
		}
	}
}

func TestCheckLanguages(t *testing.T) {
	for _, lang := range languages {
		res := checkLanguage(lang)
		assert.Equal(t, res, lang)
	}

	res := checkLanguage("badValue")
	assert.Equal(t, res, "error")

}
func TestNewRandMnemonic(t *testing.T) {
	for _, lang := range languages {
		// Good tests.
		for i := 12; i <= 24; i += 3 {
			_, err := NewRandMnemonic(lang, i)
			assert.Nil(t, err)
		}

		// Bad Values
		for i := 0; i <= 24; i++ {
			if i%3 != 0 {
				_, err := NewRandMnemonic(lang, i)
				assert.NotNil(t, err)
			}
		}
		for i := 12; i <= 24; i += 3 {
			_, err := NewRandMnemonic(lang+"badValue", i)
			assert.NotNil(t, err)
		}

	}
	// Good tests.
	for i := 12; i <= 24; i += 3 {
		_, err := NewEntropyWithMnemonicSize(i)
		assert.Nil(t, err)
	}

	// Bad Values
	for i := 0; i <= 24; i++ {
		if i%3 != 0 {
			_, err := NewEntropyWithMnemonicSize(i)
			assert.NotNil(t, err)
		}
	}
}
func TestMnemonicToByteArrayForDifferentArrayLengths(t *testing.T) {
	max := 1000
	for _, lang := range languages {
		for i := 0; i < max; i++ {
			//16, 20, 24, 28, 32
			length := 16 + (i%5)*4
			seed := make([]byte, length)

			if n, err := rand.Read(seed); err != nil {
				t.Errorf("%v", err)
			} else if n != length {
				t.Errorf("Wrong number of bytes read: %d", n)
			}

			mnemonic, err := NewMnemonic(lang, seed)
			if err != nil {
				t.Errorf("%v", err)
			}

			_, err = MnemonicToByteArray(lang, mnemonic)
			if err != nil {
				t.Errorf("Failed for %x - %v", seed, mnemonic)
			}
		}

	}

}

func TestMnemonicToByteArrayForZeroLeadingSeeds(t *testing.T) {
	for _, m := range []string{
		"00000000000000000000000000000000", "00a84c51041d49acca66e6160c1fa999",
		"00ca45df1673c76537a2020bfed1dafd", "0019d5871c7b81fd83d474ef1c1e1dae",
		"00dcb021afb35ffcdd1d032d2056fc86", "0062be7bd09a27288b6cf0eb565ec739",
		"00dc705b5efa0adf25b9734226ba60d4", "0017747418d54c6003fa64fade83374b",
		"000d44d3ee7c3dfa45e608c65384431b", "008241c1ef976b0323061affe5bf24b9",
		"0011527b8c6ddecb9d0c20beccdeb58d", "001c938c503c8f5a2bba2248ff621546",
		"0002f90aaf7a8327698f0031b6317c36", "00bff43071ed7e07f77b14f615993bac",
		"00da143e00ef17fc63b6fb22dcc2c326", "00ffc6764fb32a354cab1a3ddefb015d",
		"0062ef47e0985e8953f24760b7598cdd", "003bf9765064f71d304908d906c065f5",
		"00993851503471439d154b3613947474", "00a6aec77e4d16bea80b50a34991aaba",
		"007ad0ffe9eae753a483a76af06dfa67", "00091824db9ec19e663bee51d64c83cc",
		"00f48ac621f7e3cb39b2012ac3121543", "0072917415cdca24dfa66c4a92c885b4",
		"0027ced2b279ea8a91d29364487cdbf4", "00b9c0d37fb10ba272e55842ad812583",
		"004b3d0d2b9285946c687a5350479c8c", "00c7c12a37d3a7f8c1532b17c89b724c",
		"00f400c5545f06ae17ad00f3041e4e26", "001e290be10df4d209f247ac5878662b",
		"00bf0f74568e582a7dd1ee64f792ec8b", "00d2e43ecde6b72b847db1539ed89e23",
		"00cecba6678505bb7bfec8ed307251f6", "000aeed1a9edcbb4bc88f610d3ce84eb",
		"00d06206aadfc25c2b21805d283f15ae", "00a31789a2ab2d54f8fadd5331010287",
		"003493c5f520e8d5c0483e895a121dc9", "004706112800b76001ece2e268bc830e",
		"00ab31e28bb5305be56e38337dbfa486", "006872fe85df6b0fa945248e6f9379d1",
		"00717e5e375da6934e3cfdf57edaf3bd", "007f1b46e7b9c4c76e77c434b9bccd6b",
		"00dc93735aa35def3b9a2ff676560205", "002cd5dcd881a49c7b87714c6a570a76",
		"0013b5af9e13fac87e0c505686cfb6bf", "007ab1ec9526b0bc04b64ae65fd42631",
		"00abb4e11d8385c1cca905a6a65e9144", "00574fc62a0501ad8afada2e246708c3",
		"005207e0a815bb2da6b4c35ec1f2bf52", "00f3460f136fb9700080099cbd62bc18",
		"007a591f204c03ca7b93981237112526", "00cfe0befd428f8e5f83a5bfc801472e",
		"00987551ac7a879bf0c09b8bc474d9af", "00cadd3ce3d78e49fbc933a85682df3f",
		"00bfbf2e346c855ccc360d03281455a1", "004cdf55d429d028f715544ce22d4f31",
		"0075c84a7d15e0ac85e1e41025eed23b", "00807dddd61f71725d336cab844d2cb5",
		"00422f21b77fe20e367467ed98c18410", "00b44d0ac622907119c626c850a462fd",
		"00363f5e7f22fc49f3cd662a28956563", "000fe5837e68397bbf58db9f221bdc4e",
		"0056af33835c888ef0c22599686445d3", "00790a8647fd3dfb38b7e2b6f578f2c6",
		"00da8d9009675cb7beec930e263014fb", "00d4b384540a5bb54aa760edaa4fb2fe",
		"00be9b1479ed680fdd5d91a41eb926d0", "009182347502af97077c40a6e74b4b5c",
		"00f5c90ee1c67fa77fd821f8e9fab4f1", "005568f9a2dd6b0c0cc2f5ba3d9cac38",
		"008b481f8678577d9cf6aa3f6cd6056b", "00c4323ece5e4fe3b6cd4c5c932931af",
		"009791f7550c3798c5a214cb2d0ea773", "008a7baab22481f0ad8167dd9f90d55c",
		"00f0e601519aafdc8ff94975e64c946d", "0083b61e0daa9219df59d697c270cd31",
	} {
		seed, _ := hex.DecodeString(m)

		mnemonic, err := NewMnemonic("english", seed)
		if err != nil {
			t.Errorf("%v", err)
		}

		if _, err = MnemonicToByteArray("english", mnemonic); err != nil {
			t.Errorf("Failed for %x - %v", seed, mnemonic)
		}
	}
}
func TestEntropyFromMnemonic128(t *testing.T) {
	testEntropyFromMnemonic(t, 128)
}

func TestEntropyFromMnemonic160(t *testing.T) {
	testEntropyFromMnemonic(t, 160)
}

func TestEntropyFromMnemonic192(t *testing.T) {
	testEntropyFromMnemonic(t, 192)
}

func TestEntropyFromMnemonic224(t *testing.T) {
	testEntropyFromMnemonic(t, 224)
}

func TestEntropyFromMnemonic256(t *testing.T) {
	testEntropyFromMnemonic(t, 256)
}

func TestEntropyFromMnemonicInvalidChecksum(t *testing.T) {
	_, err := EntropyFromMnemonic("english", "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon yellow")
	assert.Equal(t, ErrChecksumIncorrect, err)
}

func TestEntropyFromMnemonicInvalidMnemonicSize(t *testing.T) {
	for _, mnemonic := range []string{
		"a a a a a a a a a a a a a a a a a a a a a a a a a", // Too many words
		"a",                           // Too few
		"a a a a a a a a a a a a a a", // Not multiple of 3
	} {
		_, err := EntropyFromMnemonic("english", mnemonic)
		assert.Equal(t, ErrInvalidMnemonic, err)
	}
}

func testEntropyFromMnemonic(t *testing.T, bitSize int) {
	for _, lang := range languages {
		for i := 0; i < 512; i++ {
			expectedEntropy, err := NewEntropy(bitSize)
			assert.Nil(t, err)
			assert.True(t, len(expectedEntropy) != 0)

			mnemonic, err := NewMnemonic(lang, expectedEntropy)
			assert.Nil(t, err)
			assert.True(t, len(mnemonic) != 0)

			actualEntropy, err := EntropyFromMnemonic(lang, mnemonic)
			assert.Nil(t, err)
			assert.EqualByteSlices(t, expectedEntropy, actualEntropy)
		}

	}

}

func TestPadByteSlice(t *testing.T) {
	assert.EqualByteSlices(t, []byte{0}, padByteSlice([]byte{}, 1))
	assert.EqualByteSlices(t, []byte{0, 1}, padByteSlice([]byte{1}, 2))
	assert.EqualByteSlices(t, []byte{1, 1}, padByteSlice([]byte{1, 1}, 2))
	assert.EqualByteSlices(t, []byte{1, 1, 1}, padByteSlice([]byte{1, 1, 1}, 2))
}

func TestCompareByteSlices(t *testing.T) {
	assert.True(t, compareByteSlices([]byte{}, []byte{}))
	assert.True(t, compareByteSlices([]byte{1}, []byte{1}))
	assert.False(t, compareByteSlices([]byte{1}, []byte{0}))
	assert.False(t, compareByteSlices([]byte{1}, []byte{}))
	assert.False(t, compareByteSlices([]byte{1}, nil))
}
