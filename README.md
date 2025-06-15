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

### Base URL
```
http://localhost:4000
```

### Health & Monitoring

#### Health Check
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

#### Debug Information
```http
GET /debug/vars
```

Returns server runtime statistics and metrics.

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
      "language": "ar",
      "duration": 1800,
      "category_name": "تقنية"
    }
  ]
}
```

#### Get Single Program
```http
GET /v1/cms/programs/{id}
```

**Parameters:**
- `id` (path, required): Program UUID

**Response:**
```json
{
  "program": {
    "id": "770e8400-e29b-41d4-a716-446655440001",
    "title": "تقنية بودكاست",
    "description": "برنامج أسبوعي يناقش أحدث التطورات في عالم التكنولوجيا",
    "language": "ar",
    "duration": 1800,
    "category_name": "تقنية"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid UUID format
- `404 Not Found`: Program not found

#### Create Program
```http
POST /v1/cms/programs
Content-Type: application/json

{
  "title": "برنامج جديد",
  "description": "وصف البرنامج",
  "category_id": "550e8400-e29b-41d4-a716-446655440001",
  "language": "ar",
  "duration": 1800
}
```

**Request Fields:**
- `title` (string, required): Program title
- `description` (string, optional): Program description
- `category_id` (string, optional): Category UUID
- `language` (string, optional): Language code (default: "ar")
- `duration` (integer, optional): Duration in seconds

**Response:** `201 Created`
```json
{
  "program": {
    "id": "generated-uuid",
    "title": "برنامج جديد",
    "description": "وصف البرنامج",
    "category_id": "550e8400-e29b-41d4-a716-446655440001",
    "language": "ar",
    "duration": 1800
  }
}
```

#### Update Program
```http
PUT /v1/cms/programs/{id}
Content-Type: application/json

{
  "title": "برنامج محدث",
  "description": "وصف محدث",
  "category_id": "550e8400-e29b-41d4-a716-446655440001",
  "language": "ar",
  "duration": 2000
}
```

**Parameters:**
- `id` (path, required): Program UUID

**Response:** `200 OK`
```json
{
  "program": {
    "id": "770e8400-e29b-41d4-a716-446655440001",
    "title": "برنامج محدث",
    "description": "وصف محدث",
    "category_id": "550e8400-e29b-41d4-a716-446655440001",
    "language": "ar",
    "duration": 2000
  }
}
```

#### Delete Program
```http
DELETE /v1/cms/programs/{id}
```

**Parameters:**
- `id` (path, required): Program UUID

**Response:** `204 No Content`

**Error Responses:**
- `400 Bad Request`: Invalid UUID format
- `404 Not Found`: Program not found

### Categories

#### List All Categories
```http
GET /v1/cms/categories
```

**Response:**
```json
{
  "categories": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "name": "تقنية"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "name": "تعليم"
    }
  ]
}
```

#### Create Category
```http
POST /v1/cms/categories
Content-Type: application/json

{
  "name": "فئة جديدة"
}
```

**Request Fields:**
- `name` (string, required): Category name

**Response:** `201 Created`
```json
{
  "category": {
    "id": "generated-uuid",
    "name": "فئة جديدة"
  }
}
```

#### Get Programs by Category
```http
GET /v1/cms/categories/{id}/programs
```

**Parameters:**
- `id` (path, required): Category UUID

**Response:**
```json
{
  "category": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "تقنية"
  },
  "programs": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440001",
      "title": "تقنية بودكاست",
      "description": "برنامج أسبوعي يناقش أحدث التطورات في عالم التكنولوجيا",
      "language": "ar",
      "duration": 1800
    }
  ]
}
```

---

## 🔍 Discovery API (Public)

### Browse Programs
```http
GET /v1/programs
```

### Search Programs
```http
GET /v1/programs?q={query}
```

**Query Parameters:**
- `q` (string, optional): Search query
- `external` (boolean, optional): Include external sources (iTunes) in search
- `import` (boolean, optional): Import external results if not found locally

**Examples:**
```http
# Basic search
GET /v1/programs?q=تقنية

# Empty query returns all programs
GET /v1/programs

# Search with external sources
GET /v1/programs?q=technology&external=true

# Search and auto-import
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
        "description": "برنامج تقني أسبوعي",
        "language": "ar",
        "duration": 1800,
        "category_name": "تقنية",
        "source": "local"
      }
    ],
    "count": 1
  }
}
```

## 🔗 External Source Integration

### iTunes Search Integration

The system automatically searches iTunes when:
1. Local search returns no results
2. `external=true` parameter is provided
3. User requests import with `import=true`

**Search Flow:**
1. Search local database first
2. If no results and `external=true`, search iTunes API
3. If `import=true`, automatically import iTunes results to local database
4. Return combined or imported results

## 📊 Error Responses

### Standard Error Format
```json
{
  "error": "Error message description"
}
```

### HTTP Status Codes
- `200 OK`: Successful request
- `201 Created`: Resource created successfully
- `204 No Content`: Successful deletion
- `400 Bad Request`: Invalid request format or parameters
- `404 Not Found`: Resource not found
- `405 Method Not Allowed`: HTTP method not supported
- `500 Internal Server Error`: Server error

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

### Python Examples

```python
import requests

def search_programs(query=None, external=False, import_results=False):
    url = "http://localhost:4000/v1/programs"
    params = {}

    if query:
        params['q'] = query
    if external:
        params['external'] = 'true'
    if import_results:
        params['import'] = 'true'

    response = requests.get(url, params=params)
    response.raise_for_status()

    data = response.json()
    return data.get('search', {}).get('results', data.get('programs', []))

# Usage
programs = search_programs('تقنية')
external_programs = search_programs('podcast', external=True, import_results=True)

def create_program(title, description=None, category_id=None, language='ar', duration=None):
    url = "http://localhost:4000/v1/cms/programs"
    data = {
        'title': title,
        'language': language
    }

    if description:
        data['description'] = description
    if category_id:
        data['category_id'] = category_id
    if duration:
        data['duration'] = duration

    response = requests.post(url, json=data)
    response.raise_for_status()

    return response.json()

# Create a category first
def create_category(name):
    url = "http://localhost:4000/v1/cms/categories"
    data = {'name': name}

    response = requests.post(url, json=data)
    response.raise_for_status()

    return response.json()

# Usage
category = create_category('تقنية جديدة')
program = create_program(
    title='برنامج جديد',
    description='وصف البرنامج',
    category_id=category['category']['id'],
    duration=1800
)
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/new-feature`)
3. Commit changes (`git commit -am 'Add new feature'`)
4. Push to branch (`git push origin feature/new-feature`)
5. Create Pull Request

## 📄 License

This project is licensed under the Apache 2 License - see the LICENSE file for details.
