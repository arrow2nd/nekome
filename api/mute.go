package api

import (
	"context"
	"fmt"
)

// Mute ユーザをミュート
func (a *API) Mute(userID string) error {
	client := a.newClient(a.CurrentUser.Token)

	if _, err := client.UserMutes(context.Background(), a.CurrentUser.ID, userID); err != nil {
		return fmt.Errorf("user mute error: %v", err)
	}

	return nil
}

// UnMute ユーザのミュートを解除
func (a *API) UnMute(userID string) error {
	client := a.newClient(a.CurrentUser.Token)

	if _, err := client.DeleteUserMutes(context.Background(), a.CurrentUser.ID, userID); err != nil {
		return fmt.Errorf("user mute error: %v", err)
	}

	return nil
}
