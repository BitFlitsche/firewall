// Test script to demonstrate cache statistics in system stats
const axios = require('axios');

async function testCacheStats() {
  console.log('=== Testing Cache Statistics in System Stats ===\n');

  try {
    // Get initial system stats
    console.log('1. Getting initial system stats...');
    const initialStats = await axios.get('http://localhost:8081/system-stats');
    console.log('Initial cache stats:', initialStats.data.cache_stats);

    // Add some test data to cache
    console.log('\n2. Adding test data to cache...');
    await axios.post('http://localhost:8081/ip', {
      address: '192.168.1.100',
      status: 'denied'
    });

    // Wait a moment for cache to be populated
    await new Promise(resolve => setTimeout(resolve, 1000));

    // Get updated system stats
    console.log('\n3. Getting updated system stats...');
    const updatedStats = await axios.get('http://localhost:8081/system-stats');
    console.log('Updated cache stats:', updatedStats.data.cache_stats);

    // Update the IP to trigger cache invalidation
    console.log('\n4. Updating IP to trigger cache invalidation...');
    const ips = await axios.get('http://localhost:8081/ips?limit=1');
    if (ips.data.data && ips.data.data.length > 0) {
      const ipId = ips.data.data[0].id;
      await axios.put(`http://localhost:8081/ip/${ipId}`, {
        address: '192.168.1.100',
        status: 'whitelisted'
      });
    }

    // Wait a moment for cache invalidation
    await new Promise(resolve => setTimeout(resolve, 1000));

    // Get final system stats
    console.log('\n5. Getting final system stats...');
    const finalStats = await axios.get('http://localhost:8081/system-stats');
    console.log('Final cache stats:', finalStats.data.cache_stats);

    console.log('\n=== Cache Statistics Test Completed! ===');
    console.log('\nYou should now see cache statistics in the System Health page:');
    console.log('- Total Items: Number of cached items');
    console.log('- Valid Items: Number of non-expired items');
    console.log('- Memory Usage: Estimated memory usage');

  } catch (error) {
    console.error('Error testing cache stats:', error.message);
  }
}

testCacheStats(); 