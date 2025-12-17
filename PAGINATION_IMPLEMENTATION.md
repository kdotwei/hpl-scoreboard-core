# LIMIT/OFFSET Pagination Implementation

## Summary

Successfully implemented LIMIT/OFFSET pagination for the HPL Scoreboard backend to support lazy loading in the React frontend.

## Changes Made

### 1. SQL Updates (internal/db/query/score.sql)
- Updated `ListTopScores` query to accept both `LIMIT $1 OFFSET $2` parameters
- Maintained `ORDER BY gflops DESC` for consistent sorting

### 2. Service Layer Updates

#### internal/service/service.go
- Updated `ListScores` interface method: `ListScores(ctx context.Context, limit int32, offset int32)`
- Updated `ListScoresParams` struct to use `Limit` and `Offset` fields instead of cursor
- Updated `PaginatedScoresResponse` to include `Limit` and `Offset` fields

#### internal/service/score.go
- Modified `ListScores` implementation to pass both limit and offset to database layer
- Updated `ListScoresWithPagination` to use LIMIT/OFFSET approach
- Improved `HasMore` calculation: `hasMore := int64(params.Offset+int32(len(scores))) < totalRecords`

### 3. Handler Layer (internal/handler/score.go)

#### ListScores Handler
- Added parsing for both `limit` and `offset` query parameters
- **Default values**: limit=10, offset=0  
- **Error handling**: Returns `400 Bad Request` for invalid parameters
- **Validation**: 
  - `limit` must be > 0
  - `offset` must be >= 0

#### ListScoresWithPagination Handler  
- Updated to use LIMIT/OFFSET instead of cursor-based pagination
- Added proper error handling for invalid limit (1-100) and offset (>=0) parameters
- Returns structured response with pagination metadata

### 4. CORS & Routing
- Maintained existing public access to `GET /api/v1/scores` 
- Kept `GET /api/v1/scores/paginated` route for enhanced pagination
- CORS middleware unchanged - continues to allow required headers

### 5. Testing Updates (internal/handler/score_test.go)

#### Updated TestListScores
- Added `expectedOffset` field to test cases
- Updated mock calls to include both limit and offset parameters  
- Added test cases for:
  - Successful offset usage (`?limit=5&offset=10`)
  - Invalid limit returns `400 Bad Request` 
  - Invalid offset returns `400 Bad Request`
  - Negative values return `400 Bad Request`

#### Updated TestListScoresWithPagination
- Replaced cursor-based tests with LIMIT/OFFSET tests
- Added comprehensive error handling test cases
- Verified different offset values return different data subsets

### 6. Mock Updates (internal/service/mocks/Service.go)
- Updated `ListScores` method signature to match new interface
- Fixed return function signatures for proper parameter handling

## API Usage Examples

### Basic Pagination
```bash
# First page (default)
GET /api/v1/scores

# First page with custom limit  
GET /api/v1/scores?limit=20

# Second page
GET /api/v1/scores?limit=20&offset=20

# Third page  
GET /api/v1/scores?limit=20&offset=40
```

### Enhanced Pagination Endpoint
```bash
# Get paginated response with metadata
GET /api/v1/scores/paginated?limit=10&offset=0
```

**Response Format:**
```json
{
  "scores": [...],
  "has_more": true,
  "total_records": 150,
  "limit": 10,
  "offset": 0
}
```

## Frontend Integration (React)

```javascript
const [scores, setScores] = useState([]);
const [offset, setOffset] = useState(0);
const [hasMore, setHasMore] = useState(true);
const limit = 20;

const loadMoreScores = async () => {
  const response = await fetch(`/api/v1/scores/paginated?limit=${limit}&offset=${offset}`);
  const data = await response.json();
  
  setScores(prev => [...prev, ...data.scores]);
  setOffset(prev => prev + limit);
  setHasMore(data.has_more);
};

// Use with Intersection Observer or onScroll for lazy loading
```

## Error Handling

- **400 Bad Request**: Invalid limit or offset parameters
- **500 Internal Server Error**: Database or service layer errors  
- **Validation**: limit must be 1-100, offset must be >= 0

## Performance Benefits

- **Consistent Performance**: LIMIT/OFFSET provides predictable query performance for reasonable page sizes
- **Simple Implementation**: Easy to understand and implement in frontend
- **Backward Compatible**: Original `/api/v1/scores` endpoint unchanged
- **Flexible**: Frontend can adjust page sizes as needed for lazy loading

## Testing

All tests pass successfully:
- Unit tests for handler layer with various pagination scenarios
- Error handling validation for invalid parameters  
- Mock service layer properly handles new interface
- Build verification confirms all components work together

The implementation is ready for production use with React frontend lazy loading.