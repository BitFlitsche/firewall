// Script to populate cache with test data
const axios = require('axios');

async function populateCache() {
  console.log('=== Populating Cache with Test Data ===\n');

  try {
    // 1. Add some IP addresses to trigger cache population
    console.log('1. Adding IP addresses...');
    const ips = [
      { address: '192.168.1.1', status: 'denied' },
      { address: '192.168.1.2', status: 'whitelisted' },
      { address: '10.0.0.1', status: 'denied' },
      { address: '172.16.0.1', status: 'allowed' }
    ];

    for (const ip of ips) {
      await axios.post('http://localhost:8081/ip', ip);
      console.log(`   Added IP: ${ip.address} (${ip.status})`);
    }

    // 2. Add some email addresses
    console.log('\n2. Adding email addresses...');
    const emails = [
      { address: 'spam@example.com', status: 'denied' },
      { address: 'admin@company.com', status: 'whitelisted' },
      { address: 'user@test.org', status: 'allowed' }
    ];

    for (const email of emails) {
      await axios.post('http://localhost:8081/email', email);
      console.log(`   Added Email: ${email.address} (${email.status})`);
    }

    // 3. Add some user agents
    console.log('\n3. Adding user agents...');
    const userAgents = [
      { user_agent: 'Mozilla/5.0 (compatible; Bot)', status: 'denied' },
      { user_agent: 'Mozilla/5.0 (iPhone; CPU iPhone OS 14_0)', status: 'allowed' },
      { user_agent: 'curl/7.68.0', status: 'denied' }
    ];

    for (const ua of userAgents) {
      await axios.post('http://localhost:8081/user-agent', ua);
      console.log(`   Added User Agent: ${ua.user_agent.substring(0, 30)}... (${ua.status})`);
    }

    // 4. Add some countries
    console.log('\n4. Adding countries...');
    const countries = [
      { code: 'US', name: 'United States', status: 'allowed' },
      { code: 'CN', name: 'China', status: 'denied' },
      { code: 'DE', name: 'Germany', status: 'allowed' }
    ];

    for (const country of countries) {
      await axios.post('http://localhost:8081/country', country);
      console.log(`   Added Country: ${country.code} (${country.name}) - ${country.status}`);
    }

    // 5. Wait for cache to be populated
    console.log('\n5. Waiting for cache to be populated...');
    await new Promise(resolve => setTimeout(resolve, 2000));

    // 6. Check cache statistics
    console.log('\n6. Checking cache statistics...');
    const stats = await axios.get('http://localhost:8081/system-stats');
    console.log('Cache Statistics:', stats.data.cache_stats);

    // 7. Trigger some cache operations
    console.log('\n7. Triggering cache operations...');
    
    // Get IPs list (should populate cache)
    await axios.get('http://localhost:8081/ips?page=1&limit=10');
    console.log('   - Retrieved IPs list');
    
    // Get emails list (should populate cache)
    await axios.get('http://localhost:8081/emails?page=1&limit=10');
    console.log('   - Retrieved emails list');
    
    // Get user agents list (should populate cache)
    await axios.get('http://localhost:8081/user-agents?page=1&limit=10');
    console.log('   - Retrieved user agents list');
    
    // Get countries list (should populate cache)
    await axios.get('http://localhost:8081/countries?page=1&limit=10');
    console.log('   - Retrieved countries list');

    // 8. Final cache statistics
    console.log('\n8. Final cache statistics...');
    await new Promise(resolve => setTimeout(resolve, 1000));
    const finalStats = await axios.get('http://localhost:8081/system-stats');
    console.log('Final Cache Statistics:', finalStats.data.cache_stats);

    console.log('\n=== Cache Population Completed! ===');
    console.log('\nYou should now see cache statistics in the System Health page:');
    console.log('- Items: Should show > 0');
    console.log('- Valid: Should show > 0');
    console.log('- Memory: Should show memory usage');

  } catch (error) {
    console.error('Error populating cache:', error.message);
    if (error.response) {
      console.error('Response data:', error.response.data);
    }
  }
}

populateCache(); 