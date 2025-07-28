package llm

import (
	"context"
	"fmt"

	"github.com/bitcs/commet/internal/config"
	"github.com/bitcs/commet/internal/prompts"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/openai"
)

type Service struct {
	llm    llms.Model
	config *config.Config
}

func NewService(cfg *config.Config) (*Service, error) {
	llmModel, err := createLLM(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM: %w", err)
	}

	return &Service{
		llm:    llmModel,
		config: cfg,
	}, nil
}

func createLLM(cfg *config.Config) (llms.Model, error) {
	model := cfg.AI.Model
	if model == "" {
		model = cfg.AI.GetDefaultModel()
	}

	switch cfg.AI.Provider {
	case config.ProviderOpenAI:
		return openai.New(
			openai.WithToken(cfg.AI.APIKey),
			openai.WithModel(model),
		)
	case config.ProviderClaude:
		return anthropic.New(
			anthropic.WithToken(cfg.AI.APIKey),
			anthropic.WithModel(model),
		)
	case config.ProviderGoogle:
		return googleai.New(context.Background(),
			googleai.WithAPIKey(cfg.AI.APIKey),
			googleai.WithDefaultModel(model),
		)
	case config.ProviderGroq:
		return openai.New(
			openai.WithToken(cfg.AI.APIKey),
			openai.WithModel(model),
			openai.WithBaseURL("https://api.groq.com/openai/v1"),
		)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.AI.Provider)
	}
}

func (s *Service) GenerateCommitMessage(ctx context.Context, gitDiff string) (string, error) {
	prompt := prompts.CommitMessagePrompt(gitDiff)

	response, err := s.llm.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate commit message: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	return response.Choices[0].Content, nil
}


func (s *Service) TestConnection(ctx context.Context) error {
	testPrompt := "Respond with 'OK' if you can understand this message."

	response, err := s.llm.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, testPrompt),
	})
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}

	if len(response.Choices) == 0 {
		return fmt.Errorf("no response from LLM")
	}

	return nil
}
