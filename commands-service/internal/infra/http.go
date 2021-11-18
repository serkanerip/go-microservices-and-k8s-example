package infra

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/serkanerip/commands-service/config"
)

type HttpRouter interface {
	SetupRoutes(router chi.Router)
}

type HttpHandler struct {
	r *chi.Mux
}

func NewHttpHandler() *HttpHandler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		message := map[string]string{
			"msg": "Pong!",
		}

		b, err := json.Marshal(&message)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	})
	return &HttpHandler{r: r}
}

func (h *HttpHandler) RegisterRoutes(routers ...HttpRouter) {
	h.r.Route("/api/c", func(chiRouter chi.Router) {
		for _, router := range routers {
			router.SetupRoutes(chiRouter)
		}
	})
}

func (h *HttpHandler) Run() {
	uri := fmt.Sprintf("0.0.0.0:%s", config.ENV.Port)
	if err := http.ListenAndServe(uri, h.r); err != nil {
		log.Fatalf("cannot start server err is :%v", err)
	}
}
