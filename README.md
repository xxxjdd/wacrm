# WACRM - WhatsApp Customer Relationship Management

A complete WhatsApp CRM system with cloud backend and Windows desktop client.

## Features

- ✅ Multi-user system with cloud authentication
- ✅ Multi-WhatsApp account management (multiple phone numbers)
- ✅ QR code login for each WhatsApp account
- ✅ Customer management with tags and notes
- ✅ Real-time messaging with conversations
- ✅ Message templates for quick replies
- ✅ Scheduled tasks for automated messaging
- ✅ Statistics and analytics dashboard
- ✅ Auto-update mechanism for desktop client
- ✅ Memory-efficient (using Tauri instead of Chrome)

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  Domain: api.dgxs.cn                                         │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                  │
│  │  Nginx   │  │  Go API  │  │  MySQL   │                  │
│  │  (SSL)   │─▶│  Server  │─▶│  (Data)  │                  │
│  └──────────┘  └──────────┘  └──────────┘                  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼ HTTP/WebSocket
┌─────────────────────────────────────────────────────────────┐
│  Windows Desktop Client (Tauri + React)                     │
│  ┌────────────────────────────────────────────────────┐    │
│  │  - Login with cloud account                        │    │
│  │  - Add WhatsApp accounts (QR scan)                │    │
│  │  - Manage customers                                │    │
│  │  - Send/receive messages                           │    │
│  │  - Scheduled tasks                                 │    │
│  └────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

## Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: MySQL 8.0
- **Auth**: JWT with session storage
- **API**: RESTful

### Frontend (Desktop)
- **Framework**: Tauri (Rust)
- **UI**: React 18 + TypeScript
- **Styling**: Tailwind CSS
- **State**: Zustand
- **HTTP Client**: Axios

## Project Structure

```
wacrm/
├── api/                    # Go backend
│   ├── config/            # Database & config
│   ├── handlers/          # API handlers
│   ├── middleware/        # Auth middleware
│   ├── models/            # Database models
│   ├── main.go            # Entry point
│   ├── go.mod             # Go dependencies
│   └── Dockerfile         # API container
├── client/                 # Tauri desktop app
│   ├── src/
│   │   ├── api/          # API clients
│   │   ├── components/   # React components
│   │   ├── pages/        # Page components
│   │   ├── store/        # State management
│   │   └── App.tsx       # Main app
│   ├── src-tauri/        # Rust/Tauri code
│   └── package.json      # Node dependencies
├── docs/
│   └── SPEC.md           # Product specification
├── docker-compose.yml    # Docker compose config
├── nginx.conf           # Nginx reverse proxy
└── README.md            # This file
```

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for local development)
- Node.js 18+ & Rust (for client development)
- Domain: api.dgxs.cn (configured with DNS)

### 1. Deploy Backend

```bash
# Clone repository
cd wacrm

# Start with Docker Compose
docker-compose up -d

# Check logs
docker-compose logs -f api
```

The API will be available at: `https://api.dgxs.cn`

### 2. Build Desktop Client

```bash
cd client

# Install dependencies
npm install

# Install Tauri CLI
npm install -g @tauri-apps/cli

# Development mode
npm run tauri dev

# Build for production
npm run tauri build
```

The installer will be created in `src-tauri/target/release/bundle/`

### 3. Configure Auto-Update

Edit `src-tauri/tauri.conf.json`:

```json
{
  "tauri": {
    "updater": {
      "active": true,
      "endpoints": ["https://api.dgxs.cn/api/updates"],
      "dialog": true,
      "pubkey": "YOUR_PUBLIC_KEY"
    }
  }
}
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| /api/auth/register | POST | User registration |
| /api/auth/login | POST | User login |
| /api/accounts | GET/POST | WhatsApp accounts |
| /api/accounts/:id/qr | GET | Get QR code |
| /api/customers | GET/POST | Customers |
| /api/messages | GET/POST | Messages |
| /api/templates | GET/POST | Templates |
| /api/tasks | GET/POST | Scheduled tasks |
| /api/stats/overview | GET | Statistics |

## Configuration

### Environment Variables

```bash
# Backend
DB_DSN=root:password@tcp(mysql:3306)/wacrm?charset=utf8mb4&parseTime=True&loc=Local
PORT=8080

# Frontend (build-time)
VITE_API_URL=https://api.dgxs.cn
```

### SSL Certificates

Place SSL certificates in:
- `/etc/nginx/ssl/fullchain.pem`
- `/etc/nginx/ssl/privkey.pem`

## Database Schema

The system uses MySQL with the following tables:
- `users` - User accounts
- `whatsapp_accounts` - WhatsApp phone connections
- `customers` - Customer contacts
- `messages` - Message history
- `message_templates` - Quick reply templates
- `scheduled_tasks` - Automated tasks

Auto-migration is enabled - tables are created automatically on first run.

## WhatsApp Protocol

This system connects directly to WhatsApp Web using WebSocket protocol.
- No Chrome browser needed
- Low memory usage (< 50MB)
- Local session storage
- Multiple account support

**Note**: WhatsApp protocol is unofficial and may be subject to restrictions.

## Development

### Backend Development

```bash
cd api
go mod tidy
go run main.go
```

### Client Development

```bash
cd client
npm install
npm run tauri dev
```

## Security

- All API requests require JWT authentication
- Passwords are hashed with bcrypt
- HTTPS enforced on production
- CORS configured for api.dgxs.cn
- Session tokens expire after 30 days

## License

MIT License - See LICENSE file

## Support

For issues and questions, please contact support.
