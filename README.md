[![GoDoc](https://godoc.org/github.com/bradylove/envstruct?status.png)](https://godoc.org/github.com/bradylove/envstruct)

# envstruct

envstruct is a simple library for populating values on structs from environment
variables.

## Usage

Export some environment variables.

```
$ export HOST_IP="127.0.0.1"
$ export HOST_PORT="443"
$ export PASSWORD="abc123"
```

*Note:* The environment variables are case
sensitive. The casing of the set environment variable must match the casing in
the struct tag.

Write some code. In this example, `Ip` requires that the `HOST_IP` environment variable is set to non empty value and
`Port` defaults to `80` if `HOST_PORT` is an empty value. Then we use the `envstruct.WriteReport()` to print a
table with a report of what fields are on the struct, the type, the environment variable where the value is read from,
whether or not it is required, and the value. If using when `envstruct.WriteReport()` you wish to omit a sensitive
value you can add `noreport` to the struct tag as shown with `Password`

```
package main

import "github.com/bradylove/envstruct"

type HostInfo struct {
	IP       string `env:"HOST_IP,required"`
	Password string `env:"PASSWORD,noreport"`
	Port     int    `env:"HOST_PORT"`
}

func main() {
	hi := HostInfo{Port: 80}

	err := envstruct.Load(&hi)
	if err != nil {
		panic(err)
	}

	envstruct.WriteReport(&hi)
}
```

Run your code and rejoice!

```
$ go run example/example.go
FIELD NAME:  TYPE:   ENV:       REQUIRED:  VALUE:
Ip           string  HOST_IP    true       127.0.0.1
Password     string  PASSWORD   false      (OMITTED)
Port         int     HOST_PORT  false      80
```

## Supported Types

- [x] string
- [x] bool (`true` and `1` results in true value, anything else results in false value)
- [x] int
- [x] int8
- [x] int16
- [x] int32
- [x] int64
- [x] uint
- [x] uint8
- [x] uint16
- [x] uint32
- [x] uint64
- [ ] float32
- [ ] float64
- [ ] complex64
- [ ] complex128
- [x] []slice (Slices of any other supported type. Environment variable should have coma separated values)
- [x] time.Duration

## Running Tests

Run tests using ginkgo.

```
$ go get github.com/apoydence/eachers
$ go get github.com/onsi/ginkgo/ginkgo
$ go get github.com/onsi/gomega
$ ginkgo
```
