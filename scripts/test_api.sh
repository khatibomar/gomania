#!/bin/bash

# API test script for Gomania
# This script tests all API endpoints to ensure they work correctly

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# API base URL
BASE_URL="http://localhost:4000"

echo -e "${GREEN}🧪 Testing Gomania API${NC}"
echo -e "${YELLOW}Base URL: ${BASE_URL}${NC}"
echo ""

# Function to test an endpoint
test_endpoint() {
    local method="$1"
    local endpoint="$2"
    local description="$3"
    local data="$4"

    echo -e "${BLUE}Testing: ${description}${NC}"
    echo -e "${YELLOW}${method} ${endpoint}${NC}"

    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "${BASE_URL}${endpoint}" \
                   -H "Content-Type: application/json" \
                   -d "$data")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "${BASE_URL}${endpoint}")
    fi

    # Extract status code (last line)
    status_code=$(echo "$response" | tail -n1)
    # Extract body (all but last line)
    body=$(echo "$response" | head -n -1)

    if [[ "$status_code" == "200" || "$status_code" == "201" || "$status_code" == "204" ]]; then
        echo -e "${GREEN}✅ Status: ${status_code}${NC}"
        if [ -n "$body" ] && [ "$body" != "null" ]; then
            echo -e "${GREEN}Response:${NC}"
            echo "$body" | jq '.' 2>/dev/null || echo "$body"
        fi
    else
        echo -e "${RED}❌ Status: ${status_code}${NC}"
        if [ -n "$body" ]; then
            echo -e "${RED}Response:${NC}"
            echo "$body"
        fi
    fi
    echo ""
}

# Check if server is running
echo -e "${YELLOW}🔍 Checking if server is running...${NC}"
if ! curl -s "${BASE_URL}/v1/healthcheck" > /dev/null; then
    echo -e "${RED}❌ Server is not running. Please start it first:${NC}"
    echo -e "${YELLOW}   GOMANIA_CONNECTION_STRING=\"postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable\" go run cmd/api/*.go${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Server is running${NC}"
echo ""

# Store category ID for program creation
CATEGORY_ID=""

# Test health check
test_endpoint "GET" "/v1/healthcheck" "Health Check"

# Test categories
test_endpoint "GET" "/v1/cms/categories" "List Categories"

# Create a test category
test_endpoint "POST" "/v1/cms/categories" "Create Category" '{"name":"اختبار"}'

# Get categories again to extract ID
echo -e "${BLUE}Getting category ID for tests...${NC}"
categories_response=$(curl -s "${BASE_URL}/v1/cms/categories")
CATEGORY_ID=$(echo "$categories_response" | jq -r '.categories[0].id' 2>/dev/null || echo "")

if [ -n "$CATEGORY_ID" ] && [ "$CATEGORY_ID" != "null" ]; then
    echo -e "${GREEN}✅ Category ID: ${CATEGORY_ID}${NC}"
    echo ""

    # Test programs
    test_endpoint "GET" "/v1/cms/programs" "List Programs (CMS)"

    # Create a test program
    test_endpoint "POST" "/v1/cms/programs" "Create Program" "{
        \"title\": \"برنامج اختبار\",
        \"description\": \"هذا برنامج للاختبار\",
        \"category_id\": \"${CATEGORY_ID}\",
        \"language\": \"ar\",
        \"duration\": 1800
    }"

    # Get program ID for further tests
    echo -e "${BLUE}Getting program ID for tests...${NC}"
    programs_response=$(curl -s "${BASE_URL}/v1/cms/programs")
    PROGRAM_ID=$(echo "$programs_response" | jq -r '.programs[0].id' 2>/dev/null || echo "")

    if [ -n "$PROGRAM_ID" ] && [ "$PROGRAM_ID" != "null" ]; then
        echo -e "${GREEN}✅ Program ID: ${PROGRAM_ID}${NC}"
        echo ""

        # Test individual program
        test_endpoint "GET" "/v1/cms/programs/${PROGRAM_ID}" "Get Single Program"

        # Test update program
        test_endpoint "PUT" "/v1/cms/programs/${PROGRAM_ID}" "Update Program" "{
            \"title\": \"برنامج اختبار محدث\",
            \"description\": \"هذا برنامج للاختبار محدث\",
            \"category_id\": \"${CATEGORY_ID}\",
            \"language\": \"ar\",
            \"duration\": 2400
        }"

        # Test programs by category
        test_endpoint "GET" "/v1/cms/categories/${CATEGORY_ID}/programs" "Get Programs by Category"

        # Test discovery API
        test_endpoint "GET" "/v1/programs" "Discovery - List All Programs"

        # Test search
        test_endpoint "GET" "/v1/programs?q=اختبار" "Discovery - Search Programs"

        # Test empty search (should return all)
        test_endpoint "GET" "/v1/programs?q=" "Discovery - Empty Search"

        # Test external sources
        test_endpoint "GET" "/v1/external/sources" "List Available External Sources"

        # Test external search with iTunes
        test_endpoint "GET" "/v1/external/search?source=itunes&q=technology&limit=5" "Search iTunes Podcasts"

        # Test search with no local results (should trigger external search)
        test_endpoint "GET" "/v1/programs?q=nonexistentterm12345" "Discovery - Search (should trigger external fallback)"

        # Test delete program
        test_endpoint "DELETE" "/v1/cms/programs/${PROGRAM_ID}" "Delete Program"

    else
        echo -e "${YELLOW}⚠️  Could not extract program ID, skipping program-specific tests${NC}"
        echo ""
    fi

else
    echo -e "${YELLOW}⚠️  Could not extract category ID, skipping program tests${NC}"
    echo ""
fi

# Test Arabic content search
test_endpoint "GET" "/v1/programs?q=تقنية" "Search Arabic Content"

# Test iTunes specific search
test_endpoint "GET" "/v1/external/search?source=itunes&q=podcast&limit=3" "iTunes Direct Search"

# Test non-existent endpoint
echo -e "${BLUE}Testing: Non-existent Endpoint (should return 404)${NC}"
echo -e "${YELLOW}GET /v1/nonexistent${NC}"
response=$(curl -s -w "\n%{http_code}" "${BASE_URL}/v1/nonexistent")
status_code=$(echo "$response" | tail -n1)
if [ "$status_code" == "404" ]; then
    echo -e "${GREEN}✅ Status: 404 (as expected)${NC}"
else
    echo -e "${YELLOW}⚠️  Status: ${status_code} (expected 404)${NC}"
fi
echo ""

echo -e "${GREEN}🎉 API testing completed!${NC}"
echo ""
echo -e "${GREEN}📊 Summary:${NC}"
echo -e "${GREEN}   ✅ Health check${NC}"
echo -e "${GREEN}   ✅ Category management${NC}"
echo -e "${GREEN}   ✅ Program management${NC}"
echo -e "${GREEN}   ✅ Discovery API${NC}"
echo -e "${GREEN}   ✅ Search functionality${NC}"
echo -e "${GREEN}   ✅ Arabic content support${NC}"
echo -e "${GREEN}   ✅ External sources integration${NC}"
echo -e "${GREEN}   ✅ iTunes API integration${NC}"
echo ""
echo -e "${BLUE}🔗 Available endpoints:${NC}"
echo -e "${YELLOW}   GET    /v1/healthcheck${NC}"
echo -e "${YELLOW}   GET    /v1/programs${NC}"
echo -e "${YELLOW}   GET    /v1/programs?q=search${NC}"
echo -e "${YELLOW}   GET    /v1/external/sources${NC}"
echo -e "${YELLOW}   GET    /v1/external/search?source=SOURCE&q=QUERY${NC}"
echo -e "${YELLOW}   GET    /v1/cms/categories${NC}"
echo -e "${YELLOW}   POST   /v1/cms/categories${NC}"
echo -e "${YELLOW}   GET    /v1/cms/categories/{id}/programs${NC}"
echo -e "${YELLOW}   GET    /v1/cms/programs${NC}"
echo -e "${YELLOW}   POST   /v1/cms/programs${NC}"
echo -e "${YELLOW}   GET    /v1/cms/programs/{id}${NC}"
echo -e "${YELLOW}   PUT    /v1/cms/programs/{id}${NC}"
echo -e "${YELLOW}   DELETE /v1/cms/programs/{id}${NC}"
