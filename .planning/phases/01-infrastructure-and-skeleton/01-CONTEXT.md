# Phase 1: Infrastructure and Skeleton - Context

**Gathered:** 2026-04-01
**Status:** Ready for planning

<domain>
## Phase Boundary

Docker Compose + Rails 8 skeleton with RSpec configured and green. The goal is `docker compose up` starts the app at localhost and `bundle exec rspec` exits green inside the container. No business logic, no custom controllers — Phase 2 builds on this foundation.

</domain>

<decisions>
## Implementation Decisions

### Rails generation

- **D-01:** Generate with `--css tailwind --asset-pipeline=propshaft` — Importmap + Propshaft, no Node.js in container
- **D-02:** Keep Hotwire (Turbo + Stimulus) — Phase 3 uses Turbo Streams explicitly
- **D-03:** Skip: Action Mailer, Action Cable, Active Storage, Jbuilder, Test::Unit (`--skip-action-mailer --skip-action-cable --skip-active-storage --skip-jbuilder --skip-test`)
- **D-04:** Database adapter: `--database=postgresql` — adapter code in Gemfile, migrations never run in V1; `db` service commented out in Compose
- **D-05:** Ruby 3.4 — target version in `.ruby-version` and Dockerfile `FROM ruby:3.4`

### Docker setup

- **D-06:** Generate Rails app locally first, then write Dockerfile/Compose around it
- **D-07:** Single dev Dockerfile — no multi-stage, production Dockerfile is V2 concern
- **D-08:** Named gem volume: `serendipity_gems`, mounted at `/usr/local/bundle`
- **D-09:** Port mapping: `3000:3000`
- **D-10:** `entrypoint.sh` removes `tmp/pids/server.pid` then `exec "$@"` — prevents stale PID crashes

### Skeleton page

- **D-11:** Rails default page — no custom controller or route in Phase 1; Phase 2 creates `SongsController`
- **D-12:** Test scope: minimal RSpec spec that confirms the suite runs (proves wiring); no routing or request specs in Phase 1

### RSpec configuration

- **D-13:** VCR + WebMock configured in Phase 1; cassettes directory: `spec/cassettes/` (already exists)
- **D-14:** No FactoryBot — no DB models in V1; add in Phase 2 when models arrive
- **D-15:** SimpleCov configured in `spec/spec_helper.rb` from Phase 1 — coverage visible from first run

### Docker Compose structure

- **D-16:** `db` service included in `docker-compose.yml` but commented out with a clear `# V2: uncomment to activate PostgreSQL` note
- **D-17:** `DATABASE_URL` env var present in the `web` service environment (commented or set to a placeholder) — zero code changes needed to activate for V2

### Claude's Discretion

- Exact `spec_helper.rb` / `rails_helper.rb` configuration details (require order, VCR record mode)
- SimpleCov minimum coverage threshold (not set in V1)
- `.ruby-version` file format

</decisions>

<specifics>
## Specific Ideas

- The `spec/cassettes/` directory already exists in the repo — VCR must be configured to use it
- `config/master.key` already exists in the repo — preserve it, do not regenerate
- Dedicated branch per phase: create `feat/phase-1-infrastructure` before any implementation commits

</specifics>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Project requirements and scope
- `.planning/REQUIREMENTS.md` — INFRA-01, INFRA-02, INFRA-03 define acceptance criteria for this phase
- `.planning/ROADMAP.md` §Phase 1 — Success criteria (4 conditions that must be TRUE)

### Prior decisions
- `.planning/STATE.md` §Accumulated Context → Decisions — Docker gem volume, No DB in V1, RSpec + mutant-rspec, service objects under `app/services/`

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `spec/cassettes/` — VCR cassette directory already exists; configure VCR to point here
- `config/master.key` — already committed; do not run `rails credentials:edit` or regenerate

### Established Patterns
- None yet — greenfield project

### Integration Points
- Phase 2 will add `SongsController`, Discogs API client, and YouTube embed logic on top of this skeleton
- Phase 3 uses Turbo Streams (kept in `rails new`) — skeleton must not skip Hotwire

</code_context>

<deferred>
## Deferred Ideas

- Production Dockerfile (multi-stage) — V2 concern
- FactoryBot setup — add in Phase 2 when `Release` / `Track` models arrive
- mutant gem configuration — post-V1 milestone (per REQUIREMENTS.md Out of Scope)

</deferred>

---

*Phase: 01-infrastructure-and-skeleton*
*Context gathered: 2026-04-01*
