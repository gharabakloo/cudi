package run

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"cudi/internal/app/config"
	docker "cudi/internal/app/docker/setup"
	_ "cudi/internal/app/init"
	"cudi/internal/pkg/domain"
	cleanupUsecase "cudi/internal/pkg/services/cleanup/usecase"
	imageRepository "cudi/internal/pkg/services/image/repository"
)

func Execute() {
	log.Println("Starting...")

	cleanup, err := SetupConfig()
	if err != nil {
		log.Panic(err)
	}

	ctx := context.Background()
	if err := SetupCleanup().CleanupImages(ctx, cleanup); err != nil {
		log.Panic(err)
	}

	log.Println("Finished")
}

func SetupConfig() (domain.Cleanup, error) {
	flag.Parse()

	config.SetVerbose(config.V)
	cleanup, err := config.ParseConfig(config.File)
	if err != nil {
		return domain.Cleanup{}, err
	}

	if config.GetVerbose() {
		fmt.Printf("Cleanup: \n")
		b, err := json.MarshalIndent(cleanup, "", "  ")
		if err != nil {
			return domain.Cleanup{}, err
		}

		fmt.Print(string(b))
	}
	return cleanup, nil
}

func SetupCleanup() domain.CleanupUsecase {
	repository := imageRepository.New(docker.Client())
	return cleanupUsecase.New(repository)
}
