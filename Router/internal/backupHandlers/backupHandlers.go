package backupHandlers

import (
	"context"
	"encoding/json"
	"fmt"
	backupv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/backup"
	commonv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/common"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func ExportToDirHandler(logger *slog.Logger, backupClient *backupv1.BackupClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawExportType := r.URL.Query().Get("export-type")
		var exportType backupv1.ExportType
		if strings.ToLower(rawExportType) == "csv" {
			exportType = backupv1.ExportType_CSV
		} else if strings.ToLower(rawExportType) == "parquet" {
			exportType = backupv1.ExportType_PARQUET
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "incorrect export type")
			return
		}

		var chatFilter commonv1.MessagesFilter
		err := json.NewDecoder(r.Body).Decode(&chatFilter)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()

		response, err := (*backupClient).ExportToDir(ctx, &backupv1.ExportToDirRequest{
			Type:          exportType,
			MessageFilter: &chatFilter,
		})

		if err != nil || !response.Ok {
			if err != nil {
				logger.Error(err.Error())
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})
}
