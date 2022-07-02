package cli

import "fmt"

// ValidateArgsFunc : 引数のバリデーション関数
type ValidateArgsFunc func(c *Command, args []string) error

// NoArgs : 引数無し
func NoArgs() ValidateArgsFunc {
	return func(c *Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("unknown command %s for %s", args[0], c.Name)
		}

		return nil
	}
}

// RequireArgs : n個の引数をとる
func RequireArgs(n int) ValidateArgsFunc {
	return func(c *Command, args []string) error {
		if len(args) != n {
			return fmt.Errorf("accepts %d arg(s), received %d", n, len(args))
		}

		return nil
	}
}

// RangeArgs : min~max個の引数をとる
func RangeArgs(min, max int) ValidateArgsFunc {
	return func(c *Command, args []string) error {
		if len(args) < min || len(args) > max {
			return fmt.Errorf("accepts between %d and %d arg(s), received %d", min, max, len(args))
		}

		return nil
	}
}
