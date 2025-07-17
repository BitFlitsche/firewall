// Test script to demonstrate the new connection pooling display
// This simulates the data structure that the backend now returns

const mockStats = {
  uptime: 3600, // 1 hour
  cpu_percent: [15.5],
  memory_used: 1073741824, // 1GB
  memory_total: 2147483648, // 2GB
  memory_percent: 50,
  disk_used: 53687091200, // 50GB
  disk_total: 107374182400, // 100GB
  disk_percent: 50,
  db_health: "ok",
  db_connections: {
    max_open_connections: 25,
    open_connections: 8,
    in_use: 3,
    idle: 5,
    wait_count: 0,
    wait_duration: "0s",
    max_idle_closed: 12,
    max_lifetime_closed: 5
  },
  es_health: "green",
  request_count: 1250,
  error_count: 0,
  go_routines: 45,
  pid: 12345
};

console.log("=== Connection Pooling Display Test ===");
console.log("Backend now returns detailed connection statistics:");
console.log(JSON.stringify(mockStats.db_connections, null, 2));

console.log("\nFrontend will now display:");
console.log(`DB Conns: ${mockStats.db_connections.open_connections}/${mockStats.db_connections.max_open_connections}`);
console.log(`Active: ${mockStats.db_connections.in_use}`);

console.log("\nTooltip information:");
console.log(`Open: ${mockStats.db_connections.open_connections}, In Use: ${mockStats.db_connections.in_use}, Idle: ${mockStats.db_connections.idle}`);

console.log("\nThis provides much more detailed information than the previous simple connection count!"); 