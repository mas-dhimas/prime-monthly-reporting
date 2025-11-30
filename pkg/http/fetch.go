package http

import (
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type HttpClient struct {
	client *http.Client
	token  string
}

func NewHttpClient(token string) *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		token: token,
	}
}

func (h *HttpClient) Fetch(uri string) ([]byte, error) {
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", h.token)

	response, err := h.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	logrus.Info("success get data from : ", uri)

	data, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Error("error io read data: ", err)
		return nil, err
	}
	logrus.Info("success io read all data")

	return data, err
}
