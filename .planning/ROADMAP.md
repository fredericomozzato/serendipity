# Roadmap: Serendipity

## Overview

Serendipity is built in four phases that flow from the dependency graph. Phase 1 establishes the Docker + Rails skeleton — Docker gem volumes and the Rails + RSpec setup must be correct before any business logic is written. Phase 2 delivers the entire value proposition: a song plays within seconds of opening the app, with metadata, autoplay compliance, and a resilient retry loop. Phase 3 adds the navigation and discovery controls that make the app usable beyond the first listen. Phase 4 layers the retro/neon aesthetic that turns a functional app into a portfolio piece. QUAL-01 (TDD) is a cross-cutting constraint applied in every phase, not a phase of its own.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [x] **Phase 1: Infrastructure and Skeleton** - Docker Compose + Rails 8 skeleton with RSpec and mutant configured and green (completed 2026-04-03)
- [ ] **Phase 2: Core Playback Loop** - Random Discogs release auto-plays on load with metadata, loading state, and resilient retry
- [ ] **Phase 3: Navigation and Discovery** - Skip, back/history, and genre + decade filter controls wired via Turbo Streams
- [ ] **Phase 4: Retro UI Polish** - Neon-on-dark aesthetic, mobile-first layout, and responsive touch targets

## Phase Details

### Phase 1: Infrastructure and Skeleton
**Goal**: Developer can `docker compose up` and hit a running Rails 8 app where `rspec` passes
**Depends on**: Nothing (first phase)
**Requirements**: INFRA-01, INFRA-02, INFRA-03
**Cross-cutting**: QUAL-01 (TDD from first commit)
**Success Criteria** (what must be TRUE):
  1. `docker compose up` starts the Rails app with no manual steps; the app responds at localhost
  2. `bundle exec rspec` runs inside the container and exits green
  3. Docker Compose file contains a commented `db` service and a `DATABASE_URL` env var that require zero code changes to activate for V2
  4. All V1 state flows through Rails session (CookieStore) — no database queries are made anywhere in the app
**Plans:** 2/2 plans complete
Plans:
- [x] 01-01-PLAN.md — Rails 8 skeleton generation + Docker infrastructure (Dockerfile, Compose, entrypoint)
- [x] 01-02-PLAN.md — RSpec/VCR/WebMock/SimpleCov configuration + smoke spec

### Phase 2: Core Playback Loop
**Goal**: Users experience the core value — a YouTube video of a random Discogs release starts playing automatically within seconds of opening the app
**Depends on**: Phase 1
**Requirements**: PLAY-01, PLAY-02, PLAY-03, PLAY-04, PLAY-05, UI-03, UI-04, QUAL-02
**Cross-cutting**: QUAL-01 (TDD)
**Success Criteria** (what must be TRUE):
  1. Opening the app on a fresh browser tab starts a YouTube video playing within ~2 seconds, without any user interaction
  2. The video plays on iOS Safari without requiring an unmute gesture (muted autoplay with `playsinline=1`)
  3. When a Discogs release has no linked YouTube video, the app silently fetches the next release — the user sees the loading indicator, never a blank player
  4. After exhausting the retry cap, the user sees an error state with a "try again" affordance rather than a broken player
  5. Artist name, track/release title, year, label, and album art are visible alongside the video for every track that plays
**Plans**: TBD

### Phase 3: Navigation and Discovery
**Goal**: Users can steer their discovery session — skipping forward, stepping back, and constraining by genre or decade — all without a page reload
**Depends on**: Phase 2
**Requirements**: NAV-01, NAV-02, DISC-01, DISC-02, DISC-03
**Cross-cutting**: QUAL-01 (TDD)
**Success Criteria** (what must be TRUE):
  1. Tapping Skip replaces the player with a new random song (no full page reload; Turbo Stream swap)
  2. Tapping Back navigates to the previously played song; the Back button is disabled on the first song of the session
  3. Session history persists up to 15–20 songs within a page session and resets cleanly on page reload
  4. User can select any of the 13 Discogs top-level genres from a filter control; subsequent skips stay within that genre
  5. User can filter by decade/era; user can clear all active filters and return to fully random discovery
**Plans**: TBD

### Phase 4: Retro UI Polish
**Goal**: The app looks and feels like a deliberate design artifact — neon-on-dark retro aesthetic, mobile-first layout with thumb-zone controls, and a fully usable desktop experience
**Depends on**: Phase 3
**Requirements**: UI-01, UI-02
**Cross-cutting**: QUAL-01 (no regressions)
**Success Criteria** (what must be TRUE):
  1. The app renders with a dark background, neon accent colors, and vintage typography (VT323 or Press Start 2P) on all viewports
  2. All primary controls (Skip, Back, genre/decade filters) have touch targets of at least 48×48px and are reachable in the bottom half of the viewport on a 375px-wide mobile screen
  3. The layout is fully usable on desktop at wider breakpoints without horizontal scrolling or broken proportions
  4. The retro aesthetic is implemented entirely with Tailwind CSS utility classes — no external CSS framework is added
**Plans**: TBD

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3 → 4

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Infrastructure and Skeleton | 2/2 | Complete   | 2026-04-03 |
| 2. Core Playback Loop | 0/TBD | Not started | - |
| 3. Navigation and Discovery | 0/TBD | Not started | - |
| 4. Retro UI Polish | 0/TBD | Not started | - |
