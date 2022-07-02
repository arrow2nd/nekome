package cli

import "fmt"

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

	// Alias
	alias := ""
	if c.Alias != "" {
		alias = fmt.Sprintf("Alias:\n  %s%s", c.Alias, newLine)
	}

	// Example
	example := ""
	if c.Example != "" {
		example = fmt.Sprintf("Example:\n  %s%s", c.Example, newLine)
	}

	// Commands
	cmds := ""
	if children := c.GetChildren(); len(children) != 0 {
		cmds += "Commands:\n"
		for _, cmd := range children {
			cmds += fmt.Sprintf("  %s\t%s\n", cmd.Name, cmd.Short)
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

	if c.HelpFunc != nil {
		c.HelpFunc(cmd, text)
		return
	}

	fmt.Print(text)
}
