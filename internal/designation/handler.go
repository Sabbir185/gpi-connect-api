package designation

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Sabbir185/gpi/pkg/httpx"
)

const (
	defaultLimit   = 6
	maxLimit       = 20
	minQueryLength = 2
)

type DesignationResponse struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Hint  *string `json:"hint"`
}

type designationHandler struct {
	repo DesignationRepository
}

func NewDesignationHandler(repo DesignationRepository) *designationHandler {
	return &designationHandler{repo: repo}
}

func (h *designationHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 || limit > maxLimit {
		limit = defaultLimit
	}

	result := []DesignationResponse{}
	if len([]rune(query)) >= minQueryLength {
		designations, err := h.repo.Search(r.Context(), query, limit)
		if err != nil {
			httpx.SendAppError(w, r, err)
			return
		}
		for _, d := range designations {
			result = append(result, DesignationResponse{ID: strconv.FormatInt(d.ID, 10), Title: d.Title, Hint: d.Hint})
		}
	}

	httpx.SendSuccess(w, http.StatusOK, httpx.CodeDataFetch, "Designations fetched", result)
}
