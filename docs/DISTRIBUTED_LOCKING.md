# Distributed Locking Implementation

## Overview

This implementation provides **configurable distributed locking** for sync operations that works seamlessly for both single and multiple instances. When enabled, it uses Redis for distributed coordination. When disabled, it uses a lightweight no-op implementation that requires no external dependencies.

## Configuration

### **Enable/Disable Distributed Locking**

In your `config.yaml`:

```yaml
locking:
  # Set to true for multi-instance deployments (requires Redis)
  # Set to false for single-instance deployments (no Redis required)
  enabled: false
  
  # Lock TTL configurations
  lock_ttl: "5m"
  incremental_ttl: "5m"
  full_sync_ttl: "30m"
  cleanup_interval: "10m"
```

### **Environment Variables**

You can also configure via environment variables:

```bash
# Enable distributed locking
export FIREWALL_LOCKING_ENABLED=true

# Configure TTL values
export FIREWALL_LOCKING_INCREMENTAL_TTL=5m
export FIREWALL_LOCKING_FULL_SYNC_TTL=30m
```

## Key Features

### 1. **Redis-based Distributed Locking**
- Uses Redis for atomic lock operations
- Automatic TTL (Time-To-Live) prevents deadlocks
- High performance with sub-millisecond operations
- Works with existing Redis infrastructure

### 2. **Instance Identification**
- Unique instance IDs generated using hostname, PID, and timestamp
- Format: `hostname-pid-timestamp`
- Enables tracking which instance holds each lock

### 3. **Lock Operations**

#### **Acquire Lock**
```go
acquired, lockInfo := distributedLock.TryAcquireLock("lock_name", 5*time.Minute)
if acquired {
    // Lock acquired successfully
    fmt.Printf("Lock acquired by instance: %s\n", lockInfo.Instance)
}
```

#### **Release Lock**
```go
released := distributedLock.ReleaseLock("lock_name")
if released {
    // Lock released successfully
}
```

#### **Check Lock Status**
```go
isLocked := distributedLock.IsLocked("lock_name")
lockInfo, err := distributedLock.GetLockInfo("lock_name")
```

### 4. **Sync Operation Integration**

#### **Incremental Sync**
- Uses lock name: `"incremental_sync"`
- TTL: 5 minutes
- Prevents multiple instances from running incremental sync simultaneously

#### **Full Sync**
- Uses lock name: `"full_sync"`
- TTL: 30 minutes (longer due to full sync duration)
- Prevents multiple instances from running full sync simultaneously

### 5. **Safety Features**

#### **Atomic Operations**
- Uses Redis Lua scripts for atomic check-and-delete operations
- Prevents race conditions during lock release

#### **Automatic Expiration**
- All locks have TTL to prevent deadlocks
- Automatic cleanup of expired locks

#### **Lock Extension**
- Allows extending lock TTL for long-running operations
- Only works if the instance owns the lock

### 6. **Monitoring & Debugging**

#### **API Endpoint**
```
GET /sync/status
```

Returns:
```json
{
  "full_sync_running": false,
  "distributed_locking_enabled": true,
  "locks": {
    "full_sync": {
      "lock_id": "full_sync",
      "instance": "hostname-12345-1234567890",
      "acquired": "2024-01-01T12:00:00Z",
      "expires_at": "2024-01-01T12:30:00Z"
    },
    "incremental_sync": null,
    "active_locks": [...]
  }
}
```

**When distributed locking is disabled:**
```json
{
  "full_sync_running": false,
  "distributed_locking_enabled": false,
  "locks": {
    "full_sync": null,
    "incremental_sync": null,
    "active_locks": []
  }
}
```

#### **Lock Information**
- Instance that acquired the lock
- When the lock was acquired
- When the lock expires
- All active locks in the system

## Implementation Details

### **Lock Key Format**
```
lock:lock_name
```

### **Lock Value Format**
```
instance_id:timestamp
```

### **Redis Commands Used**
- `SET key value NX EX seconds` - Atomic lock acquisition
- `EVAL script keys args` - Atomic lock release/extension
- `TTL key` - Check remaining time
- `KEYS pattern` - List active locks

### **Lua Scripts**

#### **Lock Release**
```lua
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end
```

#### **Lock Extension**
```lua
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("expire", KEYS[1], ARGV[2])
else
    return 0
end
```

## Deployment Scenarios

### **Single Instance Deployment**
```yaml
locking:
  enabled: false  # No Redis required
```

**Benefits:**
- No external dependencies (Redis not required)
- Zero overhead for lock operations
- Simple deployment and maintenance
- Perfect for development and small-scale production

**Use Cases:**
- Development environments
- Small production deployments
- Testing and staging environments
- Single-server setups

### **Multi-Instance Deployment**
```yaml
locking:
  enabled: true   # Requires Redis
```

**Benefits:**
- Prevents duplicate sync operations across instances
- Automatic conflict resolution
- Horizontal scaling capability
- Instance identification and monitoring

**Use Cases:**
- High-availability deployments
- Load-balanced environments
- Kubernetes deployments
- Microservices architectures

### **Migration Path**
You can easily migrate from single to multi-instance:

1. **Start with single instance:**
   ```yaml
   locking:
     enabled: false
   ```

2. **When scaling is needed, enable distributed locking:**
   ```yaml
   locking:
     enabled: true
   redis:
     host: "your-redis-server"
   ```

3. **Deploy multiple instances** - they'll automatically coordinate

## Configuration

### **Redis Requirements**
- Redis server accessible to all instances
- Same Redis instance for all application instances
- Recommended: Redis cluster for high availability

### **Lock TTLs**
- **Incremental Sync**: 5 minutes
- **Full Sync**: 30 minutes
- **Custom Locks**: Configurable per use case

## Benefits

1. **Safety**: Prevents duplicate sync operations
2. **Scalability**: Works with any number of instances
3. **Reliability**: Automatic expiration prevents deadlocks
4. **Observability**: Full lock status monitoring
5. **Performance**: Minimal overhead, fast operations
6. **Simplicity**: Transparent integration with existing code

## Testing

To test the distributed locking:

1. Start multiple application instances
2. Trigger sync operations from different instances
3. Monitor lock status via API endpoint
4. Verify only one instance runs sync at a time

## Monitoring

### **Key Metrics**
- Lock acquisition success/failure rates
- Lock duration and expiration patterns
- Instance distribution of lock ownership
- Cleanup operations frequency

### **Alerts**
- Lock acquisition failures
- Expired locks not cleaned up
- Multiple instances attempting same operation
- Redis connectivity issues

## Future Enhancements

1. **Lock Metrics**: Prometheus metrics for lock operations
2. **Lock Queuing**: Queue for lock acquisition attempts
3. **Lock Notifications**: Webhook notifications for lock events
4. **Lock History**: Audit trail of lock operations
5. **Dynamic TTL**: Adaptive TTL based on operation duration 