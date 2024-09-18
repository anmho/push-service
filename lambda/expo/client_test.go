package expo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_SendPush(t *testing.T) {
	tests := []struct {
		desc    string
		message Message

		expectedReceipt PushReceipt
	}{
		{
			desc: "happy path: valid token, title. and body",
			message: Message{
				To:    []string{"ExponentPushToken[1IoX9FHzHNRx6doi7h3nJm]"},
				Title: "Test awesome title",
				Body:  "Test awesome body",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			client := MakeClient()
			resp, err := client.SendPush(context.Background(), tc.message)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Len(t, resp.Data, 1)
			assert.Equal(t, resp.Data[0].Status, "ok")
		})
	}
}
