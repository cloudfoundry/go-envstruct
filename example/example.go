package main

import (
	"encoding/json"
	"fmt"

	envstruct "code.cloudfoundry.org/go-envstruct"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Credentials) UnmarshalEnv(data string) error {
	return json.Unmarshal([]byte(data), c)
}

type HostInfo struct {
	Credentials Credentials `env:"CREDENTIALS, required, noreport"`
	IP          string      `env:"HOST_IP,     required"`
	Port        int         `env:"HOST_PORT"`
}

func main() {
	hi := HostInfo{Port: 80}

	err := envstruct.Load(&hi)
	if err != nil {
		panic(err)
	}

	envstruct.WriteReport(&hi)

	fmt.Printf("Credentials: %+v\n", hi.Credentials)
}
