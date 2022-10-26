package api

import (
	"context"
)

// Mute : ユーザをミュート
func (a *API) Mute(userID string) error {
	_, err := a.client.UserMutes(context.Background(), a.CurrentUser.ID, userID)

	return checkError(err)
}

// UnMute : ユーザのミュートを解除
func (a *API) UnMute(userID string) error {
	_, err := a.client.DeleteUserMutes(context.Background(), a.CurrentUser.ID, userID)

	return checkError(err)
}
