## Why

The current player only supports completely random track discovery. Users want the ability to filter by music genre to discover tracks that match their mood or taste. Discogs provides genre metadata on every release, but we currently ignore it. Adding genre filtering will improve the user experience while keeping the lightweight, zero-auth nature of the app.

## What Changes

- **Parse and store Discogs genres**: Extract the canonical list of Discogs main genres into a local data file so we can display them as clickable filter tags.
- **Genre filter UI**: Add a tag-based filter component above the player. Tags render as small rounded-rectangles with an off-white background and black text when inactive, and a blue background with white text when active. Clicking toggles the filter state.
- **Additive filtering**: Selecting multiple genres (e.g., Rock and Electronic) includes results from any selected genre (OR logic).
- **Smart API usage**: With the 25 req/min Discogs limit, filtering logic must minimize wasted API calls by selecting a random genre from the active filter set per request.
- **Stimulus refactor**: Migrate the existing vanilla JavaScript player logic (history, YouTube embed, next/previous) into Stimulus controllers to align with the Rails ecosystem and enable reactive filter interactions via Turbo/Hotwire.
- **Random track endpoint update**: The `/player/random_track` endpoint will accept an optional list of active genre filters.
- **Root route redirect**: The root path (`/`) redirects to the player page so the player behaves as the app's home page.
- **No automated tests**: This is a fast-paced experimental MVP session. Do not write automated tests (unit, integration, or system tests). Manual verification only.

## Capabilities

### New Capabilities
- `genre-filtering`: Genre tag UI, toggle interactions, active filter state management, and smart genre-param querying against the Discogs API.

### Modified Capabilities
- `random-music-player`: The random track fetching requirement changes to support an optional genre filter parameter. The client-side player interactivity requirement changes from vanilla JavaScript to Stimulus controllers with Turbo stream support.

## Impact

- `PlayerController#random_track` — accepts genre filter parameters.
- `DiscogsService` — supports genre-biased random release fetching.
- `app/views/player/index.html.erb` — adds genre filter tag bar and Stimulus controller mounts.
- New `app/javascript/controllers/` files — `player_controller`, `genre_filter_controller`, `history_controller`.
- New `config/discogs_genres.yml` — canonical genre list for tag rendering.
- `config/routes.rb` — root route redirects to player.
- Browser `localStorage` — now also persists active genre filters across sessions.
