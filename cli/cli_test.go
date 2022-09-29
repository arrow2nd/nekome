package cli_test

import (
	"testing"

	"github.com/arrow2nd/nekome/cli"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func newCmd(n string) *cli.Command {
	return &cli.Command{
		Name: n,
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			return nil
		},
	}
}

func TestAddCommand(t *testing.T) {
	c := newCmd("root")

	t.Run("サブコマンドを追加できるか", func(t *testing.T) {
		assert.Equal(t, 0, len(c.GetChildren()), "追加前にサブコマンドが存在していないか")

		c.AddCommand(newCmd("test"))
		assert.Equal(t, 1, len(c.GetChildren()))
	})
}

func TestGetChildren(t *testing.T) {
	r := newCmd("root")
	r.AddCommand(newCmd("a"))
	r.AddCommand(newCmd("b"))

	t.Run("サブコマンドを取得できるか", func(t *testing.T) {
		children := r.GetChildren()

		assert.Equal(t, "a", children[0].Name)
		assert.Equal(t, "b", children[1].Name)
	})
}

func TestGetChildrenNames(t *testing.T) {
	c := newCmd("c")
	c.Hidden = true

	r := newCmd("root")
	r.AddCommand(newCmd("a"))
	r.AddCommand(newCmd("b"))
	r.AddCommand(c)

	t.Run("サブコマンド名を取得できるか", func(t *testing.T) {
		assert.Equal(t, []string{"a", "b"}, r.GetChildrenNames(false))
	})

	t.Run("非表示のコマンドは除外されているか", func(t *testing.T) {
		assert.NotEqual(t, []string{"a", "b", "c"}, r.GetChildrenNames(false))
	})
}

func TestGetChildrenNamesAll(t *testing.T) {
	b := newCmd("b")
	b.AddCommand(newCmd("a"))

	d := newCmd("d")
	d.AddCommand(b)
	d.AddCommand(newCmd("c"))

	r := newCmd("root")
	r.AddCommand(d)
	r.AddCommand(newCmd("e"))

	t.Run("全ての組み合わせを取得できるか", func(t *testing.T) {
		want := []string{"d", "d b", "d b a", "d c", "e"}
		assert.Equal(t, want, r.GetChildrenNames(true))
	})
}

func TestExecute(t *testing.T) {
	r := newCmd("root")
	c := newCmd("neko")

	c.Run = func(c *cli.Command, f *pflag.FlagSet) error {
		return nil
	}

	r.AddCommand(c)

	t.Run("実行できるか", func(t *testing.T) {
		assert.NoError(t, r.Execute([]string{"neko"}))
	})

	t.Run("引数がない場合にエラーが返るか", func(t *testing.T) {
		assert.EqualError(t, r.Execute([]string{}), "no argument")
	})

	t.Run("コマンドが存在しない場合にエラーが返るか", func(t *testing.T) {
		assert.EqualError(t, r.Execute([]string{"hoge"}), "command not found: hoge")
	})
}

func TestExecute_Args(t *testing.T) {
	r := newCmd("root")
	c := newCmd("neko")

	c.Run = func(c *cli.Command, f *pflag.FlagSet) error {
		assert.Equal(t, "neko", c.Name, "コマンド名が取得できるか")
		assert.Equal(t, "arg", f.Arg(0), "引数の取得ができるか")
		return nil
	}

	r.AddCommand(c)

	t.Run("実行できるか", func(t *testing.T) {
		assert.NoError(t, r.Execute([]string{"neko", "arg"}))
	})
}

func TestExecute_Flag(t *testing.T) {
	c := newCmd("neko")

	c.SetFlag = func(f *pflag.FlagSet) {
		f.BoolP("kawaii", "k", false, "very kawaii flag")
	}

	c.Run = func(c *cli.Command, f *pflag.FlagSet) error {
		kawaii, _ := f.GetBool("kawaii")
		assert.True(t, kawaii)
		return nil
	}

	r := newCmd("root")
	r.AddCommand(c)

	t.Run("フラグが指定できるか", func(t *testing.T) {
		assert.NoError(t, r.Execute([]string{"neko", "--kawaii"}))
	})

	c.Run = func(c *cli.Command, f *pflag.FlagSet) error {
		kawaii, _ := f.GetBool("kawaii")
		assert.False(t, kawaii)
		return nil
	}

	t.Run("フラグが初期化されているか", func(t *testing.T) {
		assert.NoError(t, r.Execute([]string{"neko"}))
	})
}

func TestExecute_Flag_Arg(t *testing.T) {
	c := newCmd("add")

	c.SetFlag = func(f *pflag.FlagSet) {
		f.String("comment", "", "comment")
	}

	c.Run = func(c *cli.Command, f *pflag.FlagSet) error {
		comment, _ := f.GetString("comment")
		assert.Equal(t, "コメント", comment, "フラグの引数が取得できるか")
		return nil
	}

	r := newCmd("root")
	r.AddCommand(c)

	t.Run("実行できるか", func(t *testing.T) {
		assert.NoError(t, r.Execute([]string{"add", "--comment", "コメント"}))
	})
}

func TestExecuteSubCommandArg(t *testing.T) {
	b := newCmd("clone")
	b.Run = func(c *cli.Command, f *pflag.FlagSet) error {
		assert.Equal(t, "nekome", f.Arg(0), "サブコマンドの引数が取得できるか")
		return nil
	}

	a := newCmd("repo")
	a.SetFlag = func(f *pflag.FlagSet) {
		f.Bool("test", false, "")
	}
	a.AddCommand(b)

	r := newCmd("root")
	r.AddCommand(a)

	t.Run("実行できるか", func(t *testing.T) {
		assert.NoError(t, r.Execute([]string{"repo", "clone", "nekome"}))
	})

	t.Run("サブコマンドを持つコマンドを実行できるか", func(t *testing.T) {
		assert.NoError(t, r.Execute([]string{"repo"}))
	})

	a.Run = func(c *cli.Command, f *pflag.FlagSet) error {
		test, _ := f.GetBool("test")
		assert.Equal(t, true, test, "フラグが取得できるか")
		return nil
	}

	t.Run("サブコマンドを持つコマンドのフラグをパースできるか", func(t *testing.T) {
		assert.NoError(t, r.Execute([]string{"repo", "--test"}))
	})
}
