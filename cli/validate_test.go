package cli_test

import (
	"testing"

	"github.com/arrow2nd/nekome/v2/cli"
	"github.com/stretchr/testify/assert"
)

func newValidateCmd(f cli.ValidateArgsFunc) *cli.Command {
	a := newCmd("test")
	a.Validate = f

	r := newCmd("root")
	r.AddCommand(a)

	return r
}

func TestVaidateNoArgs(t *testing.T) {
	r := newValidateCmd(cli.NoArgs())

	t.Run("バリデーションが正しく機能しているか", func(t *testing.T) {
		assert.NoError(t, r.Execute([]string{"test"}))
	})

	t.Run("引数の数が不正な場合にエラーが返るか", func(t *testing.T) {
		assert.EqualError(
			t,
			r.Execute([]string{"test", "a"}),
			"unknown command a for test",
		)
	})
}

func TestVaidateRequireArgs(t *testing.T) {
	r := newValidateCmd(cli.RequireArgs(2))

	t.Run("バリデーションが正しく機能しているか", func(t *testing.T) {
		assert.NoError(t, r.Execute([]string{"test", "a", "b"}))
	})

	t.Run("引数の数が少ない場合にエラーが返るか", func(t *testing.T) {
		assert.EqualError(
			t,
			r.Execute([]string{"test", "a"}),
			"accepts 2 arg(s), received 1",
		)
	})

	t.Run("引数の数が多い場合にエラーが返るか", func(t *testing.T) {
		assert.EqualError(
			t,
			r.Execute([]string{"test", "a", "b", "c"}),
			"accepts 2 arg(s), received 3",
		)
	})
}

func TestVaidateRangeArgs(t *testing.T) {
	r := newValidateCmd(cli.RangeArgs(2, 4))

	t.Run("バリデーションが正しく機能しているか", func(t *testing.T) {
		assert.NoError(t, r.Execute([]string{"test", "a", "b"}))
		assert.NoError(t, r.Execute([]string{"test", "a", "b", "c"}))
		assert.NoError(t, r.Execute([]string{"test", "a", "b", "c", "d"}))
	})

	t.Run("引数がない場合にエラーが返るか", func(t *testing.T) {
		assert.EqualError(
			t,
			r.Execute([]string{"test"}),
			"accepts between 2 and 4 arg(s), received 0",
		)
	})

	t.Run("引数の数が少ない場合にエラーが返るか", func(t *testing.T) {
		assert.EqualError(
			t,
			r.Execute([]string{"test", "a"}),
			"accepts between 2 and 4 arg(s), received 1",
		)
	})

	t.Run("引数の数が多い場合にエラーが返るか", func(t *testing.T) {
		assert.EqualError(
			t,
			r.Execute([]string{"test", "a", "b", "c", "d", "e"}),
			"accepts between 2 and 4 arg(s), received 5",
		)
	})
}
