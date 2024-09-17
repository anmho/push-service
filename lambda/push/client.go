package push

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"send-to/expo"
)

type SendReceipt struct {
}

type ScheduleReceipt struct {
}

type SchedulePushRequest struct {
	ExpoPushToken string
}

type Client interface {
	SendPush(ctx context.Context, request SendPushRequest) (SendReceipt, error)
	SchedulePush(ctx context.Context, request SchedulePushRequest) (ScheduleReceipt, error)
	UpdatePushToken(ctx context.Context, appID string, userID string, expoPushToken string) error
}

type client struct {
	dynamoClient *dynamodb.Client
	expoClient   *expo.Client
}

func MakeClient(dynamoClient *dynamodb.Client, expoClient *expo.Client) Client {
	return &client{
		dynamoClient: dynamoClient,
		expoClient:   expoClient,
	}
}

type SendPushRequest struct {
	// RecipientTokens is the list of device expo push tokens to send a notification to.
	RecipientTokens []string
	Title           string
	Body            string
	Data            map[string]any
}

func (c *client) SendPush(ctx context.Context, request SendPushRequest) (SendReceipt, error) {
	batch := make([]string, 1)
	batches := make([][]string, 1)
	for i, token := range request.RecipientTokens {
		batch = append(batch, token)
		if len(batch) == 100 || i == len(request.RecipientTokens)-1 {
			batches = append(batches, batch)
			batch = make([]string, 1)
		}
	}

	for _, batch := range batches {
		err := c.expoClient.SendPush(ctx, expo.Message{
			To:               batch,
			ContentAvailable: nil,
			Data:             request.Data,
			Title:            request.Title,
			Body:             request.Body,
			TTL:              nil,
			Expiration:       nil,
			Priority:         nil,
			Subtitle:         nil,
			Sound:            nil,
			Badge:            nil,
			ChannelID:        nil,
			CategoryID:       nil,
			MutableContent:   nil,
		})
		if err != nil {
			return SendReceipt{}, err
		}
	}

	return SendReceipt{}, nil
}
func (c *client) SchedulePush(ctx context.Context, request SchedulePushRequest) (ScheduleReceipt, error) {

	return ScheduleReceipt{}, nil
}
func (c *client) UpdatePushToken(ctx context.Context, appID string, userID string, expoPushToken string) error {
	return nil
}
