package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"test/common"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func doSignOutRequest(t *testing.T, refreshToken string) *http.Response {
	type Request struct {
		RefreshJWT string `json:"refresh_token"`
	}
	reqBody, _ := json.Marshal(Request{
		RefreshJWT: refreshToken,
	})
	req, err := http.NewRequest(http.MethodPut, common.ServiceUrl+"/v1.0/sign-out", bytes.NewReader(reqBody))
	require.NoError(t, err)
	req.Header.Add(common.HeaderIdemotencyKey, uuid.NewString())

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}
