package api

import (
	"context"
	"fmt"
)

// Mute : ユーザをミュート
func (a *API) Mute(userID string) error {
	if _, err := a.client.UserMutes(context.Background(), a.CurrentUser.ID, userID); err != nil {
		return fmt.Errorf("user mute error: %v", err)
	}

	return nil
}

// UnMute : ユーザのミュートを解除
func (a *API) UnMute(userID string) error {
	if _, err := a.client.DeleteUserMutes(context.Background(), a.CurrentUser.ID, userID); err != nil {
		return fmt.Errorf("user mute error: %v", err)
	}

	return nil
}
