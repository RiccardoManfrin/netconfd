package main

//HTTPAddr encodes an address for HTTP listen
type HTTPAddr struct {
	Port uint16 `json:"port,omitempty"`
	Host string `json:"host,omitempty"`
}

//Global global config
type Global struct {
	LogLev string   `json:"log_lev,omitempty"`
	Mgmt   HTTPAddr `json:"mgmt"`
}

//Remote encodes an HTTP addr to talk to
type Remote HTTPAddr

//Filepaths is a slice of file paths
type Filepaths []string

//Sync config
type Sync struct {
	Local     HTTPAddr  `json:"local"`
	Remote    Remote    `json:"remote"`
	Filepaths Filepaths `json:"filepaths"`
}

//Conf models application config
type Conf struct {
	Global Global `json:"global"`
	Sync   Sync   `json:"sync"`
}
