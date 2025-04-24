package generic

import (
	"encoding/json"
	"testing"
)

func TestUpdateContentMarshalJSON(t *testing.T) {
	updateContent := UpdateContent{
		TextMessage: &TextMessageContent{
			Text: "test",
		},
	}

	enc, err := json.Marshal(updateContent)
	if err != nil {
		t.Fatal(err)
	}

	if string(enc) != `{"text":"test"}` {
		t.Fatalf("got: %s", enc)
	}
}