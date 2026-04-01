# Phase 1: Infrastructure and Skeleton - Research

**Researched:** 2026-04-01
**Domain:** Rails 8 + Docker Compose + RSpec/VCR/SimpleCov setup
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

- **D-01:** Generate with `--css tailwind --asset-pipeline=propshaft` — Importmap + Propshaft, no Node.js in container
- **D-02:** Keep Hotwire (Turbo + Stimulus) — Phase 3 uses Turbo Streams explicitly
- **D-03:** Skip: Action Mailer, Action Cable, Active Storage, Jbuilder, Test::Unit (`--skip-action-mailer --skip-action-cable --skip-active-storage --skip-jbuilder --skip-test`)
- **D-04:** Database adapter: `--database=postgresql` — adapter code in Gemfile, migrations never run in V1; `db` service commented out in Compose
- **D-05:** Ruby 3.4 — target version in `.ruby-version` and Dockerfile `FROM ruby:3.4`
- **D-06:** Generate Rails app locally first, then write Dockerfile/Compose around it
- **D-07:** Single dev Dockerfile — no multi-stage, production Dockerfile is V2 concern
- **D-08:** Named gem volume: `serendipity_gems`, mounted at `/usr/local/bundle`
- **D-09:** Port mapping: `3000:3000`
- **D-10:** `entrypoint.sh` removes `tmp/pids/server.pid` then `exec "$@"` — prevents stale PID crashes
- **D-11:** Rails default page — no custom controller or route in Phase 1; Phase 2 creates `SongsController`
- **D-12:** Test scope: minimal RSpec spec that confirms the suite runs (proves wiring); no routing or request specs in Phase 1
- **D-13:** VCR + WebMock configured in Phase 1; cassettes directory: `spec/cassettes/` (already exists)
- **D-14:** No FactoryBot — no DB models in V1; add in Phase 2 when models arrive
- **D-15:** SimpleCov configured in `spec/spec_helper.rb` from Phase 1 — coverage visible from first run
- **D-16:** `db` service included in `docker-compose.yml` but commented out with a clear `# V2: uncomment to activate PostgreSQL` note
- **D-17:** `DATABASE_URL` env var present in the `web` service environment (commented or set to a placeholder) — zero code changes needed to activate for V2

### Claude's Discretion

- Exact `spec_helper.rb` / `rails_helper.rb` configuration details (require order, VCR record mode)
- SimpleCov minimum coverage threshold (not set in V1)
- `.ruby-version` file format

### Deferred Ideas (OUT OF SCOPE)

- Production Dockerfile (multi-stage) — V2 concern
- FactoryBot setup — add in Phase 2 when `Release` / `Track` models arrive
- mutant gem configuration — post-V1 milestone (per REQUIREMENTS.md Out of Scope)
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| INFRA-01 | Developer can run the full app with `docker compose up` (Rails app container with named gem volume) | Docker Compose patterns, named volume `serendipity_gems`, entrypoint.sh, port 3000 |
| INFRA-02 | Docker Compose structured for zero-code-change PostgreSQL addition in V2 (commented `db` service + `DATABASE_URL` env var pattern) | Compose service structure with commented blocks, env var placeholder pattern |
| INFRA-03 | All V1 state is session-based — no database required | Rails CookieStore default, `secret_key_base` from `master.key` (already present in repo), no DB connection needed |
</phase_requirements>

---

## Summary

Phase 1 is a pure infrastructure setup: generate a Rails 8 skeleton app, wrap it in a development Docker container, and wire up RSpec with VCR/WebMock/SimpleCov. No business logic is written; the only test is a "suite is wired correctly" smoke spec. Success means `docker compose up` serves the Rails default page at localhost:3000 and `bundle exec rspec` exits green inside the container.

The key complexity is the tailwindcss-rails v4 gem, which uses platform-specific native binaries. Because this is a dev-only Docker setup (not a cross-architecture production build), the pitfall is manageable — the Dockerfile must install the correct gem platform for `linux/amd64` or `linux/arm64` explicitly. Rails 8 generates a Dockerfile and entrypoint by default since 7.1; the generated files need to be simplified (no multi-stage, no asset precompile step) for the dev-only D-07 constraint.

The `config/master.key` already exists in the repo, so the `secret_key_base` for the CookieStore session is available without any additional secrets setup. The PostgreSQL adapter gem (`pg`) will be in the Gemfile but the `db:` Compose service is commented out — this is the zero-code-change V2 activation pattern (INFRA-02).

**Primary recommendation:** Run `rails new` locally with the locked flags, commit, then layer the Dockerfile and docker-compose.yml on top. Do not modify `config/master.key`.

---

## Standard Stack

### Core

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| rails | 8.0.5 | Web framework | Current stable, locked decision |
| ruby | 3.4 | Runtime | Locked decision D-05 |
| propshaft | bundled with Rails 8 | Asset pipeline | Replaces Sprockets; no Node needed |
| importmap-rails | bundled with Rails 8 | JS module loading | Node-free JS, locked decision D-01 |
| tailwindcss-rails | ~> 4.4 | Tailwind v4 CSS, no Node | Locked decision D-01; v4.4.0 is current (Oct 2025) |
| turbo-rails | bundled with Hotwire | Turbo Streams (Phase 3) | Kept per D-02 |
| stimulus-rails | bundled with Hotwire | Stimulus JS | Kept per D-02 |
| pg | ~> 1.5 | PostgreSQL adapter | Gemfile only in V1, DB never started |

### Testing

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| rspec-rails | ~> 8.0 | Test framework | rspec-rails 8.x required for Rails 8 (7.x is Rails 7-only) |
| vcr | ~> 6.4 | HTTP cassette recording | Configured Phase 1, used Phase 2+ for Discogs calls |
| webmock | ~> 3.26 | HTTP stubbing | VCR dependency; blocks real HTTP in specs |
| simplecov | ~> 0.22 | Coverage reporting | Configured in spec_helper; latest is 0.22.0 |

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| tailwindcss-rails v4 | pin to v3.3.1 | v3 is stable, avoids arch binary issues; v4 is current Rails 8 default |
| rspec-rails 8.x | minitest | rspec-rails 8.x is the only option for Rails 8; minitest was skipped via --skip-test |

**Installation (inside container after `bundle install`):**
```bash
# Rails new (run locally, D-06)
rails new serendipity \
  --css tailwind \
  --asset-pipeline=propshaft \
  --database=postgresql \
  --skip-action-mailer \
  --skip-action-cable \
  --skip-active-storage \
  --skip-jbuilder \
  --skip-test

# Then in Gemfile, add test group gems:
# gem "rspec-rails", "~> 8.0"
# gem "vcr", "~> 6.4"
# gem "webmock", "~> 3.26"
# gem "simplecov", "~> 0.22", require: false
```

**Version verification (confirmed against RubyGems.org 2026-04-01):**
- `rails` 8.0.5 — released 2026-03-24
- `rspec-rails` 8.0.4 — released 2026-03-11
- `tailwindcss-rails` 4.4.0 — released 2025-10-27
- `vcr` 6.4.0 — released 2025-12-22
- `webmock` 3.26.2 — released 2026-03-18
- `simplecov` 0.22.0 — released 2022-12-23 (stable, no newer version)

---

## Architecture Patterns

### Recommended Project Structure

```
serendipity/
├── app/
│   └── assets/
│       ├── builds/            # tailwind output (generated, gitignored)
│       └── tailwind/
│           └── application.css   # Tailwind v4 input file
├── config/
│   └── master.key             # ALREADY EXISTS — do not regenerate
├── spec/
│   ├── cassettes/             # ALREADY EXISTS — VCR cassette dir
│   ├── spec_helper.rb         # SimpleCov + VCR + WebMock config
│   └── rails_helper.rb        # Rails-specific RSpec config
├── Dockerfile                 # Single dev Dockerfile (D-07)
├── docker-compose.yml         # web + commented db (D-16, D-17)
└── entrypoint.sh              # Removes server.pid, then exec "$@" (D-10)
```

### Pattern 1: Dev-Only Single-Stage Dockerfile

**What:** A single Dockerfile for development. No production multi-stage build.
**When to use:** D-07 locks this — keep it simple for V1.

```dockerfile
# Source: official Docker Rails guidance + D-07 constraint
FROM ruby:3.4

RUN apt-get update -qq && apt-get install -y build-essential libpq-dev

WORKDIR /app

COPY Gemfile Gemfile.lock ./
RUN bundle install

COPY . .

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh
ENTRYPOINT ["entrypoint.sh"]

EXPOSE 3000
CMD ["bundle", "exec", "rails", "server", "-b", "0.0.0.0"]
```

### Pattern 2: Named Gem Volume in Docker Compose

**What:** Mount a named volume at `/usr/local/bundle` so gems persist across container rebuilds.
**When to use:** D-08 locks `serendipity_gems` as the volume name.

```yaml
# docker-compose.yml (simplified)
services:
  web:
    build: .
    command: bundle exec rails server -b 0.0.0.0
    volumes:
      - .:/app
      - serendipity_gems:/usr/local/bundle
    ports:
      - "3000:3000"
    environment:
      - RAILS_ENV=development
      - RAILS_MASTER_KEY=${RAILS_MASTER_KEY}
      # DATABASE_URL=postgresql://postgres:password@db:5432/serendipity_development  # V2: uncomment with db service

  # V2: uncomment to activate PostgreSQL
  # db:
  #   image: postgres:16
  #   volumes:
  #     - postgres_data:/var/lib/postgresql/data
  #   environment:
  #     - POSTGRES_PASSWORD=password

volumes:
  serendipity_gems:
  # postgres_data:  # V2: uncomment with db service
```

### Pattern 3: entrypoint.sh — Stale PID Prevention

**What:** Remove `tmp/pids/server.pid` before starting the server.
**When to use:** Always required in Rails Docker dev setups. D-10 locks this.

```bash
#!/bin/bash
set -e

# Remove stale server PID if it exists (prevents "A server is already running" crash)
rm -f /app/tmp/pids/server.pid

exec "$@"
```

### Pattern 4: RSpec + VCR + WebMock + SimpleCov Wiring

**What:** `spec_helper.rb` configures SimpleCov before any other requires; `rails_helper.rb` configures VCR and WebMock.

```ruby
# spec/spec_helper.rb
# SimpleCov MUST be required and started before everything else
require 'simplecov'
SimpleCov.start 'rails'

RSpec.configure do |config|
  # ... standard rspec config
end
```

```ruby
# spec/rails_helper.rb
require 'spec_helper'
ENV['RAILS_ENV'] ||= 'test'
require_relative '../config/environment'
require 'rspec/rails'
require 'vcr'
require 'webmock/rspec'

VCR.configure do |c|
  c.cassette_library_dir = 'spec/cassettes'  # directory already exists in repo
  c.hook_into :webmock
  c.ignore_localhost = true
  c.configure_rspec_metadata!
  c.default_cassette_options = { record: :new_episodes }
end

WebMock.disable_net_connect!(allow_localhost: true)
```

### Pattern 5: INFRA-02 — Zero-Code-Change V2 Activation

**What:** The `DATABASE_URL` env var is present but commented in Compose. The `db:` service block is present but commented. Activating PostgreSQL in V2 requires only uncommenting — no code changes.

**Critical:** The comment block structure must be maintained precisely so a `# V2:` grep makes the intent clear.

### Anti-Patterns to Avoid

- **Running `rails credentials:edit`:** `config/master.key` already exists. Regenerating it will invalidate the existing encrypted credentials.
- **Generating the app inside Docker:** D-06 locks local generation. Running `rails new` inside the container adds OS-level complexity and permission issues.
- **Using rspec-rails 7.x:** Version 7.x only supports Rails 7.0–7.2. Rails 8 requires rspec-rails 8.x.
- **Using tailwindcss-rails 3.x:** The gem's v4.x is what `rails new --css tailwind` generates when Rails 8 is current. Pinning to 3.x requires explicit override.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Stale PID on container restart | Custom health check or sleep loop | `entrypoint.sh` with `rm -f tmp/pids/server.pid` | Standard pattern; one-liner is the ecosystem answer |
| Gem caching across rebuilds | Copy gems into image on every build | Named Docker volume (`serendipity_gems`) mounted at `/usr/local/bundle` | Volumes persist across `docker compose down/up`; image layers do not |
| HTTP stubbing in tests | Custom Net::HTTP monkey-patch | WebMock (`webmock/rspec`) | Handles all HTTP adapters (Net::HTTP, Faraday, etc.) |
| Cassette-based HTTP replay | Custom fixture serialization | VCR gem | Handles all edge cases: redirects, auth headers, binary bodies |
| Coverage reporting | Manually count tests | SimpleCov with `'rails'` profile | Rails profile excludes config/, db/, vendor/ automatically |

**Key insight:** Every one of these "solve it yourself" instincts has a mature gem with edge-case handling accumulated over years. This phase is about wiring, not building.

---

## Common Pitfalls

### Pitfall 1: rspec-rails Version Mismatch

**What goes wrong:** Installing `rspec-rails ~> 7.0` with Rails 8 — bundle resolves but RSpec config hooks silently break or raise at runtime.
**Why it happens:** rspec-rails 7.x declares `railties < 8.0` as a dependency; bundler may warn but developers often ignore it.
**How to avoid:** Pin explicitly to `rspec-rails ~> 8.0` in the Gemfile. Verify with `bundle exec rspec --version`.
**Warning signs:** `NoMethodError` on Rails 8 internal APIs during spec load.

### Pitfall 2: tailwindcss-rails Native Binary Architecture Mismatch

**What goes wrong:** On Apple Silicon Macs, building a `linux/amd64` container causes the Tailwind CSS CLI binary to fail silently or with a "file not found" error even when `application.css` exists.
**Why it happens:** `tailwindcss-rails` v4 bundles a native binary via `tailwindcss-ruby`. When x86_64 emulation runs via Rosetta/QEMU, floating-point truncation causes the binary to malfunction.
**How to avoid:** For this V1 dev-only setup, the container target platform should match the host. On Apple Silicon, `FROM ruby:3.4` will pull the `linux/arm64` image and the native binary will be correct. Do NOT add `platform: linux/amd64` to the Compose service unless there's a specific reason.
**Warning signs:** `bin/rails tailwindcss:build` exits 0 but `app/assets/builds/tailwind.css` is empty or not created.

### Pitfall 3: Stale `config/master.key`

**What goes wrong:** Running `rails credentials:edit` regenerates `master.key`, invalidating all encrypted credentials and causing `ActiveSupport::MessageEncryptor::InvalidMessage` errors.
**Why it happens:** The key already exists in this repo (CONTEXT.md confirms this). `rails new` over an existing directory can prompt credential regeneration.
**How to avoid:** Do not run `rails credentials:edit`. After `rails new`, verify `config/master.key` still contains the original value. Add `config/master.key` to `.gitignore` (Rails does this by default, but verify).
**Warning signs:** `secret_key_base` resolution fails; session cookies become invalid on every restart.

### Pitfall 4: `bundle exec rspec` Fails on Missing `rails_helper` Require

**What goes wrong:** Running `bundle exec rspec` immediately after `rails generate rspec:install` shows "cannot load such file — vcr" or similar because gems are referenced in `rails_helper.rb` before being added to the Gemfile.
**Why it happens:** `rspec:install` generates a minimal `rails_helper.rb`. VCR/WebMock/SimpleCov must be manually added to both the Gemfile and the helper files.
**How to avoid:** Add all test gems to the Gemfile first, run `bundle install` inside the container, then edit the helper files.
**Warning signs:** `LoadError` on `require 'vcr'` or `require 'simplecov'`.

### Pitfall 5: SimpleCov Must Start Before Rails Loads

**What goes wrong:** SimpleCov reports 0% coverage or misses most files.
**Why it happens:** SimpleCov hooks into `require` to track which files are loaded. If Rails is required before `SimpleCov.start`, all Rails files are already loaded and not tracked.
**How to avoid:** `require 'simplecov'; SimpleCov.start 'rails'` must be the first two lines of `spec/spec_helper.rb`, before any other require.
**Warning signs:** Coverage report shows only a few files or an unrealistically low total.

### Pitfall 6: PostgreSQL Gem Requires `libpq-dev` in Dockerfile

**What goes wrong:** `bundle install` fails with a native extension compilation error for the `pg` gem.
**Why it happens:** The `pg` gem requires `libpq-dev` (Debian/Ubuntu) to compile. The base `ruby:3.4` image does not include it.
**How to avoid:** Add `libpq-dev` to the `apt-get install` line in the Dockerfile.
**Warning signs:** `Gem::Ext::BuildError: ERROR: Failed to build gem native extension` during `bundle install`.

---

## Code Examples

### Dockerfile (dev-only, D-07)

```dockerfile
# Source: official Docker Ruby guidance + project decisions D-05, D-07, D-08
FROM ruby:3.4

RUN apt-get update -qq && \
    apt-get install -y --no-install-recommends \
      build-essential \
      libpq-dev \
      curl && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Gems installed into named volume; BUNDLE_PATH must match volume mount
ENV BUNDLE_PATH=/usr/local/bundle

COPY Gemfile Gemfile.lock ./
RUN bundle install

COPY . .

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh
ENTRYPOINT ["entrypoint.sh"]

EXPOSE 3000
CMD ["bundle", "exec", "rails", "server", "-b", "0.0.0.0"]
```

### docker-compose.yml (INFRA-01, INFRA-02)

```yaml
# Source: project decisions D-08, D-09, D-16, D-17
services:
  web:
    build: .
    command: bash -c "bin/rails tailwindcss:build && bundle exec rails server -b 0.0.0.0"
    volumes:
      - .:/app
      - serendipity_gems:/usr/local/bundle
    ports:
      - "3000:3000"
    environment:
      - RAILS_ENV=development
      - RAILS_MASTER_KEY=${RAILS_MASTER_KEY}
      # V2: uncomment with db service below
      # - DATABASE_URL=postgresql://postgres:password@db:5432/serendipity_development
    tty: true

  # V2: uncomment to activate PostgreSQL
  # db:
  #   image: postgres:16
  #   volumes:
  #     - postgres_data:/var/lib/postgresql/data
  #   environment:
  #     - POSTGRES_PASSWORD=password

volumes:
  serendipity_gems:
  # postgres_data:  # V2: uncomment with db service
```

### entrypoint.sh (D-10)

```bash
#!/bin/bash
set -e

# Remove stale server PID (prevents crash on container restart)
rm -f /app/tmp/pids/server.pid

exec "$@"
```

### spec/spec_helper.rb skeleton

```ruby
# SimpleCov MUST come first, before Rails is loaded
require 'simplecov'
SimpleCov.start 'rails'

RSpec.configure do |config|
  config.expect_with :rspec do |expectations|
    expectations.include_chain_clauses_in_custom_matcher_descriptions = true
  end

  config.mock_with :rspec do |mocks|
    mocks.verify_partial_doubles = true
  end

  config.shared_context_metadata_behavior = :apply_to_host_groups
end
```

### spec/rails_helper.rb skeleton

```ruby
require 'spec_helper'
ENV['RAILS_ENV'] ||= 'test'
require_relative '../config/environment'
require 'rspec/rails'
require 'vcr'
require 'webmock/rspec'

VCR.configure do |c|
  c.cassette_library_dir = 'spec/cassettes'  # directory already exists
  c.hook_into :webmock
  c.ignore_localhost = true
  c.configure_rspec_metadata!
  c.default_cassette_options = { record: :new_episodes }
end

WebMock.disable_net_connect!(allow_localhost: true)

RSpec.configure do |config|
  config.fixture_paths = [Rails.root.join('spec/fixtures')]
  config.use_transactional_fixtures = true
  config.infer_spec_type_from_file_location!
  config.filter_rails_from_backtrace!
end
```

### Minimal smoke spec (D-12)

```ruby
# spec/smoke_spec.rb
# Source: D-12 — prove the RSpec suite wiring works, nothing more
RSpec.describe "RSpec wiring" do
  it "runs" do
    expect(true).to be true
  end
end
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Sprockets asset pipeline | Propshaft | Rails 8 default (Nov 2024) | No fingerprinting complexity; just digest and serve |
| rspec-rails 6.x/7.x | rspec-rails 8.x for Rails 8 | rspec-rails 8.0.0 released Apr 2025 | Must use 8.x; 7.x incompatible with Rails 8 |
| tailwindcss-rails 3.x (Tailwind CSS v3) | tailwindcss-rails 4.x (Tailwind CSS v4) | v4.0 released 2025 | CSS-first config; no `tailwind.config.js` needed |
| Webpacker / cssbundling-rails | tailwindcss-rails standalone executable | Progressive (2021-2024) | No Node.js required at all |
| `bin/docker-entrypoint` (generated by Rails) | Custom `entrypoint.sh` at project root | Always valid | D-06 generates locally; entrypoint must be written manually |

**Deprecated/outdated:**
- `rspec-rails ~> 7.0`: Incompatible with Rails 8 — do not use
- `tailwindcss-rails ~> 3.0`: Targets Tailwind v3; current `rails new` produces v4 config
- Sprockets: Still available but Propshaft is now the Rails 8 default for new apps

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | rspec-rails 8.0.4 |
| Config file | `.rspec` + `spec/spec_helper.rb` + `spec/rails_helper.rb` (created by `rails generate rspec:install`) |
| Quick run command | `docker compose run --rm web bundle exec rspec` |
| Full suite command | `docker compose run --rm web bundle exec rspec --format documentation` |

### Phase Requirements to Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| INFRA-01 | `docker compose up` starts app at localhost:3000 | smoke (manual verify) | `docker compose up -d && curl -f http://localhost:3000` | Wave 0 script |
| INFRA-02 | Compose has commented `db:` service and `DATABASE_URL` placeholder | structural (manual inspect) | `grep -n 'V2' docker-compose.yml` | Wave 0 |
| INFRA-03 | No DB queries on any request; session is cookie-based | unit | `bundle exec rspec spec/smoke_spec.rb` (no DB touched) | Wave 0 |
| QUAL-01 (cross-cut) | RSpec suite exists and is green | unit | `bundle exec rspec` | Wave 0 |

### Sampling Rate

- **Per task commit:** `docker compose run --rm web bundle exec rspec`
- **Per wave merge:** `docker compose run --rm web bundle exec rspec --format documentation`
- **Phase gate:** Full suite green + `curl http://localhost:3000` returns 200 before moving to Phase 2

### Wave 0 Gaps

- [ ] `spec/smoke_spec.rb` — minimal "RSpec wiring works" spec (D-12)
- [ ] `spec/spec_helper.rb` — SimpleCov config must be first lines
- [ ] `spec/rails_helper.rb` — VCR + WebMock configuration
- [ ] `.rspec` — generated by `rails generate rspec:install`
- [ ] Framework install: `bundle exec rails generate rspec:install` — once Gemfile has rspec-rails

---

## Open Questions

1. **tailwindcss-rails v4 on Apple Silicon dev machine**
   - What we know: `rails new --css tailwind` generates v4 config; native binary works correctly when container architecture matches host
   - What's unclear: Whether the developer's machine is Apple Silicon (arm64) or Intel (x86_64)
   - Recommendation: Do not add `platform: linux/amd64` to the Compose service. Let Docker select the native platform. Document this in a comment in `docker-compose.yml`.

2. **`rails new` over partial directory**
   - What we know: The repo currently has `app/assets/`, `config/master.key`, `spec/cassettes/` — a partial skeleton
   - What's unclear: Whether `rails new .` will preserve `config/master.key` or prompt to overwrite it
   - Recommendation: Run `rails new` with `--force` flag but backup `config/master.key` first and restore it immediately after generation. Verify the key value is unchanged.

3. **Tailwind CSS build in `docker compose up`**
   - What we know: `tailwindcss-rails` attaches `tailwindcss:build` to `assets:precompile` and `test:prepare` — but these don't run on `rails server` startup
   - What's unclear: Whether the CSS output file (`app/assets/builds/tailwind.css`) will be present on first `docker compose up` without a separate build step
   - Recommendation: Run `bin/rails tailwindcss:build` as part of the `command` in docker-compose.yml before `rails server`, OR include it in `entrypoint.sh`. The default page must load CSS without a manual step to satisfy INFRA-01.

---

## Sources

### Primary (HIGH confidence)

- RubyGems.org/gems/rails — version 8.0.5 confirmed, released 2026-03-24
- RubyGems.org/gems/rspec-rails — version 8.0.4 confirmed, released 2026-03-11; Rails 8 compatibility verified
- RubyGems.org/gems/tailwindcss-rails — version 4.4.0 confirmed, released 2025-10-27; targets Tailwind v4
- RubyGems.org/gems/vcr — version 6.4.0 confirmed, released 2025-12-22
- RubyGems.org/gems/webmock — version 3.26.2 confirmed, released 2026-03-18
- RubyGems.org/gems/simplecov — version 0.22.0 confirmed (stable)
- github.com/rails/tailwindcss-rails README — build process, VCR/WebMock integration, platform binary docs

### Secondary (MEDIUM confidence)

- https://radanskoric.com/articles/rails-assets-propshaft-importmaps — Propshaft + importmap-rails relationship confirmed
- https://ieftimov.com/posts/docker-compose-stray-pids-rails-beyond/ — stale PID pattern confirmed
- https://github.com/rails/tailwindcss-rails/discussions/499 — arm64/amd64 cross-platform issue, workarounds

### Tertiary (LOW confidence)

- Various DEV Community posts on RSpec/SimpleCov/VCR wiring — patterns cross-verified with official gem READMEs

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all versions verified against RubyGems.org on 2026-04-01
- Architecture: HIGH — patterns derived from locked decisions (CONTEXT.md) + official gem docs
- Pitfalls: HIGH (rspec-rails version, pg gem) / MEDIUM (tailwind arch issue, master.key) — multiple sources

**Research date:** 2026-04-01
**Valid until:** 2026-05-01 (stable ecosystem; 30 days reasonable)
