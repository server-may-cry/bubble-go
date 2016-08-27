package protocol

import "github.com/server-may-cry/bubble-go/models"

type IndexResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	I       int    `json:"i"`
}

type TestResponse struct {
	Test models.Test `json:"test"`
}

type RedisResponse struct {
	Ping string `json:"ping"`
}
