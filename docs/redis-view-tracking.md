# Redis-Based View Tracking System

This document describes the Redis-based view tracking system that prevents view count inflation through IP address throttling.

## 🎯 Overview

The view tracking system uses Redis to implement IP-based throttling, ensuring that each unique IP address can only increment a post's view count once within a 24-hour period. This prevents users from artificially inflating view counts by repeatedly refreshing pages.

## 🏗️ Architecture

### Components

1. **Redis Client** (`pkg/redis/client.go`)
   - Wrapper around go-redis/v9
   - Provides connection management and basic operations

2. **View Tracking Service** (`internal/service/view/view_tracking_service.go`)
   - Core business logic for view tracking
   - Handles IP-based throttling using Redis
   - Manages database view count updates

3. **IP Utility** (`pkg/utils/ip.go`)
   - Extracts real client IP from requests
   - Handles various proxy headers (X-Forwarded-For, X-Real-IP, etc.)
   - Validates IP addresses

4. **Enhanced Post Repository**
   - Added view tracking methods to MongoDB repository
   - Supports view count operations and analytics

## 🔧 Configuration

### Redis Configuration (`config.dev.yaml`)
```yaml
redis:
  host: "localhost"
  port: "6379"
  password: ""
  db: 0
  view_tracking_ttl: 86400  # 24 hours in seconds
```

### Environment Setup
1. **Install Redis** (if not already installed):
   ```bash
   # Windows (using Chocolatey)
   choco install redis-64
   
   # macOS (using Homebrew)
   brew install redis
   
   # Ubuntu/Debian
   sudo apt-get install redis-server
   ```

2. **Start Redis Server**:
   ```bash
   redis-server
   ```

3. **Verify Redis Connection**:
   ```bash
   redis-cli ping
   # Should return: PONG
   ```

## 🚀 How It Works

### View Tracking Flow

1. **User Requests Post**: `GET /api/v1/posts/:id`
2. **Extract IP Address**: System extracts real client IP using utility function
3. **Check Redis**: Look for key `post_view:{postID}:{ipAddress}`
4. **Decision Logic**:
   - **Key Exists**: View already counted, skip increment
   - **Key Missing**: New view, proceed with tracking
5. **Track View**: 
   - Set Redis key with 24-hour TTL
   - Increment view count in MongoDB
6. **Return Post**: Serve post data to user

### Redis Key Structure
```
post_view:{postID}:{ipAddress}
```

**Examples**:
- `post_view:507f1f77bcf86cd799439011:192.168.1.100`
- `post_view:507f1f77bcf86cd799439012:10.0.0.5`

### TTL (Time To Live)
- **Default**: 24 hours (86400 seconds)
- **Configurable**: Via `redis.view_tracking_ttl` in config
- **Purpose**: After TTL expires, same IP can increment view count again

## 📊 New API Endpoints

### View Analytics Endpoints

#### 1. Get Popular Posts
```http
GET /api/v1/posts/popular?limit=10
```

**Response**:
```json
{
  "posts": [
    {
      "id": "507f1f77bcf86cd799439011",
      "title": "Most Popular Post",
      "view_count": 1250,
      "author_id": "507f1f77bcf86cd799439001",
      "created_at": "2025-08-07T10:30:00Z"
    }
  ],
  "count": 10
}
```

#### 2. Get Total View Statistics
```http
GET /api/v1/stats/views
```

**Response**:
```json
{
  "total_views": 15420
}
```

#### 3. Get Post View Count
```http
GET /api/v1/posts/:id/views
```

**Response**:
```json
{
  "post_id": "507f1f77bcf86cd799439011",
  "view_count": 1250
}
```

## 🛡️ IP Address Handling

### Supported Headers (in order of preference)
1. `X-Forwarded-For` - Most common proxy header
2. `X-Real-IP` - Nginx proxy header
3. `X-Client-IP` - Alternative client IP header
4. `CF-Connecting-IP` - Cloudflare header
5. `RemoteAddr` - Direct connection fallback

### IP Validation
- Validates IP format using `net.ParseIP()`
- Handles both IPv4 and IPv6 addresses
- Includes private IP detection utility

### Example IP Extraction
```go
clientIP := utils.GetClientIP(c)
// Returns: "192.168.1.100" or "2001:db8::1"
```

## 🔄 Graceful Degradation

The system is designed to work gracefully even when Redis is unavailable:

1. **Redis Connection Failure**: 
   - System logs warning but continues
   - View tracking is disabled
   - Post retrieval still works normally

2. **Redis Operation Failure**:
   - Errors are logged but don't block requests
   - View count may not be perfectly accurate
   - User experience remains unaffected

3. **Fallback Behavior**:
   - Posts can still be retrieved without view tracking
   - Analytics endpoints return appropriate errors
   - System remains functional

## 📈 Database Schema Changes

### MongoDB Post Collection
The existing `view_count` field is used, with new repository methods:

```javascript
// Example post document
{
  "_id": ObjectId("507f1f77bcf86cd799439011"),
  "title": "Sample Post",
  "content": "Post content...",
  "author_id": ObjectId("507f1f77bcf86cd799439001"),
  "view_count": 1250,  // ← Tracked by Redis system
  "created_at": ISODate("2025-08-07T10:30:00Z"),
  "updated_at": ISODate("2025-08-07T10:30:00Z")
}
```

### New Repository Methods
- `IncrementViewCount(ctx, postID)` - Increment view count
- `GetViewCount(ctx, postID)` - Get current view count
- `GetTotalViews(ctx)` - Get total views across all posts
- `GetPostsByViewCount(ctx, limit)` - Get popular posts
- `ResetViewCount(ctx, postID)` - Reset view count (admin)

## 🧪 Testing

### Manual Testing with Postman

1. **Test View Tracking**:
   ```http
   GET http://localhost:8080/api/v1/posts/{post-id}
   ```
   - First request: View count increments
   - Subsequent requests: View count stays same (within 24h)

2. **Test from Different IPs**:
   - Use different proxy headers to simulate different IPs
   - Each unique IP should increment the count once

3. **Test Analytics**:
   ```http
   GET http://localhost:8080/api/v1/posts/popular
   GET http://localhost:8080/api/v1/stats/views
   ```

### Redis Monitoring
```bash
# Monitor Redis operations
redis-cli monitor

# Check specific keys
redis-cli keys "post_view:*"

# Check TTL for a key
redis-cli ttl "post_view:507f1f77bcf86cd799439011:192.168.1.100"
```

## 🔧 Configuration Options

### Redis Settings
```yaml
redis:
  host: "localhost"           # Redis server host
  port: "6379"               # Redis server port
  password: ""               # Redis password (if required)
  db: 0                      # Redis database number
  view_tracking_ttl: 86400   # TTL in seconds (24 hours)
```

### Customization Options
- **TTL Duration**: Adjust `view_tracking_ttl` for different throttling periods
- **Redis Database**: Use different `db` numbers for isolation
- **Connection Settings**: Configure host/port for different Redis instances

## 🚨 Error Handling

### Common Scenarios

1. **Redis Unavailable**:
   ```
   ⚠️ Redis connection failed: dial tcp [::1]:6379: connect: connection refused
   ```
   - System continues without view tracking
   - Posts remain accessible

2. **Invalid Post ID**:
   ```json
   {
     "error": "invalid post id"
   }
   ```

3. **Database Errors**:
   ```json
   {
     "error": "internal server error"
   }
   ```

## 📊 Performance Considerations

### Redis Performance
- **Memory Usage**: ~50 bytes per tracked view
- **Key Expiration**: Automatic cleanup after TTL
- **Scalability**: Redis can handle millions of keys

### Database Impact
- **Minimal**: Only one increment operation per unique view
- **Indexed**: View count field should be indexed for analytics
- **Aggregation**: Total views calculated via MongoDB aggregation

## 🔐 Security Considerations

### IP Spoofing Protection
- System validates IP format
- Handles proxy headers appropriately
- Logs suspicious activity

### Rate Limiting
- 24-hour TTL prevents rapid view inflation
- Per-IP throttling reduces abuse potential
- Graceful handling of invalid requests

## 🎯 Benefits

1. **Accurate Analytics**: Prevents view count inflation
2. **Performance**: Redis provides fast IP lookup
3. **Scalable**: Handles high traffic efficiently
4. **Flexible**: Configurable TTL and settings
5. **Resilient**: Works with or without Redis
6. **Real-time**: Immediate view tracking and analytics

## 🔄 Future Enhancements

Potential improvements for the system:

1. **Geographic Analytics**: Track views by country/region
2. **Time-based Analytics**: View trends over time
3. **User Agent Tracking**: Browser/device analytics
4. **Advanced Throttling**: Different TTL for different user types
5. **View Heatmaps**: Track which parts of posts are viewed
6. **A/B Testing**: Track engagement metrics

---

This Redis-based view tracking system provides accurate, scalable, and abuse-resistant view counting for the blog platform while maintaining excellent performance and user experience.