# Cline Rules for Useful Cookery Project

## Project Overview

Useful Cookery is a recipe database application providing a modern, AI-enhanced cooking experience. It stores recipes originally in TROFF format and converts them to structured JSON for use in a web application. This project is a complete rewrite with modern technology.

## Technology Stack

- **Backend**: Go
- **Database**: Multiple options supported (DGraph, PostgreSQL, in-memory)
- **API**: GraphQL
- **Frontend**: React
- **AI Integration**: External APIs initially, with plans for RAG and fine-tuning

## Project Structure

```text
useful-cookery/
├── cmd/                    # Application entry points
│   ├── api/                # GraphQL API server
│   ├── migration/          # TROFF to database migration tool
│   └── cli/                # Command-line utilities
├── internal/               # Private application code
│   ├── auth/               # Authentication services
│   ├── models/             # Data models
│   ├── parser/             # TROFF parser
│   ├── database/           # Database abstraction layer
│   ├── ai/                 # AI service integrations
│   └── graphql/            # GraphQL resolvers
├── ui/                     # Frontend React application
├── schema/                 # GraphQL schema definitions
├── scripts/                # Build and deployment scripts
└── docs/                   # Documentation
```

## Key Features

- Recipe database with comprehensive search capabilities
- Natural language query support for finding recipes
- AI-powered recommendations based on ingredient availability and user preferences
- User accounts with personalization
- Multi-category organization for recipes
- Mobile-friendly design

## Development Workflow

1. Clone the repository
2. Configure database (options include DGraph, PostgreSQL, or in-memory)
   - For DGraph setup, see docs/dgraph-setup.md
3. Configure environment (cp .env.example .env)
4. Run the API server: `go run cmd/api/main.go`
5. Run the UI: `cd ui && npm start`

## Important Notes

- The data directory contains raw recipe files in TROFF format
- The application converts TROFF recipes to structured JSON for storage in the database
- GraphQL is used for the API layer
- The project is currently in Phase 1 (Foundation) of development

## Development Phases

1. **Foundation** (Current): Basic architecture, database schema, GraphQL API, authentication, TROFF parser
2. **Core Features**: Data migration, user accounts, recipe editing, enhanced search
3. **AI Integration**: Natural language search, recommendations, ingredient substitution
4. **Mobile and Advanced Features**: Mobile app, meal planning, image recognition
5. **Optimization and Expansion**: Fine-tuning AI, internationalization, third-party integrations

## Technical Requirements

### Backend (Go)

- Go 1.20+ with modules
- gqlgen for GraphQL
- JWT authentication
- Testing with testify/assert
- Logging with slog
- Configuration with viper
- Command-line support with cobra
- Follow [Google coding standards](https://google.github.io/styleguide/go/guide)
- Use golangci-lint with the project's configuration to catch common issues
- Implement unit tests for each file, and keep them up to date
- Prefer using table-driven tests to reduce test function complexity
- Put unit tests into _test packages, and only test exported types and functions
- Avoid stuttering in type names (e.g., avoid `package.PackageThing`, prefer `package.Thing`)
- Return concrete types from functions rather than interfaces
  - When a function must return an interface (e.g., for factory functions or third-party libraries), use `//nolint:ireturn` with an explanatory comment
  - Ensure the linter name in nolint directives matches exactly what's configured in `.golangci.yml`
- Always wrap errors from external packages with additional context using `fmt.Errorf("context: %w", err)`
  - This includes errors from standard library functions like `strconv.Atoi()` and `encoding/base64.DecodeString()`
  - Never return unwrapped errors from external packages directly
- For creating new static errors, use `errors.New()` instead of `fmt.Errorf()`
  - Define package-level error variables for common error cases (e.g., `var ErrInvalidFormat = errors.New("invalid format")`)
  - Prefer these static errors over creating dynamic errors with the same message repeatedly
- When comparing errors, use `errors.Is(err, targetErr)` instead of direct comparison (`err == targetErr` or `err != targetErr`)
- In tests, use `t.Context()` instead of `context.Background()` for better test context handling
- Add `t.Parallel()` to test functions when they can safely run in parallel
- Ensure all comments end with a period for consistency
- Keep function cyclomatic complexity below 10 to maintain readability
  - For complex functions that cannot be easily refactored (e.g., GraphQL resolvers with multiple fallback strategies), use `//nolint:cyclop` with an explanatory comment
- Keep function length below 60 lines to maintain readability
  - For functions that necessarily need to be longer (e.g., database operations with many fields), use `//nolint:funlen` with an explanatory comment
  - Consider if the function can be split into smaller helper functions before adding a nolint directive
- Avoid variable names that are too short for their scope, with exceptions for standard Go idioms like `db`, `tx`, `id`, `ok`, and `err`
- Use interfaces for defining behavior, not for returning values
- Prefer dependency injection through interfaces, but have factory functions return concrete implementations
- Keep interfaces focused and small (under 10 methods) when possible; if a larger interface is necessary, document the reason
- Extract complex nested logic into separate helper functions to improve readability and reduce nesting depth
- When refactoring is not feasible, use `//nolint` directives sparingly and always with explanatory comments
  - Always specify the exact linter being disabled (e.g., `//nolint:cyclop` instead of just `//nolint`)
  - Ensure the linter name matches exactly what's configured in `.golangci.yml`
  - Add a comment explaining why the linter is being disabled for that specific case
- For test functions, it's acceptable to have higher complexity and length to ensure comprehensive test coverage
- When implementing database operations, prefer multiple smaller functions over fewer large ones
- Use constants for string literals that are used in multiple places, especially for filter keys, error messages, and other identifiers
- When handling errors that don't affect the main return value (e.g., errors from counting operations when you already have results to return), log the error and continue rather than returning nil
- Document any intentional deviations from linting rules in the code with clear explanations

### Frontend (React)

- React 18+ with hooks
- TypeScript
- Apollo Client for GraphQL
- Tailwind CSS for styling
- React Router for navigation
- React Testing Library
- Follow accepted coding standards for React and Typescript

### DevOps

- Docker for containerization
- GitHub Actions for CI/CD
- Terraform for infrastructure
- Prometheus and Grafana for monitoring

### AI Integration

- OpenAI API for initial AI features
- Vector embeddings for semantic search
- RAG system for contextual recommendations

## Success Metrics

- API response times under 100ms for non-AI endpoints
- Search results returned in under 500ms
- 95% test coverage for critical components
- Web Core Vitals meeting "Good" thresholds
