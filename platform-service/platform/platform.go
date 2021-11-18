package platform

type Platform struct {
	Id        string
	Name      string
	Publisher string
	Cost      string
}

type PlatformDTO struct {
	Id        string
	Name      string
	Publisher string
	Cost      string
}

type CreatePlatformDTO struct {
	Name      string
	Publisher string
	Cost      string
}

type PlatformPublishedDTO struct {
	Id    string
	Name  string
	Event string
}
