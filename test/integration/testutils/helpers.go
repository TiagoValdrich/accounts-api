package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

const (
	TestDocument string = "41184478007"
)

func POST(t *testing.T, app *fiber.App, path string, body any) (*http.Response, []byte) {
	t.Helper()

	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", path, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = resp.Body.Close()
	require.NoError(t, err)

	return resp, respBody
}

func GET(t *testing.T, app *fiber.App, path string) (*http.Response, []byte) {
	t.Helper()

	req := httptest.NewRequest("GET", path, nil)

	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = resp.Body.Close()
	require.NoError(t, err)

	return resp, respBody
}

func ParseJSON(t *testing.T, body []byte, dest any) {
	t.Helper()
	require.NoError(t, json.Unmarshal(body, dest))
}
