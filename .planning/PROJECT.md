# Serendipity

## What This Is

Serendipity is a random music discovery app built on Ruby on Rails. Users land on a page and a random song — sourced from the Discogs catalog via its Search API — immediately starts playing via an embedded YouTube video. Users can skip forward, navigate back through session history, and filter by genre. No account required.

## Core Value

A song is always playing within seconds of opening the app.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] Docker Compose infrastructure with Rails app container (terrain ready for PostgreSQL in V2)
- [ ] On page load, fetch a random Discogs release via Search API and auto-play its YouTube video
- [ ] Auto-skip releases with no linked YouTube video (silent retry)
- [ ] Skip button fetches the next random song
- [ ] Back button navigates session history (last 10–20 songs, resets on page reload)
- [ ] Genre filter UI lets users pick from the full Discogs genre taxonomy (~13 genres)
- [ ] When a genre is active, random songs are constrained to that genre
- [ ] Mobile-first layout with retro-inspired aesthetic (neon on dark, vintage typography)
- [ ] Full RSpec test suite following TDD

### Out of Scope

- User authentication — deferred to V2
- Database / persistence — no DB in V1; Docker Compose structured for easy PostgreSQL addition in V2
- Redis, Sidekiq, Action Cable — not needed for V1 feature set
- Curated release lists — using Discogs Search API directly
- Random ID approach — replaced by Search API random-page strategy to enable genre filtering

## Context

- **Previous attempt:** Owner built an earlier version using random Discogs IDs; worked for playback but genre filtering was not achievable with that approach. V1 switches to the Search API random-page strategy.
- **Discogs Search API randomness:** Fetch total result count → generate random page number → pick random item from page. Add `&genre=X` for genre filtering — same logic, no architectural change.
- **YouTube coverage:** Discogs `videos` field is present on a subset of releases (~20–40%). Strategy: silent retry until a release with a video is found.
- **Discogs API auth:** Owner has a personal API token (3× rate limit vs. unauthenticated).
- **Mutation testing (future):** Codebase will later be the subject of a public article series on mutation testing with the `mutant` gem — applied as a dedicated exercise after V1 ships, not as a V1 constraint.
- **Dual purpose:** Personal use tool + open-source portfolio project demonstrating AI-native development practices.

## Constraints

- **Tech Stack**: Ruby on Rails, Hotwire + Turbo (frontend), RSpec (tests), Docker + Docker Compose (infra) — no deviations
- **No DB (V1)**: All state is in-memory/session; Docker Compose must make adding PostgreSQL trivial for V2
- **Mobile-first**: UI designed for mobile viewport first, desktop secondary
- **Code quality**: TDD, clean code, well-structured service objects — this is portfolio material

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Discogs Search API (random page) over random ID | Enables genre filtering without extra API calls; random IDs can't be genre-filtered without fetching first | — Pending |
| No DB in V1 | V1 has no persistence requirements; keeps infra simple while Docker Compose is structured for easy V2 expansion | — Pending |
| RSpec over Minitest | Industry standard; best ecosystem for Rails TDD | — Pending |
| Session-only history | No auth/DB in V1; 10–20 song window is sufficient for the back-navigation UX | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd:transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd:complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-03-31 after initialization*
