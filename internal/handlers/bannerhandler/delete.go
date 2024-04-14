package bannerhandler

import (
	"banner_service/internal/services/bannerservice"
	"banner_service/pkg/http/response"
	"context"
	"errors"
	"net/http"
	"strconv"
)

func (h *bannerHandler) DeleteBannerId(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		response.JSONError(w, http.StatusBadRequest, "id must be integer and positive", h.log)
		return
	}

	err = h.service.DeleteBannerByID(r.Context(), uint32(id))
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		case errors.Is(err, bannerservice.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
			return
		}

		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
