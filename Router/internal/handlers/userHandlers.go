package handlers

import (
	"context"
	"encoding/json"
	userv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/user"
	"log/slog"
	"net/http"
	"time"
)

func GetUsersWithMessagesCountHandler(logger *slog.Logger, userClient *userv1.UserClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()

		w.Header().Set("Content-Type", "application/json")

		usersResponse, err := (*userClient).GetUsersMessagesCount(ctx, &userv1.GetUsersRequest{})

		err = json.NewEncoder(w).Encode(usersResponse)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
