package ondemand

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	typeName = "Ondemand"
)

const defaultTimeoutSeconds = 60

// Config the config that holds the service name and the timeout in seconds
type Config struct {
	Name    string
	Timeout uint64
}

// CreateConfig creates a config with its default values
func CreateConfig() *Config {
	return &Config{
		Timeout: defaultTimeoutSeconds,
	}
}

func (c *Config) validate() error {
	if len(c.Name) == 0 {
		return fmt.Errorf("name cannot be null")
	}
	return nil
}

// Ondemand holds the request for the on demand service
type Ondemand struct {
	request string
	name    string
	next    http.Handler
}

// New function creates the configuration
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	err := config.validate()

	if err != nil {
		return err
	}

	return &Ondemand{
		next:    next,
		name:    name,
		request: request,
	}, nil
}

// ServeHTTP retrieve the service status
func (e *Ondemand) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	status, err := getServiceStatus(e.request)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}

	if status == "started" {
		// Service started forward request
		e.next.ServeHTTP(rw, req)

	} else if status == "starting" {
		// Service starting, notify client
		rw.WriteHeader(http.StatusAccepted)
		rw.Write([]byte("Service is starting..."))
	} else {
		// Error
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Unexpected status answer from ondemand service"))
	}
}

func getServiceStatus(request string) (string, error) {

	// This request wakes up the service if he's scaled to 0
	resp, err := netClient.Get(request)
	if err != nil {
		return "error", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "parsing error", err
	}

	return strings.TrimSuffix(string(body), "\n"), nil
}
