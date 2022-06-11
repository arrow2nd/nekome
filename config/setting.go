package config

// Settings 設定
type Settings struct {
	// MainUser メインで使用するユーザ
	MainUser string
	// DateFormat 日付のフォーマット文字列
	DateFormat string
	// TimeFormat 時刻のフォーマット文字列
	TimeFormat string
}

func (c *Config) SaveSettings() error {
	return c.save(setingsFileName, c.Settings)
}
