package push

import (
	"context"
	"github.com/stretchr/testify/assert"
	"send-to/expo"
	"testing"
)

func TestClient_SendPush(t *testing.T) {
	tests := []struct {
		desc    string
		request NotificationRequest

		expectedReceipts int
	}{
		{
			desc: "happy path: recipient is reachable and token is valid.",
			request: NotificationRequest{
				RecipientPushTokens: []string{"ExponentPushToken[1IoX9FHzHNRx6doi7h3nJm]"},
				Title:               "Test Awesome title 2",
				Body:                "Test awesome body 2",
				Data:                nil,
			},
			expectedReceipts: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			expoClient := expo.MakeClient()
			pushService := MakeService(nil, nil, expoClient)
			receipts, err := pushService.SendPush(context.Background(), tc.request)
			assert.NotNil(t, receipts)
			assert.Equal(t, tc.expectedReceipts, len(receipts))
			assert.NoError(t, err)
		})

	}
}
