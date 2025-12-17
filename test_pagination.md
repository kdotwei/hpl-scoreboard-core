# Pagination API 测试

## 新增的 API 端点

### 1. 分页列表 API
```
GET /api/v1/scores/paginated?limit=10&cursor=<uuid>
```

#### 查询参数:
- `limit` (可选): 每页的记录数量，默认 10，最大 100
- `cursor` (可选): 游标，用于获取下一页数据，使用上一页返回的 `next_cursor` 值

#### 响应格式:
```json
{
  "scores": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "user_id": "user1",
      "gflops": 1234.56,
      "problem_size_n": 20000,
      "block_size_nb": 512,
      "submitted_at": "2023-01-01T00:00:00Z",
      "linux_username": "hpl_user",
      "n": 20000,
      "nb": 512,
      "p": 4,
      "q": 4,
      "execution_time": 125.75
    }
  ],
  "next_cursor": "550e8400-e29b-41d4-a716-446655440001",
  "has_more": true,
  "total_records": 150
}
```

## 前端使用示例

### 初始加载
```javascript
const response = await fetch('/api/v1/scores/paginated?limit=20');
const data = await response.json();
console.log(data.scores); // 显示第一页数据
```

### 加载更多 (Auto-load)
```javascript
if (data.has_more && data.next_cursor) {
  const nextResponse = await fetch(`/api/v1/scores/paginated?limit=20&cursor=${data.next_cursor}`);
  const nextData = await nextResponse.json();
  // 将 nextData.scores 追加到现有列表
}
```

### React 示例 (Infinite Scroll)
```jsx
const [scores, setScores] = useState([]);
const [nextCursor, setNextCursor] = useState(null);
const [hasMore, setHasMore] = useState(true);
const [loading, setLoading] = useState(false);

const loadMoreScores = async () => {
  if (loading) return;
  
  setLoading(true);
  try {
    const url = nextCursor 
      ? `/api/v1/scores/paginated?limit=20&cursor=${nextCursor}`
      : '/api/v1/scores/paginated?limit=20';
    
    const response = await fetch(url);
    const data = await response.json();
    
    setScores(prev => [...prev, ...data.scores]);
    setNextCursor(data.next_cursor);
    setHasMore(data.has_more);
  } catch (error) {
    console.error('Failed to load scores:', error);
  } finally {
    setLoading(false);
  }
};

// 使用 Intersection Observer 或 onScroll 事件触发 loadMoreScores
```

## 优势

1. **高效**: 使用基于 cursor 的分页，避免了 OFFSET 的性能问题
2. **一致性**: 即使有新数据插入，也不会影响分页结果
3. **适合无限滚动**: 前端可以轻松实现 auto-load 功能
4. **向后兼容**: 保留了原有的 `/api/v1/scores` 端点

## 注意事项

1. cursor 参数必须是有效的 UUID 格式
2. limit 参数最大值为 100，防止单次请求过大
3. 数据按 gflops 降序排列，相同 gflops 按 id 降序排列，保证排序一致性