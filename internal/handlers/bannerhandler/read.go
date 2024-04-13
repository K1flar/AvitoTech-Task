package bannerhandler

import (
	"banner_service/internal/domains"
	"banner_service/internal/services/bannerservice"
	"banner_service/pkg/filters"
	"banner_service/pkg/http/response"
	"context"
	"errors"
	"net/http"
	"strconv"
)

func (h *bannerHandler) GetBanner(w http.ResponseWriter, r *http.Request) {
	filter := filters.NewBannerFilterFromRequest(r)

	banners, err := h.service.GetBanners(r.Context(), filter)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			w.WriteHeader(http.StatusRequestTimeout)
			return
		}
		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	response.JSONWithMarshal(w, http.StatusOK, banners, h.log)
}

func (h *bannerHandler) GetUserBanner(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(domains.RoleKey("role")).(domains.Role)
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var isAdmin bool
	if role == domains.AdminRole {
		isAdmin = true
	}

	q := r.URL.Query()

	featureID, err := strconv.Atoi(q.Get("feature_id"))
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid feature id", h.log)
		return
	}

	tagID, err := strconv.Atoi(q.Get("tag_id"))
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid tag id", h.log)
		return
	}

	var useLastRevision bool
	if val := q.Get("use_last_revision"); val != "" {
		valb, err := strconv.ParseBool(val)
		if err != nil {
			response.JSONError(w, http.StatusBadRequest, "use_last_revision parameter must be true or false", h.log)
			return
		}
		useLastRevision = valb
	}

	banner, err := h.service.GetBannerByFeatureAndTagID(r.Context(), uint32(featureID), uint32(tagID), useLastRevision, isAdmin)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			w.WriteHeader(http.StatusRequestTimeout)
			return
		case errors.Is(err, bannerservice.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
			return
		}

		response.JSONError(w, http.StatusInternalServerError, "unknown error", h.log)
		return
	}

	response.JSON(w, http.StatusOK, banner.Content.String(), h.log)
}
