# Gomania - Arabic Podcast Management System

A simple and clean podcast content management system built with Go, focusing on essential fields and Arabic content support.

## 🚀 Features

- **Simple CMS**: Clean content management for programs with essential fields only
- **Category Management**: Organize programs by categories
- **Arabic Content**: Full Arabic language support with UTF-8 encoding
- **Smart Discovery**: Unified search API that intelligently searches local content first, then falls back to external sources when no local results are found
- **External Source Integration**:
  - iTunes API integration for podcast discovery
  - Pluggable architecture for adding new sources (Spotify, Google Podcasts, etc.)
  - Automatic fallback when local search yields no results
- **Performance Optimizations**:
  - In-memory caching system for external API responses
  - Connection pooling for database operations
  - Efficient search algorithms
- **Developer Experience**:
  - Clean layered architecture with clear separation of concerns
  - Type-safe database queries with SQLC generation
  - Comprehensive API testing suite
  - RESTful API design with consistent patterns
- **Production Ready**:
  - Structured logging with slog
  - Health check endpoints
  - CORS support
  - Environment-based configuration

## 📋 Requirements

- Go 1.24+
- dbmate (for database migrations)
- Docker & Docker Compose
- jq (for API testing script)

## 🛠️ Installation & Setup

### 1. Clone Repository
```bash
git clone https://github.com/khatibomar/gomania
cd gomania
```

### 2. Initialize Database
```bash
make setup
```

This runs: `docker-up` → `db-up` → `db-seed` automatically.

Or manually:
```bash
make docker-up
make gen
make db-up
make db-seed
```

### 3. Run Server
```bash
make build api
```

Server will start on `http://localhost:4000`

## 4. Database UI

I am using [pgweb](https://sosedoff.github.io/pgweb/)

[http://localhost:8081](http://localhost:8081)

## 📊 Database Schema

The system uses a simplified schema with only essential fields:

### Tables
- **programs**: Core podcast programs with essential fields
- **categories**: Simple category organization
- **users**: Basic user authentication for CMS

### Essential Fields (Programs)
- **title**: Program title
- **description**: Program description
- **category_id**: Foreign key to category
- **language**: Content language (default: Arabic)
- **duration**: Program duration in seconds
- **created_at**: Timestamp when record was created
- **updated_at**: Timestamp when record was last updated

### Smart Discovery System

The discovery API (`/v1/programs`) provides intelligent search capabilities:

1. **Local-First Search**: Searches your local podcast database first
2. **Automatic Fallback**: If no local results found, automatically searches external sources
3. **Unified Response**: Returns results in a consistent format regardless of source
4. **Performance Optimized**: Caches external results to reduce API calls

#### External Sources Integration
- **iTunes API**: Search and discover podcasts from iTunes Store
- **Smart Fallback**: Seamless transition to external search when local database has no matches
- **Source Management**: Extensible architecture for adding new podcast sources
- **Response Aggregation**: Combines results from multiple sources with metadata about each source
- **Caching Layer**: Intelligent caching to improve performance and reduce external API calls

## 🌐 API Documentation

[API.MD](API.md)

## 🏗️ Architecture

### Project Structure
```
gomania/
├── cmd/
│   ├── api/             # API server
│   │   ├── main.go      # Entry point
│   │   ├── routes.go    # Route definitions
│   │   ├── cms.go       # CMS handlers
│   │   ├── discovery.go # Discovery & search handlers
│   │   ├── errors.go    # Error handlers
│   │   └── ...
│   └── tools/
│       └── seed/        # Database seeding tool
├── internal/
│   ├── cache/           # Caching system
│   ├── database/        # SQLC generated code
│   ├── service/         # Business logic
│   │   └── program.go   # Program & category service
│   └── sources/         # External source integrations
│       ├── manager.go   # Source manager
│       ├── client.go    # Source client interface
│       └── itunes/      # iTunes API client
├── data/sql/
│   ├── migrations/      # Database migrations
│   └── queries/         # SQL queries
├── scripts/
│   └── test_api.sh      # API testing script
└── docker-compose.yaml  # Database setup
```

### Database Schema

#### Core Tables
- `programs` - Podcast programs with essential fields
- `categories` - Simple categories
- `users` - Basic CMS authentication

#### Relationships
- Programs → Categories (many:1)

## 📝 Configuration

### Environment Variables
- `GOMANIA_CONNECTION_STRING`: PostgreSQL connection string
- `PORT`: Server port (default: 4000)
- `ENV`: Environment (development/staging/production)

### Command Line Flags
```bash
go run cmd/api/*.go \
  -port=8080 \
  -env=production \
  -cors-trusted-origins="https://mydomain.com"
```

## 🧪 Testing

### Make Commands
```bash
# Setup - Complete initialization (recommended)
make setup              # docker-up + db-up + db-seed

# Build and run
make build              # Build API binary
make api                # Build and run API server
make debug              # Build and run with debugger

# Database
make db-up              # Run migrations
make db-drop            # Drop database
make db-seed            # Load sample data

# Docker
make docker-up          # Start database container
make docker-down        # Stop containers

# Code generation
make gen                # Generate SQLC code

# Testing
make test-api           # Run API endpoint tests
```

### Automated Testing
```bash
# Run comprehensive API tests
make test-api
```

### Manual Testing
```bash
# Health check
curl http://localhost:4000/v1/healthcheck

# Discovery API - List all programs
curl http://localhost:4000/v1/programs

# Discovery API - Smart search (searches local first, then external if no results)
curl "http://localhost:4000/v1/programs?q=تقنية"

# Discovery API - Search that triggers external fallback
curl "http://localhost:4000/v1/programs?q=nonexistentterm"

# External sources - List available sources
curl http://localhost:4000/v1/external/sources

# External sources - Direct iTunes search
curl "http://localhost:4000/v1/external/search?source=itunes&q=technology&limit=5"

# CMS - List categories
curl http://localhost:4000/v1/cms/categories

# CMS - Create category
curl -X POST http://localhost:4000/v1/cms/categories \
  -H "Content-Type: application/json" \
  -d '{"name":"تقنية"}'

# CMS - List programs
curl http://localhost:4000/v1/cms/programs

# CMS - Create program (get category_id from above first)
curl -X POST http://localhost:4000/v1/cms/programs \
  -H "Content-Type: application/json" \
  -d '{
    "title":"برنامج جديد",
    "description":"وصف البرنامج",
    "category_id":"550e8400-e29b-41d4-a716-446655440001",
    "language":"ar",
    "duration":1800
  }'

# CMS - Get programs by category
curl http://localhost:4000/v1/cms/categories/{category_id}/programs

# Debug information
curl http://localhost:4000/debug/vars
```

## 📊 Sample Data

The system includes sample Arabic categories and programs:

### Categories
- تقنية (Technology)
- تعليم (Education)
- تسلية (Entertainment)
- أخبار (News)
- رياضة (Sports)
- صحة (Health)
- تاريخ (History)
- فنون (Arts)

### Programs
- Arabic tech podcasts
- Educational content
- Entertainment shows
- News programs

### External Sources
- **iTunes**: Access to iTunes podcast directory
- **Extensible**: Architecture supports adding more sources (Spotify, Google Podcasts, etc.)

Load sample data with:
```bash
make db-seed
```

## 🚀 Deployment

### Docker Deployment
```bash
# Build image
docker build -t gomania-api .

# Run with database
docker compose up -d
```

### Production Considerations
- Use connection pooling for database
- Set up proper logging aggregation
- Configure CORS for frontend domains
- Use environment variables for secrets
- Set up health checks and monitoring

## 🔗 API Client Examples

### JavaScript/Node.js
```javascript
// Discovery API - Search programs (local + external fallback)
const searchPrograms = async (query) => {
  const response = await fetch(`http://localhost:4000/v1/programs?q=${encodeURIComponent(query)}`);
  const data = await response.json();

  if (data.search) {
    // Search response with local and potentially external results
    return {
      results: data.search.results,
      query: data.search.query,
      count: data.search.count,
      sources: data.search.sources
    };
  } else {
    // Simple list response
    return {
      results: data.programs,
      count: data.programs ? data.programs.length : 0
    };
  }
};

// External sources - Search specific source
const searchExternalSource = async (source, query, limit = 10) => {
  const response = await fetch(
    `http://localhost:4000/v1/external/search?source=${source}&q=${encodeURIComponent(query)}&limit=${limit}`
  );
  const data = await response.json();
  return data.external_search.results;
};

// External sources - List available sources
const getAvailableSources = async () => {
  const response = await fetch('http://localhost:4000/v1/external/sources');
  const data = await response.json();
  return data.external_sources.sources;
};

// Usage examples
const searchResult = await searchPrograms('تقنية');
console.log('Search Results:', searchResult.results);
console.log('Sources used:', searchResult.sources);

const itunesResults = await searchExternalSource('itunes', 'technology', 5);
console.log('iTunes Results:', itunesResults);

const sources = await getAvailableSources();
console.log('Available Sources:', sources);

// Example: Search that triggers external fallback
const fallbackResult = await searchPrograms('nonexistentterm');
if (fallbackResult.sources?.external) {
  console.log('External sources triggered:', fallbackResult.sources.external);
}

// Create category first
const createCategory = async (name) => {
  const response = await fetch('http://localhost:4000/v1/cms/categories', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name })
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  return await response.json();
};

// Create program
const createProgram = async (programData) => {
  const response = await fetch('http://localhost:4000/v1/cms/programs', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(programData)
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  return await response.json();
};

// Usage example
const categoryData = await createCategory('تقنية');
const newProgram = await createProgram({
  title: 'برنامج جديد',
  description: 'وصف البرنامج',
  category_id: categoryData.category.id,
  language: 'ar',
  duration: 1800
});
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/new-feature`)
3. Commit changes (`git commit -am 'Add new feature'`)
4. Push to branch (`git push origin feature/new-feature`)
5. Create Pull Request

## 📄 License

This project is licensed under the Apache 2 License - see the LICENSE file for details.
