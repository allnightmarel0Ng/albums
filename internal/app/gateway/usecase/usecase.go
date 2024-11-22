package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/domain/model"
)

type GatewayUseCase interface {
	Authorization(string) (model.AuthorizationResponse)
}

type gatewayUseCase struct {
	authorizationPort string
}

func NewGatewayUseCase(authorizationPort string) GatewayUseCase {
	return &gatewayUseCase{
		authorizationPort: authorizationPort,
	}
}

func (g *gatewayUseCase) interprocessCommunicationError() model.AuthorizationResponse {
	return model.AuthorizationResponse{
		Code: http.StatusInternalServerError,
		Error: "interprocess communication error",
	}
}

func (g *gatewayUseCase) Authorization(authHeader string) model.AuthorizationResponse {
	request, err := http.NewRequest("GET", fmt.Sprintf("http://authorization:%s/", g.authorizationPort), nil)
	if err != nil {
		return g.interprocessCommunicationError()
	}

	request.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return g.interprocessCommunicationError()
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return g.interprocessCommunicationError()
	}

	var data model.AuthorizationResponse

	json.Unmarshal(body, &data)
	return data
}
