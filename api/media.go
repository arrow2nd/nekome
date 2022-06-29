package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// NOTE: 現状、Twitter API v2 のメディアアップロードのエンドポイントが無いため
//       暫定的な対応として v1.1 のエンドポイントを利用

const mediaUploadEndpoint = "https://upload.twitter.com/1.1/media/upload.json"

// Image : 画像詳細
type Image struct {
	ImageType string `json:"image_type"`
	W         int    `json:"w"`
	H         int    `json:"h"`
}

// UploadImageResponse : media/uploadのレスポンス
type UploadImageResponse struct {
	MediaID          int    `json:"media_id"`
	MediaIDString    string `json:"media_id_string"`
	Size             int    `json:"size"`
	ExpiresAfterSecs int    `json:"expires_after_secs"`
	Image            Image  `json:"image"`
}

// UploadImage : 画像をアップロードする
func (a *API) UploadImage(base64Image string) (*UploadImageResponse, error) {
	v := url.Values{}
	v.Add("media_data", base64Image)

	res, err := a.client.Client.PostForm(mediaUploadEndpoint, v)
	if err != nil {
		return nil, fmt.Errorf("upload image response: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error: %s", res.Status)
	}

	decoder := json.NewDecoder(res.Body)
	raw := &UploadImageResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, fmt.Errorf("upload image decode error: %w", err)
	}

	return raw, nil
}
