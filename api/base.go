package api

import (
	"context"
	"os"

	"github.com/sapphi-red/go-traq"
)

var (
	accessToken string
	auth        context.Context
	client      *traq.APIClient
)

func init() {
	accessToken = os.Getenv("ACCESS_TOKEN")
	auth = context.WithValue(context.Background(), traq.ContextAccessToken, accessToken)
	client = traq.NewAPIClient(traq.NewConfiguration())
}
