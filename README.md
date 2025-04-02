# Useful Cookery

## Overview

Useful Cookery is a recipe database application that aims to provide users with a modern, AI-enhanced cooking experience. The application stores recipes originally in TROFF format and converts them to a structured JSON format for use in a web application.

This branch (`cd/5/rewrite`) represents a complete rewrite of the original application with modern technology and enhanced functionality.

## Key Features

- Recipe database with comprehensive search capabilities
- Natural language query support for finding recipes
- AI-powered recommendations based on ingredient availability and user preferences
- User accounts with personalization
- Multi-category organization for recipes
- Mobile-friendly design

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

## Getting Started

1. Clone the repository
2. Set up DGraph (see docs/dgraph-setup.md)
3. Configure environment (cp .env.example .env)
4. Run the API server: `go run cmd/api/main.go`
5. Run the UI: `cd ui && npm start`

## Development Roadmap

See [ROADMAP.md](ROADMAP.md) for the project development plan.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on contributing to the project.
