package push

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"log/slog"
	"os"
	"send-to/expo"
	"time"
)

const (
	AWSSchedulingExpressionFormat = "2006-01-02T15:04:05"
)

type SendReceipt struct {
}

type ScheduleReceipt struct {
}

type Service interface {
	SendPush(ctx context.Context, request NotificationRequest) ([]*expo.PushReceipt, error)
	SchedulePush(ctx context.Context, pushTime time.Time, request NotificationRequest) (*ScheduleReceipt, error)
	UpdatePushToken(ctx context.Context, appID string, userID string, expoPushToken string) error
}

type service struct {
	dynamoClient    *dynamodb.Client
	schedulerClient *scheduler.Client
	expoClient      *expo.Client
}

func MakeService(dynamoClient *dynamodb.Client, schedulerClient *scheduler.Client, expoClient *expo.Client) Service {
	return &service{
		dynamoClient:    dynamoClient,
		schedulerClient: schedulerClient,
		expoClient:      expoClient,
	}
}

type NotificationRequest struct {
	// RecipientPushTokens is the list of device expo push tokens to send a notification to.
	RecipientPushTokens []string       `json:"recipient_push_tokens"`
	Title               string         `json:"title"`
	Body                string         `json:"body"`
	Data                map[string]any `json:"data,omitempty"`
}

func (s *service) SendPush(ctx context.Context, request NotificationRequest) ([]*expo.PushReceipt, error) {
	batch := make([]string, 0)
	batches := make([][]string, 0)
	for i, token := range request.RecipientPushTokens {
		batch = append(batch, token)
		if len(batch) > 0 && (len(batch) == 100 || i == len(request.RecipientPushTokens)-1) {
			batches = append(batches, batch)
			batch = make([]string, 0)
		}
	}

	receipts := make([]*expo.PushReceipt, 0)

	for _, batch := range batches {
		pushResponse, err := s.expoClient.SendPush(ctx, expo.Message{
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
			return nil, err
		}
		for _, receipt := range pushResponse.Data {
			receipts = append(receipts, &receipt)
		}
	}

	return receipts, nil
}
func (s *service) SchedulePush(ctx context.Context, pushTime time.Time, request NotificationRequest) (*ScheduleReceipt, error) {
	notification, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	slog.Info("environment variables",
		slog.String("SCHEDULE_ROLE_ARN", os.Getenv("SCHEDULE_ROLE_ARN")),
		slog.String("EVENTBUS_ARN", os.Getenv("EVENTBUS_ARN")),
	)
	slog.Info("input passed in scheduler", slog.String("input", string(notification)))

	result, err := s.schedulerClient.CreateSchedule(ctx, &scheduler.CreateScheduleInput{
		FlexibleTimeWindow: &types.FlexibleTimeWindow{
			Mode: types.FlexibleTimeWindowModeOff,
		},
		Name:               aws.String(uuid.New().String() + "-" + "NotificationSchedule"),
		ScheduleExpression: aws.String(fmt.Sprintf(`at(%s)`, pushTime.Format(AWSSchedulingExpressionFormat))),
		Target: &types.Target{
			//Arn:     aws.String(os.Getenv("EVENTBUS_ARN")),
			//RoleArn: aws.String(os.Getenv("SCHEDULE_ROLE_ARN")),
			Arn:     aws.String(os.Getenv("SEND_SCHEDULED_NOTIFICATION_ARN")),
			RoleArn: aws.String(os.Getenv("INVOKE_SEND_SCHEDULED_NOTIFICATION_ROLE_ARN")),
			//EventBridgeParameters: &types.EventBridgeParameters{
			//	DetailType: aws.String("ScheduledNotification"),
			//	Source:     aws.String("scheduler.notifications"),
			//},
			Input: aws.String(string(notification)),
		},
		ActionAfterCompletion: types.ActionAfterCompletionDelete,
		Description:           aws.String("schedules a notification"),
		State:                 types.ScheduleStateEnabled,
	})
	slog.Info("event bridge event", slog.Any("result", result))

	if err != nil {
		slog.Error("after creating schedule", slog.Any("error", err))
		return nil, err
	}
	return &ScheduleReceipt{}, nil
}
func (s *service) UpdatePushToken(ctx context.Context, appID string, userID string, expoPushToken string) error {
	return nil
}
