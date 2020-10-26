package swarm

import (
	"github.com/docker/docker/client"
)

var dockerClient = client.NewEnvClient()

func (p *Provider) Init() error {

}

func (p *Provider) Provide() error {

}
