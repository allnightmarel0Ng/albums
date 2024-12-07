package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/admin-panel/repository"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/kafka"
	"github.com/allnightmarel0Ng/albums/internal/utils"
)

type AdminPanelUseCase interface {
	Logs(pageNumber uint, pageSize uint) api.Response
	DeleteAlbum(albumID int) api.Response
}

type adminPanelUseCase struct {
	repo        repository.AdminPanelRepository
	producer    *kafka.Producer
	profilePort string
}

func NewAdminPanelUseCase(repo repository.AdminPanelRepository, profilePort string, producer *kafka.Producer) AdminPanelUseCase {
	return &adminPanelUseCase{
		repo:        repo,
		profilePort: profilePort,
		producer:    producer,
	}
}

func (a *adminPanelUseCase) Logs(pageNumber uint, pageSize uint) api.Response {
	offset := (pageNumber - 1) * pageSize
	limit := pageSize

	ctx, cancel := utils.DeadlineContext(2)
	defer cancel()

	count, logs, err := a.repo.GetBuyLogsAndCount(ctx, offset, limit)
	if err != nil {
		return &api.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "db error",
		}
	}

	return &api.BuyLogsResponse{
		Code:      http.StatusOK,
		Logs:      logs,
		LogsCount: count,
	}
}

func (a *adminPanelUseCase) DeleteAlbum(albumID int) api.Response {
	ctx, cancel := utils.DeadlineContext(2)
	defer cancel()

	code, raw := utils.RequestAndParseResponse("GET", fmt.Sprintf("http://profile:%s/owners/%d", a.profilePort, albumID), "", nil)
	if code == http.StatusOK {
		var response api.AlbumOwnersResponse
		err := json.Unmarshal(raw, &response)
		if err == nil {
			success := false
			for i := 0; i < len(response.Ids); i++ {
				err = utils.ProduceNotificationMessage(api.NotificationKafkaMessage{
					Type:      api.Delete,
					UserID:    response.Ids[i],
					AlbumName: response.AlbumName,
					Success:   &success,
				}, a.producer)
				if err != nil {
					log.Printf("unable to produce the message: %s", err.Error())
				}
			}
		} else {
			log.Printf("got an error while parsing the response %s:", err.Error())
		}
	}

	err := a.repo.DeleteAlbum(ctx, albumID)
	if err != nil {
		return &api.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "db error",
		}
	}

	return nil
}
