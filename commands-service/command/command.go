package command

type Command struct {
	Id          string
	HowTo       string `json:"how_to"`
	CommandLine string `json:"command_line"`
	PlatformId  string
	Platform    Platform
}

type CommandService struct {
	commandRepo CommandRepo
}

func NewCommandService(repo CommandRepo) *CommandService {
	return &CommandService{commandRepo: repo}
}

func (c *CommandService) GetAllPlatforms() ([]Platform, error) {
	return c.commandRepo.GetAllPlatforms()
}

func (c *CommandService) CreatePlatform(platform Platform) error {
	return c.commandRepo.CreatePlatform(platform)
}

func (c *CommandService) PlatformExists(platformId string) (bool, error) {
	return c.commandRepo.PlatformExists(platformId)
}

func (c *CommandService) GetCommandsForPlatform(platformId string) ([]Command, error) {
	return c.commandRepo.GetCommandsForPlatform(platformId)
}

func (c *CommandService) GetCommand(platformId string, commandId string) (*Command, error) {
	return c.commandRepo.GetCommand(platformId, commandId)
}

func (c *CommandService) CreateCommand(platformId string, command Command) error {
	return c.commandRepo.CreateCommand(platformId, command)
}
