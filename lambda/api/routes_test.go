package api

// go generate

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"send-to/expo"
	"send-to/push"
	"testing"
	"time"
)

func Test_handleSendPush(t *testing.T) {
	tests := []struct {
		desc   string
		params SendPushParams

		expectedStatus int
	}{
		{
			desc: "happy path: valid title and body",
			params: SendPushParams{
				RecipientPushTokens: []string{
					"ExponentPushToken[1IoX9FHzHNRx6doi7h3nJm]",
				},
				Title: "route title",
				Body:  "route body",
			},

			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			expoClient := expo.MakeClient()
			pushService := push.MakeService(nil, nil, expoClient)

			ts := httptest.NewServer(handleSendPush(pushService))
			defer ts.Close()

			data, _ := json.Marshal(tc.params)
			res, err := http.Post(ts.URL, "application/json", bytes.NewReader(data))
			require.NoError(t, err)

			assert.Equal(t, res.StatusCode, http.StatusOK)
		})
	}
}

func Test_handleSchedulePush(t *testing.T) {
	tests := []struct {
		desc   string
		params SchedulePushParams

		expectedStatus int
	}{
		{
			desc: "happy path: valid title and body",
			params: SchedulePushParams{
				RecipientPushTokens: []string{
					"ExponentPushToken[1IoX9FHzHNRx6doi7h3nJm]",
				},
				Title:    "route title",
				Body:     "route body",
				SendTime: time.Now().Add(time.Second).Format(time.RFC3339),
			},

			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			expoClient := expo.MakeClient()
			pushService := push.MakeService(nil, nil, expoClient)

			ts := httptest.NewServer(handleSchedulePush(pushService))
			defer ts.Close()

			data, _ := json.Marshal(tc.params)
			res, err := http.Post(ts.URL, "application/json", bytes.NewReader(data))
			require.NoError(t, err)

			assert.Equal(t, res.StatusCode, http.StatusOK)
		})
	}
}
