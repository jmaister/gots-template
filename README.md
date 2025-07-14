# GOTS Template

A production-ready full-stack application template with Go backend, React frontend, and OpenAPI-generated code.

## Tools

- **Go** - Backend API server with net/http and GORM
- **TypeScript** - Frontend with React and React Router
- **Vite** - Fast frontend build tool
- **OpenAPI** - API specification and code generation
- **Tailwind 4** - Modern CSS framework
- **SQLite** - Embedded database

## How to Use

1. Click "Use this template" on GitHub
2. Clone your new repository
3. Run `./setup-template.sh` to customize for your project

## How to Run

```bash
# Generate OpenAPI code
make api-codegen

# Build everything
make build

# Run the application
make run
```

