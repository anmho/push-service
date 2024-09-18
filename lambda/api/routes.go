package api

import (
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"send-to/expo"
	"send-to/push"
	"time"
)

var (
	validate = validator.New(validator.WithRequiredStructEnabled())
)

func registerRoutes(
	mux *http.ServeMux,
	pushService push.Service,
) {
	mux.HandleFunc("POST /send-push", handleSendPush(pushService))
	mux.HandleFunc("POST /schedule-push", handleSchedulePush(pushService))
}

type SendPushParams struct {
	RecipientPushTokens []string `json:"recipient_push_tokens" validate:"required"`
	Title               string   `json:"title" validate:"required"`
	Body                string   `json:"body" validate:"required"`
}

type SendPushResponse struct {
	Receipts []*expo.PushReceipt `json:"receipts"`
}

func handleSendPush(pushService push.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := ReadJSON[SendPushParams](r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		receipts, err := pushService.SendPush(ctx, push.NotificationRequest{
			RecipientPushTokens: params.RecipientPushTokens,
			Title:               params.Title,
			Body:                params.Body,
			Data:                nil,
		})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			slog.Error("sending push", slog.Any("error", err))
			return
		}

		resp := SendPushResponse{Receipts: receipts}
		slog.Info("sending push", slog.Any("resp", resp))
		JSON(w, http.StatusOK, resp)
	}
}

type SchedulePushParams struct {
	SendTime            string   `json:"send_time" validate:"required"`
	RecipientPushTokens []string `json:"recipient_push_tokens" validate:"required,dive,min=1"`
	Title               string   `json:"title" validate:"required"`
	Body                string   `json:"body" validate:"required"`
}

func handleSchedulePush(pushService push.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		params, err := ReadJSON[SchedulePushParams](r.Body)
		slog.Info("handleSchedulePush", slog.Any("params", params))

		if err != nil {
			slog.Error("bad request", slog.Any("error", err))
			http.Error(w, http.StatusText(http.StatusBadRequest)+err.Error(), http.StatusBadRequest)
			return
		}

		sendDate, err := time.Parse(push.AWSSchedulingExpressionFormat, params.SendTime)
		if err != nil {
			slog.Error("parsing timestamp", slog.Any("error", err), slog.String("send_time", params.SendTime))
			http.Error(w, http.StatusText(http.StatusBadRequest)+err.Error(), http.StatusBadRequest)
			return
		}

		result, err := pushService.SchedulePush(ctx, sendDate, push.NotificationRequest{
			RecipientPushTokens: params.RecipientPushTokens,
			Title:               params.Title,
			Body:                params.Body,
			Data:                nil,
		})

		slog.Info("after scheduling the push", slog.Any("result", result), slog.Any("error", err))
		if err != nil {
			slog.Error("scheduling push notification", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSON(w, http.StatusOK, result)
	}
}
