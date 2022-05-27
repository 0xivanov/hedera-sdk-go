//go:build all || unit
// +build all unit

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

const testPrivateKeyStr = "302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10"

const testPublicKeyStr = "302a300506032b6570032100e0c8ec2758a5879ffac226a13c0c516b799e72e35141a0dd828f94d37988a4b7"

const testMnemonic3 = "obvious favorite remain caution remove laptop base vacant increase video erase pass sniff sausage knock grid argue salt romance way alone fever slush dune"

// generated by hedera-keygen-java, not used anywhere
const testMnemonic = "inmate flip alley wear offer often piece magnet surge toddler submit right radio absent pear floor belt raven price stove replace reduce plate home"
const testMnemonicKey = "302e020100300506032b657004220420853f15aecd22706b105da1d709b4ac05b4906170c2b9c7495dff9af49e1391da"

// backup phrase generated by the iOS wallet, not used anywhere
const iosMnemonicString = "tiny denial casual grass skull spare awkward indoor ethics dash enough flavor good daughter early hard rug staff capable swallow raise flavor empty angle"

// private key for "default account", should be index 0
const iosDefaultPrivateKey = "5f66a51931e8c99089472e0d70516b6272b94dd772b967f8221e1077f966dbda2b60cf7ee8cf10ecd5a076bffad9a7c7b97df370ad758c0f1dd4ef738e04ceb6"

// backup phrase generated by the Android wallet, also not used anywhere
const androidMnemonicString = "ramp april job flavor surround pyramid fish sea good know blame gate village viable include mixed term draft among monitor swear swing novel track"

// private key for "default account", should be index 0
const androidDefaultPrivateKey = "c284c25b3a1458b59423bc289e83703b125c8eefec4d5aa1b393c2beb9f2bae66188a344ba75c43918ab12fa2ea4a92960eca029a2320d8c6a1c3b94e06c9985"

// test pem key contests for the above testPrivateKeyStr
const pemString = `-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEINtIS4KOZLLY8SzjwKDpOguMznrxu485yXcyOUSCU44Q
-----END PRIVATE KEY-----
`

// const encryptedPem = `-----BEGIN ENCRYPTED PRIVATE KEY-----
// MIGbMFcGCSqGSIb3DQEFDTBKMCkGCSqGSIb3DQEFDDAcBAi8WY7Gy2tThQICCAAw
// DAYIKoZIhvcNAgkFADAdBglghkgBZQMEAQIEEOq46NPss58chbjUn20NoK0EQG1x
// R88hIXcWDOECttPTNlMXWJt7Wufm1YwBibrxmCq1QykIyTYhy1TZMyxyPxlYW6aV
// 9hlo4YEh3uEaCmfJzWM=
// -----END ENCRYPTED PRIVATE KEY-----`

const encryptedPem = `-----BEGIN ENCRYPTED PRIVATE KEY-----
MIGbMFcGCSqGSIb3DQEFDTBKMCkGCSqGSIb3DQEFDDAcBAi8WY7Gy2tThQICCAAw
DAYIKoZIhvcNAgkFADAdBglghkgBZQMEAQIEEOq46NPss58chbjUn20NoK0EQG1x
R88hIXcWDOECttPTNlMXWJt7Wufm1YwBibrxmCq1QykIyTYhy1TZMyxyPxlYW6aV
9hlo4YEh3uEaCmfJzWM=
-----END ENCRYPTED PRIVATE KEY-----
`

const pemPassphrase = "this is a passphrase"

func TestUnitPrivateKeyGenerate(t *testing.T) {
	key, err := GeneratePrivateKey()

	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(key.String(), ed25519PrivateKeyPrefix))
}

func TestUnitPrivateKeyExternalSerialization(t *testing.T) {
	key, err := PrivateKeyFromString(testPrivateKeyStr)

	require.NoError(t, err)
	assert.Equal(t, testPrivateKeyStr, key.String())
}

func TestUnitPrivateKeyExternalSerializationForConcatenatedHex(t *testing.T) {
	keyStr := "db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10e0c8ec2758a5879ffac226a13c0c516b799e72e35141a0dd828f94d37988a4b7"
	key, err := PrivateKeyFromStringEd25519(keyStr)

	require.NoError(t, err)
	assert.Equal(t, testPrivateKeyStr, key.String())
}

func TestUnitShouldMatchHbarWalletV1(t *testing.T) {
	mnemonic, err := MnemonicFromString("jolly kidnap tom lawn drunk chick optic lust mutter mole bride galley dense member sage neural widow decide curb aboard margin manure")
	require.NoError(t, err)

	key, err := mnemonic.ToLegacyPrivateKey()
	require.NoError(t, err)

	deriveKey, err := key.LegacyDerive(1099511627775)
	require.NoError(t, err)

	assert.Equal(t, "302a300506032b657003210045f3a673984a0b4ee404a1f4404ed058475ecd177729daa042e437702f7791e9", deriveKey.PublicKey().String())
}

func TestUnitLegacyPrivateKeyFromMnemonicDerive(t *testing.T) {
	mnemonic, err := MnemonicFromString("jolly kidnap tom lawn drunk chick optic lust mutter mole bride galley dense member sage neural widow decide curb aboard margin manure")
	require.NoError(t, err)

	key, err := mnemonic.ToLegacyPrivateKey()
	require.NoError(t, err)

	deriveKey, err := key.LegacyDerive(0)
	require.NoError(t, err)
	deriveKey2, err := key.LegacyDerive(-1)
	require.NoError(t, err)

	assert.Equal(t, "302e020100300506032b657004220420882a565ad8cb45643892b5366c1ee1c1ef4a730c5ce821a219ff49b6bf173ddf", deriveKey2.String())
	assert.Equal(t, "302e020100300506032b657004220420fae0002d2716ea3a60c9cd05ee3c4bb88723b196341b68a02d20975f9d049dc6", deriveKey.String())
}

func TestUnitPrivateKeyExternalSerializationForRawHex(t *testing.T) {
	keyStr := "db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10"
	key, err := PrivateKeyFromStringEd25519(keyStr)

	require.NoError(t, err)
	assert.Equal(t, testPrivateKeyStr, key.String())
}

func TestUnitPublicKeyExternalSerializationForDerEncodedHex(t *testing.T) {
	key, err := PublicKeyFromString(testPublicKeyStr)

	require.NoError(t, err)
	assert.Equal(t, testPublicKeyStr, key.String())
}

func TestUnitPublicKeyExternalSerializationForRawHex(t *testing.T) {
	keyStr := "e0c8ec2758a5879ffac226a13c0c516b799e72e35141a0dd828f94d37988a4b7"
	key, err := PublicKeyFromStringEd25519(keyStr)

	require.NoError(t, err)
	assert.Equal(t, testPublicKeyStr, key.String())
}

func TestUnitPrivateKeyFromMnemonic(t *testing.T) {
	mnemonic, err := MnemonicFromString(testMnemonic)
	require.NoError(t, err)

	key, err := PrivateKeyFromMnemonic(mnemonic, "")
	require.NoError(t, err)

	keyDerive, err := key.Derive(^uint32(0))
	require.NoError(t, err)

	assert.Equal(t, "302e020100300506032b657004220420e978a6407b74a0730f7aeb722ad64ab449b308e56006c8bff9aad070b9b66ddf", keyDerive.String())
	assert.Equal(t, testMnemonicKey, key.String())
}

func TestUnitMnemonicToPrivateKey(t *testing.T) {
	mnemonic, err := MnemonicFromString(testMnemonic)
	require.NoError(t, err)

	key, err := mnemonic.ToPrivateKey("")
	require.NoError(t, err)

	assert.Equal(t, testMnemonicKey, key.String())
}

func TestUnitIOSPrivateKeyFromMnemonic(t *testing.T) {
	mnemonic, err := MnemonicFromString(iosMnemonicString)
	require.NoError(t, err)

	key, err := PrivateKeyFromMnemonic(mnemonic, "")
	require.NoError(t, err)

	derivedKey, err := key.Derive(0)
	require.NoError(t, err)

	expectedKey, err := PrivateKeyFromString(iosDefaultPrivateKey)
	require.NoError(t, err)

	assert.Equal(t, expectedKey.ed25519PrivateKey.keyData, derivedKey.ed25519PrivateKey.keyData)
}

func TestUnitAndroidPrivateKeyFromMnemonic(t *testing.T) {
	mnemonic, err := MnemonicFromString(androidMnemonicString)
	require.NoError(t, err)

	key, err := PrivateKeyFromMnemonic(mnemonic, "")
	require.NoError(t, err)

	derivedKey, err := key.Derive(0)
	require.NoError(t, err)

	expectedKey, err := PrivateKeyFromString(androidDefaultPrivateKey)
	require.NoError(t, err)

	assert.Equal(t, expectedKey.ed25519PrivateKey.keyData, derivedKey.ed25519PrivateKey.keyData)
}

func TestUnitMnemonic3(t *testing.T) {
	mnemonic, err := MnemonicFromString(testMnemonic3)
	require.NoError(t, err)

	key, err := mnemonic.ToLegacyPrivateKey()
	require.NoError(t, err)

	derivedKey, err := key.LegacyDerive(0)
	require.NoError(t, err)
	derivedKey2, err := key.LegacyDerive(-1)
	require.NoError(t, err)

	assert.Equal(t, "302e020100300506032b6570042204202b7345f302a10c2a6d55bf8b7af40f125ec41d780957826006d30776f0c441fb", derivedKey.String())
	assert.Equal(t, "302e020100300506032b657004220420caffc03fdb9853e6a91a5b3c57a5c0031d164ce1c464dea88f3114786b5199e5", derivedKey2.String())
}

func TestUnitSigning(t *testing.T) {
	priKey, err := PrivateKeyFromString(testPrivateKeyStr)
	require.NoError(t, err)

	pubKey, err := PublicKeyFromString(testPublicKeyStr)
	require.NoError(t, err)

	testSignData := []byte("this is the test data to sign")
	signature := priKey.Sign(testSignData)

	assert.True(t, ed25519.Verify(pubKey.Bytes(), []byte("this is the test data to sign"), signature))
}

func TestUnitGenerated24MnemonicToWorkingPrivateKey(t *testing.T) {
	mnemonic, err := GenerateMnemonic24()

	require.NoError(t, err)

	privateKey, err := mnemonic.ToPrivateKey("")

	require.NoError(t, err)

	message := []byte("this is a test message")

	signature := privateKey.Sign(message)

	assert.True(t, ed25519.Verify(privateKey.PublicKey().Bytes(), message, signature))
}

func TestUnitGenerated12MnemonicToWorkingPrivateKey(t *testing.T) {
	mnemonic, err := GenerateMnemonic12()

	require.NoError(t, err)

	privateKey, err := mnemonic.ToPrivateKey("")

	require.NoError(t, err)

	message := []byte("this is a test message")

	signature := privateKey.Sign(message)

	assert.True(t, ed25519.Verify(privateKey.PublicKey().Bytes(), message, signature))
}

func TestUnitPrivateKeyFromKeystore(t *testing.T) {
	privatekey, err := PrivateKeyFromKeystore([]byte(testKeystore), passphrase)
	require.NoError(t, err)

	actualPrivateKey, err := PrivateKeyFromStringEd25519(testKeystoreKeyString)
	require.NoError(t, err)

	assert.Equal(t, actualPrivateKey.ed25519PrivateKey.keyData, privatekey.ed25519PrivateKey.keyData)
}

func TestUnitPrivateKeyKeystore(t *testing.T) {
	privateKey, err := PrivateKeyFromString(testPrivateKeyStr)
	require.NoError(t, err)

	keystore, err := privateKey.Keystore(passphrase)
	require.NoError(t, err)

	ksPrivateKey, err := _ParseKeystore(keystore, passphrase)
	require.NoError(t, err)

	assert.Equal(t, privateKey.ed25519PrivateKey.keyData, ksPrivateKey.ed25519PrivateKey.keyData)
}

func TestUnitPrivateKeyReadKeystore(t *testing.T) {
	actualPrivateKey, err := PrivateKeyFromStringEd25519(testKeystoreKeyString)
	require.NoError(t, err)

	keystoreReader := bytes.NewReader([]byte(testKeystore))

	privateKey, err := PrivateKeyReadKeystore(keystoreReader, passphrase)
	require.NoError(t, err)

	assert.Equal(t, actualPrivateKey.ed25519PrivateKey.keyData, privateKey.ed25519PrivateKey.keyData)
}

func TestUnitPrivateKeyFromPem(t *testing.T) {
	actualPrivateKey, err := PrivateKeyFromString(testPrivateKeyStr)
	require.NoError(t, err)

	privateKey, err := PrivateKeyFromPem([]byte(pemString), "")
	require.NoError(t, err)

	assert.Equal(t, actualPrivateKey, privateKey)
}

func TestUnitPrivateKeyFromPemInvalid(t *testing.T) {
	_, err := PrivateKeyFromPem([]byte("invalid"), "")
	assert.Error(t, err)
}

func TestUnitPrivateKeyFromPemWithPassphrase(t *testing.T) {
	actualPrivateKey, err := PrivateKeyFromString(testPrivateKeyStr)
	require.NoError(t, err)

	privateKey, err := PrivateKeyFromPem([]byte(encryptedPem), pemPassphrase)
	require.NoError(t, err)

	assert.Equal(t, actualPrivateKey, privateKey)
}

func TestUnitPrivateKeyECDSASign(t *testing.T) {
	key, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	hash := crypto.Keccak256Hash([]byte("aaa"))
	sig := key.Sign([]byte("aaa"))
	s2 := crypto.VerifySignature(key.ecdsaPrivateKey._PublicKey()._BytesRaw(), hash.Bytes(), sig)
	require.True(t, s2)
}

func DisabledTestUnitPrivateKeyECDSASign(t *testing.T) {
	message := []byte("hello world")
	key, err := PrivateKeyFromStringECSDA("8776c6b831a1b61ac10dac0304a2843de4716f54b1919bb91a2685d0fe3f3048")
	require.NoError(t, err)

	sig := key.Sign(message)

	require.Equal(t, hex.EncodeToString(sig), "f3a13a555f1f8cd6532716b8f388bd4e9d8ed0b252743e923114c0c6cbfe414cf791c8e859afd3c12009ecf2cb20dacf01636d80823bcdbd9ec1ce59afe008f0")
	require.True(t, key.PublicKey().Verify(message, sig))
}

func TestUnitPrivateKeyECDSAFromString(t *testing.T) {
	key, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	key2, err := PrivateKeyFromString(key.String())
	require.NoError(t, err)

	require.Equal(t, key2.String(), key.String())
}

func TestUnitPrivateKeyECDSAFromStringRaw(t *testing.T) {
	key, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	key2, err := PrivateKeyFromStringECSDA(key.StringRaw())
	require.NoError(t, err)

	require.Equal(t, key2.String(), key.String())
}

func TestUnitPublicKeyECDSAFromString(t *testing.T) {
	key, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := key.PublicKey()
	publicKey2, err := PublicKeyFromStringECDSA(publicKey.String())
	require.NoError(t, err)

	require.Equal(t, publicKey2.String(), publicKey.String())
}

func TestUnitPublicKeyECDSAFromStringRaw(t *testing.T) {
	key, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := key.PublicKey()
	publicKey2, err := PublicKeyFromStringECDSA(publicKey.StringRaw())
	require.NoError(t, err)

	require.Equal(t, publicKey2.String(), publicKey.String())
}

func TestUnitPrivateKeyECDSASignTransaction(t *testing.T) {
	newKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	txID := TransactionIDGenerate(AccountID{Account: 123})

	tx, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTransactionID(txID).
		SetInitialBalance(newBalance).
		SetMaxAutomaticTokenAssociations(100).
		Freeze()
	require.NoError(t, err)

	_, err = newKey.SignTransaction(&tx.Transaction)
	require.NoError(t, err)
}

func TestUnitPublicKeyFromPrivateKeyString(t *testing.T) {
	key, err := PrivateKeyFromStringECSDA("3030020100300706052b8104000a04220420d790c27a81d745ad3340e27dacedc982d1f9252c0d7a4582da9847e2094603d4")
	require.NoError(t, err)
	require.Equal(t, "302f300706052b8104000a032400042102b46925b64940f5d7d3f394aba914c05f1607fa42e9e721afee0770cb55797d99", key.PublicKey().String())
}

func TestUnitPublicKeyToEthereumAddress(t *testing.T) {
	byt, err := hex.DecodeString("03af80b90d25145da28c583359beb47b21796b2fe1a23c1511e443e7a64dfdb27d")
	require.NoError(t, err)
	key, err := PublicKeyFromBytesECDSA(byt)
	ethereumAddress := key.ToEthereumAddress()
	require.Equal(t, ethereumAddress, "627306090abab3a6e1400e9345bc60c78a8bef57")
}
