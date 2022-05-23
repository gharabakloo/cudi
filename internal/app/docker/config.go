package docker

import (
	"context"
	"log"

	docker "cudi/internal/app/docker/setup"
)

// Initialize function is for setup docker
func Initialize() {
	if err := docker.Setup(context.Background()); err != nil {
		log.Panic(err)
	}
}
