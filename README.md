# 🎓 SDU AI Education Platform - Backend

Go-based REST API backend for an AI-powered education platform for Kazakhstan school students (Grades 1-11).

## ✨ Features

### Authentication & Users
- JWT-based authentication
- User registration with grade level
- Role-based access (student/teacher)

### Quiz System
- Multiple choice questions
- Subject/category organization
- Quiz sessions with timer
- Results with scoring

### Flashcard System
- Card creation and management
- Spaced repetition tracking
- Subject organization

### Gamification (Phase 0)
- Points & XP system
- Level progression (1-20+)
- Daily streak tracking
- Robot AI Trainer with evolution stages
- Leaderboard ranking

### AI Chat Integration
- WebSocket-based real-time chat
- OpenAI-powered responses
- Multilingual support (EN, RU, KK)

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 14+
- Docker (optional)

### Local Development

```bash
# Copy environment file
cp .env.example .env
# Edit .env with your database credentials

# Run migrations
make migrate

# Start server
make run
```

### Docker

```bash
docker compose up -d
```

## 📁 Project Structure

```
backend/
├── cmd/server/           # Main entry point
├── internal/
│   ├── app/             # App initialization & routing
│   ├── config/          # Configuration
│   ├── handler/         # HTTP handlers
│   │   ├── auth/        # Login, register
│   │   ├── card/        # Flashcards
│   │   ├── chat/        # AI chat
│   │   ├── leaderboard/ # Rankings (Phase 1)
│   │   ├── progression/ # Points, XP, streaks (Phase 0)
│   │   ├── quiz/        # Quiz management
│   │   ├── session/     # Chat sessions
│   │   ├── test/        # Test/quiz taking
│   │   └── trainer/    # Robot trainer (Phase 1)
│   ├── lib/            # Shared utilities
│   ├── middleware/     # CORS, auth middleware
│   ├── models/         # Database models
│   ├── repository/      # Data access layer
│   ├── service/        # Business logic
│   └── websocket/      # WebSocket handlers
└── pkg/
    └── utils/          # Shared packages
```

## 🔌 API Endpoints

### Auth
- `POST /auth/register` - User registration
- `POST /auth/login` - Login
- `POST /auth/refresh` - Refresh token

### Quizzes
- `GET /quizzes` - List quizzes
- `GET /quizzes/:id` - Get quiz details
- `POST /quizzes/:id/start` - Start quiz session
- `POST /quizzes/:id/result` - Submit quiz result

### Flashcards
- `GET /cards` - List cards
- `POST /cards` - Create card
- `GET /cards/:id` - Get card
- `PUT /cards/:id` - Update card
- `DELETE /cards/:id` - Delete card

### Gamification (Phase 0)
- `GET /progression` - Get user progress
- `POST /progression/streak/claim` - Claim daily bonus
- `GET /progression/transactions` - Point history
- `GET /rewards` - Available rewards
- `POST /rewards/:id/redeem` - Redeem reward

### AI Trainer (Phase 1)
- `GET /trainer/profile` - Robot profile
- `PUT /trainer/profile` - Update robot name/color
- `GET /trainer/timeline` - Level-up history

### Leaderboard (Phase 1)
- `GET /leaderboard` - Rankings by level/points/streak

### Chat
- WebSocket: `ws://host/message?token=...&session_id=...`

## 📊 Grade Bands

| Grade | Band | Theme |
|-------|------|-------|
| 0-4 | 🌱 Sprouts | Cute, beginner-friendly |
| 5-9 | 🚀 Explorers | Adventure, discovery |
| 10-11 | 🎯 Champions | UNT exam prep |

## 🎮 Evolution Stages

| Level | Stage |
|-------|-------|
| 1-4 | Beginner |
| 5-9 | Thinker |
| 10-14 | Problem Solver |
| 15-19 | Scientist |
| 20+ | AI Master |

## 🐳 Docker

```bash
# Build
docker compose build

# Run
docker compose up -d

# View logs
docker compose logs -f
```

## 📝 License

MIT
