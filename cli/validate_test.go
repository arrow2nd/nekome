package cli_test

import (
	"testing"

	"github.com/arrow2nd/nekome/cli"
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

	// 正常
	assert.NoError(
		t,
		r.Execute([]string{"test"}),
		"引数の数が正しい",
	)

	// 異常
	assert.EqualError(
		t,
		r.Execute([]string{"test", "a"}),
		"unknown command a for test",
		"引数の数が不正",
	)
}

func TestVaidateRequireArgs(t *testing.T) {
	r := newValidateCmd(cli.RequireArgs(2))

	// 正常
	assert.NoError(
		t,
		r.Execute([]string{"test", "a", "b"}),
		"引数の数が正しい",
	)

	// 異常
	assert.EqualError(
		t,
		r.Execute([]string{"test", "a"}),
		"accepts 2 arg(s), received 1",
		"引数の数が少ない",
	)

	assert.EqualError(
		t,
		r.Execute([]string{"test", "a", "b", "c"}),
		"accepts 2 arg(s), received 3",
		"引数の数が多い",
	)
}

func TestVaidateRangeArgs(t *testing.T) {
	r := newValidateCmd(cli.RangeArgs(2, 4))

	// 正常
	assert.NoError(
		t,
		r.Execute([]string{"test", "a", "b"}),
		"引数の数が正しい（2つ）",
	)

	assert.NoError(
		t,
		r.Execute([]string{"test", "a", "b", "c"}),
		"引数の数が正しい（3つ）",
	)

	assert.NoError(
		t,
		r.Execute([]string{"test", "a", "b", "c", "d"}),
		"引数の数が正しい（4つ）",
	)

	// 異常
	assert.EqualError(
		t,
		r.Execute([]string{"test"}),
		"accepts between 2 and 4 arg(s), received 0",
		"引数が無い",
	)

	assert.EqualError(
		t,
		r.Execute([]string{"test", "a"}),
		"accepts between 2 and 4 arg(s), received 1",
		"引数が少ない（1つ）",
	)

	assert.EqualError(
		t,
		r.Execute([]string{"test", "a", "b", "c", "d", "e"}),
		"accepts between 2 and 4 arg(s), received 5",
		"引数の数が多い（5つ）",
	)
}
