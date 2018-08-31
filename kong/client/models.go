package client

type KongService struct {
	Id             string
	CreatedAt      int64
	UpdatedAt      int64
	ConnectTimeout int
	Name           string
	Protocol       string
	Host           string
	Port           int
	Path           string
	Retries        int
	ReadTimeout    int
	WriteTimeout   int
	Url            string
}

type KongServiceReference struct {
	Id string
}

type KongRoute struct {
	Id            string
	CreatedAt     int64
	UpdatedAt     int64
	Protocols     []string
	Methods       []string
	Hosts         []string
	Paths         []string
	RegexPriority int
	StripPath     bool
	Service       KongServiceReference
}

type KongPlugin struct {
	Id         string
	ServiceId  string
	ConsumerId string
	Name       string
	Config     map[string]interface{}
	Enabled    bool
	CreatedAt  int64
}
