# CostGuard Dashboard

Next.js dashboard for visualizing CostGuard scan results and AI decisions.

## Features

- 📊 Real-time cost analysis visualization
- 🤖 AI decision log
- 🔍 Filter by repository/project
- 📈 Cost comparison charts
- 🔄 Auto-refresh every 30 seconds
- 📝 PR tracking

## Setup

### Local Development

```bash
npm install
npm run dev
```

Open http://localhost:3000

### Production Deployment

```bash
vercel
```

## API Endpoints

- `GET /api/scan?repo=owner/repo` - Get scan data
- `GET /api/decisions?repo=owner/repo` - Get decision logs
- `GET /api/repos` - List all repositories
- `POST /api/submit` - Receive updates from CLI

## Data Storage

Data is stored in `data/` directory (gitignored):
```
data/
├── owner_repo/
│   ├── latest-scan.json
│   ├── latest-decision.json
│   ├── scan-{timestamp}.json
│   └── decision-{timestamp}.json
```

## Environment Variables

- `COSTGUARD_DASHBOARD_URL` - Dashboard URL (for CLI)
- `GITHUB_TOKEN` - Optional, for GitHub API access

## Usage

Users run:
```bash
costguard fix --dashboard-url https://your-dashboard.vercel.app
```

Dashboard automatically receives and displays the data.

