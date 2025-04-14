package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"github.com/carldunham/useful-cookery/internal/model"
)

// Error definitions.
var (
	ErrAPIKeyRequired       = errors.New("OpenAI API key is required")
	ErrNoEmbeddingsReturned = errors.New("no embeddings returned from API")
	ErrNotImplemented       = errors.New("not implemented")
	ErrNoAPIResponse        = errors.New("no response from API")
	ErrFailedToExtractJSON  = errors.New("failed to extract JSON from response")
)

// Constants for OpenAI API parameters.
const (
	SubstituteTemperature = 0.3
	SubstituteMaxTokens   = 100
	QueryTemperature      = 0.2
	QueryMaxTokens        = 500
)

// Service provides AI capabilities for the application.
type Service struct {
	openAIClient openai.Client
	config       *Config
	cache        Cache
}

// Config holds configuration for the AI service.
type Config struct {
	OpenAIAPIKey     string
	EmbeddingModel   openai.EmbeddingModel
	CompletionModel  string
	CacheEnabled     bool
	CacheTTL         time.Duration
	MaxRequestTokens int
}

// Cache interface for storing AI results.
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}

// NewAIService creates a new AI service.
func NewAIService(config *Config, cache Cache) (*Service, error) {
	if config.OpenAIAPIKey == "" {
		return nil, ErrAPIKeyRequired
	}

	openAIClient := openai.NewClient(option.WithAPIKey(config.OpenAIAPIKey))

	if config.EmbeddingModel == "" {
		config.EmbeddingModel = openai.EmbeddingModelTextEmbeddingAda002
	}

	if config.CompletionModel == "" {
		config.CompletionModel = openai.ChatModelGPT3_5Turbo
	}

	if config.MaxRequestTokens == 0 {
		config.MaxRequestTokens = 4000
	}

	return &Service{
		openAIClient: openAIClient,
		config:       config,
		cache:        cache,
	}, nil
}

// GenerateEmbedding creates vector embeddings for text using OpenAI.
//
//nolint:cyclop // TODO: simplify.
func (s *Service) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	if s.cache != nil && s.config.CacheEnabled {
		// Check cache first
		cacheKey := "embedding:" + text
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil && cached != nil {
			var embedding []float32
			if err := json.Unmarshal(cached, &embedding); err == nil {
				return embedding, nil
			}
		}
	}

	// Create input union with a string
	inputUnion := openai.EmbeddingNewParamsInputUnion{}
	inputUnion.OfString = openai.String(text)

	// Create embedding params
	params := openai.EmbeddingNewParams{
		Model: s.config.EmbeddingModel,
		Input: inputUnion,
	}

	// Create embedding
	resp, err := s.openAIClient.Embeddings.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, ErrNoEmbeddingsReturned
	}

	// Convert embedding from []float64 to []float32
	embedding64 := resp.Data[0].Embedding
	embedding32 := make([]float32, len(embedding64))
	for i, v := range embedding64 {
		embedding32[i] = float32(v)
	}

	// Cache the result
	if s.cache != nil && s.config.CacheEnabled {
		cacheKey := "embedding:" + text
		if cached, err := json.Marshal(embedding32); err == nil {
			if err := s.cache.Set(ctx, cacheKey, cached, s.config.CacheTTL); err != nil {
				log.Printf("Failed to cache embedding: %v", err)
			}
		}
	}

	return embedding32, nil
}

// SearchRecipes performs semantic search on recipes.
func (s *Service) SearchRecipes(_ context.Context, _ string, _ int) ([]model.Recipe, error) {
	// Generate embedding for the query
	// queryEmbedding, err := s.GenerateEmbedding(ctx, query)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	// }

	// In a real implementation, this would send the embedding to DGraph
	// to perform vector search. For now, return a placeholder.
	return []model.Recipe{}, ErrNotImplemented
}

// RecommendRecipes recommends recipes based on user preferences and ingredients.
func (s *Service) RecommendRecipes(
	_ context.Context,
	_ string,
	_ []string,
	_ int,
) ([]model.Recipe, error) {
	// In a real implementation, this would use a combination of
	// collaborative filtering and content-based recommendations
	return []model.Recipe{}, ErrNotImplemented
}

// GenerateSubstitutes generates ingredient substitutes.
//
//nolint:cyclop // TODO: simplify.
func (s *Service) GenerateSubstitutes(ctx context.Context, ingredient string) ([]string, error) {
	prompt := fmt.Sprintf(
		"Suggest 3 substitutes for %s in cooking recipes. Format as a comma-separated list with no explanations.",
		ingredient,
	)

	// Check cache first
	if s.cache != nil && s.config.CacheEnabled {
		cacheKey := "substitute:" + ingredient
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil && cached != nil {
			var substitutes []string
			if err := json.Unmarshal(cached, &substitutes); err == nil {
				return substitutes, nil
			}
		}
	}

	// Create messages
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are a helpful cooking assistant that suggests ingredient substitutes."),
		openai.UserMessage(prompt),
	}

	// Create chat completion params
	params := openai.ChatCompletionNewParams{
		Model:       s.config.CompletionModel,
		Messages:    messages,
		Temperature: openai.Float(SubstituteTemperature),
		MaxTokens:   openai.Int(SubstituteMaxTokens),
	}

	// Create chat completion
	resp, err := s.openAIClient.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate substitutes: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, ErrNoAPIResponse
	}

	// Parse response
	rawSubstitutes := resp.Choices[0].Message.Content
	substitutes := []string{}
	for _, sub := range strings.Split(rawSubstitutes, ",") {
		sub = strings.TrimSpace(sub)
		if sub != "" {
			substitutes = append(substitutes, sub)
		}
	}

	// Cache the result
	if s.cache != nil && s.config.CacheEnabled {
		cacheKey := "substitute:" + ingredient
		if cached, err := json.Marshal(substitutes); err == nil {
			if err := s.cache.Set(ctx, cacheKey, cached, s.config.CacheTTL); err != nil {
				log.Printf("Failed to cache substitutes: %v", err)
			}
		}
	}

	return substitutes, nil
}

// ProcessNaturalLanguageQuery processes a natural language query and translates it
// into structured search parameters.
//
//nolint:cyclop,funlen // TODO: simplify.
func (s *Service) ProcessNaturalLanguageQuery(ctx context.Context, query string) (*model.SearchParams, error) {
	prompt := fmt.Sprintf(`
Analyze this recipe search query: "%s"

Extract the following parameters in JSON format:
{
  "ingredients": [], // List of ingredients mentioned
  "excludedIngredients": [], // List of ingredients to avoid
  "categories": [], // Recipe categories (dessert, main dish, etc.)
  "cuisine": "", // Specific cuisine type
  "dietaryRestrictions": [], // Like vegetarian, vegan, gluten-free, etc.
  "maxPrepTime": 0, // Maximum preparation time in minutes (0 if not specified)
  "difficulty": "" // BEGINNER, INTERMEDIATE, ADVANCED (empty if not specified)
}
`, query)

	// Check cache
	if s.cache != nil && s.config.CacheEnabled {
		cacheKey := "nlquery:" + query
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil && cached != nil {
			var params model.SearchParams
			if err := json.Unmarshal(cached, &params); err == nil {
				return &params, nil
			}
		}
	}

	// Create messages
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are a helpful assistant that analyzes recipe search queries and extracts parameters."),
		openai.UserMessage(prompt),
	}

	// Create chat completion params
	params := openai.ChatCompletionNewParams{
		Model:       s.config.CompletionModel,
		Messages:    messages,
		Temperature: openai.Float(QueryTemperature),
		MaxTokens:   openai.Int(QueryMaxTokens),
	}

	// Create chat completion
	resp, err := s.openAIClient.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to process query: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, ErrNoAPIResponse
	}

	// Extract JSON from response
	content := resp.Choices[0].Message.Content
	jsonStr := extractJSON(content)
	if jsonStr == "" {
		return nil, ErrFailedToExtractJSON
	}

	// Parse JSON
	var searchParams model.SearchParams
	if err := json.Unmarshal([]byte(jsonStr), &searchParams); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	// Cache the result
	if s.cache != nil && s.config.CacheEnabled {
		cacheKey := "nlquery:" + query
		if cached, err := json.Marshal(searchParams); err == nil {
			if err := s.cache.Set(ctx, cacheKey, cached, s.config.CacheTTL); err != nil {
				log.Printf("Failed to cache query params: %v", err)
			}
		}
	}

	return &searchParams, nil
}

// extractJSON extracts JSON from a string.
func extractJSON(input string) string {
	startIdx := strings.Index(input, "{")
	endIdx := strings.LastIndex(input, "}")

	if startIdx == -1 || endIdx == -1 || endIdx <= startIdx {
		return ""
	}

	return input[startIdx : endIdx+1]
}
