package expo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
)

const (
	pushNotificationEndpoint = "https://exp.host/--/api/v2/push/send"
)

type Client struct {
	httpClient *http.Client
}

func MakeClient() *Client {
	return &Client{
		httpClient: http.DefaultClient,
	}
}

type Message struct {
	To               []string               `json:"to"`                          // Expo push token(s) for recipient(s)
	ContentAvailable *bool                  `json:"_contentAvailable,omitempty"` // iOS Only: true if the app should start in the background
	Data             map[string]interface{} `json:"data,omitempty"`              // JSON object delivered to the app
	Title            string                 `json:"title"`                       // Title of the notification
	Body             string                 `json:"body"`                        // Body of the notification
	TTL              *int                   `json:"ttl,omitempty"`               // Time to Live in seconds
	Expiration       *int64                 `json:"expiration,omitempty"`        // Unix timestamp of when the message expires
	Priority         *string                `json:"priority,omitempty"`          // 'default', 'normal', or 'high'
	Subtitle         *string                `json:"subtitle,omitempty"`          // iOS Only: Subtitle for the notification
	Sound            *string                `json:"sound,omitempty"`             // iOS Only: Sound to play with the notification
	Badge            *int                   `json:"badge,omitempty"`             // iOS Only: Badge number
	ChannelID        *string                `json:"channelId,omitempty"`         // Android Only: Notification Channel ID
	CategoryID       *string                `json:"categoryId,omitempty"`        // ID of the notification category
	MutableContent   *bool                  `json:"mutableContent,omitempty"`    // iOS Only: If the notification can be intercepted
}

func (c *Client) SendPush(ctx context.Context, message Message) error {
	slog.Info("SendPush invoked")
	body, err := json.Marshal(message)

	client := &http.Client{}
	r, err := http.NewRequest(http.MethodPost, pushNotificationEndpoint, bytes.NewReader(body))

	if err != nil {
		fmt.Println(err)
		return err
	}
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Host", "exp.host")
	r.Header.Set("Accept-Encoding", "gzip, deflate")
	r.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body)

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))

	return nil
}
