package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"syscall"

	"github.com/arrow2nd/nekome/v2/cli"
	"github.com/skanehira/clipboard-image/v2"
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
			f.BoolP("clipboard", "c", false, "attach the image in the clipboard (if the --image is specified, it takes precedence)")
		},
		Run: a.execTweetCmd,
	}
}

func (a *App) execTweetCmd(c *cli.Command, f *pflag.FlagSet) error {
	pref := shared.conf.Pref
	text := ""

	if f.NArg() == 0 && !term.IsTerminal(int(syscall.Stdin)) {
		// 標準入力を受け取る
		stdin, _ := io.ReadAll(os.Stdin)
		text = string(stdin)
	} else {
		// 引数を全てスペースで連結
		text = strings.Join(f.Args(), " ")
	}

	if text == "" {
		// テキストエリアを開く
		if !shared.isCommandLineMode && !pref.Feature.UseExternalEditor {
			a.view.ShowTextArea(pref.Text.TweetTextAreaHint, func(s string) {
				execPostTweet(f, s)
			})
			return nil
		}

		editor, _ := f.GetString("editor")

		// エディタを開く
		t, err := a.editTweetExternalEditor(editor)
		if err != nil {
			return err
		}

		text = t
	}

	execPostTweet(f, text)

	return nil
}

// editTweetExternalEditor : 外部エディタでツイートを編集する
func (a *App) editTweetExternalEditor(editor string) (string, error) {
	tmpFilePath := path.Join(os.TempDir(), ".nekome_tweet_tmp")
	if _, err := os.Create(tmpFilePath); err != nil {
		return "", err
	}

	if err := a.openExternalEditor(editor, tmpFilePath); err != nil {
		return "", err
	}

	bytes, err := os.ReadFile(tmpFilePath)
	if err != nil {
		return "", err
	}

	if err := os.Remove(tmpFilePath); err != nil {
		return "", err
	}

	return string(bytes), nil
}

// execPostTweet : ツイートを投稿
func execPostTweet(f *pflag.FlagSet, t string) {
	images, _ := f.GetStringSlice("image")
	text := trimEndNewline(t)

	// 文章も画像もない場合キャンセル
	if text == "" && len(images) == 0 {
		return
	}

	quoteID, _ := f.GetString("quote")
	replyID, _ := f.GetString("reply")
	existClipboardImage, _ := f.GetBool("clipboard")

	post := func() {
		mediaIDs := []string{}

		if existImages := len(images) > 0; existImages || existClipboardImage {
			var err error

			if existImages {
				mediaIDs, err = uploadImages(images)
			} else {
				mediaIDs, err = uploadImageFromClipboard()
			}

			if err != nil {
				shared.SetErrorStatus("Media", err.Error())
				return
			}
		}

		if err := shared.api.PostTweet(text, quoteID, replyID, mediaIDs); err != nil {
			shared.SetErrorStatus("Tweet", err.Error())
			return
		}

		statusLabel := "Tweeted"
		if len(mediaIDs) > 0 {
			statusLabel += fmt.Sprintf(" / %d attached images", len(mediaIDs))
		}

		shared.SetStatus(statusLabel, text)
	}

	// 確認画面不要 or コマンドラインモードならそのまま実行
	if shared.isCommandLineMode || !shared.conf.Pref.Confirm["tweet"] {
		post()
		return
	}

	// 実行しようとしている操作名
	operationType := "tweet"
	if replyID != "" {
		operationType = "reply"
	} else if quoteID != "" {
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

// uploadImageFromClipboard : クリップボードの画像をアップロード
func uploadImageFromClipboard() ([]string, error) {
	r, err := clipboard.Read()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, r); err != nil {
		return nil, err
	}

	res, err := shared.api.UploadImage(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}

	return []string{res.MediaIDString}, nil
}

// uploadImages : 複数の画像をアップロード
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
				rawImage, err := os.ReadFile(image)
				if err != nil {
					return fmt.Errorf("failed to load file (%s)", image)
				}

				res, err := shared.api.UploadImage(rawImage)
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

	mediaIDs := []string{}
	for id := range ch {
		mediaIDs = append(mediaIDs, id)
	}

	return mediaIDs, nil
}
