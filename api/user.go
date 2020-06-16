package api

import (
	"github.com/antihax/optional"
	openapi "github.com/sapphi-red/go-traq"
)

var (
	usersCache UsersMap
)

type UsersMap map[string]*openapi.User

// GET /users ユーザー一覧を取得します
func GetUsers(canUseCache bool) (UsersMap, error) {
	if canUseCache && usersCache != nil {
		return usersCache, nil
	}

	users, _, err := client.UserApi.GetUsers(auth, &openapi.UserApiGetUsersOpts{
		IncludeSuspended: optional.NewBool(true),
	})
	if err != nil {
		return nil, err
	}

	m := make(map[string]*openapi.User, len(users))
	for _, u := range users {
		user := u
		m[user.Id] = &user
	}
	usersCache = m

	return m, nil
}
