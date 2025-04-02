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
- Follow coding standards at https://google.github.io/styleguide/go/guide

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
