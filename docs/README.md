# Documentation Index

Welcome to the Firewall Application documentation. This index provides an overview of all available documentation.

## Core Documentation

### [API Reference](API_REFERENCE.md)
Complete API documentation including:
- Filter endpoint with geographic filtering
- Request/response formats
- Usage examples (cURL, JavaScript)
- Error handling
- Performance considerations

### [Installation Guide](INSTALLATION.md)
Step-by-step installation instructions:
- Prerequisites and dependencies
- Database setup
- GeoLite2 database download
- Configuration
- Production deployment
- Troubleshooting

### [Geographic Filtering](GEOGRAPHIC_FILTERING.md)
Comprehensive guide to geographic filtering:
- How automatic geolocation works
- Manual country override
- Private IP handling
- Implementation details
- Performance considerations
- Troubleshooting

## Feature Documentation

### [Validation](VALIDATION.md)
Input validation rules and security checks:
- IP address validation
- Email format validation
- Country code validation
- User agent validation
- Regex pattern validation

### [Conflict Detection](CONFLICT_DETECTION.md)
IP/CIDR conflict detection system:
- IP in CIDR detection
- CIDR overlap detection
- Conflict resolution
- API endpoints
- UI integration

### [Health Check](HEALTH_CHECK.md)
Monitoring and health check system:
- Health check endpoints
- Service monitoring
- Performance metrics
- Alerting integration

## Configuration Documentation

### [Configuration Guide](CONFIGURATION.md)
Detailed configuration options:
- Server settings
- Database configuration
- Elasticsearch setup
- Redis configuration
- Security settings

### [Environment Variables](ENVIRONMENT_VARIABLES.md)
Environment variable configuration:
- Database variables
- Elasticsearch variables
- Redis variables
- Security variables

## Development Documentation

### [Development Guide](DEVELOPMENT.md)
Development setup and guidelines:
- Local development environment
- Code structure
- Testing procedures
- Contribution guidelines

### [Architecture](ARCHITECTURE.md)
System architecture overview:
- Component diagram
- Data flow
- Service interactions
- Scalability considerations

## Deployment Documentation

### [Production Deployment](PRODUCTION.md)
Production deployment guide:
- System requirements
- Load balancing
- SSL/TLS setup
- Monitoring and alerting
- Backup strategies

### [Docker Deployment](DOCKER.md)
Docker-based deployment:
- Dockerfile
- Docker Compose
- Container orchestration
- Volume management

## Monitoring and Operations

### [Monitoring](MONITORING.md)
Monitoring and observability:
- Metrics collection
- Log aggregation
- Alerting rules
- Dashboard setup

### [Troubleshooting](TROUBLESHOOTING.md)
Common issues and solutions:
- Database issues
- Elasticsearch problems
- Geographic filtering issues
- Performance problems

## Security Documentation

### [Security Guide](SECURITY.md)
Security best practices:
- Authentication
- Authorization
- Input validation
- Network security
- Data protection

### [Compliance](COMPLIANCE.md)
Compliance considerations:
- GDPR compliance
- Data retention
- Privacy protection
- Audit logging

## API Documentation

### [Swagger UI](http://localhost:8081/swagger/index.html)
Interactive API documentation available when the application is running.

### [Postman Collection](postman_collection.json)
Postman collection for API testing (if available).

## Quick Reference

### Common Commands

```bash
# Start the application
go run main.go

# Build for production
go build -o firewall .

# Test geographic filtering
curl -X POST http://localhost:8081/api/filter \
  -H "Content-Type: application/json" \
  -d '{"ip": "8.8.8.8", "country": ""}'

# Check health
curl http://localhost:8081/api/health
```

### Key Files

- `config.yaml` - Main configuration file
- `GeoLite2-Country.mmdb` - Geographic database
- `main.go` - Application entry point
- `services/geolocation.go` - Geographic filtering service
- `services/filter_service.go` - Filter evaluation logic

### Important URLs

- **Application**: http://localhost:8081
- **API**: http://localhost:8081/api
- **Swagger UI**: http://localhost:8081/swagger/index.html
- **Health Check**: http://localhost:8081/api/health

## Getting Help

### Documentation Issues
If you find issues with the documentation:
1. Check if the information is outdated
2. Verify against the current codebase
3. Open an issue on GitHub

### Application Issues
For application problems:
1. Check the troubleshooting guide
2. Review the logs
3. Test with the health check endpoint
4. Open an issue on GitHub

### Feature Requests
For new features:
1. Check existing documentation
2. Review the architecture guide
3. Open a feature request on GitHub

## Contributing to Documentation

### Documentation Standards
- Use clear, concise language
- Include code examples
- Provide troubleshooting steps
- Keep information up-to-date

### Adding New Documentation
1. Create the new document in the `docs/` directory
2. Update this index
3. Follow the existing format and style
4. Include links to related documentation

### Updating Documentation
1. Verify information is current
2. Test all code examples
3. Update related documents
4. Review for clarity and completeness

## Version Information

This documentation covers:
- **Application Version**: 1.0.0
- **Go Version**: 1.24+
- **Node.js Version**: 18+
- **Last Updated**: 2025-07-19

## License

This documentation is licensed under the same terms as the application (MIT License). 