package hedera

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializeStakersQuery(t *testing.T) {
	query := NewAccountStakersQuery().
		SetAccountID(AccountID{Account: 3}).
		Query

	assert.Equal(t, `cryptoGetProxyStakers:{header:{}accountID:{accountNum:3}}`, strings.ReplaceAll(query.pb.String(), " ", ""))
}

func TestAccountStakersQuery_Execute(t *testing.T) {
	client := newTestClient(t)

	_, err := NewAccountStakersQuery().
		SetAccountID(client.GetOperatorAccountID()).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.Error(t, err)
}
