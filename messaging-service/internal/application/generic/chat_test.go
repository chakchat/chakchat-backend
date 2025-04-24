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

	m, err := json.Marshal(chatInfo)
	if err != nil {
		t.Fatal(err)
	}

	if string(m) != `{"blocked_by":null}` {
		t.Fatalf("Got %s", string(m))
	}
}