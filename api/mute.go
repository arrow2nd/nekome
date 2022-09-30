package api

import (
	"context"
)

// Mute : ユーザをミュート
func (a *API) Mute(userId string) error {
	_, err := a.client.UserMutes(context.Background(), a.CurrentUser.ID, userId)

	return checkError(err)
}

// UnMute : ユーザのミュートを解除
func (a *API) UnMute(userId string) error {
	_, err := a.client.DeleteUserMutes(context.Background(), a.CurrentUser.ID, userId)

	return checkError(err)
}
