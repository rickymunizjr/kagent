package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	v1alpha2 "github.com/kagent-dev/kagent/go/api/v1alpha2"
	"github.com/kagent-dev/kagent/go/internal/httpserver/handlers"
	kclient "github.com/kagent-dev/kagent/go/pkg/client"
	"github.com/kagent-dev/kagent/go/pkg/client/api"
)

func TestHandleListSupportedModels_IncludesOpenAICodexModel(t *testing.T) {
	handler := handlers.NewModelHandler(&handlers.Base{})

	req := httptest.NewRequest(http.MethodGet, "/api/models", nil)
	w := httptest.NewRecorder()

	handler.HandleListSupportedModels(&testErrorResponseWriter{w}, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp api.StandardResponse[kclient.ProviderModels]
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	openAIModels, ok := resp.Data[v1alpha2.ModelProviderOpenAI]
	require.True(t, ok, "OpenAI models list missing")

	found := false
	for _, model := range openAIModels {
		if model.Name == "gpt-5.1-codex-max" {
			found = true
			require.True(t, model.FunctionCalling, "function calling should be enabled")
			break
		}
	}

	require.True(t, found, "gpt-5.1-codex-max should be listed for OpenAI")
}
