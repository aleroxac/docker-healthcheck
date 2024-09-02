package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type healthcheckConfigs struct {
	Protocol string
	Host     string
	Port     int
	Path     string
}

func newHealthcheckConfigs(protocol string, host string, port_str string, path string) *healthcheckConfigs {
	port, err := strconv.Atoi(port_str)
	if err != nil {
		log.Fatalf("Fail to convert port '%d': %v", port, errInvalidPort)
		return nil
	}

	config := &healthcheckConfigs{
		Protocol: protocol,
		Host:     host,
		Port:     port,
		Path:     path,
	}

	if err := config.Validate(); err != nil {
		log.Fatalf("Validation failure: %v", err)
	}

	return config
}

var (
	errInvalidProtocol = errors.New("invalid protocol, only http or https are allowed")
	errInvalidHost     = errors.New("invalid host")
	errInvalidPort     = errors.New("invalid port; Network ports in TCP and UDP range from number zero up to 65535")
	errInvalidPath     = errors.New("invalid path")
)

var (
	ipv6_regex   = `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
	ipv4_regex   = `^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`
	domain_regex = `^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$`
)

func (c *healthcheckConfigs) Validate() error {
	host_is_valid, _ := regexp.MatchString("localhost"+`|`+ipv4_regex+`|`+ipv6_regex+`|`+domain_regex, c.Host)

	if c.Protocol != "http" {
		fmt.Println("protocol=", c.Protocol)
		return errInvalidProtocol
	} else if !host_is_valid {
		return errInvalidHost
	} else if c.Port < 1 || c.Port > 65535 {
		fmt.Println("port=", c.Port)
		return errInvalidPort
	}

	if c.Path == "/" {
		c.Path = ""
	} else {
		if !strings.HasPrefix(c.Path, "/") {
			return errInvalidPath

		}
	}

	return nil
}

func main() {
	needed_envs := reflect.VisibleFields(reflect.TypeOf(healthcheckConfigs{}))

	missing_envs := 0
	for _, env := range needed_envs {
		env_var_name := fmt.Sprintf("HEALTHCHECK_%v", strings.ToUpper(env.Name))
		if _, present := os.LookupEnv(env_var_name); !present {
			fmt.Printf("Please provide the %s environment variable\n", env_var_name)
			missing_envs = missing_envs + 1
		}
	}

	if missing_envs != 0 {
		return
	}

	config := newHealthcheckConfigs(
		os.Getenv("HEALTHCHECK_PROTOCOL"),
		os.Getenv("HEALTHCHECK_HOST"),
		os.Getenv("HEALTHCHECK_PORT"),
		os.Getenv("HEALTHCHECK_PATH"),
	)

	res, err := http.Get(
		fmt.Sprintf(
			"%s://%s:%d/%s",
			config.Protocol,
			config.Host,
			config.Port,
			config.Path,
		),
	)
	if res != nil {
		fmt.Println(res.StatusCode)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
