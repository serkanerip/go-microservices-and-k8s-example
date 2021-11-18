package platform

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type PlatformRouter struct {
	platformService PlatformService
}

func NewPlatformRouter(ps PlatformService) *PlatformRouter {
	return &PlatformRouter{
		platformService: ps,
	}
}

func (p *PlatformRouter) SetupRoutes(router chi.Router) {
	router.Route("/platforms", func(r chi.Router) {
		r.Get("/", p.GetAllPlatformsController)
		r.Post("/", p.CreatePlatformController)
		r.Get("/{id}", p.GetPlatformByIdController)
	})
}

func (p *PlatformRouter) GetAllPlatformsController(w http.ResponseWriter, r *http.Request) {
	platforms, err := p.platformService.GetAll()
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		log.Printf("Cannot get platforms err is: %v", err)
		w.WriteHeader(500)
		return
	}
	b, err := json.Marshal(&platforms)
	if err != nil {
		log.Printf("Cannot marshal platforms err is: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Write(b)
}

func (p *PlatformRouter) GetPlatformByIdController(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	platform, err := p.platformService.GetPlatformById(chi.URLParam(r, "id"))

	if err != nil {
		log.Printf("cannot get platform err is: %v", err)
		w.WriteHeader(500)
		return
	}
	b, err := json.Marshal(&platform)
	if err != nil {
		log.Printf("cannot marshal platform err is: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Write(b)
}

func (p *PlatformRouter) CreatePlatformController(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("cannot read body err is: %v", err)
		w.WriteHeader(500)
		return
	}
	var dto CreatePlatformDTO
	if err = json.Unmarshal(b, &dto); err != nil {
		log.Printf("cannot unmarshal body to dto err is: %v", err)
		w.WriteHeader(500)
		return
	}
	platform, err := p.platformService.CreatePlatform(dto)
	if err != nil {
		log.Printf("cannot create platform err is: %v", err)
		w.WriteHeader(500)
		return
	}

	b, err = json.Marshal(&platform)
	if err != nil {
		log.Printf("cannot marshal platform err is: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Write(b)
}
