package repository

import (
	"gitlab.com/mas-dhimas/xlsx-prime-monthly-reporting/config"
	"gitlab.com/mas-dhimas/xlsx-prime-monthly-reporting/pkg/http"
)

type Repository struct {
	httpClient *http.HttpClient
	configURI  config.NetworkSourceData
}

func NewRepository(httpClient *http.HttpClient, configUri config.NetworkSourceData) Repository {
	return Repository{
		httpClient,
		configUri,
	}
}
