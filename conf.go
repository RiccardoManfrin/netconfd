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

//Conf models application config
type Conf struct {
	Global Global `json:"global"`
}
