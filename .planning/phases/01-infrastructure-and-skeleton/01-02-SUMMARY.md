---
phase: 01-infrastructure-and-skeleton
plan: 02
subsystem: testing
tags: [rspec, rspec-rails, vcr, webmock, simplecov, docker, rails, postgresql]

# Dependency graph
requires:
  - phase: 01-01
    provides: Rails 8.1.3 skeleton with all test gems in Gemfile (rspec-rails 8.0, vcr, webmock, simplecov)
provides:
  - RSpec installed and configured inside Docker container
  - spec/spec_helper.rb with SimpleCov as first require (SimpleCov.start 'rails')
  - spec/rails_helper.rb with VCR + WebMock configuration and no-DB fix
  - spec/smoke_spec.rb proving 2 examples, 0 failures
  - .rspec with --require spec_helper and --format documentation
  - SimpleCov coverage report generated at coverage/index.html on each run
  - No PostgreSQL connection during any spec execution
affects:
  - phase-02
  - all subsequent phases that write specs

# Tech tracking
tech-stack:
  added:
    - rspec 3.13 (via rspec-rails 8.0.4)
    - vcr 6.4.0
    - webmock 3.26.2
    - simplecov 0.22.0
  patterns:
    - SimpleCov must be first require in spec_helper.rb (before Rails loads) — prevents 0% coverage
    - NoDbFixtureSetup module overrides before_setup/after_teardown to prevent PG connections in V1
    - rspec-rails 8 MinitestLifecycleAdapter runs before_setup in around(:each) — must be no-op in no-DB setup
    - coverage/ directory gitignored — SimpleCov output never committed

key-files:
  created:
    - spec/spec_helper.rb
    - spec/rails_helper.rb
    - spec/smoke_spec.rb
    - .rspec
  modified:
    - .gitignore

key-decisions:
  - "NoDbFixtureSetup module overrides before_setup + after_teardown in rails_helper.rb — rspec-rails 8 always includes ActiveRecord::TestFixtures via FixtureSupport regardless of use_active_record setting; both lifecycle hooks must be no-ops in V1"
  - "use_active_record = false kept in config — belt-and-suspenders, prevents fixture paths from being set"
  - "SimpleCov start 'rails' profile — excludes config/, db/, vendor/ automatically; correct default for Rails apps"

patterns-established:
  - "NoDbFixtureSetup: override before_setup + after_teardown when Rails app has no DB in test env"
  - "SimpleCov first: require 'simplecov'; SimpleCov.start 'rails' must be lines 1-2 of spec_helper.rb"
  - "VCR default record mode :new_episodes — records real HTTP on first run, replays from cassette thereafter"

requirements-completed: [INFRA-01, INFRA-03]

# Metrics
duration: 15min
completed: 2026-04-01
---

# Phase 01 Plan 02: RSpec + VCR + WebMock + SimpleCov Summary

**RSpec 3.13 running green inside Docker with VCR/WebMock/SimpleCov wired, smoke spec passing 2 examples 0 failures, no PostgreSQL connection in any test**

## Performance

- **Duration:** ~15 min
- **Started:** 2026-04-01T12:45:00Z
- **Completed:** 2026-04-01T12:57:00Z
- **Tasks:** 2
- **Files modified:** 5 (spec_helper.rb, rails_helper.rb, smoke_spec.rb, .rspec, .gitignore)

## Accomplishments

- Installed RSpec via `rails generate rspec:install` inside Docker, then replaced generated files with exact plan content
- Configured SimpleCov as first require in spec_helper.rb (Pitfall 5 prevention — coverage 10% on minimal app confirms reporting works)
- Configured VCR with `spec/cassettes/` library dir, WebMock disabled net connect, `record: :new_episodes` default
- Fixed rspec-rails 8 + no-DB incompatibility: `MinitestLifecycleAdapter#before_setup` always runs regardless of `use_active_record = false` setting; added `NoDbFixtureSetup` module override
- Verified `docker compose run --rm web bundle exec rspec` exits 0 with 2 examples, 0 failures, coverage/index.html generated

## Task Commits

Each task was committed atomically:

1. **Task 1: Install RSpec and generate configuration files** - `c7d59b9` (feat) — installed via prior session
2. **Task 2: Create smoke spec and verify green suite** - `51650d3` (feat)

**Plan metadata:** (pending)

## Files Created/Modified

- `spec/spec_helper.rb` - SimpleCov.start 'rails' first, standard RSpec config
- `spec/rails_helper.rb` - VCR + WebMock config + NoDbFixtureSetup override for no-DB V1
- `spec/smoke_spec.rb` - 2-example smoke test: `expect(true).to be true` + `expect(defined?(Rails)).to eq("constant")`
- `.rspec` - `--require spec_helper --format documentation`
- `.gitignore` - Added `/coverage/` entry

## Decisions Made

- **NoDbFixtureSetup override:** rspec-rails 8.0.4's `FixtureSupport` unconditionally includes `ActiveRecord::TestFixtures` (line 6 in fixture_support.rb, outside the `use_active_record?` guard). The `MinitestLifecycleAdapter`'s `around(:each)` hook calls `before_setup` on the example instance, which triggers `setup_transactional_fixtures` → `pin_connection!` → PG::ConnectionBad. Fix: override both `before_setup` and `after_teardown` to no-ops via a concern included after `FixtureSupport`.

- **`use_active_record = false` retained:** Belt-and-suspenders — still prevents fixture paths from being configured and `use_transactional_tests` from being set to true in the `included do` block of `FixtureSupport`.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed rspec-rails 8 ActiveRecord::TestFixtures PG::ConnectionBad on every spec**
- **Found during:** Task 2 (running the smoke spec for the first time)
- **Issue:** `rspec-rails 8.0.4` unconditionally includes `ActiveRecord::TestFixtures` into all example groups via `FixtureSupport`. This module overrides `before_setup` to call `setup_transactional_fixtures` → `pin_connection!`, which attempts a PostgreSQL connection. With no DB running (INFRA-03 design), every spec fails with `ActiveRecord::ConnectionNotEstablished`.
- **Fix:** Added `NoDbFixtureSetup` module to `rails_helper.rb` that overrides `before_setup` and `after_teardown` to no-ops. Included after rspec-rails includes `FixtureSupport` so the override wins. Also kept `config.use_active_record = false` and `config.use_transactional_fixtures = false` for belt-and-suspenders protection.
- **Files modified:** `spec/rails_helper.rb`
- **Verification:** `docker compose run --rm web bundle exec rspec` exits 0 with 2 examples, 0 failures; `grep -i "database\|postgres\|PG::" output` returns nothing.
- **Committed in:** `51650d3` (Task 2 commit)

---

**Total deviations:** 1 auto-fixed (Rule 1 - Bug)
**Impact on plan:** Fix was necessary for INFRA-03 compliance (no DB connection in tests). Minimal scope — 12 lines added to rails_helper.rb. No architectural changes.

## Issues Encountered

- rspec-rails 8.0.4 bug/design gap: `use_active_record = false` config option does NOT prevent `ActiveRecord::TestFixtures#before_setup` from running. The `FixtureSupport` module includes `ActiveRecord::TestFixtures` unconditionally at line 6, outside the `if RSpec.configuration.use_active_record?` guard at line 20. This means any Rails app using rspec-rails 8 without a database must use the `NoDbFixtureSetup` override pattern established here.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- `docker compose run --rm web bundle exec rspec` exits 0 with 2 examples, 0 failures
- Phase 1 success criteria all satisfied: app boots at localhost:3000, RSpec suite green, no DB required
- Phase 2 can write specs from day one: VCR cassettes dir ready, WebMock blocks external HTTP, SimpleCov tracks coverage
- The `NoDbFixtureSetup` pattern in `rails_helper.rb` is established for all future phases in V1

## Self-Check: PASSED

- spec/spec_helper.rb: FOUND
- spec/rails_helper.rb: FOUND
- spec/smoke_spec.rb: FOUND
- .rspec: FOUND
- .gitignore /coverage/ entry: FOUND
- Commit c7d59b9 (Task 1): FOUND
- Commit 51650d3 (Task 2): FOUND
- 2 examples, 0 failures: VERIFIED
- coverage/index.html generated: VERIFIED
- No PG connection errors: VERIFIED

---
*Phase: 01-infrastructure-and-skeleton*
*Completed: 2026-04-01*
