## 1. Genre Data and Backend Support

- [x] 1.1 Create `config/discogs_genres.yml` with the canonical list of Discogs main genres (e.g., Rock, Electronic, Jazz, Hip Hop, Blues, Classical, Folk, World & Country, Funk / Soul, Latin, Reggae, Pop, Stage & Screen, Non-Music, Children's, Brass & Military)
- [x] 1.2 Add `DiscogsService.genres` class method that reads and returns the genre list from `config/discogs_genres.yml`
- [x] 1.3 Update `DiscogsService#fetch_random_release` to accept an optional `genre` parameter and pass it to the Discogs API search query string when present
- [x] 1.4 Update `DiscogsService#random_playable_track` to accept an optional array of `genres` and, when provided, select one at random to pass to `fetch_random_release`
- [x] 1.5 Update `PlayerController#random_track` to read `genres` from query params (as an array) and pass them to `DiscogsService.random_playable_track`

## 2. Stimulus Controllers

- [x] 2.1 Generate `app/javascript/controllers/player_controller.js` to manage the YouTube embed, track info display, and next/previous navigation (refactored from existing vanilla JS)
- [x] 2.2 Generate `app/javascript/controllers/genre_filter_controller.js` to manage genre tag toggling, active filter state, and emitting filter-change events
- [x] 2.3 Generate `app/javascript/controllers/history_controller.js` to manage `localStorage` read/write for both play history and active genre filters
- [x] 2.4 Ensure `history_controller` restores active filters from `localStorage` on page load and emits an event so `genre_filter_controller` can render tags correctly
- [x] 2.5 Ensure `player_controller` listens for filter changes and includes active genres as query params when calling `/player/random_track`
- [x] 2.6 Remove or deprecate existing vanilla JavaScript player logic from `app/views/player/index.html.erb`

## 3. Player View and UI

- [x] 3.1 Update `app/views/player/index.html.erb` to render the genre filter tag bar above the player using the genres from `DiscogsService.genres`
- [x] 3.2 Style inactive genre tags as small rounded-rectangles with an off-white background and black text; active genre tags with a blue background and white text (use Tailwind or existing CSS framework)
- [x] 3.3 Add Stimulus `data-controller` attributes to mount `player_controller`, `genre_filter_controller`, and `history_controller` on the player page
- [x] 3.4 Ensure the YouTube iframe, track info, and next/previous buttons remain functional after the Stimulus refactor

## 4. Routes and Configuration

- [x] 4.1 Add a root route redirect in `config/routes.rb` so `/` redirects to `/player`
- [x] 4.2 Verify `config/routes.rb` does not require other changes (existing `player/random_track` route supports query params)
- [x] 4.3 Verify Stimulus controllers are loaded correctly via `config/importmap.rb` or existing Rails 7+ default setup

## 5. Manual Verification (No Automated Tests)

- [x] 5.1 Manually verify `DiscogsService` correctly appends the `genre` parameter to the Discogs API search URL when a genre is provided
- [x] 5.2 Manually verify `PlayerController#random_track` correctly parses and passes multiple genres from query params
- [x] 5.3 Manually verify genre tags render in the correct active/inactive state on page load
- [x] 5.4 Manually verify active genre filters persist across page refreshes via `localStorage`
- [x] 5.5 Manually verify the "Next" button fetches tracks respecting active genre filters
- [x] 5.6 Manually verify play history (previous/next) continues to work after the Stimulus refactor
- [x] 5.7 Do NOT run or write automated tests. This is an MVP session — manual checks only.

## 6. Post-Verification Fixes

- [x] 6.1 Fix `history_controller.js` connection order bug: defer `history:restored` dispatch with `requestAnimationFrame` so sibling controllers (`player_controller`, `genre_filter_controller`) finish connecting before the event fires
- [x] 6.2 Fix `player_controller.js` autoplay: append `?autoplay=1` to YouTube embed URL so videos start playing immediately when loaded
- [x] 6.3 Fix `app/views/player/index.html.erb` iframe `src=""`: remove empty `src` attribute to prevent loading the current page inside the iframe before any track is fetched
- [x] 6.4 Fix `player_controller.js` prev button state on page load: update `prevButtonTarget.disabled` in `handleHistoryRestored` so the Previous button correctly reflects history state after refresh
