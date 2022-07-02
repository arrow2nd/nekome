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

	"github.com/arrow2nd/nekome/cli"
	"github.com/arrow2nd/nekome/config"
	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
)

// newTweetCmd : tweetコマンド生成
func (a *App) newTweetCmd() *cli.Command {
	return &cli.Command{
		Name:      "tweet",
		Shorthand: "t",
		Short:     "Post a tweet",
		Long: `Post a tweet.
If you omit the tweet statement, the editor will be activated.
You cannot tweet only images.`,
		Example: "tweet にゃーん --image cute_cat.png,very_cute_cat.png",
		SetFlag: func(f *pflag.FlagSet) {
			f.StringP("quote", "q", "", "specify the ID of the tweet to quote")
			f.StringP("reply", "r", "", "specify the ID of the tweet to which you are replying")
			f.StringP("editor", "e", os.Getenv("EDITOR"), "specify the editor to start (default is $EDITOR)")
			f.StringSliceP("image", "i", nil, "image to be attached (if there is more than one comma separated)")
		},
		Run: a.execTweetCmd,
	}
}

// execTweetCmd : tweetコマンド処理
func (a *App) execTweetCmd(c *cli.Command, f *pflag.FlagSet) error {
	text := f.Arg(0)

	// ツイート文が無いなら、エディタを起動
	if text == "" {
		editor, _ := f.GetString("editor")

		t, err := a.editTweet(editor)
		if err != nil {
			return err
		}

		text = t
	}

	// 末尾の改行を削除
	text = trimEndNewline(text)
	if text == "" {
		return nil
	}

	quote, _ := f.GetString("quote")
	reply, _ := f.GetString("reply")
	images, _ := f.GetStringSlice("image")

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

	// 確認画面不要 or コマンドラインモードならそのまま実行
	if shared.isCommandLineMode || !shared.conf.Settings.Feature.Confirm["Tweet"] {
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

// editTweet : エディタを起動してツイートを編集
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

	if shared.isCommandLineMode {
		err = cmd.Run()
	} else {
		a.app.Suspend(func() {
			err = cmd.Run()
		})
	}

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
