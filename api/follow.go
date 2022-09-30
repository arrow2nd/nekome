package api

import (
	"context"
)

// Follow : ユーザをフォロー
func (a *API) Follow(userId string) error {
	_, err := a.client.UserFollows(context.Background(), a.CurrentUser.ID, userId)

	return checkError(err)
}

// UnFollow : ユーザのフォローを解除
func (a *API) UnFollow(userId string) error {
	_, err := a.client.DeleteUserFollows(context.Background(), a.CurrentUser.ID, userId)

	return checkError(err)
}
