package main

import (
	"crypto/tls"

	envstruct "code.cloudfoundry.org/go-envstruct"
)

type HostInfo struct {
	IP       string `env:"HOST_IP,required"`
	Password string `env:"PASSWORD,noreport"`
	Port     int    `env:"HOST_PORT"`

	TLSConfig *tls.Config
}

func main() {
	hi := HostInfo{Port: 80}

	err := envstruct.Load(&hi)
	if err != nil {
		panic(err)
	}

	envstruct.WriteReport(&hi)
}
