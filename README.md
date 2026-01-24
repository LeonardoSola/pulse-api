# Pulse API

A WhatsApp API service built with Go, providing RESTful endpoints for managing WhatsApp devices and sending messages.

## ⚠️ Work in Progress

**This API is currently under active development and does not include all planned features.** Some functionality may be incomplete, and the API structure may change in future releases. Contributions and feedback are welcome!

## Features

- 🔐 API Key authentication
- 📱 WhatsApp multi-device management (create, read, update, delete)
- 📨 Send text and image messages via WhatsApp
- 🔄 RabbitMQ integration for message queuing
- 🗄️ PostgreSQL database with migrations
- 🐳 Docker support for easy deployment
- 🔌 WAHA (WhatsApp HTTP API) adapter support

## Roadmap

The following features are planned for future releases:

### Message Features
- ✏️ Message update and delete
- 👁️ Send presence (typing, recording, etc.)
- ✅ Read message receipts
- 📎 Additional message types:
  - Media messages (video, documents)
  - Stickers
  - Location sharing
  - Contact cards
  - Audio messages
  - Polls

### Event Handling
- 🔔 Webhook support for event listening
- 🔌 SSE (Server-Sent Events) or WebSocket for real-time event streaming
- 📡 Additional event handlers for various WhatsApp events

### Contact & Chat Management
- 👥 Contacts management
- 💬 Chat management
- 📜 Chat message history retrieval

### Group Management
- 👨‍👩‍👧‍👦 Group creation and management
- 👤 Group member management
- ⚙️ Group settings and configuration

### Profile Management
- 👤 User profile management
- 🖼️ Profile picture updates
- 📝 Status updates

### API Adapter Support
- 🔌 **WAHA (WhatsApp HTTP API)** - Expand current partial support to include all message types and operations
- 🚀 **Evolution API** - Full adapter implementation for Evolution API integration
- 📡 **WHAPI** - Complete adapter support for WHAPI (WhatsApp API) integration
- 🔄 Unified adapter interface to seamlessly switch between different WhatsApp API providers
- 📦 Support for all adapter operations: messaging, device management, webhooks, and more

## Tech Stack

- **Language**: Go 1.25.2
- **Framework**: Chi Router
- **Database**: PostgreSQL
- **Message Queue**: RabbitMQ
- **WhatsApp Library**: whatsmeow
- **Containerization**: Docker & Docker Compose

## Prerequisites

- Go 1.25.2 or higher
- Docker and Docker Compose
- PostgreSQL 16 (or use Docker Compose)
- RabbitMQ (optional, for message queuing)

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/LeonardoSola/pulse-api.git
cd pulse-api
```

### 2. Environment Configuration

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASS=your_password
DB_NAME=pulse_api
DB_SSL_MODE=disable
# OR use connection string
# DB_CONNECTION_STRING=postgres://user:password@localhost:5432/dbname

# RabbitMQ
RABBITMQ_ENABLED=false
RABBITMQ_HOST=
RABBITMQ_PORT=
RABBITMQ_USER=
RABBITMQ_PASSWORD=
RABBITMQ_VHOST=
RABBITMQ_EXCHANGE=
# OR use connection string
RABBITMQ_URL=amqp://user:password@host:port/

# Server
SERVER_PORT=8080
API_KEY=your_api_key_here  # Generate with: openssl rand -hex 32
```

### 3. Generate API Key

Generate a secure API key:

```bash
openssl rand -hex 32
```

Add the generated key to your `.env` file.

### 4. Run with Docker Compose

The easiest way to get started:

```bash
docker-compose up -d
```

This will start:
- PostgreSQL database
- The API service

### 5. Run Locally (Development)

For local development:

```bash
# Install dependencies
go mod download

# Run migrations (if not using Docker)
# Make sure PostgreSQL is running and configured

# Run the application (if you are using bash)
make run
```

Or manually:

```bash
source .env
air .
```

## API Endpoints

All endpoints require API key authentication via the `X-API-Key` header. (expect /health)
You can find the Postman collection [here](Pulse API.postman_collection.json).

### Health Check

- **GET** `/health` - Check API health status

### Devices

- **GET** `/devices` - List all devices
- **GET** `/devices/{id}` - Get device by ID
- **POST** `/devices` - Create a new device
- **PUT** `/devices/{id}` - Update device
- **DELETE** `/devices/{id}` - Delete device
- **GET** `/devices/{id}/qr` - Get QR code for device authentication

### Messages

- **POST** `/devices/{deviceId}/messages/sendText` - Send a text message
- **POST** `/devices/{deviceId}/messages/sendImage` - Send an image message

### Adapters

- **GET** `/adapters/waha/*` - WAHA adapter endpoints


## Example Usage

### Send a Text Message

```bash
curl -X POST http://localhost:8080/devices/{deviceId}/messages/sendText \
  -H "X-API-Key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Hello, World!",
    "recipient": "5511999999999",
    "recipientType": "individual"
  }'
```

### Create a Device

```bash
curl -X POST http://localhost:8080/devices \
  -H "X-API-Key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Device"
  }'
```

## Project Structure

```
pulse-api/
├── adapters/          # External adapters (HTTP, WhatsApp, Database, Messaging)
├── internal/          # Internal business logic
│   ├── config/       # Configuration management
│   ├── entities/     # Domain entities
│   ├── ports/        # Interface definitions
│   └── services/     # Business logic services
├── migrations/        # Database migrations
├── main.go           # Application entry point
├── docker-compose.yaml
├── Dockerfile
└── Makefile
```

## Development

### Running Migrations

Database migrations are handled automatically on startup using Goose.

### Available Make Commands

```bash
make run    # Run the application in development mode
make help   # Show available commands
```

## Contributing

Contributions are welcome! Since this project is in active development, please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is open source. Please check the LICENSE file for details.

## Disclaimer

This API is provided as-is for development purposes. Use at your own risk. Make sure to comply with WhatsApp's Terms of Service when using this API.

## Support

For issues, questions, or contributions, please open an issue on the GitHub repository.

---

**Note**: This API is a work in progress. Features may be incomplete or subject to change.
