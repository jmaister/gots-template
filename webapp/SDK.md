# TypeScript API Client Generation

This project automatically generates a complete TypeScript API client from the OpenAPI specification using `@hey-api/openapi-ts` with the `@hey-api/client-fetch` plugin.

## How it works

1. The OpenAPI specification is defined in `../api/openapi-spec.yaml`
2. The `@hey-api/openapi-ts` package generates TypeScript types and client functions
3. The `@hey-api/client-fetch` plugin creates a fetch-based HTTP client
4. All generated files are saved to `src/apiclient/` directory
5. The client provides type-safe API functions based on your OpenAPI spec

## Generated Files Structure

```
src/apiclient/
├── index.ts          # Main exports
├── types.gen.ts      # Generated TypeScript types
├── sdk.gen.ts        # Generated API functions
├── client.gen.ts     # Pre-configured client instance
├── client/           # Client implementation
└── core/             # Core utilities
```

## Commands

### Generate TypeScript API client
```bash
npm run generate-api
```

### Build with client generation
```bash
npm run build
```

### From the root project (includes Go codegen)
```bash
make api-codegen        # Generate both Go and TypeScript code
make api-codegen-ts     # Generate only TypeScript API client
make build              # Full build including client generation
```

## Usage

### Recommended Approach (Pre-configured)
```typescript
import { healthCheck, type HealthResponse } from '../apiclient/configured';

// No configuration needed - already set up with correct baseUrl
const healthData = await healthCheck();
console.log(healthData.data); // Fully typed HealthResponse
```

### Manual Configuration
```typescript
import { healthCheck, type HealthResponse } from '../apiclient';
import { client } from '../apiclient/client.gen';

// Configure the client manually
client.setConfig({
    baseUrl: '/_/api'
});

const healthData = await healthCheck();
```

### Using Types
```typescript
import { type HealthResponse } from '../apiclient';

// All types are automatically generated from OpenAPI spec
const processHealth = (health: HealthResponse) => {
    console.log(`Status: ${health.status}`);
    console.log(`Timestamp: ${health.timestamp}`);
};
```

## Generated Files

- `src/apiclient/` - Complete generated API client (excluded from git)
  - `index.ts` - Main exports
  - `types.gen.ts` - All TypeScript types from OpenAPI
  - `sdk.gen.ts` - Generated API functions
  - `client.gen.ts` - Pre-configured client instance
  - `client/` & `core/` - Client implementation

## Adding New API Endpoints

1. Update the OpenAPI specification in `../api/openapi-spec.yaml`
2. Run `make api-codegen` or `npm run generate-api`
3. **All new endpoints are automatically available** with full type safety
4. Import and use the new functions directly from the generated client

Example: If you add a `/users` endpoint to OpenAPI spec:
```typescript
// After regeneration, this will be automatically available:
import { getUsers, createUser } from '../apiclient';

const users = await getUsers();
const newUser = await createUser({ 
    body: { name: 'John', email: 'john@example.com' } 
});
```

## Best Practices

- **Everything is generated** - no manual API code needed
- Regenerate the client whenever the OpenAPI spec changes
- Configure the client once in your app initialization
- Use the generated types for all API-related data structures
- Import functions directly from `../apiclient` for cleaner code
