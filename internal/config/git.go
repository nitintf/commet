package config

type GitConfig struct {
	AutoStage   bool `mapstructure:"auto_stage" yaml:"auto_stage"`
	ShowDiff    bool `mapstructure:"show_diff" yaml:"show_diff"`
	ConfirmPush bool `mapstructure:"confirm_push" yaml:"confirm_push"`
}

func (c *GitConfig) SetDefaults() {
	// Only set defaults for uninitialized values
	// For booleans, we can't distinguish between false and unset,
	// so we'll let the YAML file or explicit user settings take precedence
	// The defaults are now just documentation of the recommended settings
}
