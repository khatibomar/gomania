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

#### Search with External Sources
```http
GET /v1/programs?q=technology&external=true
```

#### Search and Auto-Import
```http
GET /v1/programs?q=podcast&external=true&import=true
```

**Search Response:**
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
    "count": 1
  }
}
```

---

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

#### Search with iTunes Integration
```bash
curl -X GET "http://localhost:4000/v1/programs?q=technology&external=true&import=true"
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

### Python Examples

#### Search and List Programs
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
```

#### Create Program
```python
import requests

def create_program(title, description=None, category=None, language='ar', duration=None):
    url = "http://localhost:4000/v1/cms/programs"
    data = {
        'title': title,
        'language': language
    }
    
    if description:
        data['description'] = description
    if category:
        data['category'] = category
    if duration:
        data['duration'] = duration
    
    response = requests.post(url, json=data)
    response.raise_for_status()
    
    return response.json()

# Usage
program = create_program(
    title='برنامج جديد',
    description='وصف البرنامج',
    category='تقنية',
    duration=1800
)
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

### Future Endpoints
The following endpoints are planned for future releases:
- Episode management (`/v1/cms/episodes/*`)
- Category management (`/v1/cms/categories/*`)
- User management (`/v1/cms/users/*`)
- Tag management (`/v1/cms/tags/*`)
- Direct iTunes import (`/v1/cms/import/itunes/{id}`)
- Analytics and statistics (`/v1/cms/analytics/*`)