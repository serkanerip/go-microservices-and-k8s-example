package platform

import "github.com/serkanerip/platform-service/internal/mongodb"

func PlatformToDTO(p Platform) *PlatformDTO {
	return &PlatformDTO{
		Id:        p.Id,
		Name:      p.Name,
		Publisher: p.Publisher,
		Cost:      p.Cost,
	}
}

func PlatformDTOToModel(dto PlatformDTO) *Platform {
	return &Platform{
		Id:        dto.Id,
		Name:      dto.Name,
		Publisher: dto.Publisher,
		Cost:      dto.Cost,
	}
}

func PlatformToPersistence(p Platform) *mongodb.Platform {
	return &mongodb.Platform{
		Name:      p.Name,
		Publisher: p.Publisher,
		Cost:      p.Cost,
	}
}

func PersistenceToPlatformModel(p mongodb.Platform) *Platform {
	return &Platform{
		Id:        p.Id.Hex(),
		Name:      p.Name,
		Publisher: p.Publisher,
		Cost:      p.Cost,
	}
}

func PersistenceToPlatformDTO(p mongodb.Platform) *PlatformDTO {
	return &PlatformDTO{
		Id:        p.Id.Hex(),
		Name:      p.Name,
		Publisher: p.Publisher,
		Cost:      p.Cost,
	}
}
