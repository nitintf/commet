package config

import (
	"fmt"
	"strings"
)

type Provider string

const (
	ProviderOpenAI Provider = "openai"
	ProviderClaude Provider = "claude"
	ProviderGoogle Provider = "google"
)

type AIConfig struct {
	Provider Provider `mapstructure:"provider" yaml:"provider"`
	APIKey   string   `mapstructure:"api_key" yaml:"api_key"`
	Model    string   `mapstructure:"model" yaml:"model"`
}

func (p Provider) String() string {
	return string(p)
}

func (p Provider) IsValid() bool {
	switch p {
	case ProviderOpenAI, ProviderClaude, ProviderGoogle:
		return true
	default:
		return false
	}
}

func ParseProvider(s string) (Provider, error) {
	provider := Provider(strings.ToLower(s))
	if !provider.IsValid() {
		return "", fmt.Errorf("invalid provider: %s (valid options: openai, claude, google)", s)
	}
	return provider, nil
}

func (c *AIConfig) GetDefaultModel() string {
	switch c.Provider {
	case ProviderOpenAI:
		return "gpt-4"
	case ProviderClaude:
		return "claude-3-sonnet-20240229"
	case ProviderGoogle:
		return "gemini-pro"
	default:
		return ""
	}
}

func (c *AIConfig) SetDefaults() {
	if c.Provider == "" {
		c.Provider = ProviderOpenAI
	}
}

func (c *AIConfig) MaskAPIKey() string {
	if c.APIKey == "" {
		return ""
	}
	if len(c.APIKey) <= 8 {
		return strings.Repeat("*", len(c.APIKey))
	}
	return c.APIKey[:4] + strings.Repeat("*", len(c.APIKey)-8) + c.APIKey[len(c.APIKey)-4:]
}