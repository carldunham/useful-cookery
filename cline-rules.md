# Cline Rules for Useful Cookery Project

## Project Overview

Useful Cookery is a recipe database application providing a modern, AI-enhanced cooking experience. It stores recipes originally in TROFF format and converts them to structured JSON for use in a web application. This project is a complete rewrite with modern technology (branch `cd/5/rewrite`).

## Technology Stack

- **Backend**: Go
- **Database**: DGraph
- **API**: GraphQL
- **Frontend**: React
- **AI Integration**: External APIs initially, with plans for RAG and fine-tuning

## Project Structure

```text
useful-cookery/
├── cmd/                    # Application entry points
│   ├── api/                # GraphQL API server
│   ├── migration/          # TROFF to DGraph migration tool
│   └── cli/                # Command-line utilities
├── internal/               # Private application code
│   ├── auth/               # Authentication services
│   ├── models/             # Data models
│   ├── parser/             # TROFF parser
│   ├── database/           # DGraph interface
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
2. Set up DGraph (see docs/dgraph-setup.md)
3. Configure environment (cp .env.example .env)
4. Run the API server: `go run cmd/api/main.go`
5. Run the UI: `cd ui && npm start`

## Important Notes

- The data directory contains raw recipe files in TROFF format
- The application converts TROFF recipes to structured JSON
- GraphQL is used for the API layer
- The project is currently in Phase 1 (Foundation) of development

## Development Phases

1. **Foundation** (Current): Basic architecture, DGraph schema, GraphQL API, authentication, TROFF parser
2. **Core Features**: Data migration, user accounts, recipe editing, enhanced search
3. **AI Integration**: Natural language search, recommendations, ingredient substitution
4. **Mobile and Advanced Features**: Mobile app, meal planning, image recognition
5. **Optimization and Expansion**: Fine-tuning AI, internationalization, third-party integrations

## Technical Requirements

### Backend (Go)

- Go 1.20+ with modules
- gqlgen for GraphQL
- DGraph Go client
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
- Always wrap errors from external packages with additional context using `fmt.Errorf("context: %w", err)`
- For creating new static errors, use `errors.New()` instead of `fmt.Errorf()`
- When comparing errors, use `errors.Is(err, targetErr)` instead of direct comparison (`err == targetErr` or `err != targetErr`)
- In tests, use `t.Context()` instead of `context.Background()` for better test context handling
- Add `t.Parallel()` to test functions when they can safely run in parallel
- Ensure all comments end with a period for consistency
- Keep function cyclomatic complexity below 10 to maintain readability
- Keep function length below 60 lines to maintain readability
- Avoid variable names that are too short for their scope, with exceptions for standard Go idioms like `db`, `tx`, `id`, `ok`, and `err`
- Use interfaces for defining behavior, not for returning values
- Prefer dependency injection through interfaces, but have factory functions return concrete implementations
- Keep interfaces focused and small (under 10 methods) when possible; if a larger interface is necessary, document the reason
- Extract complex nested logic into separate helper functions to improve readability and reduce nesting depth
- When refactoring is not feasible, use `//nolint` directives sparingly and always with explanatory comments
- For test functions, it's acceptable to have higher complexity and length to ensure comprehensive test coverage
- When implementing database operations, prefer multiple smaller functions over fewer large ones
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
