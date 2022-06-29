package app

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/arrow2nd/nekome/config"
	flag "github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
)

// postTweet : ツイートを投稿
func (a *App) postTweet(args []string) error {
	var (
		quote  string
		reply  string
		editor string
		images []string
	)

	// フラグを設定
	f := flag.NewFlagSet("tweet", flag.ContinueOnError)
	f.StringVarP(&quote, "quote", "q", "", "Specify the ID of the tweet to quote")
	f.StringVarP(&reply, "reply", "r", "", "Specify the ID of the tweet to which you are replying")
	f.StringVarP(&editor, "editor", "e", os.Getenv("EDITOR"), "Specify the editor to start (Default is $EDITOR)")
	f.StringSliceVarP(&images, "image", "i", nil, "Image to be attached")

	if err := f.Parse(args); err != nil {
		return err
	}

	text := f.Arg(1)

	// エディタを起動
	if text == "" {
		t, err := a.editTweet(editor)
		if err != nil {
			return err
		}

		text = t
	}

	text = trimEndNewline(text)
	if text == "" {
		return nil
	}

	post := func() {
		var mediaIDs []string = nil

		// 画像をアップロード
		if images != nil {
			ids, err := a.uploadImages(images)
			if err != nil {
				shared.SetErrorStatus("Upload Image", err.Error())
				return
			}

			mediaIDs = ids
		}

		// 投稿
		if err := shared.api.PostTweet(text, quote, reply, mediaIDs); err != nil {
			shared.SetErrorStatus("Tweet", err.Error())
			return
		}

		shared.SetStatus("Tweeted", text)
	}

	// 確認画面が不要ならそのままツイート
	if !shared.conf.Settings.Feature.Confirm["Tweet"] {
		post()
		return nil
	}

	shared.ReqestPopupModal(&ModalOpt{
		title:  "Do you want to tweet?",
		text:   text,
		onDone: post,
	})

	return nil
}

// uploadImages : 画像をアップロード
func (a *App) uploadImages(images []string) ([]string, error) {
	imagesCount := len(images)

	containsGIF := find(images, func(v string) bool {
		return strings.HasSuffix(strings.ToLower(v), ".gif")
	})

	// 複数の画像と一緒にGIFをアップロードしようとしていないか
	if containsGIF && imagesCount > 1 {
		return nil, errors.New("gif images cannot be attached with other images")
	}

	// 画像の枚数チェック
	if imagesCount > 4 {
		return nil, errors.New("you can attach up to 4 images")
	}

	eg, ctx := errgroup.WithContext(context.Background())
	ch := make(chan string, imagesCount)

	for _, image := range images {
		// 拡張子のチェック
		ext := strings.ToLower(path.Ext(image))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
			return nil, fmt.Errorf("unsupported extensions (%s)", image)
		}

		image := image

		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				rawImage, err := ioutil.ReadFile(image)
				if err != nil {
					return fmt.Errorf("failed to load file (%s)", image)
				}

				base64Image := base64.StdEncoding.EncodeToString(rawImage)
				res, err := shared.api.UploadImage(base64Image)
				if err != nil {
					return fmt.Errorf("upload failed (%s): %w", image, err)
				}

				ch <- res.MediaIDString

				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	close(ch)

	mediaIds := []string{}
	for id := range ch {
		mediaIds = append(mediaIds, id)
	}

	return mediaIds, nil
}

// editTweet ツイートをエディタで編集
func (a *App) editTweet(editor string) (string, error) {
	dir, err := config.GetConfigDir()
	if err != nil {
		return "", err
	}

	tmpFile := path.Join(dir, ".tmp")
	if _, err := os.Create(tmpFile); err != nil {
		return "", err
	}

	// エディタを起動
	cmd := exec.Command(editor, tmpFile)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	a.app.Suspend(func() {
		err = cmd.Run()
	})

	if err != nil {
		return "", fmt.Errorf("failed to open editor (%s) : %w", editor, err)
	}

	// 一時ファイル読み込み
	bytes, err := ioutil.ReadFile(tmpFile)
	if err != nil {
		return "", err
	}

	// 一時ファイル削除
	if err := os.Remove(tmpFile); err != nil {
		return "", err
	}

	return string(bytes), nil
}
