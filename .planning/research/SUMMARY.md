# Project Research Summary

**Project:** Serendipity — Random Music Discovery App
**Domain:** No-auth web jukebox (Discogs API + YouTube embed + Rails + Hotwire)
**Researched:** 2026-03-31
**Confidence:** HIGH (stack and architecture), MEDIUM (pitfalls — API behaviour)

## Executive Summary

Serendipity is a server-rendered Rails 8 jukebox that plays random music by combining the Discogs Search API (17M+ releases) with YouTube embeds, with no database or user accounts in V1. Experts in this space build this kind of app with thin Rails controllers delegating to pure service objects, Hotwire Turbo Streams for partial DOM updates (no WebSockets), and cookie-based session storage for play history. The recommended stack — Rails 8.0.5 / Ruby 3.4.9 / Faraday / Turbo-Rails / Stimulus / no Node build pipeline — is a coherent whole that avoids unnecessary complexity while remaining trivially extensible to V2 (PostgreSQL, auth, favourites).

The core value proposition ("a song is always playing within 2 seconds of opening the app") depends on three things working correctly from day one: autoplay policy compliance (muted embed with `playsinline=1`), correct Discogs pagination handling (hard cap at page 100 regardless of reported total), and a rate-aware retry loop (60 req/min authenticated limit with a hard retry cap). These are not polish concerns — they are load-bearing infrastructure that must be built before any UI work begins. The mutation testing requirement (mutant-rspec) further constrains architecture: all business logic must live in pure service objects under `app/services/`, not in controllers.

The recommended build order flows from dependencies: value objects first (no dependencies), then the Discogs HTTP adapter, then pure services (RandomPageStrategy, VideoFilter), then the orchestrating SearchService, then controllers and views, and finally Turbo Stream wiring and genre filter UI. This is an unusually well-understood problem domain with clear, well-documented patterns. The risk surface is narrow but concentrated in the Discogs API's underdocumented pagination behaviour and browser autoplay policy quirks.

---

## Key Findings

### Recommended Stack

Rails 8 ships a generated Dockerfile, Propshaft, and importmap-rails by default — eliminating all tooling overhead for a no-Node project. The combination of Turbo Frames and Stimulus is exactly the right tool for this use case: server renders HTML, Turbo swaps partials on skip/back, Stimulus manages the YouTube IFrame Player API lifecycle. No JavaScript framework, no Node build step.

**Core technologies:**
- Ruby 3.4.9 / Rails 8.0.5 — current stable, Rails 8 default asset pipeline eliminates build overhead
- Faraday ~> 2.14 + faraday-retry ~> 2.3 — middleware-based retry and backoff; the silent-retry-until-video pattern maps directly to Faraday middleware
- turbo-rails ~> 2.0 + stimulus-rails — partial DOM swaps and YouTube player lifecycle management; no SPA needed
- importmap-rails + propshaft — both bundled with Rails 8; zero configuration asset serving; no Node
- rspec-rails ~> 8.0 + mutant-rspec ~> 0.14 — project constraint; mutant requires RSpec; mutant-rspec is open-source licensed for public repos
- webmock ~> 3.x — HTTP stubbing for Faraday; prevents real Discogs calls in test suite
- Docker + Docker Compose v2 — project constraint; `ruby:3.4.9-slim-bookworm` base image; multi-stage Dockerfile; DB service pre-stubbed in compose file for V2

**Avoid:** `discogs-wrapper` gem (opaque, unclear maintenance), HTTParty (no middleware layer), hotwire-rails gem (deprecated), any Node/esbuild pipeline, NES.css (wrong aesthetic, mobile issues), Tailwind (requires Node).

### Expected Features

Research confirms the core feature set is well-understood from the music discovery / radio app domain. The retro neon aesthetic is a genuine differentiator at near-zero code cost (pure CSS).

**Must have (table stakes):**
- Immediate autoplay on load — the entire value proposition
- Loading indicator — covers 0.5–2s API fetch; without it users assume the app is broken
- Error state + silent retry (capped at 5 attempts) — ~60–80% of Discogs releases have no YouTube video
- Track metadata + album art — artist, title, year, label, cover image from Discogs response
- Skip button — always visible, in thumb zone, triggers new Discogs fetch
- Back / session history (10–15 entry window, no DB) — "what was that song?" is the #1 post-discovery frustration
- Genre filter (13 Discogs top-level genres, single-select) — expected for any catalog-backed product
- Mobile-first touch targets (48×48px minimum, controls in bottom-half viewport)

**Should have (differentiators):**
- Retro neon-on-dark aesthetic with VT323/Press Start 2P typography — personality at zero code cost
- "Serendipity" framing in UX copy — makes skipping feel like discovery, not failure
- "Now playing" link to Discogs release — one anchor tag; high value for power users
- Animated "searching crates" loading state — turns API latency into anticipation
- Surface the catalog depth ("pulling from 17 million releases") in empty/loading state

**Defer to V2+:**
- Favourites / likes (requires DB + auth)
- Visible session history as scrollable track list (back button covers the core need)
- Decade / era filter (low cost but not V1 blocking)
- Social / shareable links (requires real-time infra)
- User accounts, playlists, queue management, volume control, equalizer — all anti-features for V1

### Architecture Approach

The architecture is a classic thin-controller / fat-service Rails app, structured specifically for mutation testing. Controllers are restricted to: extract params/session, call one service, write session, render. All decisions live in service objects under `app/services/` namespaced as `Discogs::` and `History::`. Value objects (`Release`, `HistoryWindow`) use `Data.define` — immutable, zero logic, zero mutant targets. Turbo Streams are used over request-response (no ActionCable, no Redis). Session state (CookieStore, 4KB limit) stores only the 5–6 minimum fields per history entry.

**Major components:**
1. `Discogs::ApiClient` — Faraday HTTP adapter; auth headers; rate limit header inspection; single point of external API contact
2. `Discogs::RandomPageStrategy` — pure function; given total page count, returns `rand(1..[pages, 100].min)`; critical pagination cap logic lives here
3. `Discogs::VideoFilter` — pure function; given a list of releases, returns the first with a YouTube video URL (or nil)
4. `Discogs::SearchService` — orchestrator; calls ApiClient, RandomPageStrategy, VideoFilter; the retry loop lives here
5. `History::Reader` / `History::Writer` — pure functions on arrays; never touch session directly; controller reads/writes session and passes plain arrays
6. `SongsController` — three actions only: `show`, `skip`, `back`; no business logic
7. Turbo Stream templates (`skip.turbo_stream.erb`, `back.turbo_stream.erb`) — replace player and history-nav frames without full page reload

### Critical Pitfalls

1. **YouTube autoplay blocked by browser policy** — Embed must include `?autoplay=1&mute=1&playsinline=1` and `allow="autoplay; encrypted-media"` attribute. Silent failure: no JS error, just no playback. iOS Safari requires `playsinline=1`. Test on real iOS early; desktop Chrome is more permissive. Address in Phase 1.

2. **Discogs pagination hard cap at page 100** — `pagination.items` reports the true total (potentially millions) but only pages 1–100 are reachable. `RandomPageStrategy` must cap at `[pagination.pages, 100].min`, never use `items` for page calculation. Silent failure: retry loop exhausts rate limit budget. Address in Phase 1.

3. **Discogs rate limit exhaustion in retry loop** — 60 req/min with authenticated token (25 unauthenticated). A mashing user exhausts the budget in under 10 seconds. Hard-cap retries at 5 per action; inspect `X-Discogs-Ratelimit-Remaining`; always send a descriptive `User-Agent` header (required by Discogs ToS). Address in Phase 1.

4. **mutant discovers 0 subjects — load path and naming mismatch** — Most common first-run failure. Spec descriptions must match Ruby constant names exactly. Remove `require 'rspec/autorun'` from `spec_helper.rb`. Use `--include app --require config/environment` in `.mutant.yml`. Run `--zombie` as a sanity check. Architecture must be service-object-first from day one; retrofitting is painful. Address in Phase 0 setup.

5. **Docker Compose gem volume conflicts** — Host-mounted app volume can shadow gems installed into the image at build time. Use a named volume (`bundle_cache:/usr/local/bundle`) separate from the app code volume. Do not mount the app directory at the bundle path. Address in Phase 0.

---

## Implications for Roadmap

Based on the dependency graph in ARCHITECTURE.md and phase warnings in PITFALLS.md, the natural phase structure is:

### Phase 0: Infrastructure and Skeleton
**Rationale:** Docker gem-volume conflicts and mutant load-path configuration are foundational — getting these wrong means debugging environment issues instead of building features. Set up once, never revisit.
**Delivers:** Working Docker Compose environment, Rails 8 app skeleton with RSpec + mutant configured and returning non-zero subjects, `.mutant.yml` scoped to `app/services/`, CI-ready test run in Docker.
**Avoids:** Pitfall 4 (mutant zero subjects), Pitfall 6 (Docker gem volume conflicts).
**Research flag:** Standard patterns — no additional research needed. Rails 8 generated Dockerfile + Compose is well-documented.

### Phase 1: Core Playback (Discogs + YouTube)
**Rationale:** This is the entire value proposition. Every other feature depends on a song playing. The retry loop, rate-limit handling, and autoplay compliance are foundational — not polish. Build order: value objects → ApiClient → RandomPageStrategy + VideoFilter → SearchService → SongsController#show → basic view.
**Delivers:** App loads, fetches a random Discogs release, plays the YouTube video, displays metadata and album art. No skip, no history, no genre filter.
**Features:** Immediate autoplay on load, loading indicator, error state + silent retry (capped), track metadata + album art.
**Avoids:** Pitfall 1 (autoplay), Pitfall 2 (pagination cap), Pitfall 3 (rate limit exhaustion), Pitfall 7 (master vs. release video field).
**Research flag:** Discogs API pagination behaviour is community-documented, not official — verify the page 100 cap empirically in Phase 1 with a direct API call.

### Phase 2: Navigation and Discovery Controls
**Rationale:** Once a song is playing, the next user need is skip and back. Genre filter is table stakes and architecturally simple (a param on the existing `show` action). Session history requires careful data structure design to stay within the 4KB cookie limit and to avoid Turbo Drive cache desync.
**Delivers:** Skip button, back navigation (10–15 entry session history), genre filter (13 genres, single-select). Full core feature set complete.
**Features:** Skip, back/history, genre filter, mobile touch targets.
**Uses:** Turbo Streams for skip/back (request-response, no ActionCable); Turbo Drive for genre change (full page nav resets history cleanly).
**Avoids:** Pitfall 4 (Turbo history desync — use `data-turbo-action="advance"`, let Turbo own history), Pitfall 8 (4KB session limit — store only 5 minimum fields, cap at 15 entries), Pitfall 9 (low-coverage genres — cap retries at 3 for genre-filtered requests).
**Research flag:** Standard Hotwire patterns. No additional research needed; ARCHITECTURE.md covers data flow in detail.

### Phase 3: Retro Polish and UX Refinement
**Rationale:** The aesthetic is a differentiator at near-zero code cost. Once the core is working, retro CSS and UX copy can be layered on without touching business logic. Mobile responsiveness and thumb-zone validation happen here.
**Delivers:** Neon-on-dark aesthetic, VT323/Press Start 2P typography, animated loading state, "now playing" Discogs link, mobile-optimised touch targets, UX copy ("searching the crates...").
**Features:** Retro aesthetic, Discogs link, animated loading state, mobile UX.
**Research flag:** Custom CSS — no authoritative guide exists. Recommendation is evidence-based (what CSS frameworks cannot do). Implementation is low-risk.

### Phase Ordering Rationale

- Phase 0 before everything: mutant configuration must work before any business logic is written; retrofitting is documented as the most common and painful mutant failure mode.
- Phase 1 before Phase 2: skip and back cannot exist without a playing song; the retry loop and rate-limit wrapper must be in place before navigation features add more API calls.
- Phase 3 last: the retro CSS is purely additive; deferring it until Phase 3 keeps early development focused on correctness over appearance and avoids CSS rework if layout changes during Phase 2.
- Genre filter in Phase 2 (not Phase 1): it is architecturally trivial (a param on `show`) but depends on the session history data structure being settled to handle the history-reset-on-genre-change behaviour correctly.

### Research Flags

Phases needing deeper research during planning:
- **Phase 1:** Discogs pagination cap (page 100) is community-documented, not officially confirmed — validate empirically with a real API call against a broad genre. Discogs `videos` field on master vs. release type also needs empirical confirmation of hit rate.
- **Phase 1:** `faraday-retry` exact version — sourced from libraries.io, not RubyGems.org. Verify before pinning.
- **Phase 1:** `stimulus-rails` latest version — not confirmed in research. Check RubyGems.org before pinning.

Phases with standard patterns (skip research-phase):
- **Phase 0:** Rails 8 Docker setup is well-documented; mutant RSpec integration has clear official docs.
- **Phase 2:** Hotwire Turbo Streams request-response pattern is fully documented in official Turbo handbook.
- **Phase 3:** Custom CSS — no research needed; implementation is direct.

---

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH | Core versions verified against official release pages (Rails, Ruby, RubyGems). Two gem versions (faraday-retry, stimulus-rails) are MEDIUM — verify before pinning. |
| Features | HIGH | Table stakes derived from multiple music app domain sources; anti-features rationale is well-reasoned and consistent across research. |
| Architecture | HIGH | Standard Rails service-object patterns, official Turbo/Hotwire docs, mutant official docs. Component boundaries and data flow are unusually well-specified for research output. |
| Pitfalls | MEDIUM | YouTube autoplay: HIGH (official Google docs). Discogs pagination cap: MEDIUM (community-confirmed, not officially documented). Rate limits: MEDIUM (community forum). mutant config: HIGH (official repo). Docker volumes: MEDIUM (community). |

**Overall confidence:** HIGH

### Gaps to Address

- **Discogs pagination hard cap:** Community-documented but not in official API docs. Verify with a live API call in Phase 1 before building RandomPageStrategy. Cap defensively at `[pages, 100].min` regardless.
- **faraday-retry exact version:** Confirm on RubyGems.org before Gemfile lock. Listed as ~> 2.3 from libraries.io.
- **stimulus-rails latest version:** Not confirmed in research. Check RubyGems.org; pin to verified version.
- **Master release video hit rate:** The actual proportion of Discogs search results that are type `"master"` vs. `"release"` and their respective video hit rates are unknown. Log this in Phase 1 integration and adjust retry cap if needed.
- **Custom CSS retro aesthetic:** No authoritative reference exists. The approach (hand-written CSS with `text-shadow` glow effects, CSS custom properties, Google Fonts) is sound but the exact implementation is left to Phase 3.
- **4KB session cookie limit in practice:** Architecture caps history at 15 entries × ~200 bytes = ~3KB. Verify actual serialised size with real Discogs data in Phase 2; adjust cap constant if needed.

---

## Sources

### Primary (HIGH confidence)
- [Rails 8.0.5 release announcement](https://rubyonrails.org/2026/3/24/Rails-Versions-8-0-5-and-8-1-3-have-been-released) — version confirmed
- [Ruby 3.4.9 release](https://www.ruby-lang.org/en/news/2026/03/11/ruby-3-4-9-released/) — version confirmed
- [YouTube Embedded Players — Google Developers](https://developers.google.com/youtube/player_parameters) — autoplay, mute, playsinline parameters
- [Hotwire Turbo Handbook](https://turbo.hotwired.dev/handbook/streams) — Turbo Streams request-response pattern
- [Turbo Drive Handbook](https://turbo.hotwired.dev/handbook/drive) — restoration visits, history management
- [mutant GitHub](https://github.com/mbj/mutant) — scope configuration, RSpec integration
- [mutant-rspec docs](https://github.com/mbj/mutant/blob/main/docs/mutant-rspec.md) — integration setup
- [faraday on RubyGems.org](https://rubygems.org/gems/faraday) — version 2.14.1
- [turbo-rails on RubyGems.org](https://rubygems.org/gems/turbo-rails/versions) — version 2.0.21
- [rspec-rails on RubyGems.org](https://rubygems.org/gems/rspec-rails) — version 8.0.3
- [mutant-rspec on RubyGems.org](https://rubygems.org/gems/mutant-rspec) — version 0.14.1
- [ruby:3.4-slim-bookworm Dockerfile](https://github.com/docker-library/ruby/blob/main/3.4/slim-bookworm/Dockerfile) — base image confirmed
- [Rails CookieStore — Rails API](https://api.rubyonrails.org/classes/ActionDispatch/Session/CookieStore.html) — 4096-byte limit
- [Turbo history.pushState issue #792](https://github.com/hotwired/turbo/issues/792) — manual pushState conflict confirmed

### Secondary (MEDIUM confidence)
- [Discogs rate limiting — community forum](https://www.discogs.com/forum/thread/997721) — 60 req/min authenticated
- [Discogs 10,000 result cap — community forum](https://www.discogs.com/forum/thread/996056) — page 100 hard cap
- [Discogs User-Agent requirement — community forum](https://www.discogs.com/forum/thread/740651) — ToS requirement
- [faraday-retry on GitHub](https://github.com/lostisland/faraday-retry) — version ~> 2.3
- [Rails Service Objects — OneUptime Blog 2026](https://oneuptime.com/blog/post/2026-01-26-rails-service-objects/view) — `.call` pattern conventions
- [Docker Compose bundle install — bradgessler.com](https://bradgessler.com/articles/docker-bundler/) — named volume pattern
- [Best Ruby HTTP Clients 2026 — Scrapingdog](https://www.scrapingdog.com/blog/ruby-http-clients/) — Faraday vs HTTParty benchmark
- [NES.css mobile limitation](https://github.com/nostalgic-css/NES.css/) — mobile responsiveness issues confirmed

### Tertiary (LOW confidence)
- [Turbo and Rails Common Pitfalls — Reintech](https://reintech.io/blog/turbo-rails-ajax-application-pitfalls) — single source; verify against Turbo official docs
- [YouTube iframe iOS autoplay fix — technetexperts](https://www.technetexperts.com/youtube-iframe-ios-autoplay-fix/) — cross-checked against official Google docs; use official docs as authority

---

*Research completed: 2026-03-31*
*Ready for roadmap: yes*
