package setup

import (
	"context"
	"log"

	"github.com/docker/docker/client"
)

var c *client.Client

func Setup(ctx context.Context) error {
	var err error
	if c, err = newClient(); err != nil {
		return err
	}

	if resp, err := c.Ping(ctx); err != nil {
		return err
	} else {
		log.Printf("Ping: %+v\n", resp)
	}
	return nil
}

// Client returns the dockerClient
func Client() *client.Client {
	return c
}

func newClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}
