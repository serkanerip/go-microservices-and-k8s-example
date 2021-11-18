package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type PlatformRouter struct {
	commandService CommandService
}

func NewPlatformRouter(cs CommandService) *PlatformRouter {
	return &PlatformRouter{
		commandService: cs,
	}
}

func (c *PlatformRouter) SetupRoutes(router chi.Router) {
	router.Route("/platforms", func(r chi.Router) {
		r.Get("/", c.GetAllPlatforms)
		r.Post("/", c.CreatePlatform)
		r.Get("/{platformId}/commands", c.GetCommandsForPlatform)
		r.Post("/{platformId}/commands", c.CreateCommand)
		r.Get("/{platformId}/commands/{commandId}", c.GetCommandOfPlatform)
	})
}

func (c *PlatformRouter) GetAllPlatforms(w http.ResponseWriter, r *http.Request) {
	platforms, err := c.commandService.GetAllPlatforms()
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

func (c *PlatformRouter) GetCommandsForPlatform(w http.ResponseWriter, r *http.Request) {
	var platformId = chi.URLParam(r, "platformId")
	w.Header().Add("Content-Type", "application/json")

	fmt.Printf("--> Get commands of platform with id: %s\n", platformId)
	if ok, err := c.commandService.PlatformExists(platformId); !ok || err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	commands, err := c.commandService.GetCommandsForPlatform(platformId)
	if err != nil {
		log.Printf("Cannot get commands err is: %v", err)
		w.WriteHeader(500)
		return
	}
	b, err := json.Marshal(&commands)
	if err != nil {
		log.Printf("Cannot marshal commands err is: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Write(b)
}

func (c *PlatformRouter) CreatePlatform(w http.ResponseWriter, r *http.Request) {
	log.Print("--> Inbound POST # Command Service")
	w.Header().Add("Content-Type", "application/json")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("cannot read body err is: %v", err)
		w.WriteHeader(500)
		return
	}
	var platform Platform
	if err = json.Unmarshal(b, &platform); err != nil {
		log.Printf("cannot unmarshal body to dto err is: %v", err)
		w.WriteHeader(500)
		return
	}
	err = c.commandService.commandRepo.CreatePlatform(platform)
	if err != nil {
		log.Printf("cannot create platform err is: %v", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (p *PlatformRouter) GetCommandOfPlatform(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var platformId, commandId = chi.URLParam(r, "platformId"), chi.URLParam(r, "commandId")

	if ok, err := p.commandService.PlatformExists(platformId); !ok || err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	command, err := p.commandService.commandRepo.GetCommand(platformId, commandId)
	if err != nil {
		log.Printf("cannot get command err is: %v", err)
		w.WriteHeader(500)
		return
	}
	b, err := json.Marshal(command)
	if err != nil {
		log.Printf("cannot marshal command err is: %v", err)
		w.WriteHeader(500)
		return
	}
	w.Write(b)
}

func (p *PlatformRouter) CreateCommand(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var platformId = chi.URLParam(r, "platformId")

	if ok, err := p.commandService.PlatformExists(platformId); !ok || err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var command Command
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		internalError(w, "cannot read request body", err)
		return
	}
	if err := json.Unmarshal(b, &command); err != nil {
		internalError(w, "cannot unmarshal body", err)
		return
	}
	if err := p.commandService.CreateCommand(platformId, command); err != nil {
		internalError(w, "cannot create command", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func internalError(w http.ResponseWriter, errMsg string, err error) {
	log.Printf("%s err is: %v", errMsg, err)
	w.WriteHeader(500)
}
