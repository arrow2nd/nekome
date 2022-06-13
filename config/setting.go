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

func defaultSettings() *Settings {
	return &Settings{
		MainUser:   "",
		DateFormat: "2006/01/02",
		TimeFormat: "15:04:05",
	}
}

func (c *Config) SaveSettings() error {
	return c.save(setingsFileName, c.Settings)
}
