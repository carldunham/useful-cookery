package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/carldunham/useful-cookery/internal/models"
	"github.com/sashabaranov/go-openai"
)

// AIService provides AI capabilities for the application
type AIService struct {
	openAIClient *openai.Client
	config       *Config
	cache        Cache
}

// Config holds configuration for the AI service
type Config struct {
	OpenAIAPIKey     string
	EmbeddingModel   string
	CompletionModel  string
	CacheEnabled     bool
	CacheTTL         time.Duration
	MaxRequestTokens int
}

// Cache interface for storing AI results
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}

// NewAIService creates a new AI service
func NewAIService(config *Config, cache Cache) (*AIService, error) {
	if config.OpenAIAPIKey == "" {
		return nil, errors.New("OpenAI API key is required")
	}

	openAIClient := openai.NewClient(config.OpenAIAPIKey)

	if config.EmbeddingModel == "" {
		config.EmbeddingModel = openai.AdaEmbeddingV2
	}

	if config.CompletionModel == "" {
		config.CompletionModel = openai.GPT3Dot5Turbo
	}

	if config.MaxRequestTokens == 0 {
		config.MaxRequestTokens = 4000
	}

	return &AIService{
		openAIClient: openAIClient,
		config:       config,
		cache:        cache,
	}, nil
}

// GenerateEmbedding creates vector embeddings for text using OpenAI
func (s *AIService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	if s.cache != nil && s.config.CacheEnabled {
		// Check cache first
		cacheKey := fmt.Sprintf("embedding:%s", text)
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil && cached != nil {
			var embedding []float32
			if err := json.Unmarshal(cached, &embedding); err == nil {
				return embedding, nil
			}
		}
	}

	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: s.config.EmbeddingModel,
	}

	resp, err := s.openAIClient.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, errors.New("no embeddings returned from API")
	}

	embedding := resp.Data[0].Embedding

	// Cache the result
	if s.cache != nil && s.config.CacheEnabled {
		cacheKey := fmt.Sprintf("embedding:%s", text)
		if cached, err := json.Marshal(embedding); err == nil {
			if err := s.cache.Set(ctx, cacheKey, cached, s.config.CacheTTL); err != nil {
				log.Printf("Failed to cache embedding: %v", err)
			}
		}
	}

	return embedding, nil
}

// SearchRecipes performs semantic search on recipes
func (s *AIService) SearchRecipes(ctx context.Context, query string, limit int) ([]models.Recipe, error) {
	// Generate embedding for the query
	queryEmbedding, err := s.GenerateEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// In a real implementation, this would send the embedding to DGraph
	// to perform vector search. For now, return a placeholder.
	return []models.Recipe{}, errors.New("not implemented")
}

// RecommendRecipes recommends recipes based on user preferences and ingredients
func (s *AIService) RecommendRecipes(
	ctx context.Context,
	userID string,
	availableIngredients []string,
	limit int,
) ([]models.Recipe, error) {
	// In a real implementation, this would use a combination of
	// collaborative filtering and content-based recommendations
	return []models.Recipe{}, errors.New("not implemented")
}

// GenerateSubstitutes generates ingredient substitutes
func (s *AIService) GenerateSubstitutes(ctx context.Context, ingredient string) ([]string, error) {
	prompt := fmt.Sprintf(
		"Suggest 3 substitutes for %s in cooking recipes. Format as a comma-separated list with no explanations.",
		ingredient,
	)

	// Check cache first
	if s.cache != nil && s.config.CacheEnabled {
		cacheKey := fmt.Sprintf("substitute:%s", ingredient)
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil && cached != nil {
			var substitutes []string
			if err := json.Unmarshal(cached, &substitutes); err == nil {
				return substitutes, nil
			}
		}
	}

	req := openai.ChatCompletionRequest{
		Model: s.config.CompletionModel,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful cooking assistant that suggests ingredient substitutes.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.3,
		MaxTokens:   100,
	}

	resp, err := s.openAIClient.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate substitutes: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("no response from API")
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
		cacheKey := fmt.Sprintf("substitute:%s", ingredient)
		if cached, err := json.Marshal(substitutes); err == nil {
			if err := s.cache.Set(ctx, cacheKey, cached, s.config.CacheTTL); err != nil {
				log.Printf("Failed to cache substitutes: %v", err)
			}
		}
	}

	return substitutes, nil
}

// ProcessNaturalLanguageQuery processes a natural language query and translates it
// into structured search parameters
func (s *AIService) ProcessNaturalLanguageQuery(ctx context.Context, query string) (*models.SearchParams, error) {
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
		cacheKey := fmt.Sprintf("nlquery:%s", query)
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil && cached != nil {
			var params models.SearchParams
			if err := json.Unmarshal(cached, &params); err == nil {
				return &params, nil
			}
		}
	}

	req := openai.ChatCompletionRequest{
		Model: s.config.CompletionModel,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful assistant that analyzes recipe search queries and extracts structured parameters.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.2,
		MaxTokens:   500,
	}

	resp, err := s.openAIClient.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to process query: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("no response from API")
	}

	// Extract JSON from response
	content := resp.Choices[0].Message.Content
	jsonStr := extractJSON(content)
	if jsonStr == "" {
		return nil, errors.New("failed to extract JSON from response")
	}

	// Parse JSON
	var params models.SearchParams
	if err := json.Unmarshal([]byte(jsonStr), &params); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	// Cache the result
	if s.cache != nil && s.config.CacheEnabled {
		cacheKey := fmt.Sprintf("nlquery:%s", query)
		if cached, err := json.Marshal(params); err == nil {
			if err := s.cache.Set(ctx, cacheKey, cached, s.config.CacheTTL); err != nil {
				log.Printf("Failed to cache query params: %v", err)
			}
		}
	}

	return &params, nil
}

// extractJSON extracts JSON from a string
func extractJSON(input string) string {
	startIdx := strings.Index(input, "{")
	endIdx := strings.LastIndex(input, "}")
	
	if startIdx == -1 || endIdx == -1 || endIdx <= startIdx {
		return ""
	}
	
	return input[startIdx : endIdx+1]
}