package bybit

import (
	"encoding/json"
	"testing"

	"github.com/hirokisan/bybit/v2/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebsocketV5Public_Bookticker(t *testing.T) {
	t.Run("spot", func(t *testing.T) {
		respBody := map[string]interface{}{
			"topic": "bookticker.BTCUSDT",
			"ts":    1673437259336,
			"type":  "snapshot",
			"data": map[string]interface{}{
				"s":  "BTCUSDT",
				"bp": "17440",
				"bq": "0.0002",
				"ap": "17440.01",
				"aq": "0.21302",
				"t":  1673437259336,
			},
		}
		bytesBody, err := json.Marshal(respBody)
		require.NoError(t, err)

		category := CategoryV5Spot

		server, teardown := testhelper.NewWebsocketServer(
			testhelper.WithWebsocketHandlerOption(V5WebsocketPublicPathFor(category), bytesBody),
		)
		defer teardown()

		wsClient := NewTestWebsocketClient().
			WithBaseURL(server.URL)

		svc, err := wsClient.V5().Public(category)
		require.NoError(t, err)

		{
			_, err := svc.SubscribeBookticker(
				V5WebsocketPublicBooktickerParamKey{
					Symbol: SymbolV5BTCUSDT,
				},
				func(response V5WebsocketPublicBooktickerResponse) error {
					assert.Equal(t, respBody["topic"], response.Topic)
					testhelper.Compare(t, respBody["data"], response.Data)
					return nil
				},
			)
			require.NoError(t, err)
		}

		assert.NoError(t, svc.Run())
		assert.NoError(t, svc.Ping())
		assert.NoError(t, svc.Close())
	})
}
