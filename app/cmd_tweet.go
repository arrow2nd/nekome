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
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func (a *App) newTweetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tweet",
		Long:    "Post a tweet",
		Example: "tweet [text] [--quote <tweet id>] [--reply <tweet id>] [--editor <editor>] [--image <path...>]",
		RunE:    a.execTweetCmd,
	}

	flags := cmd.Flags()
	flags.StringP("quote", "q", "", "Specify the ID of the tweet to quote")
	flags.StringP("reply", "r", "", "Specify the ID of the tweet to which you are replying")
	flags.StringP("editor", "e", os.Getenv("EDITOR"), "Specify the editor to start (Default is $EDITOR)")
	flags.StringSliceP("image", "i", nil, "Image to be attached")

	return cmd
}

func (a *App) execTweetCmd(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()
	text := ""

	if len(args) == 0 {
		editor, _ := flags.GetString("editor")

		// エディタを起動
		t, err := a.editTweet(editor)
		if err != nil {
			return err
		}

		text = t
	} else {
		text = args[0]
	}

	text = trimEndNewline(text)
	if text == "" {
		return nil
	}

	quote, _ := flags.GetString("quote")
	reply, _ := flags.GetString("reply")
	images, _ := flags.GetStringSlice("image")

	post := func() {
		var mediaIDs []string

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
