# Tournament Maker — Frontend

A Vue 3 + Vite + Tailwind + PrimeVue frontend for the Tournament Maker API.

## Stack

- **Vue 3** (Composition API + `<script setup>`)
- **Vite** — build tool
- **Tailwind CSS** — utility-first styling with a custom dark theme
- **PrimeVue 4** — dialog, toast, confirm dialog components
- **Pinia** — state management
- **Vue Router 4** — client-side routing
- **Axios** — API client with auto token-refresh interceptor
- **qrcode** — QR code generation for tournament invite
- **date-fns** — date formatting

## Setup

```bash
# Install dependencies
npm install

# Copy env and set API URL if needed
cp .env.example .env

# Start dev server
npm run dev

# Build for production
npm run build
```

## Pages

| Route | Description |
|---|---|
| `/login` | Sign in or register |
| `/` | Dashboard — all tournaments |
| `/tournaments/new` | Create tournament |
| `/tournaments/:id` | Participant view — encounters, scoreboard |
| `/tournaments/:id/manage` | Admin view — participants, encounters, QR invite |
| `/profile` | Profile, password change, account upgrade |

## Design

Dark theme with an industrial/editorial aesthetic:
- **Syne** display font (headings, labels)
- **DM Sans** body font
- **JetBrains Mono** for IDs/code
- Brand color: `#FF4D00` (orange-red)
- Accent: `#00C896` (jade) for success/wins
- Grid background with subtle noise texture

## Auth Flow

- Register with just a `visible_name` → participant-only account (gets token immediately)
- Register with `username` + `password` → full account that can create tournaments
- JWT stored in `localStorage`; auto-refreshed via axios interceptor

## API

Configured via `VITE_API_URL` env variable, defaults to `https://tournament-maker.baeriswyl.dev/api`.

See `src/api/index.js` for all endpoint wrappers.
