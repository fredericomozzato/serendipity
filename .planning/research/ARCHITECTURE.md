# Architecture Patterns

**Domain:** Random music discovery — Rails app with Discogs API, Hotwire/Turbo, session-based state
**Project:** Serendipity
**Researched:** 2026-03-31
**Overall confidence:** HIGH (standard Rails patterns, well-documented)

---

## Recommended Architecture

### System Overview

```
Browser
  |
  |-- Turbo Drive (full page nav)
  |-- Turbo Frames (partial swaps: player, history nav)
  |-- Turbo Streams (targeted DOM updates after form/button actions)
  |
Rails App (Docker container)
  |
  |-- Controllers (thin — format detection only, delegate everything)
  |       SongsController
  |       HistoryController (future)
  |
  |-- Services (pure business logic — mutation-testing targets)
  |       Discogs::SearchService         # Fetch random release from Search API
  |       Discogs::RandomPageStrategy    # Compute random page number from total count
  |       Discogs::VideoFilter           # Find first release that has a YouTube video URL
  |       Discogs::ApiClient             # HTTP adapter wrapping Faraday/Net::HTTP
  |
  |-- Query Objects (session reads — pure input/output)
  |       History::Reader                # Read and deserialize session history
  |       History::Writer                # Append/pop within 10-20 item window
  |
  |-- Value Objects (immutable data containers)
  |       Release                        # id, title, artist, year, genre, video_url
  |       HistoryWindow                  # ordered list of Release, max_size boundary
  |
  |-- Views / Partials
  |       songs/show.html.erb            # Shell: player frame + history nav frame
  |       songs/_player.html.erb         # Turbo Frame: YouTube embed + skip button
  |       songs/_history_nav.html.erb    # Turbo Frame: back button + genre selector
  |       songs/show.turbo_stream.erb    # Stream template: update player + history on skip
  |
  |-- Helpers / Presenters
  |       ReleasePresenter               # Formatting only, no logic
```

---

## Component Boundaries

| Component | Responsibility | Communicates With | Must NOT Touch |
|-----------|---------------|-------------------|----------------|
| `SongsController` | Resolve format, call service, render | Services, session, views | Discogs API directly, business logic |
| `Discogs::ApiClient` | HTTP to Discogs, auth headers, error wrapping | External API only | Rails session, models |
| `Discogs::RandomPageStrategy` | Given total count and optional genre, return random page + index | `ApiClient` (via service) | HTTP, session |
| `Discogs::SearchService` | Orchestrate: fetch count → compute page → fetch page → pick item | `ApiClient`, `RandomPageStrategy`, `VideoFilter` | Session, views |
| `Discogs::VideoFilter` | Given a list of releases, return first with a video URL (or nil) | Pure data only | HTTP, session |
| `History::Writer` | Given current history array + new release, return updated array (capped) | Pure data only | Session directly (controller writes result) |
| `History::Reader` | Given raw session value, return `HistoryWindow` | Pure data only | HTTP, API |
| `Release` (value object) | Hold release data — immutable struct | Nothing (data container) | Everything |
| `HistoryWindow` (value object) | Hold ordered list of releases with max size | Nothing (data container) | Everything |

**Key boundary rule:** Services and query objects never read from `session` or `params` directly. Controllers extract what is needed and pass it as plain arguments. This makes every service trivially unit-testable and maximises mutant coverage.

---

## Data Flow

### 1. Initial Page Load

```
GET /
  └─> SongsController#show
        ├─ History::Reader.call(session[:history])  → HistoryWindow
        ├─ Discogs::SearchService.call(genre: nil)  → Release (or retry loop)
        ├─ History::Writer.call(history, release)   → new_history (Array)
        ├─ session[:history] = new_history           (controller writes)
        └─ render songs/show.html.erb
                ├─ <turbo-frame id="player">   → _player partial
                └─ <turbo-frame id="history-nav"> → _history_nav partial
```

### 2. Skip (Next Song)

```
POST /songs/skip
  └─> SongsController#skip
        ├─ History::Reader.call(session[:history])     → HistoryWindow
        ├─ Discogs::SearchService.call(genre: params[:genre])  → Release
        ├─ History::Writer.call(history, release)      → new_history
        ├─ session[:history] = new_history
        └─ respond_to :turbo_stream
              └─ render songs/show.turbo_stream.erb
                    ├─ turbo_stream.replace("player", partial: "_player")
                    └─ turbo_stream.replace("history-nav", partial: "_history_nav")
```

No WebSockets. No ActionCable. Pure request-response Turbo Streams — the controller responds with `Content-Type: text/vnd.turbo-stream.html` and the browser applies targeted DOM updates.

### 3. Back Navigation

```
POST /songs/back
  └─> SongsController#back
        ├─ History::Reader.call(session[:history])  → HistoryWindow
        ├─ HistoryWindow#pop                        → [prev_release, new_history]
        ├─ session[:history] = new_history
        └─ respond_to :turbo_stream
              └─ render songs/show.turbo_stream.erb (same template, different data)
```

### 4. Genre Filter Change

```
GET /songs?genre=jazz  (Turbo Drive navigation, full page replace)
  └─> SongsController#show
        ├─ genre = params[:genre]   (whitelist against Discogs taxonomy)
        ├─ Discogs::SearchService.call(genre: genre) → Release
        ├─ session[:history] = []   (genre change resets history)
        └─ render songs/show.html.erb
```

Genre change is a clean page navigation, not a stream. This avoids managing history state across genre boundaries and keeps the controller simple.

---

## Session Data Shape

Rails cookie store (CookieStore, default) is used. No additional session backend needed for V1.

```ruby
# session[:history] shape
[
  {
    "id"        => "12345",
    "title"     => "Kind of Blue",
    "artist"    => "Miles Davis",
    "year"      => 1959,
    "genre"     => "Jazz",
    "video_url" => "https://www.youtube.com/watch?v=..."
  },
  # ... up to 20 entries
]
```

**Cookie size constraint:** Rails CookieStore is capped at 4096 bytes. With 20 entries at ~200 bytes each (serialized) that is ~4000 bytes — tight. Mitigation: store only the 5 fields above (no full Discogs JSON), cap history at 15 entries in V1, monitor session size in tests.

`History::Writer` enforces the cap as a pure method. The maximum size is a constant passed in, not hardcoded, making it trivially testable and mutable.

---

## Patterns to Follow

### Pattern 1: Thin Controller, Fat Service

**What:** Controller does format detection, session read/write, and delegates to services. No conditional logic, no API calls.

**Why it matters for mutant:** Mutant targets instance and singleton methods. Logic in controllers lives inside `before_action` callbacks and action methods that are hard to isolate. Services with `.call` class methods are easy to test in isolation with no Rails boot.

**Structure:**
```ruby
# app/services/discogs/search_service.rb
module Discogs
  class SearchService
    def self.call(genre: nil)
      new(genre: genre).call
    end

    def initialize(genre:)
      @genre = genre
    end

    def call
      total   = api_client.fetch_total(genre: @genre)
      page    = RandomPageStrategy.call(total: total)
      results = api_client.fetch_page(page: page, genre: @genre)
      VideoFilter.call(releases: results)
    end

    private

    def api_client
      @api_client ||= ApiClient.new
    end
  end
end
```

### Pattern 2: Value Objects for Release Data

**What:** `Release` is a plain Ruby struct (or `Data.define` in Ruby 3.2+) with no methods beyond accessors.

**Why it matters for mutant:** Value objects have zero branching logic by design, keeping mutant scope focused on services and query objects where decisions actually live.

```ruby
# app/value_objects/release.rb
Release = Data.define(:id, :title, :artist, :year, :genre, :video_url)
```

### Pattern 3: Pure Methods for Session History

**What:** `History::Writer` and `History::Reader` take plain data in, return plain data out. They never touch `session` directly.

**Why it matters for mutant:** Purity means the full behaviour is expressible in unit tests with no mocks. Every mutation in the append/cap/pop logic is caught by a simple array assertion.

```ruby
# app/services/history/writer.rb
module History
  class Writer
    MAX_SIZE = 15

    def self.call(history:, release:, max_size: MAX_SIZE)
      entries = history.dup
      entries.unshift(release.to_h.stringify_keys)
      entries.first(max_size)
    end
  end
end
```

### Pattern 4: Adapter Pattern for Discogs HTTP

**What:** `Discogs::ApiClient` is the single place that knows about HTTP, authentication headers, and the Discogs base URL. All other services call it as a dependency.

**Why:** Swappable in tests (inject a stub client), easy to mock for mutant runs, and isolated when the Discogs API shape changes.

```ruby
# app/services/discogs/api_client.rb
module Discogs
  class ApiClient
    BASE_URL = "https://api.discogs.com"

    def initialize(token: ENV.fetch("DISCOGS_TOKEN"))
      @token = token
    end

    def fetch_total(genre: nil)
      response = get("/database/search", genre: genre)
      response.dig("pagination", "items")
    end

    def fetch_page(page:, genre: nil)
      response = get("/database/search", page: page, per_page: 50, genre: genre)
      response["results"].map { |r| Release.from_discogs(r) }
    end

    private

    def get(path, params = {})
      # Faraday or Net::HTTP call, raise on non-200
    end
  end
end
```

### Pattern 5: Turbo Stream Templates (No Inline Rendering)

**What:** Skip and back actions render a `.turbo_stream.erb` template, not inline `render turbo_stream:` in the controller.

**Why:** Keeps controllers free of view DSL. Templates are the right place for multiple stream operations. Easier to test with `assert_turbo_stream` helpers in RSpec.

```ruby
# app/controllers/songs_controller.rb
def skip
  # ... service call, session write ...
  respond_to do |format|
    format.turbo_stream   # renders songs/skip.turbo_stream.erb
    format.html { redirect_to root_path }
  end
end
```

---

## Anti-Patterns to Avoid

### Anti-Pattern 1: Logic in Controllers

**What:** Putting retry loops, genre filtering, or page calculation inside a controller action.

**Why bad:** Untestable without full request stack. Mutant cannot get meaningful coverage. Controllers balloon in size.

**Instead:** Every decision lives in a service or value object. Controllers only call `.call` and assign the result.

### Anti-Pattern 2: Reading session Inside a Service

**What:** Passing `session` object into a service so it can read history itself.

**Why bad:** Couples pure logic to Rails internals. Makes unit tests require a Rails session mock. Kills mutant coverage efficiency.

**Instead:** Controller reads `session[:history]`, passes the raw array to the service, service returns new array, controller writes it back.

### Anti-Pattern 3: God Service (DiscogsService Doing Everything)

**What:** One service class that fetches count, computes page, fetches results, filters by video, and builds history.

**Why bad:** Every public method becomes a mutation target sharing the same collaborator setup. Tests become integration tests by default.

**Instead:** Decompose into `SearchService` (orchestrator), `RandomPageStrategy`, `VideoFilter`, `ApiClient`. Each is independently testable.

### Anti-Pattern 4: ActionCable for Skip/Back

**What:** Using WebSockets and Action Cable to push new song data.

**Why bad:** Adds Redis dependency, increases infrastructure complexity, unnecessary for a synchronous user-initiated action.

**Instead:** Turbo Streams over request-response (the controller responds with `turbo_stream` format). The skip/back button submits a form; the response is a stream of DOM patches. No WebSockets needed.

### Anti-Pattern 5: Storing Full Discogs JSON in Session

**What:** Dumping the raw API response object into `session[:history]`.

**Why bad:** Discogs release objects contain large nested structures. The 4096-byte cookie limit is exceeded quickly.

**Instead:** Serialize only the 5–6 fields needed for playback and display into the session. `Release#to_session_hash` returns only the required keys.

---

## Directory Structure

```
app/
  controllers/
    songs_controller.rb          # show, skip, back
    application_controller.rb
  services/
    discogs/
      api_client.rb
      search_service.rb
      random_page_strategy.rb
      video_filter.rb
    history/
      reader.rb
      writer.rb
  value_objects/
    release.rb
    history_window.rb
  views/
    songs/
      show.html.erb
      skip.turbo_stream.erb
      back.turbo_stream.erb
      _player.html.erb
      _history_nav.html.erb
      _genre_selector.html.erb
  helpers/
    release_presenter.rb         # formatting only

spec/
  services/
    discogs/
      api_client_spec.rb
      search_service_spec.rb
      random_page_strategy_spec.rb
      video_filter_spec.rb
    history/
      reader_spec.rb
      writer_spec.rb
  value_objects/
    release_spec.rb
    history_window_spec.rb
  requests/
    songs_spec.rb                # request specs for controller integration

.mutant.yml                      # scope: Services::, ValueObjects::, History::
```

---

## Suggested Build Order

Dependencies dictate this sequence. Each step produces something the next step consumes.

### Step 1: Value Objects
`Release` and `HistoryWindow` have zero dependencies. Build first. All other components depend on them.

### Step 2: ApiClient
Wraps HTTP, knows Discogs auth, returns raw hashes. No application logic. Can be tested with VCR cassettes or stub HTTP.

### Step 3: RandomPageStrategy + VideoFilter
Pure functions, depend only on value objects. Build and fully test in isolation before touching orchestration.

### Step 4: SearchService
Orchestrates `ApiClient` + `RandomPageStrategy` + `VideoFilter`. Integration-test with a stubbed `ApiClient`.

### Step 5: History::Reader + History::Writer
Pure functions on arrays. Zero external dependencies. Build alongside or after value objects.

### Step 6: SongsController + Views
Can only be assembled once services exist. Thin shell: reads session, calls service, writes session, renders.

### Step 7: Turbo Stream Templates
Wire up `skip.turbo_stream.erb` and `back.turbo_stream.erb` once the controller actions are functional.

### Step 8: Genre Filter UI
Build genre selector as a form that drives a Turbo Drive navigation to `/?genre=X`. No new controller logic needed — `show` already accepts the param.

---

## Scalability Considerations

| Concern | V1 (No DB, Session Only) | V2 Triggers |
|---------|--------------------------|-------------|
| Discogs rate limits | Personal token (3× rate limit); retry on 429 in `ApiClient` | Cache popular genre counts in Redis |
| Session size | 15 entries × ~200 bytes = ~3 KB (within 4 KB limit) | Move history to DB-backed session when auth added |
| PostgreSQL readiness | Docker Compose defines `db` service (commented out); `database.yml` has `adapter: postgresql` with `DATABASE_URL` env var; `DATABASE_URL` not set in V1 | Un-comment db service, set env var |
| Test suite speed | Pure unit tests + request specs; no DB means no fixtures to load | Add factory_bot + database_cleaner when DB introduced |

---

## Docker Compose: V1 → V2 Transition Pattern

```yaml
# docker-compose.yml (V1 — db service present but inactive)
version: "3.9"
services:
  web:
    build: .
    ports:
      - "3000:3000"
    environment:
      DISCOGS_TOKEN: ${DISCOGS_TOKEN}
      RAILS_ENV: development
      # DATABASE_URL: postgres://postgres:password@db:5432/serendipity_development
    # depends_on:
    #   - db

  # db:                          # Un-comment for V2
  #   image: postgres:16
  #   environment:
  #     POSTGRES_PASSWORD: password
  #   volumes:
  #     - postgres_data:/var/lib/postgresql/data

# volumes:
#   postgres_data:
```

`config/database.yml` always references `DATABASE_URL`, which is simply unset in V1. When V2 adds the `db` service, no application code changes — only Docker Compose and environment variable configuration changes.

---

## Mutation-Testing-Friendly Architecture Summary

The mutant gem targets instance and singleton methods. Coverage is maximised by:

| Principle | Implementation |
|-----------|---------------|
| No logic in controllers | Controllers only call `.call`, write session, render |
| Small methods with one decision | `RandomPageStrategy.call` is 2-3 lines; `VideoFilter.call` is one `find` |
| Pure input/output | History services take arrays, return arrays; no side effects |
| `.mutant.yml` scope | Target `Discogs::`, `History::`, `ValueObjects::` — exclude controllers and views |
| No conditionals in value objects | `Release = Data.define(...)` has nothing for mutant to kill |
| Explicit max constants | `MAX_SIZE = 15` passed as argument enables boundary mutation tests |

Mutant configuration scope:
```yaml
# .mutant.yml
integration: rspec
includes:
  - lib
requires:
  - serendipity
subjects:
  - Discogs*
  - History*
  - Release
  - HistoryWindow
```

---

## Sources

- [Hotwire Turbo Streams Handbook](https://turbo.hotwired.dev/handbook/streams) — official, authoritative
- [turbo-rails GitHub](https://github.com/hotwired/turbo-rails) — request-response stream pattern without ActionCable
- [Turbo Streams on Rails — Colby.so](https://www.colby.so/posts/turbo-streams-on-rails) — controller respond_to pattern
- [mutant GitHub](https://github.com/mbj/mutant) — scope configuration, RSpec integration
- [Rails Service Objects — OneUptime Blog 2026](https://oneuptime.com/blog/post/2026-01-26-rails-service-objects/view) — current naming/structure conventions
- [Query Object Pattern in Rails — DEV Community](https://dev.to/vladhilko/how-to-implement-query-object-pattern-in-ruby-on-rails-59fn) — `.call` pattern
- [Discogs Ruby Wrapper](https://github.com/buntine/discogs) — existing Ruby gem for Discogs API
- [Discogs API Documentation](https://www.discogs.com/developers) — Search API endpoint structure
- [Rails CookieStore — Rails API](https://api.rubyonrails.org/classes/ActionDispatch/Session/CookieStore.html) — 4096-byte limit
- [TurboCable — no external dependencies](https://intertwingly.net/blog/2025/11/04/TurboCable.html) — confirms request-response streams viable without ActionCable
