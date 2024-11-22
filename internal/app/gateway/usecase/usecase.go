package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type GatewayUseCase interface {
	Authorization(string) (string, int, error)
}

type gatewayUseCase struct {
	authorizationPort string
}

func NewGatewayUseCase(authorizationPort string) GatewayUseCase {
	return &gatewayUseCase{
		authorizationPort: authorizationPort,
	}
}

func (g *gatewayUseCase) interprocessCommunicationError() (string, int, error) {
	return "", http.StatusInternalServerError, errors.New("interprocess communication error")
}

func (g *gatewayUseCase) Authorization(authHeader string) (string, int, error) {
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

	var data struct {
		Code int `json:"code"`
		Jwt string `json:"jwt"`
	}

	json.Unmarshal(body, &data)
	return data.Jwt, http.StatusOK, nil
}