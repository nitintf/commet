package config

type GitConfig struct {
	AutoStage     bool `mapstructure:"auto_stage" yaml:"auto_stage"`
	ShowDiff      bool `mapstructure:"show_diff" yaml:"show_diff"`
	ConfirmPush   bool `mapstructure:"confirm_push" yaml:"confirm_push"`
	DirectCommit  bool `mapstructure:"direct_commit" yaml:"direct_commit"`
	UseAI         bool `mapstructure:"use_ai" yaml:"use_ai"`
	Interactive   bool `mapstructure:"interactive" yaml:"interactive"`
}

func (c *GitConfig) SetDefaults() {
	// Since we can't distinguish between false and unset for boolean fields in Go,
	// we rely on the config file or explicit user settings for the actual values.
	// The application should handle defaults at runtime when needed.
}
