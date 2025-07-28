package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	AI  AIConfig  `mapstructure:"ai" yaml:"ai"`
	Git GitConfig `mapstructure:"git" yaml:"git"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Set UseAI default to true if not explicitly set in config
	if !viper.IsSet("git.use_ai") {
		cfg.Git.UseAI = true
	}

	cfg.SetDefaults()

	return &cfg, nil
}

func (c *Config) Save() error {
	viper.Set("ai.provider", c.AI.Provider)
	viper.Set("ai.api_key", c.AI.APIKey)
	viper.Set("ai.model", c.AI.Model)
	viper.Set("git.auto_stage", c.Git.AutoStage)
	viper.Set("git.show_diff", c.Git.ShowDiff)
	viper.Set("git.confirm_push", c.Git.ConfirmPush)
	viper.Set("git.direct_commit", c.Git.DirectCommit)
	viper.Set("git.use_ai", c.Git.UseAI)
	viper.Set("git.interactive", c.Git.Interactive)

	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("could not get home directory: %w", err)
		}
		configPath = fmt.Sprintf("%s/.commet.yaml", home)
	}

	return viper.WriteConfigAs(configPath)
}

func (c *Config) MaskAPIKey() string {
	return c.AI.MaskAPIKey()
}

func (c *Config) GetDefaultModel() string {
	return c.AI.GetDefaultModel()
}

func (c *Config) SetDefaults() {
	c.AI.SetDefaults()
	c.Git.SetDefaults()
}
