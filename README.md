# Gomania - Arabic Podcast Management System

A simple and clean podcast content management system built with Go, focusing on essential fields and Arabic content support.

## 🚀 Features

- **Simple CMS**: Clean content management for programs with essential fields only
- **Category Management**: Organize programs by categories
- **Arabic Content**: Full Arabic language support
- **External Source Integration**: Optional iTunes search integration with auto-import capability
- **Clean Architecture**: Simple layered design with clear separation of concerns
- **Type Safety**: SQLC-generated database queries
- **RESTful API**: Simple and intuitive API endpoints

## 📋 Requirements

- Go 1.24+
- dbmate
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
make docker-up
make gen
make db-up
```

### 5. Run Server
```bash
make build api
```

Server will start on `http://localhost:4000`

## 6. Database UI

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

## 🌐 API Documentation

[API.MD](API.md)

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
│   └── service/         # Business logic
│       └── program.go   # Program & category service
├── data/sql/
│   ├── migrations/      # Database migrations
│   ├── queries/         # SQL queries
│   └── seed.sql         # Sample data
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

### Database Commands
```bash
# Initialize database (migrations + sample data)
make db-init

# Reset database completely
make db-reset

# Run migrations only
make db-migrate

# Load sample data only
make db-seed
```

### Manual Testing
```bash
# Health check
curl http://localhost:4000/v1/healthcheck

# List programs
curl http://localhost:4000/v1/programs

# Search
curl "http://localhost:4000/v1/programs?q=تقنية"

# Search with external sources
curl "http://localhost:4000/v1/programs?q=technology&external=true"

# List categories
curl http://localhost:4000/v1/cms/categories

# Create category
curl -X POST http://localhost:4000/v1/cms/categories \
  -H "Content-Type: application/json" \
  -d '{"name":"تقنية"}'

# Create program
curl -X POST http://localhost:4000/v1/cms/programs \
  -H "Content-Type: application/json" \
  -d '{
    "title":"برنامج جديد",
    "description":"وصف البرنامج",
    "category_id":"550e8400-e29b-41d4-a716-446655440001",
    "language":"ar",
    "duration":1800
  }'
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

Load sample data with:
```bash
make db-seed
# OR
./scripts/init_db.sh
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
// Search programs
const searchPrograms = async (query) => {
  const response = await fetch(`http://localhost:4000/v1/programs?q=${encodeURIComponent(query)}`);
  const data = await response.json();
  return data.search ? data.search.results : data.programs;
};

// Usage
const programs = await searchPrograms('تقنية');
console.log(programs);

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
