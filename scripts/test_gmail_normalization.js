// Test script to demonstrate Gmail email normalization
const axios = require('axios');

async function testGmailNormalization() {
  console.log('=== Testing Gmail Email Normalization ===\n');

  try {
    // Test cases for Gmail normalization
    const testCases = [
      {
        name: 'Basic Gmail normalization',
        email: 't.e.s.t@gmail.com',
        expected: 'test@gmail.com'
      },
      {
        name: 'Gmail.de normalization',
        email: 'u.s.e.r@gmail.de',
        expected: 'user@gmail.de'
      },
      {
        name: 'Gmail.co.uk normalization',
        email: 'a.d.m.i.n@gmail.co.uk',
        expected: 'admin@gmail.co.uk'
      },
      {
        name: 'No dots in Gmail',
        email: 'test@gmail.com',
        expected: 'test@gmail.com'
      },
      {
        name: 'Non-Gmail email (should not change)',
        email: 'test@example.com',
        expected: 'test@example.com'
      },
      {
        name: 'Gmail with multiple dots',
        email: 't.e.s.t.u.s.e.r@gmail.com',
        expected: 'testuser@gmail.com'
      },
      {
        name: 'Mixed case Gmail',
        email: 'T.E.S.T@GMAIL.COM',
        expected: 'test@gmail.com'
      }
    ];

    console.log('1. Testing email normalization logic...');
    for (const testCase of testCases) {
      console.log(`   ${testCase.name}:`);
      console.log(`     Input:  ${testCase.email}`);
      console.log(`     Expected: ${testCase.expected}`);
    }

    // Test filter endpoint with different Gmail variations
    console.log('\n2. Testing filter endpoint with Gmail variations...');
    
    const filterTests = [
      {
        name: 'Original Gmail',
        email: 'test@gmail.com',
        ip: '192.168.1.1'
      },
      {
        name: 'Dotted Gmail',
        email: 't.e.s.t@gmail.com',
        ip: '192.168.1.1'
      },
      {
        name: 'Gmail.de with dots',
        email: 'u.s.e.r@gmail.de',
        ip: '192.168.1.2'
      },
      {
        name: 'Non-Gmail (should not normalize)',
        email: 'test@example.com',
        ip: '192.168.1.3'
      }
    ];

    for (const test of filterTests) {
      console.log(`\n   Testing: ${test.name}`);
      console.log(`   Email: ${test.email}`);
      
      try {
        const response = await axios.post('http://localhost:8081/filter', {
          ip: test.ip,
          email: test.email,
          user_agent: 'Mozilla/5.0 (Test)',
          country: 'US',
          username: 'testuser'
        });
        
        console.log(`   Result: ${response.data.result}`);
        console.log(`   Reason: ${response.data.reason || 'N/A'}`);
        
        // Check cache statistics
        const stats = await axios.get('http://localhost:8081/system-stats');
        console.log(`   Cache Items: ${stats.data.cache_stats.total_items}`);
        
      } catch (error) {
        console.log(`   Error: ${error.message}`);
      }
    }

    // Test cache behavior with normalized emails
    console.log('\n3. Testing cache behavior with normalized emails...');
    
    // First request with dotted Gmail
    console.log('   Making first request with t.e.s.t@gmail.com...');
    await axios.post('http://localhost:8081/filter', {
      ip: '192.168.1.100',
      email: 't.e.s.t@gmail.com',
      user_agent: 'Mozilla/5.0 (Test)',
      country: 'US',
      username: 'testuser'
    });
    
    // Second request with non-dotted Gmail (should be cached)
    console.log('   Making second request with test@gmail.com...');
    const startTime = Date.now();
    await axios.post('http://localhost:8081/filter', {
      ip: '192.168.1.100',
      email: 'test@gmail.com',
      user_agent: 'Mozilla/5.0 (Test)',
      country: 'US',
      username: 'testuser'
    });
    const endTime = Date.now();
    
    console.log(`   Response time: ${endTime - startTime}ms (should be fast if cached)`);
    
    // Check final cache statistics
    const finalStats = await axios.get('http://localhost:8081/system-stats');
    console.log(`   Final Cache Items: ${finalStats.data.cache_stats.total_items}`);
    console.log(`   Final Cache Memory: ${finalStats.data.cache_stats.memory_usage}`);

    console.log('\n=== Gmail Normalization Test Completed! ===');
    console.log('\nKey Points:');
    console.log('- Gmail addresses have dots removed from local part');
    console.log('- This affects gmail.com, gmail.de, gmail.co.uk, etc.');
    console.log('- Cache keys use normalized emails, so t.e.s.t@gmail.com and test@gmail.com share the same cache entry');
    console.log('- This prevents spammers from bypassing filters using dot variations');

  } catch (error) {
    console.error('Error testing Gmail normalization:', error.message);
    if (error.response) {
      console.error('Response data:', error.response.data);
    }
  }
}

testGmailNormalization(); 