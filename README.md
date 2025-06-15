# Gomania - Arabic Podcast Management System

A comprehensive podcast content management and discovery system built with Go, featuring Arabic content support and external source integration (iTunes API).

## 🚀 Features

- **CMS System**: Internal content management for programs, episodes, and metadata
- **Discovery API**: Public search and browsing interface  
- **External Integration**: iTunes API import with extensible architecture
- **Arabic Content**: Full Arabic language support with RTL content
- **Structured Logging**: Comprehensive logging with slog
- **Type Safety**: SQLC-generated database queries
- **Clean Architecture**: Layered design with clear separation of concerns

## 📋 Requirements

- Go 1.24+
- PostgreSQL 17
- Docker & Docker Compose

## 🛠️ Installation & Setup

### 1. Clone Repository
```bash
git clone <repository-url>
cd gomania
```

### 2. Start Database
```bash
docker compose up -d database
```

### 3. Set Environment Variables
```bash
export GOMANIA_CONNECTION_STRING="postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable"
```

### 4. Initialize Database
```bash
# Generate database code
sqlc generate

# Initialize database with existing migrations
go run init_db.go
```

### 5. Run Server
```bash
go run cmd/api/*.go
```

Server will start on `http://localhost:4000`

## 📊 Sample Data

The system comes pre-loaded with 10 Arabic podcast programs:

1. **تقنية بودكاست** - Technology discussions
2. **ريادة الأعمال العربية** - Arabic entrepreneurship  
3. **علوم المستقبل** - Future sciences
4. **كوميديا الشارع** - Street comedy
5. **أخبار التقنية اليومية** - Daily tech news
6. **تعلم البرمجة** - Programming tutorials
7. **صوت الشباب** - Youth voices
8. **مستثمر ذكي** - Smart investing
9. **تاريخ وحضارة** - History and civilization
10. **صحة ولياقة** - Health and fitness

## 🌐 API Documentation

### Base URL
```
http://localhost:4000
```

### Health Check
```http
GET /v1/healthcheck
```

**Response:**
```json
{
  "status": "available",
  "system_info": {
    "environment": "development",
    "version": "1.0.0"
  }
}
```

---

## 🔒 CMS API (Content Management)

### Programs

#### List All Programs
```http
GET /v1/cms/programs
```

**Response:**
```json
{
  "programs": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440001",
      "title": "تقنية بودكاست",
      "description": "برنامج أسبوعي يناقش أحدث التطورات في عالم التكنولوجيا",
      "summary": "برنامج تقني أسبوعي",
      "language": "ar",
      "country": "SA",
      "author": "أحمد محمد",
      "publisher": "شبكة تقنية",
      "status": "active",
      "total_episodes": 25,
      "rating": 4.5,
      "source": "local",
      "published_at": "2024-01-15T10:00:00Z"
    }
  ]
}
```

#### Get Single Program
```http
GET /v1/cms/programs/{id}
```

#### Create Program
```http
POST /v1/cms/programs
Content-Type: application/json

{
  "title": "برنامج جديد",
  "description": "وصف البرنامج",
  "category": "تقنية",
  "language": "ar",
  "duration": 1800,
  "published_at": "2024-01-15T10:00:00Z"
}
```

#### Update Program
```http
PUT /v1/cms/programs/{id}
Content-Type: application/json

{
  "title": "برنامج محدث",
  "description": "وصف محدث",
  "category": "تقنية",
  "language": "ar",
  "duration": 2000
}
```

#### Delete Program
```http
DELETE /v1/cms/programs/{id}
```

### Episodes, Categories, and Import Features

*Note: Episode management, category management, and iTunes import features are planned for future releases. Currently, only basic program management is implemented.*

---

## 🔍 Discovery API (Public)

### Browse Programs
```http
GET /v1/programs
```

**Response:**
```json
{
  "programs": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440001",
      "title": "تقنية بودكاست",
      "description": "برنامج أسبوعي يناقش أحدث التطورات",
      "language": "ar",
      "country": "SA",
      "author": "أحمد محمد",
      "rating": 4.5,
      "total_episodes": 25
    }
  ]
}
```

### Search Programs
```http
GET /v1/programs?q={query}
```

**Parameters:**
- `q` (required): Search query
- `external=true`: Search external sources (iTunes)
- `import=true`: Import results if not found locally

**Examples:**
```http
# Basic search
GET /v1/programs?q=تقنية

# Search with external sources (iTunes integration)
GET /v1/programs?q=technology&external=true

# Search and import if not found locally
GET /v1/programs?q=podcast&external=true&import=true
```

**Response:**
```json
{
  "search": {
    "query": "تقنية",
    "results": [
      {
        "id": "770e8400-e29b-41d4-a716-446655440001",
        "title": "تقنية بودكاست",
        "description": "برنامج تقني أسبوعي"
      }
    ],
    "count": 1
  }
}
```

## 🏗️ Architecture

### Project Structure
```
gomania/
├── cmd/api/              # API server
│   ├── main.go          # Entry point
│   ├── routes.go        # Route definitions
│   ├── cms.go           # CMS handlers
│   ├── errors.go        # Error handlers
│   └── ...
├── internal/
│   ├── database/        # SQLC generated code
│   ├── service/         # Business logic
│   │   └── program.go   # Program service
│   └── sources/         # External sources
│       ├── client.go    # Interface
│       ├── manager.go   # Source manager
│       └── itunes/      # iTunes client
├── data/sql/
│   ├── migrations/      # Database migrations
│   └── queries/         # SQL queries
└── docker-compose.yaml  # Database setup
```

### Database Schema

#### Core Tables
- `programs` - Podcast programs
- `episodes` - Individual episodes  
- `categories` - Program categories
- `external_sources` - External source tracking
- `users` - CMS users
- `tags` - Flexible tagging

#### Relationships
- Programs → Episodes (1:many)
- Programs → Categories (many:1)  
- Programs ↔ External Sources (1:many)
- Programs ↔ Tags (many:many)

## 🔌 External Sources

### Current Integrations
- **iTunes API**: Podcast search and import functionality

### Adding New Sources

The system is designed to support multiple external sources. To add a new source:

1. **Implement the Client Interface:**
```go
// internal/sources/spotify/client.go
type SpotifyClient struct{}

func (s *SpotifyClient) SearchPodcasts(term string, limit int) ([]sources.Podcast, error) {
    // Your implementation here
    return podcasts, nil
}

func (s *SpotifyClient) GetSourceName() string {
    return "spotify"
}
```

2. **Register the Client:**
```go
// In internal/service/program.go NewProgramService function
sourcesManager.RegisterClient(&spotify.SpotifyClient{})
```

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
  -cors-trusted-origins="https://mydomain.com https://anotherdomain.com"
```

## 🧪 Testing

### Manual Testing
```bash
# Health check
curl http://localhost:4000/v1/healthcheck

# List programs
curl http://localhost:4000/v1/programs

# Search
curl "http://localhost:4000/v1/programs?q=تقنية"

# Create program (CMS)
curl -X POST http://localhost:4000/v1/cms/programs \
  -H "Content-Type: application/json" \
  -d '{"title":"برنامج جديد","description":"وصف","category":"تقنية"}'
```

## 📊 Logging

The system uses structured logging with `slog`:

```
level=INFO msg="Searching programs" query=تقنية external=false found=2
level=INFO msg="Creating new program" title="My Podcast"
level=INFO msg="Program created successfully" id=abc123...
level=ERROR msg="Failed to create program" title="Bad Program" error="validation failed"
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

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/new-feature`)
3. Commit changes (`git commit -am 'Add new feature'`)
4. Push to branch (`git push origin feature/new-feature`)
5. Create Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🔗 API Client Examples

### JavaScript/Node.js
```javascript
// Search programs
const response = await fetch('http://localhost:4000/v1/programs?q=تقنية');
const data = await response.json();
console.log(data.programs);

// Create program
const program = await fetch('http://localhost:4000/v1/cms/programs', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    title: 'برنامج جديد',
    description: 'وصف البرنامج',
    category: 'تقنية'
  })
});
```

### cURL Examples
```bash
# Get all programs
curl -X GET "http://localhost:4000/v1/programs"

# Search with external sources and auto-import
curl -X GET "http://localhost:4000/v1/programs?q=podcast&external=true&import=true"

# Create new program
curl -X POST "http://localhost:4000/v1/cms/programs" \
     -H "Content-Type: application/json" \
     -d '{
       "title": "برنامج جديد",
       "description": "وصف مفصل للبرنامج",
       "category": "تقنية",
       "language": "ar",
       "duration": 1800
     }'

# Update existing program
curl -X PUT "http://localhost:4000/v1/cms/programs/{id}" \
     -H "Content-Type: application/json" \
     -d '{
       "title": "برنامج محدث",
       "description": "وصف محدث"
     }'

# Delete program
curl -X DELETE "http://localhost:4000/v1/cms/programs/{id}"
```
