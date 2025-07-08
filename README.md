# Firewall Application

A full-stack firewall application that provides IP, email, user agent, and country-based filtering capabilities. The application consists of a Go backend API and a React frontend.

## Features

- IP address filtering
- Email address filtering
- User agent filtering
- Country-based filtering
- Real-time request filtering
- Elasticsearch integration for search capabilities
- MariaDB for data persistence
- Redis for caching
- Kibana for data visualization

## Prerequisites

- Go 1.24 or higher
- Node.js 14 or higher
- Docker and Docker Compose
- Git

## Project Structure

```
firewall/
├── config/         # Configuration files
├── controllers/    # API controllers
├── firewall-app/   # React frontend
├── geoip/         # GeoIP integration
├── middleware/    # HTTP middleware
├── migrations/    # Database migrations
├── models/        # Data models
├── routes/        # API routes
├── services/      # Business logic
├── docker-compose.yml
├── go.mod
├── main.go
└── README.md
```

## Setup and Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd firewall
```

2. Start the backend services using Docker Compose:
```bash
docker-compose up -d
```
This will start:
- MariaDB (port 3306)
- Redis (port 6379)
- Elasticsearch (port 9200)
- Kibana (port 5601)

3. Start the Go backend server:
```bash
go run main.go
```
The backend server will start on port 8081.

4. Start the React frontend:
```bash
cd firewall-app
npm install
npm start
```
The frontend will be available at http://localhost:3000

## API Endpoints

### IP Management
- `POST /ip` - Create new IP address
- `GET /ips` - List all IP addresses

### Email Management
- `POST /email` - Create new email
- `GET /emails` - List all emails

### User Agent Management
- `POST /user-agent` - Create new user agent
- `GET /user-agents` - List all user agents

### Country Management
- `POST /country` - Create new country
- `GET /countries` - List all countries

### Filtering
- `POST /filter` - Filter requests based on multiple criteria

## API-Dokumentation (Swagger)

Dieses Projekt nutzt [Swaggo](https://github.com/swaggo/swag) zur automatischen Generierung einer Swagger/OpenAPI-Dokumentation.

### Beispiel für Handler-Kommentare
```go
// @Summary      Filtert IP, E-Mail, User-Agent und Land
// @Description  Prüft, ob die angegebenen Werte erlaubt oder blockiert sind
// @Tags         filter
// @Accept       json
// @Produce      json
// @Param        filter  body      FilterRequest  true  "Filterdaten"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]string
// @Failure      504     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /filter [post]
```

### Swagger-Dokumentation generieren

1. Installiere swag (falls noch nicht geschehen):
   ```sh
   go install github.com/swaggo/swag/cmd/swag@latest
   ```
2. Generiere die Swagger-Dokumentation:
   ```sh
   swag init
   ```
   Dadurch wird das Verzeichnis `docs/` mit der OpenAPI-Dokumentation erstellt.

3. Starte das Backend und rufe die Swagger-UI im Browser auf:
   ```
   http://localhost:8081/swagger/index.html
   ```

Weitere Infos: [Swaggo Doku](https://github.com/swaggo/swag)

## Usage

1. Access the frontend at http://localhost:3000
2. Use the interface to:
   - Add new IP addresses, emails, user agents, or countries
   - View existing entries
   - Filter requests based on multiple criteria

## Development

### Backend Development
- The backend is written in Go using the Gin framework
- Database migrations are handled through GORM
- CORS is configured to allow frontend access

### Frontend Development
- Built with React
- Uses Axios for API communication
- Development server with hot reloading

## Stopping the Application

1. Stop the frontend:
   - Press `Ctrl+C` in the frontend terminal

2. Stop the backend:
   - Press `Ctrl+C` in the backend terminal

3. Stop Docker containers:
```bash
docker-compose down
```

## Troubleshooting

### Common Issues

1. Port conflicts:
   - If port 8081 is already in use, you can change it in `main.go`
   - If port 3000 is in use, React will automatically suggest an alternative port

2. Database connection issues:
   - Ensure Docker containers are running
   - Check database credentials in configuration

3. CORS issues:
   - Verify CORS configuration in `routes/routes.go`
   - Check browser console for specific error messages

## Security Considerations

- The application is configured for development mode
- For production:
  - Enable Gin release mode
  - Configure proper CORS settings
  - Set up proper authentication
  - Configure secure proxy settings

## Lizenz

Dieses Projekt ist Open Source und steht unter der MIT License.  
Copyright (c) 2024 github.com/BitFlitsche

Drittanbieter-Bibliotheken, wie z.B. der ElasticSearch Go Client, können unter anderen Lizenzen (z.B. Apache 2.0) stehen. Die entsprechenden Lizenztexte sind diesem Projekt beigelegt.

## Contributing

Beiträge sind willkommen!  
Wenn du Fehler findest, Ideen hast oder neue Features beitragen möchtest, erstelle gerne ein Issue oder einen Pull Request.

Vielen Dank für deine Unterstützung! 