package init

import (
	"cudi/internal/app/config"
	"log"

	"cudi/internal/app/docker"
)

func init() {
	config.InitializeArg()
	docker.Initialize()

	log.Println("Initialization was done.")
}
