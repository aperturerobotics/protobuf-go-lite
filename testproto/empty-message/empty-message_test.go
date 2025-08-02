package empty_message

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

var testData = Parent{
	Empty: &Parent_Empty{},
}

func TestEmptyMessageProto(t *testing.T) {
	b, err := testData.MarshalVT()
	require.NoError(t, err)

	newData := Parent{}
	err = newData.UnmarshalVT(b)
	require.NoError(t, err)

	require.Condition(t, func() bool { return testData.EqualVT(&newData) })
}

func TestEmptyMessageJson(t *testing.T) {
	b, err := json.Marshal(&testData)
	require.NoError(t, err)

	newData := Parent{}
	err = json.Unmarshal(b, &newData)
	require.NoError(t, err)

	require.Condition(t, func() bool { return testData.EqualVT(&newData) })
}

func TestEmptyMessageText(t *testing.T) {
	// _ = testData.MarshalProtoText()
	t.Skip("MessageText seems to not have Unmarshal for now. skipped")
}
