package entity

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/aleroxac/docker-healthcheck/internal/errors"
)

type HealthcheckConfigs struct {
	Protocol string
	Host     string
	Port     int
	Path     string
}

func NewHealthcheckConfigs(protocol string, host string, port_str string, path string) *HealthcheckConfigs {
	port, err := strconv.Atoi(port_str)
	if err != nil {
		log.Fatalf("Fail to convert port '%d': %v", port, errors.ErrInvalidPort)
		return nil
	}

	config := &HealthcheckConfigs{
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

func (c *HealthcheckConfigs) Validate() error {
	var (
		ipv6_regex   = `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
		ipv4_regex   = `^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`
		domain_regex = `^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$`
	)
	host_is_valid, _ := regexp.MatchString("localhost"+`|`+ipv4_regex+`|`+ipv6_regex+`|`+domain_regex, c.Host)

	if c.Protocol != "http" {
		fmt.Println("protocol=", c.Protocol)
		return errors.ErrInvalidProtocol
	} else if !host_is_valid {
		return errors.ErrInvalidHost
	} else if c.Port < 1 || c.Port > 65535 {
		fmt.Println("port=", c.Port)
		return errors.ErrInvalidPort
	}

	if c.Path == "/" {
		c.Path = ""
	} else {
		if !strings.HasPrefix(c.Path, "/") {
			return errors.ErrInvalidPath

		}
	}

	return nil
}
