# Domain Pitfalls

**Domain:** Rails music discovery app — Discogs API + YouTube embed + Hotwire + mutant + Docker Compose
**Researched:** 2026-03-31
**Confidence:** MEDIUM (WebSearch-verified; Discogs/Hotwire official docs and mutant GitHub confirmed patterns)

---

## Critical Pitfalls

Mistakes that cause rewrites or behaviour that fundamentally breaks the core value proposition ("a song is always playing within seconds").

---

### Pitfall 1: YouTube Autoplay Blocked by Browser Policy

**What goes wrong:**
The app's core UX is instant playback on page load. Modern browsers silently block autoplay with sound unless the user has previously interacted with the page. iOS Safari blocks all autoplay (even muted) unless `playsinline` is also set. The result is a blank player — the "song is always playing" promise fails on first visit for most users.

**Why it happens:**
Browser autoplay policies have tightened since 2018 and are now enforced by Chrome, Firefox, Edge, and Safari. The YouTube iframe will render without error but the video will not start. The `allow="autoplay"` iframe attribute is also required at the browser permissions level — missing it silently blocks playback even when URL parameters are correct.

**Consequences:**
- Silent failure: no JS error thrown, just no playback
- iOS users (primary mobile audience) see a fully broken experience unless `playsinline=1` is in the embed URL
- Testing in a desktop browser during development often masks the issue because desktop Chrome is more permissive

**Prevention:**
- Always embed with: `?autoplay=1&mute=1&playsinline=1` in the src URL
- Always set the iframe attribute: `allow="autoplay; encrypted-media"`
- On mobile the video will play silently on load (muted); design the UI to show a clear "tap to unmute" affordance
- Test on real iOS Safari early — simulator behaviour can differ

**Warning signs:**
- Desktop dev works but mobile QA shows static player
- Playwright/Capybara tests pass but manual iOS test fails
- No JS console errors despite no playback (silent browser policy enforcement)

**Phase:** Address in Phase 1 (core playback). Do not defer to polish.

---

### Pitfall 2: Discogs Search API 10,000-Result Pagination Cap with Skewed Randomness

**What goes wrong:**
The strategy is: fetch total result count → pick random page → pick random item. Discogs caps pageable results at 10,000 items (100 pages × 100 per page), regardless of how many total matches exist. If a broad search returns 500,000 releases, `pagination.items` will report the full count but pages beyond 100 are inaccessible (404 or empty). Generating a random page from the full `items` count will frequently land on unreachable pages.

**Why it happens:**
Discogs applies a hard pagination wall for search endpoints at page 100 (with `per_page=100`). The `pagination.items` field reports the true total, not the reachable total. This is documented in community threads but not prominently in the official API docs.

**Consequences:**
- Random page calculation using raw `items` count silently fails ~99% of the time for popular genres (e.g., Rock has millions of entries)
- The app silently retries forever (or exhausts rate limit budget) rather than finding a song
- Even with correct page capping, only the first 10,000 results of any genre are accessible, skewing selection toward Discogs-ranked "top" results

**Prevention:**
- Cap the random page calculation: `rand(1..[pagination.pages, 100].min)` — never trust `pagination.items` directly for page selection
- Verify `per_page=100` is sent on every search request to maximise reachable randomness
- Accept the top-10,000 bias as a product reality in V1; document it explicitly in code comments

**Warning signs:**
- Retry loop hits rate limit before finding a video
- Manual API calls to page 101+ return 404 or empty results
- `pagination.pages` value exceeds 100 in API response

**Phase:** Address in Phase 1 (Discogs integration). Critical path item.

---

### Pitfall 3: Discogs Rate Limit Exhaustion During Silent Retry Loop

**What goes wrong:**
The project requires silently retrying when a release has no YouTube video (~60–80% of releases). Each retry costs one API call. With authenticated token, the limit is 60 requests/minute. An aggressive retry loop can exhaust the budget in under 10 seconds for a single user, returning 429 errors. Without backoff and `X-Discogs-Ratelimit-Remaining` header inspection, the app will throw uncaught errors or hang.

**Why it happens:**
- Unauthenticated requests: 25/min. Token-authenticated: 60/min (3× boost)
- Rate limiting is IP-based for unauthenticated, token-based for authenticated
- The retry loop for "no video found" releases has no natural upper bound unless explicitly constrained
- 429 responses are not always handled gracefully in naive HTTP client wrappers

**Consequences:**
- Rapid successive skips (user mashing skip) can exhaust the rate budget entirely
- In production with multiple concurrent users sharing one token, one heavy user can deny service to all
- Unhandled 429s crash the request or return a 500 to the user

**Prevention:**
- Always send a unique, descriptive `User-Agent` header (required by Discogs ToS; generic agents get lower limits)
- Inspect `X-Discogs-Ratelimit-Remaining` header on every response; pause or return an error when it drops below a safe threshold (e.g., 5)
- Hard-cap the retry loop at 5 attempts per user action; surface a graceful "couldn't find a track, try again" message rather than looping infinitely
- Wrap all Discogs HTTP calls in a dedicated service object that enforces these rules — never call the API from a controller directly

**Warning signs:**
- 429 responses in Rails logs during testing with rapid skip clicks
- Rate limit remaining header drops to zero in development even with one user
- Retry loop causes request timeouts (>10 seconds) for skip actions

**Phase:** Address in Phase 1. The rate-limiting wrapper is foundational — build it before building features on top.

---

## Moderate Pitfalls

---

### Pitfall 4: Hotwire Turbo + Browser Back Button Breaks Session History

**What goes wrong:**
The back-navigation requirement (10–20 songs in session history) is implemented via server-side session state. When the user hits the browser's native Back button, Turbo Drive performs a "restoration visit" — it may serve a cached version of the page from Turbo's page cache rather than fetching the correct historical track from session. The displayed song can be wrong, or the session array pointer can desync from the browser history stack.

**Why it happens:**
- Turbo Drive maintains its own page cache separate from session state
- Manually calling `history.pushState` (to update the URL on skip) conflicts with Turbo's internal restoration identifiers; Turbo's documentation explicitly warns against this
- Turbo Frames do not push history entries by default — using a frame for the player without `data-turbo-action="advance"` means the back button never navigates within the player at all

**Consequences:**
- Back button shows stale UI (cached page) but the session cursor is in a different position
- Rapid skip + back sequences can corrupt the session history array
- URL may not reflect the current track, making deep-linking or sharing impossible

**Prevention:**
- Use Turbo Drive (full page visits with `data-turbo-action="advance"`) for skip/back navigation, not Turbo Frames, so Turbo handles all history entries natively
- Store history as an ordered array in the Rails session; use `params` or the URL to signal "go back N steps" rather than relying on Turbo's cache to restore the right state
- If using Turbo Frames, explicitly set `data-turbo-action="advance"` on the frame wrapper and test browser Back extensively
- Never call `history.pushState` manually; let Turbo own all history manipulation

**Warning signs:**
- Back button shows the correct URL but wrong song content
- Session array length grows unexpectedly (duplicate pushes)
- Browser forward/back works in development but breaks after deploy (caching difference)

**Phase:** Address in Phase 2 (navigation UX). Design the session history data structure before implementing skip/back.

---

### Pitfall 5: mutant Cannot Discover Subjects — Naming and Load Path Mismatches

**What goes wrong:**
mutant identifies subjects by matching Ruby constant names to example group description strings in RSpec. If the subject constant path does not match the spec description exactly, mutant reports "0 subjects" and exits silently without testing anything. This is the most common first-run failure when integrating mutant into a Rails app.

**Why it happens:**
- Rails autoloading means constants may be defined at runtime but mutant needs to `require` them explicitly via `--require` flags
- Rails apps do not use a simple `lib/` structure; mutant's default `--include lib` path may not cover `app/services`, `app/models`, etc.
- Controllers, helpers, and views are rarely good mutant subjects — they are too tightly coupled to Rails machinery; mutant works best against pure domain objects
- Leaving `require 'rspec/autorun'` in `spec_helper.rb` causes double-run conflicts with mutant's test runner

**Consequences:**
- `bundle exec mutant run` exits with 0 mutations killed, giving false confidence
- CI passes mutation coverage gates because 0/0 = 100% by some interpretations
- Time is wasted debugging mutant configuration rather than writing tests

**Prevention:**
- Remove or guard `require 'rspec/autorun'` from `spec_helper.rb` before running mutant
- Use explicit `--include app --require config/environment` (or a dedicated mutant bootstrap file) for Rails apps
- Scope mutant subjects to domain service objects in `app/services/` and plain Ruby classes — not controllers or views
- Design the codebase with mutant in mind from the start: small, single-responsibility, pure methods with no side effects are easy mutant subjects; stateful, I/O-heavy methods are not
- Run `bundle exec mutant run --zombie` as a sanity check: if it exits with 0 subjects, the load path or constant naming is wrong

**Warning signs:**
- `mutant run` output shows "0 subjects selected"
- Spec descriptions use informal names ("describes the searcher") rather than matching the Ruby constant (`Serendipity::Discogs::Searcher`)
- All logic lives in controllers rather than service objects

**Phase:** Address in Phase 1 setup. The architecture must support mutant before any logic is written. Retrofitting is painful.

---

### Pitfall 6: Docker Compose Gem Volume Conflicts Break bundle install

**What goes wrong:**
When a host-mounted volume (e.g., `.:/app`) overlays the container's gem install path, `bundle install` during image build is invisible at runtime — the volume mount wipes out everything installed into the image. Alternatively, mounting `/usr/local/bundle` directly destroys Bundler itself (it installs to that exact path). The result is `bundle exec rails server` failing with `cannot load such file` errors even though the image built successfully.

**Why it happens:**
- Docker volumes are applied at container start, after image build — gems installed during `docker build` can be shadowed by a mounted empty directory
- The standard pattern of `VOLUME /usr/local/bundle` in Dockerfile or mounting the bundle path directly is a documented footgun
- Rails Docker Compose setups without a database are less commonly documented, making it easy to copy patterns from Postgres-inclusive examples that handle volumes differently

**Consequences:**
- `docker compose up` works on first pull but breaks after `git clean` or on a fresh clone
- New gems added to Gemfile appear to install but are missing at runtime
- Debugging requires understanding the build vs runtime volume layering model

**Prevention:**
- Use a **named volume for gems** (`bundle_cache:/usr/local/bundle`) separate from the application code volume (`.:app`)
- Do NOT mount the app directory at `/usr/local/bundle` — keep gem storage separate
- Structure `docker-compose.yml` so the database service can be added in V2 by simply appending a `db` service and a `DATABASE_URL` env var — document this explicitly in a comment in the compose file
- Add a `BUNDLE_PATH` env var pointing to the named volume path so Bundler is consistent between build and runtime

**Warning signs:**
- `bundle exec rails` works inside `docker compose run` but not `docker compose up`
- `gem list` shows fewer gems than Gemfile requires
- Adding a new gem requires `docker compose down -v` + full rebuild to take effect

**Phase:** Address in Phase 0 (infrastructure setup). Getting this wrong means debugging environment issues instead of building features.

---

## Minor Pitfalls

---

### Pitfall 7: Discogs `videos` Field Absent on Master Releases vs. Release Versions

**What goes wrong:**
The Discogs search API returns a mix of master releases and specific release versions. The `videos` field is more likely to be populated on specific release versions than on master releases. Fetching a master release and checking `videos` will frequently return an empty array even when YouTube links exist on sub-releases. This increases the effective retry rate above the expected ~60–80%.

**Prevention:**
- When the search result's `type` field is `"master"`, prefer fetching the master's main release version (`/masters/{id}/main_release`) to find video links, or treat master results as lower priority and retry for a `type: "release"` result
- Log `type` in retry telemetry during development to understand the actual video hit rate

**Phase:** Phase 1 refinement, after initial integration is working.

---

### Pitfall 8: Session Size Limits Break History Array on Long Sessions

**What goes wrong:**
Rails cookie-based sessions have a 4KB limit. Storing 20 full track objects (title, artist, YouTube ID, Discogs URL, genre) in session can exceed this limit. The session silently fails to write — no error is raised, but the history array stops growing, breaking back navigation.

**Prevention:**
- Store only the minimum required data per history entry: `{ youtube_id:, title:, artist: }` — not full API response objects
- Add a session size guard in the history service: measure before writing, drop oldest entries if approaching 3KB
- Consider storing history as a compact array of IDs only and re-fetching on back navigation (one extra API call, but safe)

**Phase:** Phase 2 (back navigation). Design the session data structure with this constraint in mind.

---

### Pitfall 9: Genre Filter UI Shows All 13 Discogs Genres But Not All Have Equal Coverage

**What goes wrong:**
Discogs has ~13 official genres but their catalog coverage is highly uneven. Selecting "Brass & Military" or "Stage & Screen" with a random-page strategy may consistently return releases with no YouTube video, causing the retry loop to exhaust its budget. From a UX perspective, this looks like a broken genre filter.

**Prevention:**
- Document the known low-coverage genres in code comments
- Consider capping retries lower (3 instead of 5) for genre-filtered requests and surfacing a "Couldn't find a video for this genre, try again" message faster
- V1 ships all 13 genres but this is a known rough edge

**Phase:** Phase 2 (genre filter). Note in code, not a blocker.

---

## Phase-Specific Warnings

| Phase Topic | Likely Pitfall | Mitigation |
|---|---|---|
| Phase 0: Docker infra | Gem volume conflicts (`bundle_cache` vs app mount) | Use named volume for gems; add DB-readiness comment to compose file |
| Phase 1: Discogs integration | Pagination cap at page 100; rate limit exhaustion in retry loop; missing User-Agent | Cap page calculation; wrap API calls in dedicated rate-aware service; set User-Agent header |
| Phase 1: YouTube embed | Autoplay blocked on mobile; missing `playsinline` for iOS | Embed with `mute=1&playsinline=1`; `allow="autoplay"` attribute; test on real iOS early |
| Phase 1: mutant setup | Zero subjects discovered; `rspec/autorun` conflict; load path wrong | Remove autorun; use `--include app`; design for service objects from day one |
| Phase 2: Skip/back navigation | Turbo cache desync with session history; manual `pushState` conflicts with Turbo | Use Turbo Drive `advance` action; let Turbo own history; store only IDs in session |
| Phase 2: Genre filter | Low-video-coverage genres exhaust retry budget | Cap retries at 3 for genre requests; surface friendly error faster |
| Phase 2: Session history | 4KB cookie limit breaks history array | Store minimal data per entry; add size guard in history service |

---

## Sources

- [Discogs API Rate Limiting — Community Forum](https://www.discogs.com/forum/thread/997721) (MEDIUM confidence — community-verified)
- [Discogs API Rate Limits Thread](https://www.discogs.com/forum/thread/1104957) (MEDIUM confidence)
- [Discogs 10,000 Result Cap Thread](https://www.discogs.com/forum/thread/996056) (MEDIUM confidence — multiple community confirmations)
- [Discogs User-Agent Rate Limiting Thread](https://www.discogs.com/forum/thread/740651) (MEDIUM confidence)
- [YouTube Embedded Players and Player Parameters — Google Developers](https://developers.google.com/youtube/player_parameters) (HIGH confidence — official docs)
- [YouTube Iframe Autoplay Fix — technetexperts](https://www.technetexperts.com/youtube-iframe-ios-autoplay-fix/) (MEDIUM confidence — verified against official docs)
- [Hotwire Turbo Drive Handbook](https://turbo.hotwired.dev/handbook/drive) (HIGH confidence — official docs)
- [Turbo history.pushState manual override issue #792](https://github.com/hotwired/turbo/issues/792) (HIGH confidence — official repo)
- [mutant — GitHub](https://github.com/mbj/mutant) (HIGH confidence — official source)
- [mutant-rspec docs](https://github.com/mbj/mutant/blob/main/docs/mutant-rspec.md) (HIGH confidence — official source)
- [Docker Compose bundle install volume pitfalls — bradgessler.com](https://bradgessler.com/articles/docker-bundler/) (MEDIUM confidence)
- [Rails Docker Compose bundle install — anonoz.github.io](https://anonoz.github.io/tech/2021/01/10/rails-docker-compose-yml.html) (MEDIUM confidence)
- [Turbo and Rails Common Pitfalls — Reintech](https://reintech.io/blog/turbo-rails-ajax-application-pitfalls) (LOW confidence — single source, verify against Turbo docs)
- [Mutation Testing Introduction — Nebulab](https://nebulab.com/blog/an-introduction-to-mutation-testing) (MEDIUM confidence)
