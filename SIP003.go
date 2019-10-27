package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
)

func parseEnv() {
	SSRemoteHost := os.Getenv("SS_REMOTE_HOST")
	SSRemotePort := os.Getenv("SS_REMOTE_PORT")
	SSLocalHost := os.Getenv("SS_LOCAL_HOST")
	SSLocalPort := os.Getenv("SS_LOCAL_PORT")
	SSPluginOptions := os.Getenv("SS_PLUGIN_OPTIONS")
	if len(SSRemoteHost) > 0 {
		*remoteAddress = SSRemoteHost
	}
	if port, err := strconv.Atoi(SSRemotePort); err == nil && port > 0 && port < 65536 {
		*remotePort = port
	}
	if len(SSLocalHost) > 0 {
		*localAddress = SSLocalHost
	}
	if port, err := strconv.Atoi(SSLocalPort); err == nil && port > 0 && port < 65536 {
		*localPort = port
	}
	if len(SSPluginOptions) > 0 {
		parseOptions(SSPluginOptions)
	}
}

func parseOptions(options string) {
	reg := regexp.MustCompile(`\\(.)`)
	i := 0
	length := len(options)
	for {
		if i >= length {
			break
		}

		escape := false

		start := i
		for {
			if i >= length {
				break
			}
			if escape {
				escape = false
				i++
				continue
			}
			if options[i] == '=' {
				break
			}
			if options[i] == '\\' {
				escape = true
			}
			i++
		}
		key := reg.ReplaceAllString(options[start:i], "$1")

		i++ // Skip '='

		if i >= length {
			log.Fatalf("Parse error: %s has no value\n", key)
		}

		start = i
		for {
			if i >= length {
				break
			}
			if escape {
				escape = false
				i++
				continue
			}
			if options[i] == ';' {
				break
			}
			if options[i] == '\\' {
				escape = true
			}
			i++
		}
		value := reg.ReplaceAllString(options[start:i], "$1")

		i++ // Skip ';'

		if key == "socks5Address" || key == "address" {
			*socks5Address = value
		}
		if key == "socks5Port" || key == "port" {
			if port, err := strconv.Atoi(value); err == nil && port > 0 && port < 65536 {
				*socks5Port = port
			}
		}
	}
}
