package api

import (
	"log/slog"
	"net/http"
	"send-to/push"
)

func registerRoutes(
	mux *http.ServeMux,
	pushClient push.Client,
) {
	mux.HandleFunc("POST /send-push", handleSendPush(pushClient))
	mux.HandleFunc("POST /schedule-push", func(w http.ResponseWriter, r *http.Request) {
	})
}

type SendPushParams struct {
	RecipientPushTokens []string `json:"recipient_push_tokens" validate:"required"`
	Title               string   `json:"title" validate:"required"`
	Body                string   `json:"body" validate:"required"`
}

func handleSendPush(pushClient push.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := ReadJSON[SendPushParams](r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		resp, err := pushClient.SendPush(ctx, push.SendPushRequest{
			RecipientTokens: params.RecipientPushTokens,
			Title:           params.Title,
			Body:            params.Body,
			Data:            nil,
		})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			slog.Info("sending push", slog.Any("error", err))
			return
		}
		slog.Info("sending push", slog.Any("resp", resp))
		JSON(w, http.StatusOK, resp)
	}
}

func handleSchedulePush() {

}
