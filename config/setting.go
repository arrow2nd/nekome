package config

// Settings 設定
type Settings struct {
	MainUser string
}

func (c *Config) SaveSettings() error {
	return c.save(setingsFileName, c.Settings)
}
