## Context

We have an existing Rails application and want to add a zero-auth, single-page music discovery feature. Users will land on a dedicated page that fetches random releases from the Discogs API, filters for tracks with YouTube video links, embeds a YouTube player, and allows basic next/previous navigation. No user accounts, database persistence, or advanced filtering is required.

## Goals / Non-Goals

**Goals:**
- Provide a dedicated `/player` route that renders the music player page.
- Fetch random releases from Discogs API and extract a playable track with a valid YouTube video link.
- Automatically skip releases/tracks that lack YouTube links without user intervention.
- Embed a YouTube video player via iframe using the extracted video ID.
- Display the current track title and artist name.
- Allow users to skip to the next random song and return to the previous song.
- Store recent song history in browser `localStorage` so navigation persists across sessions.
- Keep the UI minimal: video player, song info, and navigation buttons only.

**Non-Goals:**
- User authentication or authorization.
- Persistent database storage of play history, favorites, or preferences.
- Search, filtering, or genre selection.
- Volume control, playlists, or shuffle logic beyond pure randomness.
- Handling non-YouTube video sources (e.g., Vimeo, Spotify).
- Caching or proxying Discogs API responses server-side.
- Deploy-specific configuration (e.g., Heroku setup).

## Decisions

### Use server-side Discogs API calls from Rails
- **Rationale**: Discogs API does not require authentication for public read endpoints. Keeping API calls on the server avoids CORS issues and keeps any future API key centralized.
- **Alternative considered**: Client-side fetch directly from the browser. Rejected because Discogs does not send CORS headers for all endpoints, and server-side is easier to extend with caching or auth later.

### Use vanilla JavaScript for player interactivity
- **Rationale**: The feature is small and scoped to a single page. Adding a full JS framework is unnecessary overhead.
- **Alternative considered**: Stimulus (Rails default) or React. Rejected to keep dependencies minimal and match the user's preference for simplicity.

### Store history in `localStorage`
- **Rationale**: Persists song history across browser sessions and tabs, so users can return to previously played songs even after closing the browser. No backend storage needed.
- **Alternative considered**: `sessionStorage`. Rejected because it only survives within a single tab session and is lost when the browser is closed, which limits the "go back" experience.

### Parse YouTube links from Discogs release `videos` array
- **Rationale**: Discogs release data includes a `videos` array with `uri` fields containing YouTube URLs. We will extract the `v` parameter or handle youtu.be short links to get the video ID.
- **Alternative considered**: Using a separate YouTube Data API search. Rejected because it adds complexity and another API key requirement; Discogs `videos` field is sufficient when present.

### Single controller action with JSON endpoint for random track
- **Rationale**: The page loads once and fetches new tracks via AJAX. This allows seamless "Next" transitions without full page reloads.
- **Alternative considered**: Server-rendered page that reloads on each song change. Rejected for worse UX and slower transitions.

## Risks / Trade-offs

- **[Risk]** Discogs API has rate limits (unauthenticated: 25 requests per minute). Rapid skipping could hit the limit.
  - **Mitigation**: Implement a small client-side debounce on the "Next" button and consider a simple server-side cache (e.g., Rails cache) for release metadata if needed in the future.
- **[Risk]** Many Discogs releases do not have `videos` entries, leading to multiple API calls before finding a playable track.
  - **Mitigation**: Fetch releases in batches or use pagination parameters to increase the candidate pool. Add a reasonable retry limit (e.g., 5 attempts) and show a user-friendly message if no video is found.
- **[Risk]** YouTube links in Discogs may be dead, region-blocked, or removed.
  - **Mitigation**: This is accepted as a known limitation. The user can skip to the next random song if a video fails to load.
- **[Risk]** YouTube iframe embedding may be blocked by some browser privacy settings or extensions.
  - **Mitigation**: Display a fallback message directing the user to open the link on YouTube directly if the iframe fails to load.
