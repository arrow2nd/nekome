package api

import (
	"context"
)

// Follow : ユーザをフォロー
func (a *API) Follow(userID string) error {
	_, err := a.client.UserFollows(context.Background(), a.CurrentUser.ID, userID)

	return checkError(err)
}

// UnFollow : ユーザのフォローを解除
func (a *API) UnFollow(userID string) error {
	_, err := a.client.DeleteUserFollows(context.Background(), a.CurrentUser.ID, userID)

	return checkError(err)
}
