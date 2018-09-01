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
	Id string `json:"id"`
}

type KongRoute struct {
	Id            string               `json:"id,omitempty"`
	CreatedAt     int64                `json:"created_at,omitempty"`
	UpdatedAt     int64                `json:"updated_at,omitempty"`
	Protocols     []string             `json:"protocols,omitempty"`
	Methods       []string             `json:"methods,omitempty"`
	Hosts         []string             `json:"hosts,omitempty"`
	Paths         []string             `json:"paths,omitempty"`
	RegexPriority int                  `json:"regex_priority"`
	StripPath     bool                 `json:"strip_path"`
	PreserveHost  bool                 `json:"preserve_host"`
	Service       KongServiceReference `json:"service"`
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
