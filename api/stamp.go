package api

import (
	"github.com/antihax/optional"
	openapi "github.com/sapphi-red/go-traq"
)

// GET /stamps スタンプ一覧を取得します
func GetStamps() (stamps []openapi.Stamp, err error) {
	stamps, _, err = client.StampApi.GetStamps(auth, &openapi.StampApiGetStampsOpts{
		IncludeUnicode: optional.NewBool(true),
	})
	return
}
