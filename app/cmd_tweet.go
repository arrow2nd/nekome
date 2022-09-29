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

	"github.com/arrow2nd/nekome/cli"
	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
	"golang.org/x/term"
)

func (a *App) newTweetCmd() *cli.Command {
	longHelp := `Post a tweet.

If the tweet statement is omitted, the internal editor is invoked if from the TUI, or the external editor if from the CLI.
Tips: If 'feature.use_external_editor' in preferences.toml is true, an external editor will be launched even from the TUI.

When specifying multiple images, please separate them with commas.
You may attach up to four images at a time.`

	example := `tweet にゃーん --image cute_cat.png,very_cute_cat.png
echo "にゃーん" | nekome tweet`

	return &cli.Command{
		Name:      "tweet",
		Shorthand: "t",
		Short:     "Post a tweet",
		Long:      longHelp,
		UsageArgs: "[text]",
		Example:   example,
		SetFlag: func(f *pflag.FlagSet) {
			f.StringP("quote", "q", "", "specify the ID of the tweet to quote")
			f.StringP("reply", "r", "", "specify the ID of the tweet to which you are replying")
			f.StringP("editor", "e", os.Getenv("EDITOR"), "specify which editor to use (default is $EDITOR)")
			f.StringSliceP("image", "i", nil, "specify the image to attach (if there is more than one comma separated)")
		},
		Run: a.execTweetCmd,
	}
}

func (a *App) execTweetCmd(c *cli.Command, f *pflag.FlagSet) error {
	isTerm := term.IsTerminal(int(syscall.Stdin))
	text := f.Arg(0)

	// 標準入力を受け取る
	if f.NArg() == 0 && !isTerm {
		stdin, _ := ioutil.ReadAll(os.Stdin)
		text = string(stdin)
	}

	pref := shared.conf.Pref

	editor, _ := f.GetString("editor")
	quoteId, _ := f.GetString("quote")
	replyId, _ := f.GetString("reply")
	images, _ := f.GetStringSlice("image")

	if text == "" {
		// テキストエリアを開く
		if isTerm && !pref.Feature.UseExternalEditor {
			a.view.ShowTextArea(pref.Text.TweetTextAreaHint, func(s string) {
				execPostTweet(s, quoteId, replyId, images)
			})
			return nil
		}

		// エディタを開く
		var err error
		text, err = a.editTweetExternalEditor(editor)
		if err != nil {
			return err
		}
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

	_, containsGif := find(images, func(v string) bool {
		return strings.HasSuffix(strings.ToLower(v), ".gif")
	})

	// GIFと複数画像の同時アップロードを防止（GIFは動画扱い）
	if containsGif && imagesCount > 1 {
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
