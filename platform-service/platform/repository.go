package platform

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

type PlatformRepository interface {
	GetAll() ([]PlatformDTO, error)
	GetPlatformById(string) (*PlatformDTO, error)
	// CreatePlatform returns id, and error
	CreatePlatform(CreatePlatformDTO) (string, error)
	Seed()
}

var db = make([]Platform, 0)

type PlatformInMemoryRepository struct {
}

func (p *PlatformInMemoryRepository) GetAll() ([]PlatformDTO, error) {
	dtos := make([]PlatformDTO, 0)
	for _, platform := range db {
		dtos = append(dtos, *PlatformToDTO(platform))
	}
	return dtos, nil
}

func (p *PlatformInMemoryRepository) GetPlatformById(id string) (*PlatformDTO, error) {
	for _, platform := range db {
		if platform.Id == id {
			return PlatformToDTO(platform), nil
		}
	}
	return nil, fmt.Errorf("platform with id: (%s) not found", id)
}

func (p *PlatformInMemoryRepository) CreatePlatform(dto CreatePlatformDTO) (string, error) {
	platform := Platform{
		Id:        uuid.NewString(),
		Name:      dto.Name,
		Publisher: dto.Publisher,
		Cost:      dto.Cost,
	}

	db = append(db, platform)

	return platform.Id, nil
}

func (p *PlatformInMemoryRepository) Seed() {
	if len(db) > 0 {
		log.Println("--> We already have data")
		return
	}

	log.Println("--> Seeding Data...")
	db = append(db,
		Platform{Id: uuid.NewString(), Name: "Go", Publisher: "Google", Cost: "Free"},
		Platform{Id: uuid.NewString(), Name: "Mysql", Publisher: "Mysql", Cost: "Free"},
		Platform{Id: uuid.NewString(), Name: "Kubernetes", Publisher: "Cloud Native Computing Foundation", Cost: "Free"},
	)
}
