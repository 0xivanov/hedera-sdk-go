//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitTransactionReceiptQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkyk")
	transactionID := TransactionIDGenerate(accountID)
	assert.NoError(t, err)

	receiptQuery := NewTransactionReceiptQuery().
		SetTransactionID(transactionID)

	err = receiptQuery._ValidateNetworkOnIDs(client)
	assert.NoError(t, err)
}

func TestUnitTransactionReceiptQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	transactionID := TransactionIDGenerate(accountID)
	assert.NoError(t, err)

	receiptQuery := NewTransactionReceiptQuery().
		SetTransactionID(transactionID)

	err = receiptQuery._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}