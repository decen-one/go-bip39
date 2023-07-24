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

type vector struct {
	entropy  string
	mnemonic string
	seed     string
	password string
	lang     string
}

var languages = [10]string{"chinesesimplified", "chinesetraditional", "czech", "english", "french", "italian", "japanese", "korean", "portuguese", "spanish"}

func TestGetWordList(t *testing.T) {
	for _, lang := range languages {
		fmt.Println(lang)
		switch lang {
		case "chinesesimplified":
			assert.EqualStringsSlices(t, wordlist.ChineseSimplified, wordList[lang])
		case "chinesetraditional":
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

func testVectors() []vector {
	return []vector{
		{
			entropy:  "00000000000000000000000000000000",
			mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
			seed:     "c55257c360c07c72029aebc1b53c05ed0362ada38ead3e3e9efa3708e53495531f09a6987599d18264c1e1c92f2cf141630c7a3c4ab7c81b2f001698e7463b04",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
			mnemonic: "legal winner thank year wave sausage worth useful legal winner thank yellow",
			seed:     "2e8905819b8723fe2c1d161860e5ee1830318dbf49a83bd451cfb8440c28bd6fa457fe1296106559a3c80937a1c1069be3a3a5bd381ee6260e8d9739fce1f607",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "80808080808080808080808080808080",
			mnemonic: "letter advice cage absurd amount doctor acoustic avoid letter advice cage above",
			seed:     "d71de856f81a8acc65e6fc851a38d4d7ec216fd0796d0a6827a3ad6ed5511a30fa280f12eb2e47ed2ac03b5c462a0358d18d69fe4f985ec81778c1b370b652a8",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "ffffffffffffffffffffffffffffffff",
			mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo wrong",
			seed:     "ac27495480225222079d7be181583751e86f571027b0497b5b5d11218e0a8a13332572917f0f8e5a589620c6f15b11c61dee327651a14c34e18231052e48c069",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "000000000000000000000000000000000000000000000000",
			mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon agent",
			seed:     "035895f2f481b1b0f01fcf8c289c794660b289981a78f8106447707fdd9666ca06da5a9a565181599b79f53b844d8a71dd9f439c52a3d7b3e8a79c906ac845fa",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
			mnemonic: "legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth useful legal will",
			seed:     "f2b94508732bcbacbcc020faefecfc89feafa6649a5491b8c952cede496c214a0c7b3c392d168748f2d4a612bada0753b52a1c7ac53c1e93abd5c6320b9e95dd",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "808080808080808080808080808080808080808080808080",
			mnemonic: "letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter always",
			seed:     "107d7c02a5aa6f38c58083ff74f04c607c2d2c0ecc55501dadd72d025b751bc27fe913ffb796f841c49b1d33b610cf0e91d3aa239027f5e99fe4ce9e5088cd65",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "ffffffffffffffffffffffffffffffffffffffffffffffff",
			mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo when",
			seed:     "0cd6e5d827bb62eb8fc1e262254223817fd068a74b5b449cc2f667c3f1f985a76379b43348d952e2265b4cd129090758b3e3c2c49103b5051aac2eaeb890a528",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "0000000000000000000000000000000000000000000000000000000000000000",
			mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
			seed:     "bda85446c68413707090a52022edd26a1c9462295029f2e60cd7c4f2bbd3097170af7a4d73245cafa9c3cca8d561a7c3de6f5d4a10be8ed2a5e608d68f92fcc8",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
			mnemonic: "legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth title",
			seed:     "bc09fca1804f7e69da93c2f2028eb238c227f2e9dda30cd63699232578480a4021b146ad717fbb7e451ce9eb835f43620bf5c514db0f8add49f5d121449d3e87",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "8080808080808080808080808080808080808080808080808080808080808080",
			mnemonic: "letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic bless",
			seed:     "c0c519bd0e91a2ed54357d9d1ebef6f5af218a153624cf4f2da911a0ed8f7a09e2ef61af0aca007096df430022f7a2b6fb91661a9589097069720d015e4e982f",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo vote",
			seed:     "dd48c104698c30cfe2b6142103248622fb7bb0ff692eebb00089b32d22484e1613912f0a5b694407be899ffd31ed3992c456cdf60f5d4564b8ba3f05a69890ad",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "77c2b00716cec7213839159e404db50d",
			mnemonic: "jelly better achieve collect unaware mountain thought cargo oxygen act hood bridge",
			seed:     "b5b6d0127db1a9d2226af0c3346031d77af31e918dba64287a1b44b8ebf63cdd52676f672a290aae502472cf2d602c051f3e6f18055e84e4c43897fc4e51a6ff",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "b63a9c59a6e641f288ebc103017f1da9f8290b3da6bdef7b",
			mnemonic: "renew stay biology evidence goat welcome casual join adapt armor shuffle fault little machine walk stumble urge swap",
			seed:     "9248d83e06f4cd98debf5b6f010542760df925ce46cf38a1bdb4e4de7d21f5c39366941c69e1bdbf2966e0f6e6dbece898a0e2f0a4c2b3e640953dfe8b7bbdc5",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "3e141609b97933b66a060dcddc71fad1d91677db872031e85f4c015c5e7e8982",
			mnemonic: "dignity pass list indicate nasty swamp pool script soccer toe leaf photo multiply desk host tomato cradle drill spread actor shine dismiss champion exotic",
			seed:     "ff7f3184df8696d8bef94b6c03114dbee0ef89ff938712301d27ed8336ca89ef9635da20af07d4175f2bf5f3de130f39c9d9e8dd0472489c19b1a020a940da67",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "0460ef47585604c5660618db2e6a7e7f",
			mnemonic: "afford alter spike radar gate glance object seek swamp infant panel yellow",
			seed:     "65f93a9f36b6c85cbe634ffc1f99f2b82cbb10b31edc7f087b4f6cb9e976e9faf76ff41f8f27c99afdf38f7a303ba1136ee48a4c1e7fcd3dba7aa876113a36e4",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "72f60ebac5dd8add8d2a25a797102c3ce21bc029c200076f",
			mnemonic: "indicate race push merry suffer human cruise dwarf pole review arch keep canvas theme poem divorce alter left",
			seed:     "3bbf9daa0dfad8229786ace5ddb4e00fa98a044ae4c4975ffd5e094dba9e0bb289349dbe2091761f30f382d4e35c4a670ee8ab50758d2c55881be69e327117ba",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "2c85efc7f24ee4573d2b81a6ec66cee209b2dcbd09d8eddc51e0215b0b68e416",
			mnemonic: "clutch control vehicle tonight unusual clog visa ice plunge glimpse recipe series open hour vintage deposit universe tip job dress radar refuse motion taste",
			seed:     "fe908f96f46668b2d5b37d82f558c77ed0d69dd0e7e043a5b0511c48c2f1064694a956f86360c93dd04052a8899497ce9e985ebe0c8c52b955e6ae86d4ff4449",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "eaebabb2383351fd31d703840b32e9e2",
			mnemonic: "turtle front uncle idea crush write shrug there lottery flower risk shell",
			seed:     "bdfb76a0759f301b0b899a1e3985227e53b3f51e67e3f2a65363caedf3e32fde42a66c404f18d7b05818c95ef3ca1e5146646856c461c073169467511680876c",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "7ac45cfe7722ee6c7ba84fbc2d5bd61b45cb2fe5eb65aa78",
			mnemonic: "kiss carry display unusual confirm curtain upgrade antique rotate hello void custom frequent obey nut hole price segment",
			seed:     "ed56ff6c833c07982eb7119a8f48fd363c4a9b1601cd2de736b01045c5eb8ab4f57b079403485d1c4924f0790dc10a971763337cb9f9c62226f64fff26397c79",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "4fa1a8bc3e6d80ee1316050e862c1812031493212b7ec3f3bb1b08f168cabeef",
			mnemonic: "exile ask congress lamp submit jacket era scheme attend cousin alcohol catch course end lucky hurt sentence oven short ball bird grab wing top",
			seed:     "095ee6f817b4c2cb30a5a797360a81a40ab0f9a4e25ecd672a3f58a0b5ba0687c096a6b14d2c0deb3bdefce4f61d01ae07417d502429352e27695163f7447a8c",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "18ab19a9f54a9274f03e5209a2ac8a91",
			mnemonic: "board flee heavy tunnel powder denial science ski answer betray cargo cat",
			seed:     "6eff1bb21562918509c73cb990260db07c0ce34ff0e3cc4a8cb3276129fbcb300bddfe005831350efd633909f476c45c88253276d9fd0df6ef48609e8bb7dca8",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "18a2e1d81b8ecfb2a333adcb0c17a5b9eb76cc5d05db91a4",
			mnemonic: "board blade invite damage undo sun mimic interest slam gaze truly inherit resist great inject rocket museum chief",
			seed:     "f84521c777a13b61564234bf8f8b62b3afce27fc4062b51bb5e62bdfecb23864ee6ecf07c1d5a97c0834307c5c852d8ceb88e7c97923c0a3b496bedd4e5f88a9",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "15da872c95a13dd738fbf50e427583ad61f18fd99f628c417a61cf8343c90419",
			mnemonic: "beyond stage sleep clip because twist token leaf atom beauty genius food business side grid unable middle armed observe pair crouch tonight away coconut",
			seed:     "b15509eaa2d09d3efd3e006ef42151b30367dc6e3aa5e44caba3fe4d3e352e65101fbdb86a96776b91946ff06f8eac594dc6ee1d3e82a42dfe1b40fef6bcc3fd",
			password: "TREZOR",
			lang:     "english",
		},
		{
			entropy:  "00000000000000000000000000000000",
			mnemonic: "あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あおぞら",
			seed:     "a262d6fb6122ecf45be09c50492b31f92e9beb7d9a845987a02cefda57a15f9c467a17872029a9e92299b5cbdf306e3a0ee620245cbd508959b6cb7ca637bd55",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
			mnemonic: "そつう れきだい ほんやく わかす りくつ ばいか ろせん やちん そつう れきだい ほんやく わかめ",
			seed:     "aee025cbe6ca256862f889e48110a6a382365142f7d16f2b9545285b3af64e542143a577e9c144e101a6bdca18f8d97ec3366ebf5b088b1c1af9bc31346e60d9",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "80808080808080808080808080808080",
			mnemonic: "そとづら あまど おおう あこがれる いくぶん けいけん あたえる いよく そとづら あまど おおう あかちゃん",
			seed:     "e51736736ebdf77eda23fa17e31475fa1d9509c78f1deb6b4aacfbd760a7e2ad769c714352c95143b5c1241985bcb407df36d64e75dd5a2b78ca5d2ba82a3544",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "ffffffffffffffffffffffffffffffff",
			mnemonic: "われる われる われる われる われる われる われる われる われる われる われる ろんぶん",
			seed:     "4cd2ef49b479af5e1efbbd1e0bdc117f6a29b1010211df4f78e2ed40082865793e57949236c43b9fe591ec70e5bb4298b8b71dc4b267bb96ed4ed282c8f7761c",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "000000000000000000000000000000000000000000000000",
			mnemonic: "あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あらいぐま",
			seed:     "d99e8f1ce2d4288d30b9c815ae981edd923c01aa4ffdc5dee1ab5fe0d4a3e13966023324d119105aff266dac32e5cd11431eeca23bbd7202ff423f30d6776d69",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
			mnemonic: "そつう れきだい ほんやく わかす りくつ ばいか ろせん やちん そつう れきだい ほんやく わかす りくつ ばいか ろせん やちん そつう れいぎ",
			seed:     "eaaf171efa5de4838c758a93d6c86d2677d4ccda4a064a7136344e975f91fe61340ec8a615464b461d67baaf12b62ab5e742f944c7bd4ab6c341fbafba435716",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "808080808080808080808080808080808080808080808080",
			mnemonic: "そとづら あまど おおう あこがれる いくぶん けいけん あたえる いよく そとづら あまど おおう あこがれる いくぶん けいけん あたえる いよく そとづら いきなり",
			seed:     "aec0f8d3167a10683374c222e6e632f2940c0826587ea0a73ac5d0493b6a632590179a6538287641a9fc9df8e6f24e01bf1be548e1f74fd7407ccd72ecebe425",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "ffffffffffffffffffffffffffffffffffffffffffffffff",
			mnemonic: "われる われる われる われる われる われる われる われる われる われる われる われる われる われる われる われる われる りんご",
			seed:     "f0f738128a65b8d1854d68de50ed97ac1831fc3a978c569e415bbcb431a6a671d4377e3b56abd518daa861676c4da75a19ccb41e00c37d086941e471a4374b95",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "0000000000000000000000000000000000000000000000000000000000000000",
			mnemonic: "あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん あいこくしん いってい",
			seed:     "23f500eec4a563bf90cfda87b3e590b211b959985c555d17e88f46f7183590cd5793458b094a4dccc8f05807ec7bd2d19ce269e20568936a751f6f1ec7c14ddd",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
			mnemonic: "そつう れきだい ほんやく わかす りくつ ばいか ろせん やちん そつう れきだい ほんやく わかす りくつ ばいか ろせん やちん そつう れきだい ほんやく わかす りくつ ばいか ろせん まんきつ",
			seed:     "cd354a40aa2e241e8f306b3b752781b70dfd1c69190e510bc1297a9c5738e833bcdc179e81707d57263fb7564466f73d30bf979725ff783fb3eb4baa86560b05",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "8080808080808080808080808080808080808080808080808080808080808080",
			mnemonic: "そとづら あまど おおう あこがれる いくぶん けいけん あたえる いよく そとづら あまど おおう あこがれる いくぶん けいけん あたえる いよく そとづら あまど おおう あこがれる いくぶん けいけん あたえる うめる",
			seed:     "6b7cd1b2cdfeeef8615077cadd6a0625f417f287652991c80206dbd82db17bf317d5c50a80bd9edd836b39daa1b6973359944c46d3fcc0129198dc7dc5cd0e68",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			mnemonic: "われる われる われる われる われる われる われる われる われる われる われる われる われる われる われる われる われる われる われる われる われる われる われる らいう",
			seed:     "a44ba7054ac2f9226929d56505a51e13acdaa8a9097923ca07ea465c4c7e294c038f3f4e7e4b373726ba0057191aced6e48ac8d183f3a11569c426f0de414623",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "77c2b00716cec7213839159e404db50d",
			mnemonic: "せまい うちがわ あずき かろう めずらしい だんち ますく おさめる ていぼう あたる すあな えしゃく",
			seed:     "344cef9efc37d0cb36d89def03d09144dd51167923487eec42c487f7428908546fa31a3c26b7391a2b3afe7db81b9f8c5007336b58e269ea0bd10749a87e0193",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "b63a9c59a6e641f288ebc103017f1da9f8290b3da6bdef7b",
			mnemonic: "ぬすむ ふっかつ うどん こうりつ しつじ りょうり おたがい せもたれ あつめる いちりゅう はんしゃ ごますり そんけい たいちょう らしんばん ぶんせき やすみ ほいく",
			seed:     "b14e7d35904cb8569af0d6a016cee7066335a21c1c67891b01b83033cadb3e8a034a726e3909139ecd8b2eb9e9b05245684558f329b38480e262c1d6bc20ecc4",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "3e141609b97933b66a060dcddc71fad1d91677db872031e85f4c015c5e7e8982",
			mnemonic: "くのう てぬぐい そんかい すろっと ちきゅう ほあん とさか はくしゅ ひびく みえる そざい てんすう たんぴん くしょう すいようび みけん きさらぎ げざん ふくざつ あつかう はやい くろう おやゆび こすう",
			seed:     "32e78dce2aff5db25aa7a4a32b493b5d10b4089923f3320c8b287a77e512455443298351beb3f7eb2390c4662a2e566eec5217e1a37467af43b46668d515e41b",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "0460ef47585604c5660618db2e6a7e7f",
			mnemonic: "あみもの いきおい ふいうち にげる ざんしょ じかん ついか はたん ほあん すんぽう てちがい わかめ",
			seed:     "0acf902cd391e30f3f5cb0605d72a4c849342f62bd6a360298c7013d714d7e58ddf9c7fdf141d0949f17a2c9c37ced1d8cb2edabab97c4199b142c829850154b",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "72f60ebac5dd8add8d2a25a797102c3ce21bc029c200076f",
			mnemonic: "すろっと にくしみ なやむ たとえる へいこう すくう きない けってい とくべつ ねっしん いたみ せんせい おくりがな まかい とくい けあな いきおい そそぐ",
			seed:     "9869e220bec09b6f0c0011f46e1f9032b269f096344028f5006a6e69ea5b0b8afabbb6944a23e11ebd021f182dd056d96e4e3657df241ca40babda532d364f73",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "2c85efc7f24ee4573d2b81a6ec66cee209b2dcbd09d8eddc51e0215b0b68e416",
			mnemonic: "かほご きうい ゆたか みすえる もらう がっこう よそう ずっと ときどき したうけ にんか はっこう つみき すうじつ よけい くげん もくてき まわり せめる げざい にげる にんたい たんそく ほそく",
			seed:     "713b7e70c9fbc18c831bfd1f03302422822c3727a93a5efb9659bec6ad8d6f2c1b5c8ed8b0b77775feaf606e9d1cc0a84ac416a85514ad59f5541ff5e0382481",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "eaebabb2383351fd31d703840b32e9e2",
			mnemonic: "めいえん さのう めだつ すてる きぬごし ろんぱ はんこ まける たいおう さかいし ねんいり はぶらし",
			seed:     "06e1d5289a97bcc95cb4a6360719131a786aba057d8efd603a547bd254261c2a97fcd3e8a4e766d5416437e956b388336d36c7ad2dba4ee6796f0249b10ee961",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "7ac45cfe7722ee6c7ba84fbc2d5bd61b45cb2fe5eb65aa78",
			mnemonic: "せんぱい おしえる ぐんかん もらう きあい きぼう やおや いせえび のいず じゅしん よゆう きみつ さといも ちんもく ちわわ しんせいじ とめる はちみつ",
			seed:     "1fef28785d08cbf41d7a20a3a6891043395779ed74503a5652760ee8c24dfe60972105ee71d5168071a35ab7b5bd2f8831f75488078a90f0926c8e9171b2bc4a",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "4fa1a8bc3e6d80ee1316050e862c1812031493212b7ec3f3bb1b08f168cabeef",
			mnemonic: "こころ いどう きあつ そうがんきょう へいあん せつりつ ごうせい はいち いびき きこく あんい おちつく きこえる けんとう たいこ すすめる はっけん ていど はんおん いんさつ うなぎ しねま れいぼう みつかる",
			seed:     "43de99b502e152d4c198542624511db3007c8f8f126a30818e856b2d8a20400d29e7a7e3fdd21f909e23be5e3c8d9aee3a739b0b65041ff0b8637276703f65c2",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "18ab19a9f54a9274f03e5209a2ac8a91",
			mnemonic: "うりきれ さいせい じゆう むろん とどける ぐうたら はいれつ ひけつ いずれ うちあわせ おさめる おたく",
			seed:     "3d711f075ee44d8b535bb4561ad76d7d5350ea0b1f5d2eac054e869ff7963cdce9581097a477d697a2a9433a0c6884bea10a2193647677977c9820dd0921cbde",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "18a2e1d81b8ecfb2a333adcb0c17a5b9eb76cc5d05db91a4",
			mnemonic: "うりきれ うねる せっさたくま きもち めんきょ へいたく たまご ぜっく びじゅつかん さんそ むせる せいじ ねくたい しはらい せおう ねんど たんまつ がいけん",
			seed:     "753ec9e333e616e9471482b4b70a18d413241f1e335c65cd7996f32b66cf95546612c51dcf12ead6f805f9ee3d965846b894ae99b24204954be80810d292fcdd",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
		{
			entropy:  "15da872c95a13dd738fbf50e427583ad61f18fd99f628c417a61cf8343c90419",
			mnemonic: "うちゅう ふそく ひしょ がちょう うけもつ めいそう みかん そざい いばる うけとる さんま さこつ おうさま ぱんつ しひょう めした たはつ いちぶ つうじょう てさぎょう きつね みすえる いりぐち かめれおん",
			seed:     "346b7321d8c04f6f37b49fdf062a2fddc8e1bf8f1d33171b65074531ec546d1d3469974beccb1a09263440fc92e1042580a557fdce314e27ee4eabb25fa5e5fe",
			lang:     "japanese",
			password: "㍍ガバヴァぱばぐゞちぢ十人十色",
		},
	}
}

func badMnemonicSentences() []vector {
	return []vector{
		{mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon", lang: "english"},
		{mnemonic: "legal winner thank year wave sausage worth useful legal winner thank yellow yellow", lang: "english"},
		{mnemonic: "letter advice cage absurd amount doctor acoustic avoid letter advice caged above", lang: "english"},
		{mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo, wrong", lang: "english"},
		{mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon", lang: "english"},
		{mnemonic: "legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth useful legal will will will", lang: "english"},
		{mnemonic: "letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter always.", lang: "english"},
		{mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo why", lang: "english"},
		{mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art art", lang: "english"},
		{mnemonic: "legal winner thank year wave sausage worth useful legal winner thanks year wave worth useful legal winner thank year wave sausage worth title", lang: "english"},
		{mnemonic: "letter advice cage absurd amount doctor acoustic avoid letters advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic bless", lang: "english"},
		{mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo voted", lang: "english"},
		{mnemonic: "jello better achieve collect unaware mountain thought cargo oxygen act hood bridge", lang: "english"},
		{mnemonic: "renew, stay, biology, evidence, goat, welcome, casual, join, adapt, armor, shuffle, fault, little, machine, walk, stumble, urge, swap", lang: "english"},
		{mnemonic: "dignity pass list indicate nasty", lang: "english"},

		// From issue 32
		{mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon letter", lang: "english"},
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
