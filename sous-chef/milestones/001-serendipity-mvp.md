---
id: "001"
name: serendipity-mvp
status: IN_PROGRESS
---

# Serendipity MVP

Core music discovery experience: a random-release player backed by the Discogs catalog with YouTube embedding, session history, genre filtering, Google sign-in, authenticated history management, and a liked collection. Mixtapes are scoped to a separate milestone.

All anonymous features work without an account. Authenticated features layer on top of the same player without changing the core experience.

## Slices

### Slice 001 — Bootstrap the application
STATUS: PENDING

Delivers: Rails app and Docker environment ready for development

Scope:
- Run /chef:bootstrap to complete this slice — do not scaffold manually

---

### Slice 002 — Player skeleton with hardcoded videos
STATUS: IN_PROGRESS

Delivers: user can load the app and cycle through a hardcoded set of YouTube videos using Next and Previous controls

Scope:
- Single-page layout with video embed area, track metadata area, and player controls
- Next and Previous buttons cycle through a hardcoded list of YouTube video IDs
- YouTube IFrame API loaded and controlled via a Stimulus controller
- Player UI structure matches the PRD layout (player-centric, mobile-first)
- No Discogs integration, no persistence — purely frontend skeleton

---

### Slice 003 — Discogs integration
STATUS: PENDING

Delivers: player fetches and displays a real random release from the Discogs API, retrying automatically if the release has no video

Scope:
- Discogs service object fetches a random release and retries until one with a video URL is found
- Release is cached on first fetch (find-or-create by Discogs ID); only persisted if a video URL is present
- Player displays release metadata: title, artist, year, cover art, Discogs genre
- Next button triggers a new Discogs fetch; Previous still cycles through session state
- No genre filtering yet — full catalog random

---

### Slice 004 — Anonymous history (localStorage)
STATUS: PENDING

Delivers: anonymous user's last 20 played releases are tracked in the browser; Previous button navigates through that history

Scope:
- Each played release is appended to a localStorage history list (capped at 20, FIFO)
- Previous button loads and plays the most recent history entry
- History list powers the Previous button for all users in this slice
- No server-side storage; no authenticated history yet

---

### Slice 005 — Genre filtering
STATUS: PENDING

Delivers: user can narrow the random draw by one or more Discogs genres; filter state persists across Next/Previous navigations within the session

Scope:
- Multi-select genre filter UI populated from the fixed Discogs genre taxonomy
- Active filters pass the selected genre(s) to the Discogs search API before randomizing (OR logic — one request per genre, results merged and sampled)
- No filter = full catalog random
- Filter state held in the Stimulus player controller; persists for the session
- Clearing all filters returns to full-catalog random

---

### Slice 006 — Google OAuth sign-in
STATUS: PENDING

Delivers: user can sign in with Google and remain authenticated across page loads; sign-out returns them to anonymous mode

Scope:
- Devise with OmniAuth Google provider; no password registration or reset flows
- User record created on first sign-in (email, name, avatar_url, google_uid)
- Sign-in button visible to anonymous users; sign-out available to authenticated users
- current_user available throughout the app via Devise helper
- No feature changes to the player yet — authentication infrastructure only

---

### Slice 007 — Authenticated history
STATUS: PENDING

Delivers: signed-in user's listening history is saved to the database; they can remove individual entries, pause recording, and clear all history

Scope:
- Each played release appended to history_entries for the current user (played_at timestamp)
- Authenticated users see their DB history instead of the localStorage list
- Previous button loads from DB history for authenticated users
- History management UI: remove individual entry, pause/resume recording, clear all
- history_paused flag on the user record; when paused, DB writes stop but in-memory session list still powers Previous
- Anonymous localStorage history unchanged

---

### Slice 008 — Liked collection
STATUS: PENDING

Delivers: signed-in user can like or unlike the currently playing release, and browse their private liked collection

Scope:
- Like/unlike button on the player, visible only to authenticated users
- Like state reflected immediately in the player UI (Turbo Stream update)
- Likes stored in the likes join table (unique per user + release)
- Collection view: chronological list of liked releases with cover art, title, artist, and unlike button
- Unlike removes the entry and updates the UI immediately
