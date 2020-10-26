package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func (p *Provider) ScaleToZeroAfter(service string, timeout uint64) error {

}

func ScaleToZero(service string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()

	if err != nil {
		return fmt.Errorf("%+v", "Could not connect to docker API")
	}

	cli.ContainerList(ctx)
}

func (p *Provider) scaleToZero(container types.Container) error {

	container
}

func (p *Provider) Scale(number uint64) {
	if number != 1 {
		return fmt.Errorf("Scaling for Docker is only supported for number=1. Got number=%d", number)
	}


}

func getContainerByName(name string, cli *client) types.Container {
	ctx := context.Background()

	containers = cli.ContainerList(ctx)

	containers.
}