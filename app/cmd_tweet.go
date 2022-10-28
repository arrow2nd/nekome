package app

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"syscall"

	"github.com/arrow2nd/nekome/v2/cli"
	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
	"golang.org/x/term"
)

func (a *App) newTweetCmd() *cli.Command {
	longHelp := `Post a tweet.

If the tweet statement is omitted, the internal editor is invoked if from the TUI, or the external editor if from the CLI.
Also, setting 'feature.use_external_editor' to true in preferences.toml will launch the external editor even from the TUI.`

	example := `tweet にゃーん --image cat.png,dog.png
  echo "にゃーん" | nekome tweet`

	return &cli.Command{
		Name:      "tweet",
		Shorthand: "t",
		Short:     "Post a tweet",
		Long:      longHelp,
		UsageArgs: "[text]",
		Example:   example,
		SetFlag: func(f *pflag.FlagSet) {
			f.StringP("quote", "q", "", "quotes the tweet with the specified ID")
			f.StringP("reply", "r", "", "send a reply to the tweet with the specified ID")
			f.StringP("editor", "e", os.Getenv("EDITOR"), "specify the editor to use for editing")
			f.StringSliceP("image", "i", nil, "attach the image (if there is more than one comma separated)")
			f.BoolP("clipboard", "c", false, "attach the image in the clipboard")
		},
		Run: a.execTweetCmd,
	}
}

func (a *App) execTweetCmd(c *cli.Command, f *pflag.FlagSet) error {
	pref := shared.conf.Pref
	text := ""

	// 標準入力を受け取る
	if f.NArg() == 0 && !term.IsTerminal(int(syscall.Stdin)) {
		stdin, _ := ioutil.ReadAll(os.Stdin)
		text = string(stdin)
	} else {
		text = strings.Join(f.Args(), " ")
	}

	editor, _ := f.GetString("editor")
	quoteId, _ := f.GetString("quote")
	replyId, _ := f.GetString("reply")
	images, _ := f.GetStringSlice("image")

	if text == "" {
		// テキストエリアを開く
		if !shared.isCommandLineMode && !pref.Feature.UseExternalEditor {
			a.view.ShowTextArea(pref.Text.TweetTextAreaHint, func(s string) {
				execPostTweet(s, quoteId, replyId, images)
			})
			return nil
		}

		// エディタを開く
		t, err := a.editTweetExternalEditor(editor)
		if err != nil {
			return err
		}

		text = t
	}

	execPostTweet(text, quoteId, replyId, images)

	return nil
}

// editTweetExternalEditor : 外部エディタでツイートを編集する
func (a *App) editTweetExternalEditor(editor string) (string, error) {
	// 一時ファイル作成
	tmpFilePath := path.Join(os.TempDir(), ".nekome_tweet_tmp")
	if _, err := os.Create(tmpFilePath); err != nil {
		return "", err
	}

	if err := a.openExternalEditor(editor, tmpFilePath); err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		return "", err
	}

	if err := os.Remove(tmpFilePath); err != nil {
		return "", err
	}

	return string(bytes), nil
}

// execPostTweet : ツイートを投稿
func execPostTweet(text, quoteId, replyId string, images []string) {
	text = trimEndNewline(text)

	// 文章も画像もない場合キャンセル
	if text == "" && len(images) == 0 {
		return
	}

	post := func() {
		var mediaIids []string

		// 画像をアップロード
		if images != nil {
			ids, err := uploadImages(images)
			if err != nil {
				shared.SetErrorStatus("Upload Image", err.Error())
				return
			}

			mediaIids = ids
		}

		if err := shared.api.PostTweet(text, quoteId, replyId, mediaIids); err != nil {
			shared.SetErrorStatus("Tweet", err.Error())
			return
		}

		shared.SetStatus("Tweeted", text)
	}

	// 確認画面不要 or コマンドラインモードならそのまま実行
	if shared.isCommandLineMode || !shared.conf.Pref.Confirm["tweet"] {
		post()
		return
	}

	operationType := "tweet"

	if replyId != "" {
		operationType = "reply"
	} else if quoteId != "" {
		operationType = "quote tweet"
	}

	shared.ReqestPopupModal(&ModalOpt{
		title: fmt.Sprintf(
			"Do you want to post a [%s]%s[-:-:-]?",
			shared.conf.Style.App.EmphasisText,
			operationType,
		),
		text:   text,
		onDone: post,
	})
}

// uploadImages : 画像をアップロード
func uploadImages(images []string) ([]string, error) {
	imagesCount := len(images)

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
