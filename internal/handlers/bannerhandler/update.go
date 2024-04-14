package bannerhandler

import (
	"banner_service/internal/domains"
	"banner_service/internal/models"
	"banner_service/internal/services/bannerservice"
	"banner_service/pkg/http/response"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

func (h *bannerHandler) PatchBannerId(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		response.JSONError(w, http.StatusBadRequest, "id must be integer and positive", h.log)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	var bannerInfo models.Banner
	err = json.Unmarshal(body, &bannerInfo)
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "bad request", h.log)
		return
	}

	bannerWithTagIDs := &domains.BannerWithTagIDs{
		Banner: domains.Banner{Content: bannerInfo.Content, IsActive: bannerInfo.IsActive, FeatureID: bannerInfo.FeatureID},
		TagIDs: bannerInfo.TagIDs,
	}

	err = h.service.UpdateBannerByID(r.Context(), uint32(id), bannerWithTagIDs)
	if err != nil {
		switch {
		case errors.Is(err, bannerservice.ErrInvalidFeatureID):
			response.JSONError(w, http.StatusBadRequest, "feature id must be positive", h.log)
			return
		case errors.Is(err, bannerservice.ErrInvalidTagID):
			response.JSONError(w, http.StatusBadRequest, "tag id must be positive", h.log)
			return
		case errors.Is(err, bannerservice.ErrNoTagIDs):
			response.JSONError(w, http.StatusBadRequest, "no tag ids", h.log)
			return
		case errors.Is(err, bannerservice.ErrNotJSON):
			response.JSONError(w, http.StatusBadRequest, "content is not json", h.log)
			return
		case errors.Is(err, bannerservice.ErrAlreadyExists):
			response.JSONError(w, http.StatusBadRequest, "there is already a pair that uniquely identifies the banner", h.log)
			return
		case errors.Is(err, bannerservice.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
			return
		case errors.Is(err, context.DeadlineExceeded):
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}

		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	w.WriteHeader(http.StatusOK)
}
