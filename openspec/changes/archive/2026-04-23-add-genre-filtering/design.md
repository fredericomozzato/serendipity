## Context

We have a lightweight Rails music discovery app that fetches random releases from the Discogs API and plays YouTube-embedded tracks. The current implementation uses vanilla JavaScript for client-side interactivity and stores play history in browser `localStorage`. There is no database persistence. The Discogs API has a hard rate limit of 25 unauthenticated requests per minute.

## Goals / Non-Goals

**Goals:**
- Display a row of clickable genre filter tags above the player, sourced from a canonical Discogs genre list.
- Allow users to toggle genres on/off. Active tags render as small rounded-rectangles with a blue background and white text; inactive tags render with an off-white background and black text.
- Support additive (OR) filtering: selecting Rock and Electronic includes results from either genre.
- Refactor the player from vanilla JavaScript to Stimulus controllers.
- Persist active genre filters alongside play history in `localStorage`.
- Minimize Discogs API request waste by passing a single genre param per request when filters are active.

**Non-Goals:**
- Database persistence of releases, genres, or user preferences.
- Sub-genre / style filtering (e.g., "Death Metal"). We stick to Discogs main genres only.
- Genre search, autocomplete, or free-text filtering.
- Server-side caching of Discogs responses.
- Changing the YouTube embedding strategy or supporting additional video sources.
- Writing automated tests (unit, integration, or system). This is a fast-paced MVP session; verify features manually only.

## Decisions

### Use a static YAML file for Discogs genres
- **Rationale**: Discogs main genres are stable and well-documented (e.g., Rock, Electronic, Jazz, Hip Hop). Hard-coding them in `config/discogs_genres.yml` avoids an API call, guarantees availability at boot, and simplifies tag rendering.
- **Alternative considered**: Fetching genres dynamically from Discogs API. Rejected because Discogs does not expose a dedicated genre endpoint, and scraping/searching is wasteful given the 25 req/min limit.

### Additive filtering via random genre selection
- **Rationale**: The Discogs search endpoint accepts a single `genre` parameter. To achieve additive OR logic with multiple selected genres, the server picks one genre at random from the active set and queries Discogs with it. Over multiple skips, the user sees a mix from all selected genres.
- **Alternative considered**: Fetching random releases and filtering server-side by genre. Rejected because it wastes API calls on releases that don't match the filter, which is dangerous at 25 req/min.

### Refactor player to Stimulus controllers
- **Rationale**: The user explicitly requested Turbo, Hotwire, and Stimulus for reactive components. Stimulus provides a lightweight, maintainable way to handle DOM events (tag toggles, next/previous buttons, YouTube embed updates) without a heavy JS framework.
- **Alternative considered**: Keeping vanilla JS and only using Stimulus for filters. Rejected because mixing paradigms creates technical debt; a full Stimulus refactor keeps the codebase consistent.

### Persist filters in `localStorage`
- **Rationale**: Filters should survive page refreshes and browser restarts, just like play history. Storing them in the same `localStorage` key space keeps the persistence strategy uniform and avoids backend state.
- **Alternative considered**: Storing filters in the URL query string. Rejected because localStorage is already used for history and provides a cleaner UX without cluttering URLs.

### Root route redirects to player
- **Rationale**: The player is the only meaningful page in the app. Redirecting `/` to `/player` removes friction for users visiting the base URL and makes the player feel like the home page.
- **Alternative considered**: Rendering the player directly at `/` without a redirect. Rejected because it would require renaming the controller/view paths and complicating the route structure; a redirect is simpler and preserves the existing `/player` route.

## Risks / Trade-offs

- **[Risk]** Highly restrictive genre combinations (e.g., one obscure genre) may cause many skipped releases before finding a playable track, burning API requests.
  - **Mitigation**: Cap retries at 5 per request (existing behavior). If no track is found, show a user-friendly message suggesting broader filters.
- **[Risk]** Stimulus refactor could introduce regressions in history navigation or YouTube embedding.
  - **Mitigation**: Maintain the same `localStorage` schema and event handlers. Test next/previous/history-restore scenarios thoroughly.
- **[Risk]** Discogs API rate limit (25 req/min) is shared across all users if deployed to a single server IP.
  - **Mitigation**: This is accepted as a known limitation. The single-genre param approach actually reduces average API usage compared to unfiltered random search because genre-constrained searches have a higher hit rate for releases with videos in popular genres.

## Migration Plan

No migration required. This is a purely additive client-side and controller change with no database schema changes.

## Open Questions

None.
