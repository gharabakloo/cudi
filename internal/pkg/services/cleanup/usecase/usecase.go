package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"cudi/internal/app/config"
	"cudi/internal/pkg/domain"
)

type usecase struct {
	repository domain.ImageRepository
}

var _ domain.CleanupUsecase = &usecase{}

func New(repository domain.ImageRepository) domain.CleanupUsecase {
	return &usecase{
		repository: repository,
	}
}

func (u *usecase) CleanupImages(ctx context.Context, cleanup domain.Cleanup) error {
	for _, image := range cleanup.Images {
		if image.Type = u.ValidateType(image.Type); image.Type == "" {
			return errors.New("clean up type must be 'separately' or 'together'")
		}

		if image.Type == domain.TypeTogether {
			if err := u.CleanupImage(ctx, image); err != nil {
				return err
			}
			continue
		}

		allRepository, err := u.parseRepository(ctx, image.Repository)
		if err != nil {
			return err
		}

		for repository := range allRepository {
			image.Repository = repository
			if err := u.CleanupImage(ctx, image); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *usecase) ValidateType(imageType string) string {
	if imageType == "" {
		imageType = domain.TypeTogether
	}

	imageType = strings.ToLower(imageType)
	if imageType != domain.TypeTogether && imageType != domain.TypeSeparately {
		return ""
	}
	return imageType
}

func (u *usecase) parseRepository(ctx context.Context, repository string) (map[string]struct{}, error) {
	if config.GetVerbose() {
		fmt.Printf("\n*************************************************************************\n")
		fmt.Printf("FetchImagesByTags: AllTags\n")
	}

	allRepoMap, err := u.FetchImagesByTags(ctx, repository, []string{domain.AllTags})
	if err != nil {
		return nil, err
	}

	allRepository := make(map[string]struct{})
	for key := range allRepoMap {
		repo := strings.Split(key, domain.Colon)[0]
		allRepository[repo] = struct{}{}
	}

	if config.GetVerbose() {
		fmt.Printf("ParseRepository: repository: %+v\n", repository)
		b, err := json.MarshalIndent(allRepository, "", "  ")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Print(string(b))
	}
	return allRepository, nil
}

func (u *usecase) CleanupImage(ctx context.Context, image domain.Image) error {
	if config.GetVerbose() {
		fmt.Printf("\n*************************************************************************\n")
		fmt.Printf("CleanupImage: repository: %+v\nremoveTags: %+v\n\n", image.Repository, image.RemoveTags)
	}

	filteredTags, err := u.FetchFilteredTags(ctx, image)
	if err != nil {
		return err
	}

	for _, filteredTag := range filteredTags {
		if err := u.RemoveImage(ctx, filteredTag); err != nil {
			fmt.Println(err)
			return nil
		}
	}
	return nil
}

func (u *usecase) FetchFilteredTags(ctx context.Context, image domain.Image) ([]string, error) {
	if config.GetVerbose() {
		fmt.Printf("FetchImagesByTags: RemoveTags: %+v\n", image.RemoveTags)
	}

	removeRepoMap, err := u.FetchImagesByTags(ctx, image.Repository, image.RemoveTags)
	if err != nil {
		return nil, err
	}

	if len(removeRepoMap) == 0 {
		return nil, nil
	}

	if config.GetVerbose() {
		fmt.Printf("FetchImagesByTags: KeepTags: %+v\n", image.KeepTags)
	}

	keepRepoMap, err := u.FetchImagesByTags(ctx, image.Repository, image.KeepTags)
	if err != nil {
		return nil, err
	}

	removeTime, err := u.fetchRemoveTime(image.OlderThan)
	if err != nil {
		return nil, err
	}

	keepCounter := 0
	var filteredTags []string
	var createdTags []int64
	for repoTag, created := range removeRepoMap {
		if created > removeTime {
			keepCounter++
			continue
		}

		if _, ok := keepRepoMap[repoTag]; ok {
			keepCounter++
			continue
		}

		filteredTags = append(filteredTags, repoTag)
		createdTags = append(createdTags, created)
	}

	if keepCounter >= image.KeepNumber {
		return filteredTags, nil
	}

	count := image.KeepNumber - keepCounter
	for i := 0; i < count; i++ {
		if len(filteredTags) == 0 {
			break
		}

		u.popLatestTag(&filteredTags, &createdTags)
	}
	return filteredTags, nil
}

func (u *usecase) FetchImagesByTags(ctx context.Context, repository string, tags []string) (map[string]int64, error) {
	repoMap := make(map[string]int64)
	for _, tag := range tags {
		reference := fmt.Sprintf("%s:%s", repository, tag)
		if err := u.FetchImagesByTag(ctx, reference, repoMap); err != nil {
			return nil, err
		}
	}
	return repoMap, nil
}

func (u *usecase) FetchImagesByTag(ctx context.Context, reference string, repoMap map[string]int64) error {
	images, err := u.repository.ReadImages(ctx, reference)
	if err != nil {
		return err
	}

	for _, image := range images {
		for _, repoTag := range image.RepoTags {
			repoMap[repoTag] = image.Created
			if config.GetVerbose() {
				fmt.Printf("FetchImages: created: %v ,repoTag: %+v\n", image.Created, repoTag)
			}
		}
	}

	if config.GetVerbose() {
		fmt.Println()
	}
	return nil
}

func (u *usecase) fetchRemoveTime(olderThan string) (int64, error) {
	if olderThan == "" {
		return time.Now().Unix(), nil
	}

	sep := " "
	timeUnit := strings.Split(olderThan, sep)
	if len(timeUnit) != 2 {
		return 0, errors.New("olderThan format is invalid")
	}

	ts, err := strconv.Atoi(timeUnit[0])
	if err != nil {
		return 0, errors.New("olderThan format is invalid")
	}

	duration := u.unitDuration(timeUnit[1])
	if duration == 0 {
		return 0, errors.New("olderThan format is invalid")
	}
	return time.Now().Add(-time.Duration(ts) * duration).Unix(), nil
}

func (u *usecase) unitDuration(unitSign string) time.Duration {
	switch unitSign {
	case domain.YearSign:
		return domain.Year
	case domain.MonthSign:
		return domain.Month
	case domain.WeekSign:
		return domain.Week
	case domain.DaySign:
		return domain.Day
	case domain.HourSign:
		return domain.Hour
	case domain.MinuteSign:
		return domain.Minute
	case domain.SecondSign:
		return domain.Second
	default:
		return 0
	}
}

func (u *usecase) popLatestTag(filteredTags *[]string, createdTags *[]int64) {
	maxIndex := 0
	maxCreated := int64(-1)
	for index, created := range *createdTags {
		if created > maxCreated {
			maxCreated = created
			maxIndex = index
		}
	}

	*filteredTags = append((*filteredTags)[:maxIndex], (*filteredTags)[maxIndex+1:]...)
	*createdTags = append((*createdTags)[:maxIndex], (*createdTags)[maxIndex+1:]...)
}

func (u *usecase) RemoveImage(ctx context.Context, imageID string) error {
	resp, err := u.repository.DeleteImage(ctx, imageID)
	if err != nil {
		return err
	}

	if config.GetVerbose() {
		fmt.Printf("RemoveImage: imageID: %+v\n", imageID)
		b, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Print(string(b))
	}
	return nil
}
