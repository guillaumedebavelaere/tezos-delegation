package tezos_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/tezos"
	"github.com/guillaumedebavelaere/tezos-delegation/pkg/http"
)

type underTest struct {
	client        tezos.API
	mockTransport *httpmock.MockTransport
}

func setupTest(t *testing.T, unmarshalFunc func(data []byte, v any) error) *underTest {
	t.Helper()

	unmarshaler := json.Unmarshal
	if unmarshalFunc != nil {
		unmarshaler = unmarshalFunc
	}

	ut := &underTest{}

	ut.mockTransport = httpmock.NewMockTransport()

	ut.client = tezos.NewClient(&tezos.Config{
		HTTP: http.ClientConfig{
			Debug:   false,
			BaseURL: "https://api.tezos.test/v1",
			Timeout: 5 * time.Second,
		},
	},
		http.WithTransport(ut.mockTransport),
		http.WithUnmarshaller(unmarshaler),
	)

	require.NoError(t, ut.client.Init())

	return ut
}

func TestTezos_NewClient(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, tezos.NewClient(&tezos.Config{}))
	})
}

func TestTezos_Init(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		ut := tezos.NewClient(&tezos.Config{})

		require.NoError(t, ut.Init())
	})
}
