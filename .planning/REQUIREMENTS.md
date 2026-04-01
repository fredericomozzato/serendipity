# Requirements: Serendipity

**Defined:** 2026-03-31
**Core Value:** A song is always playing within seconds of opening the app.

## v1 Requirements

### Infrastructure

- [x] **INFRA-01**: Developer can run the full app with `docker compose up` (Rails app container with named gem volume)
- [x] **INFRA-02**: Docker Compose is structured for zero-code-change PostgreSQL addition in V2 (commented `db` service + `DATABASE_URL` env var pattern)
- [x] **INFRA-03**: All V1 state is session-based — no database required

### Playback

- [ ] **PLAY-01**: User sees a YouTube video playing automatically within seconds of opening the app
- [ ] **PLAY-02**: YouTube embed is configured for mobile autoplay compliance (`mute=1`, `playsinline=1`, `allow="autoplay"` attribute)
- [ ] **PLAY-03**: When a fetched Discogs release has no linked YouTube video, the app silently fetches the next release (retries are capped to prevent infinite loops)
- [ ] **PLAY-04**: User can play and pause the current track
- [ ] **PLAY-05**: Artist name, track/release title, year, label, and album cover art are displayed for the current track

### Navigation

- [ ] **NAV-01**: User can skip to the next random song via a Skip button
- [ ] **NAV-02**: User can navigate back to a previously played song via a Back button (session history of last 10–20 songs, resets on page reload)

### Discovery

- [ ] **DISC-01**: User can filter songs by genre using the full Discogs genre taxonomy (~13 genres)
- [ ] **DISC-02**: User can filter songs by decade/era (maps to Discogs `year` range parameter)
- [ ] **DISC-03**: User can clear all active filters to return to unfiltered random discovery

### UI/UX

- [ ] **UI-01**: Layout is mobile-first with touch targets ≥ 48×48px and primary controls in thumb reach; fully usable on desktop at wider breakpoints
- [ ] **UI-02**: Retro/neon aesthetic implemented with Tailwind CSS (dark background, neon accents, vintage typography)
- [ ] **UI-03**: A retro-themed loading indicator is shown during all API fetches (e.g. "searching crates...")
- [ ] **UI-04**: An error/retry state is shown when no playable video is found after the retry cap is reached

### Quality

- [ ] **QUAL-01**: All domain logic has RSpec test coverage written before implementation (TDD)
- [ ] **QUAL-02**: Discogs API client handles rate limiting and transient errors with a retry strategy

## v2 Requirements

### Authentication

- **AUTH-01**: User can sign up with email and password
- **AUTH-02**: User session persists across browser refresh
- **AUTH-03**: User can reset password via email link

### Persistence

- **PERS-01**: User can favourite a track and retrieve it later
- **PERS-02**: User's genre/decade filter preferences are saved across sessions
- **PERS-03**: Played track history persists across page reloads

### Discovery

- **DISC-04**: User can view their full play history as a scrollable list
- **DISC-05**: User can click a "Now Playing" link to view the full Discogs release page

## Out of Scope

| Feature | Reason |
|---------|--------|
| mutant gem / mutation testing | Applied as a dedicated future milestone after V1 ships — not a V1 constraint |
| Volume control | YouTube embed + OS/browser handle volume; in-app control adds no value |
| Playlists / queue | Inverts the product's serendipity identity; skip IS the queue |
| Social features (share, vote, collaborative) | Requires real-time infra (Action Cable/Redis) not in V1 stack |
| Search by artist or album | Defeats the purpose; randomness is the feature |
| Recommendation engine / ML | Discogs genre + decade filter IS the discovery mechanism |
| Offline / PWA | Audio streams through YouTube; caching is not applicable |
| OAuth / magic link login | Email/password sufficient for V2; not needed in V1 |
| Podcast or non-music content | Discogs is a music catalog; scope is music only |

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| INFRA-01 | Phase 1 | Complete |
| INFRA-02 | Phase 1 | Complete |
| INFRA-03 | Phase 1 | Complete |
| PLAY-01 | Phase 2 | Pending |
| PLAY-02 | Phase 2 | Pending |
| PLAY-03 | Phase 2 | Pending |
| PLAY-04 | Phase 2 | Pending |
| PLAY-05 | Phase 2 | Pending |
| NAV-01 | Phase 3 | Pending |
| NAV-02 | Phase 3 | Pending |
| DISC-01 | Phase 3 | Pending |
| DISC-02 | Phase 3 | Pending |
| DISC-03 | Phase 3 | Pending |
| UI-01 | Phase 4 | Pending |
| UI-02 | Phase 4 | Pending |
| UI-03 | Phase 2 | Pending |
| UI-04 | Phase 2 | Pending |
| QUAL-01 | All phases (cross-cutting) | Pending |
| QUAL-02 | Phase 2 | Pending |

**Coverage:**
- v1 requirements: 19 total
- Mapped to phases: 19
- Unmapped: 0 ✓

---
*Requirements defined: 2026-03-31*
*Last updated: 2026-03-31 after roadmap creation*
