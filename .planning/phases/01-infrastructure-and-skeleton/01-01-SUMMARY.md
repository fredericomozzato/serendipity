---
phase: 01-infrastructure-and-skeleton
plan: 01
subsystem: infra
tags: [rails, docker, ruby, tailwindcss, propshaft, rspec, vcr, webmock, simplecov, postgresql]

# Dependency graph
requires: []
provides:
  - Rails 8.1.3 app skeleton with Importmap + Propshaft + Hotwire (no Node.js)
  - Dockerfile (FROM ruby:3.4, libpq-dev, named gem volume at /usr/local/bundle)
  - docker-compose.yml with web service, commented db + DATABASE_URL (zero-code V2 activation)
  - entrypoint.sh removing stale PID before server start
  - Gemfile with all V1 gems: pg, tailwindcss-rails, rspec-rails ~> 8.0, vcr, webmock, simplecov
  - App boots via docker compose up at localhost:3000, no database required
affects:
  - 01-02
  - phase-02
  - all subsequent phases

# Tech tracking
tech-stack:
  added:
    - rails 8.1.3
    - ruby 3.4.9 (aarch64-linux in container)
    - propshaft 1.3.1
    - importmap-rails 2.2.3
    - tailwindcss-rails 4.4.0
    - turbo-rails 2.0.23
    - stimulus-rails 1.3.4
    - pg 1.6.3
    - rspec-rails 8.0.4
    - vcr 6.4.0
    - webmock 3.26.2
    - simplecov 0.22.0
    - puma 7.2.0
  patterns:
    - Single-stage dev Dockerfile (no multi-stage)
    - Named Docker volume (serendipity_gems) at /usr/local/bundle for gem persistence
    - entrypoint.sh stale PID removal pattern
    - Commented db service + DATABASE_URL in Compose for zero-code V2 PostgreSQL activation
    - migration_error disabled in development (no DB in V1)

key-files:
  created:
    - Dockerfile
    - docker-compose.yml
    - entrypoint.sh
    - .dockerignore
    - Gemfile
    - Gemfile.lock
    - config/application.rb
    - config/routes.rb
    - config/puma.rb
    - app/views/layouts/application.html.erb
    - app/assets/tailwind/application.css
    - bin/rails
    - bin/rake
    - bin/setup
    - config/importmap.rb
    - .ruby-version
    - Rakefile
    - config.ru
  modified:
    - config/environments/development.rb

key-decisions:
  - "Rails 8.1.3 used (latest available) rather than 8.0.5 from research — satisfies ~> 8.x requirement"
  - "Disabled config.active_record.migration_error in development — Rails 8.1 checks migrations on page load; no DB in V1 causes 500 without this"
  - "Rails 8.1 generates solid_cache/solid_queue — only active in production, not an issue for dev"

patterns-established:
  - "Docker-first: all dev runs inside Docker containers (CLAUDE.md constraint)"
  - "Named gem volume serendipity_gems at /usr/local/bundle persists gems across container rebuilds"
  - "Commented db service block with V2 markers enables PostgreSQL via uncomment-only change"
  - "entrypoint.sh removes tmp/pids/server.pid before exec to prevent stale PID crashes"

requirements-completed: [INFRA-01, INFRA-02, INFRA-03]

# Metrics
duration: 10min
completed: 2026-04-01
---

# Phase 01 Plan 01: Infrastructure and Skeleton Summary

**Rails 8.1.3 skeleton app running in Docker at localhost:3000 with named gem volume, commented PostgreSQL service for zero-code V2 activation, and full test gem set (rspec-rails 8.0, vcr, webmock, simplecov)**

## Performance

- **Duration:** ~10 min
- **Started:** 2026-04-01T10:51:59Z
- **Completed:** 2026-04-01T10:57:00Z
- **Tasks:** 2
- **Files modified:** ~95 (91 from rails new + 5 Docker files + 1 development.rb fix)

## Accomplishments

- Generated Rails 8.1.3 skeleton with `--css tailwind --asset-pipeline=propshaft --database=postgresql` and all required skip flags (D-01 through D-06)
- Added all V1 test gems: rspec-rails ~> 8.0, vcr ~> 6.4, webmock ~> 3.26, simplecov ~> 0.22
- Created custom Dockerfile, docker-compose.yml, entrypoint.sh, .dockerignore satisfying INFRA-01/02/03
- Verified `docker compose up` + `curl localhost:3000` returns HTTP 200 (Rails welcome page)
- Preserved config/master.key (original 32-byte value) and spec/cassettes/ directory

## Task Commits

Each task was committed atomically:

1. **Task 1: Generate Rails 8 skeleton and add test gems** - `a26fd26` (feat)
2. **Task 2: Create Docker infrastructure files** - `cdfbc39` (feat)

**Plan metadata:** (pending)

## Files Created/Modified

- `Gemfile` - Rails 8.1.3 with pg, tailwindcss-rails, propshaft, turbo-rails, stimulus-rails, rspec-rails, vcr, webmock, simplecov
- `Gemfile.lock` - Fully resolved with all 128 gem dependencies
- `Dockerfile` - FROM ruby:3.4, apt-get libpq-dev, BUNDLE_PATH=/usr/local/bundle, ENTRYPOINT entrypoint.sh
- `docker-compose.yml` - web service with serendipity_gems volume, commented db + DATABASE_URL (INFRA-02)
- `entrypoint.sh` - Removes tmp/pids/server.pid then exec "$@" (D-10)
- `.dockerignore` - Excludes .git, .planning/, .ruby-lsp/, node_modules/, log/*, tmp/*
- `config/environments/development.rb` - migration_error check disabled (no DB in V1)
- `.ruby-version` - ruby-3.4.5
- `app/assets/tailwind/application.css` - Tailwind v4 CSS input file

## Decisions Made

- **Rails 8.1.3 vs 8.0.5:** Used latest locally available Rails (8.1.3); the locked version 8.0.5 from research was not available locally. Satisfies `~> 8.x` requirement.
- **migration_error disabled:** Rails 8.1 introduced `CheckPending` middleware that runs on every page load to detect pending migrations. With no database in V1, this causes a 500 error. Commented out with a V1 note. This is the correct approach for a no-DB V1.
- **solid_queue/solid_cache:** Rails 8.1 includes these by default but they are production-only in the generated config. No action needed for V1.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Disabled migration check preventing app boot**
- **Found during:** Task 2 (Docker verification — curl returned HTTP 500)
- **Issue:** Rails 8.1 `CheckPending` middleware tries to connect to PostgreSQL on every page load to check for pending migrations. With no database service running (INFRA-03 design), the app responded 500 on every request.
- **Fix:** Commented out `config.active_record.migration_error = :page_load` in `config/environments/development.rb` with a `# V1: No database in V1` comment.
- **Files modified:** `config/environments/development.rb`
- **Verification:** `curl -sf http://localhost:3000` returned HTTP 200 with Rails welcome page after fix and container rebuild.
- **Committed in:** `cdfbc39` (Task 2 commit)

---

**Total deviations:** 1 auto-fixed (Rule 1 - Bug)
**Impact on plan:** Fix was necessary for INFRA-03 compliance (no DB connection on boot). No scope creep.

## Issues Encountered

- Rails 8.1 (available locally) vs 8.0.5 (researched version) — no issues, both satisfy `~> 8.x`; Rails 8.1 generates slightly more files (solid_queue, solid_cache, kamal) but these do not affect V1 functionality

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- Rails 8.1.3 app is fully booted in Docker at localhost:3000
- All V1 gem dependencies locked in Gemfile.lock (128 gems)
- Docker infrastructure ready: web service, named gem volume, entrypoint.sh
- PostgreSQL infrastructure stubbed for zero-code V2 activation (commented db service + DATABASE_URL)
- Phase 01-02 can proceed: install RSpec, configure spec_helper.rb/rails_helper.rb/VCR/SimpleCov, write smoke spec

## Self-Check: PASSED

- Dockerfile: FOUND
- docker-compose.yml: FOUND
- entrypoint.sh: FOUND
- .dockerignore: FOUND
- Gemfile: FOUND
- Gemfile.lock: FOUND
- config/master.key: FOUND
- spec/cassettes/: FOUND
- 01-01-SUMMARY.md: FOUND
- Commit a26fd26: FOUND
- Commit cdfbc39: FOUND

---
*Phase: 01-infrastructure-and-skeleton*
*Completed: 2026-04-01*
