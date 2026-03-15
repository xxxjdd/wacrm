# WACRM - WhatsApp Customer Relationship Management

## Product Overview

- **Product Name**: WACRM
- **Type**: Windows Desktop CRM Application
- **Core Functionality**: Multi-account WhatsApp customer management with automated messaging
- **Target Users**: Sales teams, customer support, marketing teams in Africa (Nigeria, Ghana, Cameroon, India)

## Technical Stack

| Component | Technology |
|----------|------------|
| Backend API | Go (Gin) |
| Database | MySQL |
| Desktop Client | Tauri (Rust) |
| Frontend | React + TypeScript |
| WhatsApp Protocol | Web Socket Direct Connection |
| Auto-update | Tauri Updater |

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        api.dgxs.cn                              │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐          │
│  │  Auth   │  │ Users   │  │Accounts │  │Customers│          │
│  └─────────┘  └─────────┘  └─────────┘  └─────────┘          │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐                       │
│  │ Messages│  │  Tasks  │  │ Stats   │                       │
│  └─────────┘  └─────────┘  └─────────┘                       │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Windows Client (Tauri)                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐       │
│  │  Login   │  │Dashboard │  │Accounts  │  │Customers │       │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘       │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐       │
│  │ Messages  │  │ Tasks    │  │ Stats    │  │ Settings │       │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘       │
└─────────────────────────────────────────────────────────────────┘
```

## Feature List

### Phase 1: Authentication & Core
- [x] User registration (cloud)
- [x] User login (cloud)
- [x] Auto-update
- [x] Local WhatsApp session management

### Phase 2: WhatsApp Connection
- [x] QR code scanning for WhatsApp login
- [x] Multi-account support (switch between accounts)
- [x] Real-time message receiving
- [x] Message sending
- [x] Contact/chat list

### Phase 3: Customer Management
- [x] Customer profile (from WhatsApp contacts)
- [x] Customer tags/labels
- [x] Customer notes
- [x] Customer search
- [x] Import customers

### Phase 4: Automation & Stats
- [x] Scheduled message sending
- [x] Auto-reply rules
- [x] Message templates
- [x] Conversion statistics
- [x] Daily/weekly reports

## Database Schema

### users
| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT | Primary key |
| username | VARCHAR(100) | Login username |
| password_hash | VARCHAR(255) | Bcrypt hash |
| email | VARCHAR(255) | Email |
| role | ENUM | admin, user |
| created_at | DATETIME | Creation time |
| updated_at | DATETIME | Update time |

### whatsapp_accounts
| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT | Primary key |
| user_id | BIGINT | Owner user |
| phone | VARCHAR(20) | WhatsApp phone |
| session_data | TEXT | Encrypted session |
| nickname | VARCHAR(100) | Display name |
| status | ENUM | online, offline, connecting |
| created_at | DATETIME | Creation time |

### customers
| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT | Primary key |
| account_id | BIGINT | WhatsApp account |
| phone | VARCHAR(20) | Customer phone |
| name | VARCHAR(100) | Customer name |
| avatar | VARCHAR(500) | Profile picture |
| tags | JSON | Tags array |
| notes | TEXT | Notes |
| created_at | DATETIME | Creation time |
| updated_at | DATETIME | Update time |

### messages
| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT | Primary key |
| account_id | BIGINT | WhatsApp account |
| customer_id | BIGINT | Customer |
| direction | ENUM | inbound, outbound |
| content | TEXT | Message content |
| message_type | ENUM | text, image, video, etc |
| sent_at | DATETIME | Send time |
| created_at | DATETIME | Creation time |

### message_templates
| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT | Primary key |
| user_id | BIGINT | Owner |
| name | VARCHAR(100) | Template name |
| content | TEXT | Template content |
| variables | JSON | Variables like {{name}} |
| created_at | DATETIME | Creation time |

### scheduled_tasks
| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT | Primary key |
| user_id | BIGINT | Owner |
| account_id | BIGINT | WhatsApp account |
| name | VARCHAR(100) | Task name |
| template_id | BIGINT | Message template |
| customer_ids | JSON | Target customers |
| scheduled_at | DATETIME | When to send |
| status | ENUM | pending, running, completed, failed |
| created_at | DATETIME | Creation time |

## API Endpoints

### Auth
- POST /api/auth/register - Register
- POST /api/auth/login - Login
- POST /api/auth/logout - Logout
- GET /api/auth/me - Current user

### Users
- GET /api/users - List users (admin)
- PUT /api/users/:id - Update user

### WhatsApp Accounts
- GET /api/accounts - List accounts
- POST /api/accounts - Add account
- DELETE /api/accounts/:id - Remove account
- POST /api/accounts/:id/logout - Logout account

### Customers
- GET /api/customers - List customers
- POST /api/customers - Create customer
- PUT /api/customers/:id - Update customer
- DELETE /api/customers/:id - Delete customer
- POST /api/customers/import - Import customers

### Messages
- GET /api/messages - List messages
- POST /api/messages/send - Send message

### Templates
- GET /api/templates - List templates
- POST /api/templates - Create template
- PUT /api/templates/:id - Update template
- DELETE /api/templates/:id - Delete template

### Scheduled Tasks
- GET /api/tasks - List tasks
- POST /api/tasks - Create task
- PUT /api/tasks/:id - Update task
- DELETE /api/tasks/:id - Delete task
- POST /api/tasks/:id/run - Run task now

### Stats
- GET /api/stats/overview - Overview stats
- GET /api/stats/messages - Message stats
- GET /api/stats/conversion - Conversion stats

## UI/UX Requirements

### Design System
- Color scheme: Blue primary (#2563EB), White background
- Dark mode support
- Responsive layout
- Native Windows look

### Performance
- Memory usage < 50MB
- Startup time < 3 seconds
- Smooth scrolling
- Lazy loading for lists

### Key Screens
1. Login/Register
2. Dashboard (stats overview)
3. Account management (list, add, QR scan)
4. Chat/Messages (conversations)
5. Customer list (with search, filter, tags)
6. Customer detail
7. Message templates
8. Scheduled tasks
9. Settings
