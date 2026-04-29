# Architecture — Serendipity

---

## App

- Name: `serendipity`
- Rails name: `Serendipity`

---

## Stack

| Layer | Choice | Notes |
|-------|--------|-------|
| Ruby on Rails | 8.x (latest stable) | |
| Database | PostgreSQL | |
| Auth | Devise + OmniAuth (Google) | No password login; Google OAuth only |
| Authorization | None (Pundit) | Simple owner-checks inline at MVP; add Pundit if complexity grows |
| Background jobs | None at MVP | Discogs retry loop is synchronous; revisit if latency becomes an issue |
| Frontend | Hotwire (Turbo + Stimulus) | Stimulus manages player state and YouTube IFrame API |
| CSS | Tailwind CSS | Dark mode via `dark:` variant; visual system deferred post-MVP |
| Component library | Deferred post-MVP | Tailwind utilities only at launch |
| File uploads | None | Cover art and videos sourced from Discogs/YouTube URLs |
| Tests | RSpec + FactoryBot + SimpleCov + Mutant | |
| Linting | RuboCop + rubocop-rails + rubocop-rspec | |
| Code quality | RubyCritic | |
| Security | Brakeman + bundler-audit | |
| DB integrity | database_consistency + strong_migrations | |

---

## External APIs

| Service | Purpose | Auth | Rate limit |
|---------|---------|------|------------|
| Discogs API | Random release fetch, genre-filtered search | Personal access token (`DISCOGS_API_KEY`) | 60 req/min (authenticated) |
| YouTube IFrame API | Embed and control video playback | None (client-side JS embed) | — |
| Google OAuth | User authentication | `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET` | — |

---

## Conventions

**Authentication**
- Google OAuth only. No password registration or password reset flows.
- `current_user` helper from Devise. Controllers that require auth use `before_action :authenticate_user!`.
- Anonymous sessions use Rails session store (no Devise session needed). Anonymous history lives in the client (localStorage); the server never stores anonymous data.

**Authorization**
- Owner-check pattern at MVP: `current_user == resource.user` before mutating mixtapes, history, likes.
- No policy objects at MVP unless multiple roles emerge.

**Frontend**
- Stimulus controllers own all interactive state: player, genre filter, dark/light mode toggle.
- Turbo Frames for replacing player content (next/prev) without full-page reload.
- Turbo Streams for updating like state, history count, mixtape track list.
- Drag-to-reorder in mixtapes: use a touch-compatible JS library (e.g. Sortable.js) wired via Stimulus.
- No React, no build step beyond Tailwind CLI (or Propshaft + importmap).

**Discogs integration**
- All Discogs API calls go through `app/services/discogs/` service objects.
- `Discogs::RandomReleaseService` encapsulates the fetch-and-retry loop; callers receive a fully hydrated `Release` record or an error.
- Releases are cached on first fetch (`find_or_create_by(discogs_id:)`). A release is only persisted if `video_url` is present.
- Genre filter passed as Discogs search param, not applied post-fetch. Multiple genres → OR query (one Discogs request per genre, results merged and sampled).

**YouTube embedding**
- YouTube IFrame API loaded once via Stimulus lifecycle. The player Stimulus controller manages play, pause, and swap (load new video by ID).
- Video ID extracted from the stored `video_url` at render time.

**Dark mode**
- Tailwind `darkMode: 'class'` strategy. `<html class="dark">` toggled by a Stimulus controller.
- Preference stored in localStorage key `serendipity_theme`. On authenticated users, also persisted to `users.theme` column (optional, add if needed).

**Testing**
- Request specs for all API/controller actions.
- System specs for critical user flows: discovery play, like/unlike, mixtape creation and sharing.
- FactoryBot factories for all models; no fixtures.
- VCR or WebMock to stub Discogs API calls in tests.

**Service objects**
- Cross-model or external-API logic in `app/services/`, named `VerbNounService`.
- Return a plain result object with `success?`, `value`, and `error` accessors.

---

## Decisions

**Hotwire over React** — Chosen: Hotwire. Rejected: React (Inertia), API-only + React SPA. The player's state (current track, history cursor, filter selection) is manageable in a single Stimulus controller without a virtual DOM. Keeps the stack Rails-native and avoids a separate frontend build pipeline at MVP.

**Devise + OmniAuth over Rodauth/custom** — Chosen: Devise + omniauth-google-oauth2. Rejected: Rodauth, bare OmniAuth. Best-documented path for Rails + Google OAuth; Devise's session and model helpers integrate cleanly with the rest of the app even if password auth is never enabled.

**Synchronous Discogs retry over pre-warmed cache** — Chosen: retry loop on demand. Rejected: background job pre-warming a release pool. Simpler infra at MVP, no stale-cache problem, no Redis/queue dependency. Acceptable latency for a discovery use case where the user expects a brief "loading" moment.

**Releases cached locally** — Discogs releases that have been fetched and confirmed playable are stored in the local `releases` table. This reduces redundant API calls, enables the history/likes/mixtapes foreign keys, and survives Discogs API downtime for re-playing known tracks.

**Share token over sequential ID for mixtapes** — Mixtape public URLs use a random token (`/m/:share_token`) rather than the DB primary key. Prevents enumeration of other users' mixtapes.
