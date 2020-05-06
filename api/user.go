package api

import (
	"github.com/antihax/optional"
	openapi "github.com/sapphi-red/go-traq"
)

// GET /users ユーザー一覧を取得します
func GetUsers() (users []openapi.User, err error) {
	users, _, err = client.UserApi.GetUsers(auth, &openapi.UserApiGetUsersOpts{
		IncludeSuspended: optional.NewBool(true),
	})
	return
}
