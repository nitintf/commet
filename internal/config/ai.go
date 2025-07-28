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
	ProviderGroq   Provider = "groq"
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
	case ProviderOpenAI, ProviderClaude, ProviderGoogle, ProviderGroq:
		return true
	default:
		return false
	}
}

func ParseProvider(s string) (Provider, error) {
	provider := Provider(strings.ToLower(s))
	if !provider.IsValid() {
		return "", fmt.Errorf("invalid provider: %s (valid options: openai, claude, google, groq)", s)
	}
	return provider, nil
}

// GetAvailableModels returns the available models for each provider
func GetAvailableModels(provider Provider) []string {
	switch provider {
	case ProviderOpenAI:
		return []string{
			"gpt-4o",
			"gpt-4-turbo",
			"gpt-4",
			"gpt-3.5-turbo",
		}
	case ProviderClaude:
		return []string{
			"claude-3-5-sonnet-20241022",
			"claude-3-opus-20240229",
			"claude-3-sonnet-20240229",
			"claude-3-haiku-20240307",
		}
	case ProviderGoogle:
		return []string{
			"gemini-1.5-pro",
			"gemini-1.5-flash",
			"gemini-pro",
			"gemini-pro-vision",
		}
	case ProviderGroq:
		return []string{
			"llama-3.1-70b-versatile",
			"llama-3.1-8b-instant",
			"mixtral-8x7b-32768",
			"gemma-7b-it",
		}
	default:
		return []string{}
	}
}

func (c *AIConfig) GetDefaultModel() string {
	models := GetAvailableModels(c.Provider)
	if len(models) > 0 {
		return models[0]
	}
	return ""
}

func (c *AIConfig) GetAvailableModels() []string {
	return GetAvailableModels(c.Provider)
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
