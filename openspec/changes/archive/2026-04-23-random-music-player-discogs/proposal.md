## Why

We want a lightweight, zero-auth music discovery experience that lets anyone explore random songs via the Discogs catalog and immediately play them through embedded YouTube videos. This gives users a fun, serendipitous way to discover music without needing accounts, playlists, or complex UI.

## What Changes

- Add a new Rails controller (`PlayerController`) with a single `index` action that renders the music player page.
- Integrate the Discogs API to fetch random releases and extract tracks with valid YouTube video links.
- Add backend logic to filter out releases/tracks that lack YouTube links, automatically skipping until a playable track is found.
- Embed an iframe YouTube player on the page using the video ID extracted from Discogs release data.
- Add client-side JavaScript for "Next Random Song" and "Previous Song" navigation.
- Store the user's recent song history in browser `localStorage` so it persists across sessions (no database required).
- Display the current song title and artist name below the video player.

## Capabilities

### New Capabilities
- `random-music-player`: Core music discovery and playback capability. Covers Discogs API integration, random track selection with YouTube filtering, player UI, and localStorage-persisted navigation history.

### Modified Capabilities
<!-- No existing capabilities are being modified -->

## Impact

- New controller, view, and route in the Rails application.
- New JavaScript file for player interactivity and localStorage-persisted history.
- Dependency on external Discogs API (no API key required for public endpoints).
- No database migrations or user authentication changes.
