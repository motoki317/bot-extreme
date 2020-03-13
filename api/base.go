package api

import (
	"context"
	"github.com/sapphi-red/go-traq"
	"os"
)

var (
	accessToken string
	auth        context.Context
	client      *openapi.APIClient
)

func init() {
	accessToken = os.Getenv("ACCESS_TOKEN")
	auth = context.WithValue(context.Background(), openapi.ContextAccessToken, accessToken)
	client = openapi.NewAPIClient(openapi.NewConfiguration())
}
