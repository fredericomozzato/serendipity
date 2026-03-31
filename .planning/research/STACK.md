# Technology Stack

**Project:** Serendipity — Random Music Discovery App
**Researched:** 2026-03-31
**Overall confidence:** HIGH (core stack), MEDIUM (CSS approach)

---

## Recommended Stack

### Core Framework

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| Ruby | 3.4.9 | Runtime | Latest stable 3.4 patch (released 2026-03-11); Rails 8 explicitly supports 3.2+; 3.4 brings performance gains (YJIT improvements, Prism parser default). Use 3.4 not 3.3 — 3.4 is now the actively maintained branch. |
| Rails | 8.0.5 | Web framework | Latest 8.0 patch (released 2026-03-24). Prefer 8.0.x over 8.1.x for this project — 8.0 is stable and battle-tested; 8.1 (released Oct 2025) adds features (job continuations, structured events) irrelevant to this app. Bug fixes until May 2026, security fixes beyond that. |

**Do NOT use Rails 7.x.** Rails 8.0 ships a generated Dockerfile, Propshaft as the default asset pipeline, and importmap-rails by default — all of which eliminate tooling overhead for this no-Node project.

---

### HTTP Client (Discogs API)

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| faraday | ~> 2.14 | HTTP client for Discogs Search API | Middleware-based architecture enables clean retry logic (critical for the "silent retry until a video is found" strategy), response parsing, and timeout configuration — all as composable layers. 87% faster than HTTParty in benchmarks. Actively maintained (2.14.1 released 2026-02-07). |
| faraday-retry | ~> 2.3 | Automatic retry with backoff | First-class Faraday middleware for exponential backoff and jitter. Handles Discogs 429 rate-limit responses and transient network errors. The silent-retry-until-video pattern maps directly to this middleware. |

**Do NOT use the `discogs-wrapper` gem.** The `discogs-wrapper` gem (latest: 2.5.0) wraps only specific endpoints and its last meaningful update is unclear. The Discogs Search API is a straightforward REST endpoint; raw Faraday calls keep the code simple, transparent, and mutation-testing-friendly (fewer opaque dependencies = more directly testable code). Building a thin, purpose-built `DiscogsClient` service object is the right architecture here.

**Do NOT use HTTParty.** HTTParty has no middleware layer. Implementing retry, timeout, and instrumentation requires monkey-patching or wrapper boilerplate. Faraday handles all of this cleanly.

---

### Frontend

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| turbo-rails | ~> 2.0 | Turbo Frames + Turbo Streams | Latest stable is 2.0.21 (released 2026-01-16). Skip/Back navigation and genre filter UI are exactly the use case for Turbo Frames — server renders partial HTML, Turbo swaps it in without a full page reload. No JavaScript framework needed. |
| stimulus-rails | ~> 1.x | Minimal JS controllers | Handles YouTube IFrame Player API lifecycle (player ready event, autoplay trigger). Stimulus keeps JS scoped and testable; avoid inline `<script>` tags. |
| importmap-rails | bundled with Rails 8 | JavaScript module loading | Ships with Rails 8 by default. No Node.js, no Webpack, no bundler. Pins @hotwired/turbo and @hotwired/stimulus directly from CDN. Appropriate for this app — no complex JS build pipeline needed. |
| propshaft | bundled with Rails 8 | Asset pipeline | Rails 8 default. Replaces Sprockets. Simple digest fingerprinting with no transpilation. CSS and static assets served with zero configuration overhead. |

**Do NOT use hotwire-rails gem.** It is officially deprecated. Use `turbo-rails` and `stimulus-rails` directly.

**Do NOT use jsbundling-rails or Node/esbuild.** This app has no complex JavaScript requirements. Adding a Node build step for a retro CSS app with Stimulus controllers would be unnecessary complexity.

---

### CSS / Styling

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| Custom hand-written CSS | N/A | Retro neon-on-dark aesthetic | No retro CSS framework fits this exact aesthetic well. NES.css targets 8-bit pixel art (not neon/vintage), has known mobile responsiveness issues, and is not maintained for modern Rails. The required aesthetic — neon glow on dark background, vintage typography — is 50–100 lines of custom CSS using `text-shadow` for glow effects and Google Fonts (Press Start 2P or VT323 for terminals, Cinzel for vintage). Custom CSS is also mutation-test-irrelevant (not Ruby code) so adds zero overhead to the testing strategy. |

**Recommended fonts (load via `<link>` in layout, no npm required):**
- `VT323` — terminal/retro monospace, excellent mobile legibility
- `Press Start 2P` — 8-bit pixel font for accent text / headings

**CSS architecture:** Single `application.css` with CSS custom properties (`--neon-pink`, `--dark-bg`, etc.) for the color palette. Mobile-first media queries. No Sass, no PostCSS — Propshaft serves plain CSS files directly.

---

### Testing

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| rspec-rails | ~> 8.0 | Test framework | Latest is 8.0.3 (released 2026-02-17). Requires Rails 7.2+, fully Rails 8 compatible. The mutant gem's primary supported integration is RSpec (mutant-rspec). This is a project constraint, not a choice. |
| mutant | ~> 0.14 | Mutation testing engine | Latest is 0.14.1 (released 2026-01-06). Open source license is free for public repositories — this project qualifies as open-source. |
| mutant-rspec | ~> 0.14 | Mutant + RSpec integration | Latest is 0.14.1 (same release as mutant). Install `mutant-rspec`; it pulls the correct `mutant` version as a dependency. Do NOT pin `mutant` directly — let `mutant-rspec` control the version constraint. |
| webmock | ~> 3.x | HTTP call stubbing in tests | Required for testing Faraday client. Prevents real Discogs API calls in test suite. Pairs naturally with Faraday's test adapter. |
| factory_bot_rails | ~> 6.x | Test data | Standard; useful even without a DB for building domain objects in specs. |
| shoulda-matchers | ~> 6.x | Matcher helpers | Reduces boilerplate in model/controller specs. |

**Architecture requirement for mutation testing:** The mutant gem selects tests via namespace prefix matching (e.g., `Serendipity::DiscogsClient#fetch_random_release` runs the `Serendipity::DiscogsClient` describe block). This requires:
1. All business logic lives in service objects under `app/services/`, not in controllers or views
2. Classes use module namespacing (e.g., `Serendipity::DiscogsClient`)
3. Methods are small, pure, and have a single responsibility — fat methods produce many surviving mutants

---

### Infrastructure

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| Docker | Latest stable | Container runtime | Project constraint. Use official Ruby image, not a Rails-specific image (the official Rails image is deprecated since 2016). |
| Docker Compose | v2 (compose v2 CLI) | Local dev orchestration | Single-service in V1 (Rails app only). Structured to trivially add `db` service for PostgreSQL in V2 by adding a `postgres:16-alpine` service and an `environment` block — no architectural change needed. |

**Docker base image:** `ruby:3.4.9-slim-bookworm`

Rationale:
- `slim` variant: 219MB vs 1.01GB for full bookworm. Sufficient for Rails — only needs `curl`, `libjemalloc2`, `libvips` added via apt.
- `bookworm` (Debian 12): Current stable Debian. Specify the release explicitly to avoid breakage on new Debian releases.
- `jemalloc` (`libjemalloc2`): Significant Ruby memory allocation performance improvement; trivial to add.

**Multi-stage Dockerfile pattern (Rails 8 generated default):**
- Stage 1 (`base`): Install system dependencies
- Stage 2 (`build`): Install gems, precompile assets
- Stage 3 (`production`): Copy artifacts, non-root user, minimal footprint

**V1 `docker-compose.yml` structure:**
```yaml
services:
  web:
    build: .
    ports: ["3000:3000"]
    environment:
      - DISCOGS_API_TOKEN
    volumes:
      - .:/rails  # dev only; remove for production image

# V2 addition (pre-stub for easy expansion):
# db:
#   image: postgres:16-alpine
#   environment:
#     POSTGRES_DB: serendipity_development
```

---

## Alternatives Considered

| Category | Recommended | Alternative | Why Not |
|----------|-------------|-------------|---------|
| Rails version | 8.0.5 | 8.1.x | 8.1 adds features irrelevant to V1 (job continuations, structured events); adds V1 risk for no gain. Upgrade to 8.1 after V1 ships. |
| HTTP client | Faraday | HTTParty | No middleware layer; retry/timeout/instrumentation require boilerplate; slower. |
| HTTP client | Faraday | Net::HTTP (stdlib) | No middleware composability; retry logic must be hand-rolled; poor DX for API clients. |
| Discogs client | Custom service object | `discogs-wrapper` gem | Opaque wrapper over a simple REST API; unclear maintenance status; harms mutation coverage by hiding logic in gem code. |
| CSS | Custom CSS | NES.css | 8-bit pixel art aesthetic, not neon/vintage; documented mobile responsiveness issues; unmaintained for modern Rails. |
| CSS | Custom CSS | Tailwind CSS | Requires Node.js build pipeline; overkill for a focused retro aesthetic; contradicts #NOBUILD philosophy. |
| CSS | Custom CSS | Bootstrap | Generic look, not retro; adds 30KB+ of unused styles; no neon aesthetic. |
| Frontend JS | Stimulus + Turbo | React/Vue | Contradicts Rails 8 #NOBUILD and Hotwire philosophy; adds Node dependency; no SPA benefits for a server-rendered music player. |
| Asset pipeline | Propshaft (Rails 8 default) | Sprockets | Sprockets is legacy; Propshaft is Rails 8 default; simpler, faster, no compilation step. |
| Ruby version | 3.4.9 | 3.3.x | 3.3 receives security-only updates; 3.4 is actively maintained with performance improvements. |

---

## Installation

```bash
# Create new Rails 8 app (no default test framework, will add RSpec manually)
gem install rails -v '~> 8.0.5'
rails new serendipity --skip-test --skip-bundle

# Key Gemfile additions (beyond Rails defaults)
# HTTP client
gem 'faraday', '~> 2.14'
gem 'faraday-retry', '~> 2.3'

# Testing
group :development, :test do
  gem 'rspec-rails', '~> 8.0'
  gem 'factory_bot_rails', '~> 6.0'
  gem 'webmock', '~> 3.0'
end

group :test do
  gem 'mutant-rspec', '~> 0.14'
  gem 'shoulda-matchers', '~> 6.0'
end

# Install
bundle install
rails generate rspec:install
```

**mutant configuration (`.mutant.yml`):**
```yaml
integration: rspec
includes:
  - app/services
  - app/models
requires:
  - config/environment
subjects:
  - Serendipity*  # namespace-prefix matching
```

**Run mutation testing:**
```bash
RAILS_ENV=test bundle exec mutant run --integration rspec
```

---

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Rails 8.0.5 | HIGH | Verified via rubyonrails.org release announcement (2026-03-24) |
| Ruby 3.4.9 | HIGH | Verified via ruby-lang.org release (2026-03-11) |
| Faraday 2.14.1 | HIGH | Verified via RubyGems.org (released 2026-02-07) |
| faraday-retry 2.3.x | MEDIUM | Version sourced from libraries.io; confirm exact version on RubyGems before pinning |
| turbo-rails 2.0.21 | HIGH | Verified via RubyGems.org (released 2026-01-16) |
| stimulus-rails version | MEDIUM | Latest version not confirmed from search results; check RubyGems.org before pinning |
| rspec-rails 8.0.3 | HIGH | Verified via RubyGems.org (released 2026-02-17) |
| mutant-rspec 0.14.1 | HIGH | Verified via RubyGems.org (released 2026-01-06) |
| Docker base image | HIGH | ruby:3.4.9-slim-bookworm confirmed from official Docker library |
| Custom CSS approach | MEDIUM | Best fit for aesthetic requirements; no authoritative "retro neon Rails" guide — recommendation based on evidence of what CSS frameworks cannot do, not a single definitive source |

---

## Sources

- [Rails 8.0.5 + 8.1.3 release announcement](https://rubyonrails.org/2026/3/24/Rails-Versions-8-0-5-and-8-1-3-have-been-released)
- [Rails 8.0 Release Notes](https://guides.rubyonrails.org/8_0_release_notes.html)
- [Ruby 3.4.9 Released](https://www.ruby-lang.org/en/news/2026/03/11/ruby-3-4-9-released/)
- [faraday on RubyGems.org](https://rubygems.org/gems/faraday)
- [faraday-retry on GitHub](https://github.com/lostisland/faraday-retry)
- [turbo-rails versions on RubyGems.org](https://rubygems.org/gems/turbo-rails/versions)
- [hotwire-rails deprecation notice](https://github.com/hotwired/hotwire-rails)
- [rspec-rails on RubyGems.org](https://rubygems.org/gems/rspec-rails)
- [mutant on RubyGems.org](https://rubygems.org/gems/mutant)
- [mutant-rspec on RubyGems.org](https://rubygems.org/gems/mutant-rspec)
- [Propshaft + importmap Rails 8 asset pipeline](https://radanskoric.com/articles/rails-assets-propshaft-importmaps)
- [NES.css mobile limitation disclosure](https://github.com/nostalgic-css/NES.css/)
- [ruby:3.4/slim-bookworm Dockerfile](https://github.com/docker-library/ruby/blob/main/3.4/slim-bookworm/Dockerfile)
- [Rails 8 Docker best practices](https://jetthoughts.com/blog/rails-8-docker-deployment-production-guide/)
- [Best Ruby HTTP Clients 2026](https://www.scrapingdog.com/blog/ruby-http-clients/)
