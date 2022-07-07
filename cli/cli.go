package cli

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/pflag"
)

// Command : コマンド
type Command struct {
	// Name : コマンド名
	Name string
	// Shorthand : ショートハンド
	Shorthand string
	// Short : 短いヘルプ文
	Short string
	// Long : 長いヘルプ文
	Long string
	// UsageArgs : 使い方（引数）
	UsageArgs string
	// Example : サンプル
	Example string
	// Hidden : コマンドを表示しない
	Hidden bool

	// Validate : 引数のバリデーション関数
	Validate ValidateArgsFunc
	// SetFlag : フラグの設定
	SetFlag func(f *pflag.FlagSet)
	// Run : コマンドの処理
	Run func(c *Command, f *pflag.FlagSet) error
	// Help : ヘルプ関数（オーバーライド）
	Help func(c *Command, h string)

	// children : サブコマンド
	children map[string]*Command
}

// AddCommand : コマンドを追加
func (c *Command) AddCommand(newCmds ...*Command) {
	if c.children == nil {
		c.children = make(map[string]*Command)
	}

	for _, cmd := range newCmds {
		c.children[cmd.Name] = cmd
	}
}

// GetChildren : サブコマンドを取得
func (c *Command) GetChildren() []*Command {
	ls := []*Command{}

	for _, cmd := range c.children {
		if !cmd.Hidden {
			ls = append(ls, cmd)
		}
	}

	sort.Slice(ls, func(i, j int) bool {
		return ls[i].Name < ls[j].Name
	})

	return ls
}

// getAllChidrenCombinations : 全てのサブコマンドの組み合わせ一覧を取得
func getAllChidrenCombinations(prefix string, parent *Command) []string {
	ls := []string{}

	for _, c := range parent.GetChildren() {
		p := prefix + " " + c.Name
		ls = append(ls, p)

		if c.children != nil {
			ls = append(ls, getAllChidrenCombinations(p, c)...)
		}
	}

	return ls
}

// GetChildrenNames : サブコマンド名の一覧を取得
func (c *Command) GetChildrenNames(all bool) []string {
	ls := []string{}

	for _, cmd := range c.GetChildren() {
		ls = append(ls, cmd.Name)

		if all {
			ls = append(ls, getAllChidrenCombinations(cmd.Name, cmd)...)
		}
	}

	return ls
}

// NewFlagSet : flagsetを生成
func (c *Command) NewFlagSet() *pflag.FlagSet {
	f := pflag.NewFlagSet(c.Name, pflag.ContinueOnError)

	if c.SetFlag != nil {
		c.SetFlag(f)
	}

	f.BoolP("help", "h", false, fmt.Sprintf("help for %s", c.Name))

	return f
}

// find : サブコマンドを再帰的に検索
func find(cmd *Command, args []string) (*Command, []string) {
	// 先頭がフラグなら検索終了
	if strings.HasPrefix(args[0], "-") {
		return cmd, args
	}

	for _, c := range cmd.GetChildren() {
		if args[0] != c.Name && args[0] != c.Shorthand {
			continue
		}

		// サブコマンドを持たない, 後ろにコマンドが無いなら検索終了
		if c.children == nil || len(args) <= 1 {
			return c, args[1:]
		}

		return find(c, args[1:])
	}

	return nil, args
}

// Execute : 実行
func (c *Command) Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("no argument")
	}

	cmd := c

	// 先頭がフラグでないなら、該当するコマンドを検索
	if !strings.HasPrefix(args[0], "-") {
		fCmd, fArgs := find(cmd, args)

		if fCmd == nil {
			return fmt.Errorf("command not found: %s", fArgs[0])
		}

		cmd, args = fCmd, fArgs
	}

	// パース
	f := cmd.NewFlagSet()
	if err := f.Parse(args); err != nil {
		return err
	}

	// ヘルプフラグが指定されているなら、ヘルプを表示
	if help, _ := f.GetBool("help"); help {
		c.help(cmd)
		return nil
	}

	// 引数のバリデーション
	if cmd.Validate != nil {
		if err := cmd.Validate(cmd, f.Args()); err != nil {
			return err
		}
	}

	// 実行関数が無い場合、ヘルプを表示
	if cmd.Run == nil {
		c.help(cmd)
		return nil
	}

	return cmd.Run(cmd, f)
}
