# Gomania API Documentation

## Base URL
```
http://localhost:4000
```

## Authentication
Currently, no authentication is required for any endpoints. Authentication will be added in future versions for CMS endpoints.

---

## 🏥 Health & Monitoring

### Health Check
**GET** `/v1/healthcheck`

Check if the API server is running and healthy.

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

### Debug Information
**GET** `/debug/vars`

Returns server runtime statistics and metrics.

---

## 🔒 CMS API (Content Management System)

### Programs

#### List All Programs
**GET** `/v1/cms/programs`

Retrieve all programs in the system.

**Response:**
```json
{
  "programs": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440001",
      "title": "تقنية بودكاست",
      "description": "برنامج أسبوعي يناقش أحدث التطورات في عالم التكنولوجيا والبرمجة",
      "category": "تقنية",
      "language": "ar",
      "duration": 1800,
      "published_at": "2024-01-15T10:00:00Z",
      "source": "local"
    }
  ]
}
```

#### Get Single Program
**GET** `/v1/cms/programs/{id}`

Retrieve a specific program by its UUID.

**Parameters:**
- `id` (path, required): Program UUID

**Response:**
```json
{
  "program": {
    "id": "770e8400-e29b-41d4-a716-446655440001",
    "title": "تقنية بودكاست",
    "description": "برنامج أسبوعي يناقش أحدث التطورات في عالم التكنولوجيا",
    "category": "تقنية",
    "language": "ar",
    "duration": 1800,
    "published_at": "2024-01-15T10:00:00Z",
    "source": "local"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid UUID format
- `404 Not Found`: Program not found

#### Create Program
**POST** `/v1/cms/programs`

Create a new program.

**Request Body:**
```json
{
  "title": "برنامج جديد",
  "description": "وصف مفصل للبرنامج الجديد",
  "category": "تقنية",
  "language": "ar",
  "duration": 1800,
  "published_at": "2024-01-15T10:00:00Z"
}
```

**Request Fields:**
- `title` (string, required): Program title
- `description` (string, optional): Program description
- `category` (string, optional): Program category
- `language` (string, optional): Language code (default: "ar")
- `duration` (integer, optional): Duration in seconds
- `published_at` (string, optional): ISO 8601 timestamp

**Response:** `201 Created`
```json
{
  "program": {
    "id": "generated-uuid",
    "title": "برنامج جديد",
    "description": "وصف مفصل للبرنامج الجديد",
    "category": "تقنية",
    "language": "ar",
    "duration": 1800,
    "published_at": "2024-01-15T10:00:00Z",
    "source": "local"
  }
}
```

#### Update Program
**PUT** `/v1/cms/programs/{id}`

Update an existing program.

**Parameters:**
- `id` (path, required): Program UUID

**Request Body:**
```json
{
  "title": "برنامج محدث",
  "description": "وصف محدث للبرنامج",
  "category": "تقنية محدثة",
  "language": "ar",
  "duration": 2000
}
```

**Response:** `200 OK`
```json
{
  "program": {
    "id": "770e8400-e29b-41d4-a716-446655440001",
    "title": "برنامج محدث",
    "description": "وصف محدث للبرنامج",
    "category": "تقنية محدثة",
    "language": "ar",
    "duration": 2000,
    "published_at": "2024-01-15T10:00:00Z",
    "source": "local"
  }
}
```

#### Delete Program
**DELETE** `/v1/cms/programs/{id}`

Delete a program permanently.

**Parameters:**
- `id` (path, required): Program UUID

**Response:** `204 No Content`

**Error Responses:**
- `400 Bad Request`: Invalid UUID format
- `404 Not Found`: Program not found

---

## 🔍 Discovery API (Public)

### Browse Programs
**GET** `/v1/programs`

Browse all available programs. This endpoint serves both as a simple listing and as a search endpoint when query parameters are provided.

**Query Parameters:**
- `q` (string, optional): Search query
- `external` (boolean, optional): Include external sources (iTunes) in search
- `import` (boolean, optional): Import external results if not found locally

**Examples:**

#### Simple Browse
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
      "description": "برنامج أسبوعي يناقش أحدث التطورات في عالم التكنولوجيا",
      "category": "تقنية",
      "language": "ar",
      "duration": 1800,
      "published_at": "2024-01-15T10:00:00Z",
      "source": "local"
    }
  ]
}
```

#### Search Local Programs
```http
GET /v1/programs?q=تقنية
```

#### Search with Automatic External Fallback
```http
GET /v1/programs?q=technology
```

When no local results are found, the system automatically searches external sources.

**Search Response with External Fallback:**
```json
{
  "search": {
    "query": "technology",
    "results": [],
    "count": 0,
    "sources": {
      "local": {
        "count": 0
      },
      "external": {
        "itunes": [
          {
            "id": "12345",
            "title": "Tech Talk Podcast",
            "description": "Latest technology discussions",
            "host": "John Doe",
            "genre": "Technology",
            "country": "US",
            "duration": 3600,
            "published_at": "2024-01-15T10:00:00Z",
            "artwork_url": "https://example.com/artwork.jpg"
          }
        ]
      }
    },
    "external_count": 1
  }
}
```

**Local Search Response:**
```json
{
  "search": {
    "query": "تقنية",
    "results": [
      {
        "id": "770e8400-e29b-41d4-a716-446655440001",
        "title": "تقنية بودكاست",
        "description": "برنامج أسبوعي يناقش أحدث التطورات",
        "category": "تقنية",
        "language": "ar",
        "duration": 1800,
        "source": "local"
      }
    ],
    "count": 1,
    "sources": {
      "local": {
        "count": 1
      }
    }
  }
}
```

---

## 🔗 External Sources API

### List Available External Sources
**GET** `/v1/external/sources`

Get a list of all available external sources.

**Response:**
```json
{
  "external_sources": {
    "sources": ["itunes"],
    "count": 1
  }
}
```

### Search Specific External Source
**GET** `/v1/external/search`

Search a specific external source directly.

**Query Parameters:**
- `source` (string, required): Source name (e.g., "itunes")
- `q` (string, required): Search query
- `limit` (integer, optional): Maximum results to return (default: 10)

**Examples:**

#### Search iTunes
```http
GET /v1/external/search?source=itunes&q=technology&limit=5
```

**Response:**
```json
{
  "external_search": {
    "query": "technology",
    "source": "itunes",
    "results": [
      {
        "id": "12345",
        "title": "Tech Talk Podcast",
        "description": "Latest technology discussions and trends",
        "host": "John Doe",
        "genre": "Technology",
        "country": "US",
        "duration": 3600,
        "published_at": "2024-01-15T10:00:00Z",
        "artwork_url": "https://example.com/artwork.jpg"
      }
    ],
    "count": 1
  }
}
```

**Error Responses:**
- `400 Bad Request`: Missing required parameters
- `500 Internal Server Error`: External source unavailable

### iTunes Search Integration

The discovery endpoint (`/v1/programs`) automatically searches iTunes when:
1. Local search returns no results
2. Search query is provided with `?q=` parameter

**Automatic Search Flow:**
1. Search local database first
2. If no local results found, automatically search iTunes
3. Return combined response with both local and external results

---

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

### Common Error Examples

#### Invalid UUID Format
```json
{
  "error": "invalid UUID length: 5"
}
```

#### Resource Not Found
```json
{
  "error": "the requested resource could not be found"
}
```

#### Invalid JSON
```json
{
  "error": "invalid character '}' looking for beginning of object key string"
}
```

---

## 🧪 Testing Examples

### cURL Examples

#### Health Check
```bash
curl -X GET "http://localhost:4000/v1/healthcheck"
```

#### List All Programs
```bash
curl -X GET "http://localhost:4000/v1/cms/programs"
```

#### Create Program
```bash
curl -X POST "http://localhost:4000/v1/cms/programs" \
     -H "Content-Type: application/json" \
     -d '{
       "title": "برنامج تجريبي",
       "description": "هذا برنامج للاختبار",
       "category": "تقنية",
       "language": "ar",
       "duration": 1800
     }'
```

#### Search Programs
```bash
curl -X GET "http://localhost:4000/v1/programs?q=تقنية"
```

#### Search External Sources Directly
```bash
curl -X GET "http://localhost:4000/v1/external/search?source=itunes&q=technology&limit=5"
```

#### List Available External Sources
```bash
curl -X GET "http://localhost:4000/v1/external/sources"
```

#### Update Program
```bash
curl -X PUT "http://localhost:4000/v1/cms/programs/770e8400-e29b-41d4-a716-446655440001" \
     -H "Content-Type: application/json" \
     -d '{
       "title": "برنامج محدث",
       "description": "وصف جديد"
     }'
```

#### Delete Program
```bash
curl -X DELETE "http://localhost:4000/v1/cms/programs/770e8400-e29b-41d4-a716-446655440001"
```

### JavaScript Examples

#### Search Programs
```javascript
const searchPrograms = async (query) => {
  const response = await fetch(`http://localhost:4000/v1/programs?q=${encodeURIComponent(query)}`);
  const data = await response.json();
  return data.search ? data.search.results : data.programs;
};

// Usage
const programs = await searchPrograms('تقنية');
console.log(programs);
```

#### Create Program
```javascript
const createProgram = async (programData) => {
  const response = await fetch('http://localhost:4000/v1/cms/programs', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(programData)
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  return await response.json();
};

// Usage
const newProgram = await createProgram({
  title: 'برنامج جديد',
  description: 'وصف البرنامج',
  category: 'تقنية'
});
```

---

## 📝 Notes

### Arabic Content Support
- All text fields support Arabic content with proper UTF-8 encoding
- RTL (Right-to-Left) text is properly handled
- Arabic search queries are fully supported

### Rate Limiting
Currently, no rate limiting is implemented. This will be added in future versions.

### Pagination
Currently, all endpoints return complete result sets. Pagination will be added for large datasets in future versions.

### Categories

#### List All Categories
**GET** `/v1/cms/categories`

Retrieve all categories in the system.

**Response:**
```json
{
  "categories": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "name": "تقنية",
      "created_at": "2024-01-15T10:00:00Z"
    }
  ]
}
```

#### Create Category
**POST** `/v1/cms/categories`

Create a new category.

**Request Body:**
```json
{
  "name": "تقنية"
}
```

**Response:** `201 Created`
```json
{
  "category": {
    "id": "generated-uuid",
    "name": "تقنية",
    "created_at": "2024-01-15T10:00:00Z"
  }
}
```

**Error Responses:**
- `409 Conflict`: Category with this name already exists

#### Get Programs by Category
**GET** `/v1/cms/categories/{id}/programs`

Retrieve all programs in a specific category.

**Parameters:**
- `id` (path, required): Category UUID

**Response:**
```json
{
  "programs": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440001",
      "title": "تقنية بودكاست",
      "description": "برنامج أسبوعي يناقش أحدث التطورات",
      "category": "تقنية",
      "language": "ar",
      "duration": 1800,
      "source": "local"
    }
  ]
}
```

---

## 📋 Complete API Endpoints Summary

### Health & Monitoring
- `GET /v1/healthcheck` - Health check
- `GET /debug/vars` - Debug information

### Discovery API (Public)
- `GET /v1/programs` - Browse/search programs with automatic external fallback
- `GET /v1/programs?q={query}` - Search programs (auto-searches iTunes if no local results)

### External Sources
- `GET /v1/external/sources` - List available external sources
- `GET /v1/external/search?source={source}&q={query}&limit={limit}` - Search specific external source

### CMS - Programs
- `GET /v1/cms/programs` - List all programs
- `POST /v1/cms/programs` - Create new program
- `GET /v1/cms/programs/{id}` - Get single program
- `PUT /v1/cms/programs/{id}` - Update program
- `DELETE /v1/cms/programs/{id}` - Delete program

### CMS - Categories
- `GET /v1/cms/categories` - List all categories
- `POST /v1/cms/categories` - Create new category
- `GET /v1/cms/categories/{id}/programs` - Get programs by category

---

### Future Endpoints
The following endpoints are planned for future releases:
- Episode management (`/v1/cms/episodes/*`)
- User management (`/v1/cms/users/*`)
- Tag management (`/v1/cms/tags/*`)
- Direct iTunes import (`/v1/cms/import/itunes/{id}`)
- Analytics and statistics (`/v1/cms/analytics/*`)
- Bulk import from external sources (`/v1/cms/import/bulk`)
- Subscription management (`/v1/cms/subscriptions/*`)
