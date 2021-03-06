package traefik_ondemand_plugin

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	typeName = "Ondemand"
)

const defaultTimeoutSeconds = 60

// Net client is a custom client to timeout after 20 seconds if the service is not ready
var netClient = &http.Client{
	Timeout: time.Second * 20,
}

// Config the plugin configuration
type Config struct {
	Names      []string
	ServiceUrl string
	Timeout    uint64
}

// CreateConfig creates a config with its default values
func CreateConfig() *Config {
	return &Config{
		Timeout: defaultTimeoutSeconds,
	}
}

// Ondemand holds the request for the on demand service
type Ondemand struct {
	request string
	name    string
	next    http.Handler
}

func buildRequest(url string, names []string, timeout uint64) (string, error) {
	// TODO: Check url validity
	request := fmt.Sprintf("%s?names=%s&timeout=%d", url, strings.Join(names, ","), timeout)
	return request, nil
}

// New function creates the configuration
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.ServiceUrl) == 0 {
		return nil, fmt.Errorf("serviceUrl cannot be null")
	}

	if config.Names == nil || len(config.Names) == 0 {
		return nil, fmt.Errorf("names cannot be empty")
	}

	request, err := buildRequest(config.ServiceUrl, config.Names, config.Timeout)

	if err != nil {
		return nil, fmt.Errorf("error while building request")
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

	i := 0

	for i = 0; status == "starting" && i < 40; {
		time.Sleep(250 * time.Millisecond)
		status, err = getServiceStatus(e.request)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}

		i++
	}

	if status == "started" {
		// Service started, forward the request
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
