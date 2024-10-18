package helpers

import (
	"log/slog"
	"net/http"
)

type serverErrorResponse struct {
	Error string `json:"error"`
}

func ServerError(w http.ResponseWriter, logger *slog.Logger, err error) {
	response := serverErrorResponse{
		Error: "something went wrong",
	}

	logger.Error("server error", "err", err)
	_ = WriteJSON(w, http.StatusInternalServerError, response)
}
