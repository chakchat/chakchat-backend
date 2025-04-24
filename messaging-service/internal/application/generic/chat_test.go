package generic

import (
	"encoding/json"
	"testing"
)

func TestChatMarshalJSON(t *testing.T) {
	chatInfo := ChatInfo{
		Personal: &PersonalInfo{
			BlockedBy: nil,
		},
	}

	m, err := json.MarshalIndent(chatInfo, "", "	")
	if err != nil {
		t.Fatal(err)
	}

	t.Fatal(string(m))
}