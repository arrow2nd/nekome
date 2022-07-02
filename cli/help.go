package cli

import (
	"fmt"
	"strings"
)

// createHelpText : ヘルプ文を作成
func (c *Command) createHelpText() string {
	newLine := "\n\n"

	// 詳細
	desc := c.Long
	if desc == "" {
		desc = c.Short
	}
	desc += newLine

	// Usage
	usage := fmt.Sprintf("Usage:\n  %s", c.Name)
	if c.children != nil {
		usage += " [command]"
	}
	usage += " [flags]" + newLine

	// Shorthand
	alias := ""
	if c.Shorthand != "" {
		alias = fmt.Sprintf("Shorthand:\n  %s%s", c.Shorthand, newLine)
	}

	// Example
	example := ""
	if c.Example != "" {
		example = fmt.Sprintf("Example:\n  %s%s", c.Example, newLine)
	}

	// Commands
	cmds := ""
	if children := c.GetChildren(); len(children) != 0 {
		maxLen := 0
		for _, cmd := range children {
			if l := len(cmd.Name); maxLen < l {
				maxLen = l
			}
		}

		cmds += "Commands:\n"
		for _, cmd := range children {
			space := strings.Repeat(" ", maxLen-len(cmd.Name)+3)
			cmds += fmt.Sprintf("  %s%s%s\n", cmd.Name, space, cmd.Short)
		}
		cmds += "\n"
	}

	// Flags
	fs := c.newFlagSet().FlagUsages()
	flags := fmt.Sprintf("Flags:\n%s", fs)

	return desc + usage + alias + example + cmds + flags
}

// help : ヘルプを表示
func (c *Command) help(cmd *Command) {
	text := cmd.createHelpText()

	if c.Help != nil {
		c.Help(cmd, text)
		return
	}

	fmt.Print(text)
}
