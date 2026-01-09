# Teams Notification Center

A full-stack application for managing Microsoft Teams notifications with a Go backend and React frontend.

## Project Structure

```
TeamsNotifyClaude/
├── Backend/          # Go backend service
│   ├── services/     # Business logic
│   ├── handlers/     # HTTP handlers
│   ├── libs/         # Shared libraries
│   └── ...
├── Frontend/         # React frontend application
│   ├── pages/        # Page components
│   ├── services/     # API services
│   ├── components/   # Reusable components
│   └── ...
└── start-dev.sh      # Development startup script
```

## Quick Start

### Development

Run both backend and frontend in development mode:

```bash
./start-dev.sh
```

### Backend Only

```bash
cd Backend
make run
```

### Frontend Only

```bash
cd Frontend
npm install
npm run dev
```

## Architecture

- **Backend**: Go (Gin framework) with PostgreSQL database
- **Frontend**: React (Vite) with TypeScript
- **Communication**: RESTful API between frontend and backend

## Documentation

- Integration progress: See `INTEGRATION_PROGRESS.md`
- Integration fixes: See `INTEGRATION_FIXES.md`
- Integration analysis: See `INTEGRATION_ANALYSIS.md`

## Git Workflow

This project uses a monorepo structure with conventional commits:

- `feat:` - New features
- `fix:` - Bug fixes
- `chore:` - Maintenance tasks
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
