package destroy

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang-api-1/internal/lib/api/response"
	"golang-api-1/internal/lib/logger/sl"
	"golang.org/x/exp/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=URLDestroyer
type URLDestroyer interface {
	DestroyURL(alias string) error
}

type Request struct {
	Alias string `json:"alias" validate:"required,alias"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

func New(log *slog.Logger, urlDestroyer URLDestroyer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.destroy.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		err := urlDestroyer.DestroyURL(alias)
		if err != nil {
			log.Info("failed to remove url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to remove url"))

			return
		}

		log.Info("url removed", slog.String("alias", alias))

		responseOk(w, r, alias)
	}
}

func responseOk(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
