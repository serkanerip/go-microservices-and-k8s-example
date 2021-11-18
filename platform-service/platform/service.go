package platform

import "log"

type PlatformService struct {
	platformRepository PlatformRepository
	commandDataClient  CommandDataClient
	messageBusClient   MessageBusClient
}

func NewPlatformService(pr PlatformRepository, c CommandDataClient, mbc MessageBusClient) *PlatformService {
	return &PlatformService{
		platformRepository: pr,
		commandDataClient:  c,
		messageBusClient:   mbc,
	}
}

func (p *PlatformService) GetAll() (platforms []PlatformDTO, err error) {
	platforms, err = p.platformRepository.GetAll()
	return
}

func (p *PlatformService) GetPlatformById(id string) (dto *PlatformDTO, err error) {
	dto, err = p.platformRepository.GetPlatformById(id)
	return
}

func (p *PlatformService) CreatePlatform(dto CreatePlatformDTO) (*PlatformDTO, error) {
	id, err := p.platformRepository.CreatePlatform(dto)
	if err != nil {
		return nil, err
	}
	outDto := &PlatformDTO{
		Id:        id,
		Name:      dto.Name,
		Publisher: dto.Publisher,
		Cost:      dto.Cost,
	}

	// sync message
	p.commandDataClient.SendPlatformToCommand(*outDto)

	// async message
	go func() {
		err := p.messageBusClient.PublishNewPlatform(PlatformPublishedDTO{
			Id:    outDto.Id,
			Name:  outDto.Name,
			Event: "Platform_Published",
		})
		if err != nil {
			log.Printf("cannot publish new platform err is: %err", err)
		}
		log.Println("platform published to message bus")
	}()

	return outDto, nil
}
