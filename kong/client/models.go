package client

type KongService struct {
	Id             string `json:"id,omitempty"`
	CreatedAt      int64  `json:"created_at,omitempty"`
	UpdatedAt      int64  `json:"updated_at,omitempty"`
	ConnectTimeout int    `json:"connect_timeout,omitempty"`
	Name           string `json:"name,omitempty"`
	Retries        int    `json:"retries,omitempty"`
	ReadTimeout    int    `json:"read_timeout,omitempty"`
	WriteTimeout   int    `json:"write_timeout,omitempty"`

	// Kong's api treats `url` as a write-only property.
	// This is useful for creating or updating a service (simply supply the url instead of four other fields),
	// However, in the interest of a consistent model, this package only exposes a url field, for both reading and writing.
	// The other fields (protocol, host, port, path) are in the struct only so that the url field can be populated.
	Url      string `json:"url,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Path     string `json:"path,omitempty"`
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
	Id         string                 `json:"id,omitempty"`
	ServiceId  string                 `json:"service_id,omitempty"`
	RouteId    string                 `json:"route_id,omitempty"`
	ConsumerId string                 `json:"consumer_id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty"`
	Enabled    bool                   `json:"enabled,omitempty"`
	CreatedAt  int64                  `json:"created_at,omitempty"`
}
