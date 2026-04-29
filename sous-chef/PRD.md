# PRD — Serendipity

Serendipity is a true random music discovery tool that surfaces releases from the Discogs catalog and plays them as embedded YouTube videos. It solves the paradox-of-choice problem in streaming by removing all browsing — users just press play and let the catalog surprise them. Anonymous visitors get the full discovery experience immediately; registered users (Google OAuth only) unlock a liked collection, full history management, and shareable mixtapes.

---

## Users

**Anonymous visitors** — no account required. Access the full discovery player, genre filtering, and a session history of the last 20 songs (stored in localStorage). Cannot like songs, create mixtapes, or persist data across devices.

**Authenticated users** — sign in with Google. Everything anonymous users can do, plus: liked collection, unlimited listening history with management controls, and mixtapes. Authentication exists solely to persist and sync data — there are no social or public profile features.

**Mixtape guests** — anyone with a mixtape share link. Can listen to the mixtape using the full player. No account required, no interaction beyond playback.

---

## Features

### Random Discovery Player
STATUS: PLANNED

The core experience. Fetches a random Discogs release, extracts the first YouTube URL from its `videos` array, and renders it as an embedded YouTube player. Releases with an empty `videos` array are silently skipped and a new random release is fetched (retry loop). The player exposes three controls: **Next** (new random release), **Previous** (load most recent history entry), and **Genre Filter** (multi-select, OR logic).

**Scope:**
- Fetch a random release from Discogs, retry until a release with a non-empty `videos` array is found
- Embed the first `videos[0]` YouTube URL using the YouTube IFrame API
- Display release metadata: title, artist, year, cover art, Discogs genre
- Next and Previous controls work for all users (anonymous and authenticated)
- Genre filter accepts multiple selections; OR logic — a release matching any selected genre qualifies
- Filter state persists across next/previous navigations within a session

**Out of scope:**
- Multiple video selection per release
- Audio-only fallback if the YouTube video is unavailable
- Autoplay (user must initiate first play)

---

### Genre Filtering
STATUS: PLANNED

Users can filter the random draw by one or more Discogs top-level genres. When filters are active, the genre parameter is passed to the Discogs search API before randomizing — not applied after the fact. Multiple selected genres use OR logic.

**Scope:**
- Multi-select filter UI, populated from the fixed Discogs genre taxonomy
- Active filters narrow the Discogs query; no filter = full catalog
- Filter selections persist in the current session (localStorage for anonymous, DB preference for authenticated users)
- Clearing all filters returns to full-catalog random

**Out of scope:**
- Discogs "styles" (sub-genres)
- Saving named filter presets

---

### Listening History
STATUS: PLANNED

Every played release is recorded as a history entry. Anonymous users get a rolling 20-entry list in localStorage. Authenticated users get an unlimited history stored in the database.

**Scope:**
- Anonymous: localStorage history, capped at 20 entries (FIFO)
- Authenticated: unlimited DB history; entries record `played_at` timestamp
- Previous track button loads and plays the most recent history entry
- Authenticated users can: remove individual entries, pause history recording, clear all history
- When history is paused: DB writes stop, but the current session's in-memory list still powers the Previous button

**Out of scope:**
- Exporting history
- History analytics or charts

---

### Liked Collection
STATUS: PLANNED

Authenticated users can like any currently playing release. Likes are stored as a private collection — no public visibility. The collection view is a flat list with the ability to unlike (remove) entries.

**Scope:**
- Like/unlike button on the player (authenticated only)
- Collection view: chronological list of liked releases with metadata
- Unlike removes the entry from the collection
- Like state is reflected immediately in the player UI

**Out of scope:**
- Playing the collection as a queue
- Sharing the collection
- Notes, tags, or ratings on liked items

---

### Mixtapes
STATUS: PLANNED

Authenticated users can create named, ordered playlists of up to 50 songs. A mixtape can be shared with anyone via a permanent link. Guests with the link can listen using the full player but cannot interact. The creator can edit the mixtape at any time; the share link always reflects the current state.

**Scope:**
- Create a mixtape with a name
- Add a currently playing release to a mixtape via a picker modal (supports multiple mixtapes)
- Remove songs from a mixtape
- Reorder songs via drag-and-drop
- Maximum 50 songs per mixtape
- Generate a unique share token; share link is permanent and always current
- Mixtape player: same embedded player, navigates within the mixtape's ordered track list
- Guests with share link: full playback, no liking or saving
- Creator can rename or delete a mixtape

**Out of scope:**
- Collaborative editing (only the creator can modify)
- Mixtape comments or reactions
- Making a mixtape private again after sharing (the link just stops working if the mixtape is deleted)

---

## UI / UX

**Layout**: Single-page, player-centric. The player dominates the viewport. Controls (next, previous, genre filter) are immediately accessible without scrolling. On mobile, the layout is a vertical stack: cover art / video embed → track metadata → player controls → genre filter chips.

**Key screens:**
1. **Discovery player** — the home screen for all users. Video embed front and center, controls below, genre filter as a collapsible chip row.
2. **Mixtape player** — same player shell, track list visible below or as a slide-up drawer. Distinguishable from discovery mode (header shows mixtape name).
3. **Collection** (authenticated) — scrollable list of liked releases with cover art, title, artist, and unlike button.
4. **History** (authenticated) — scrollable list with played-at timestamps, remove button per entry, and global controls (pause, clear all).
5. **Mixtape management** (authenticated) — list of mixtapes with create button; tap into a mixtape to reorder/remove tracks.
6. **Sign-in** — single button: "Continue with Google". No email/password fields.

**Mobile-first requirements:**
- Touch-friendly tap targets (minimum 44×44px)
- No hover-dependent interactions for core controls
- Drag-to-reorder in mixtapes uses touch-compatible library

**Dark/light mode:**
- System default on first visit; user can toggle and preference is persisted (localStorage for anonymous, DB for authenticated)
- Tailwind `dark:` variant throughout; no separate stylesheet

**Visual design**: Full palette, typography, and component decisions are deferred to post-MVP. Tailwind is in place from day one to make the transition seamless.

---

## Design

### Visual direction

Deferred post-MVP — full design system (palette, typography, components) to be defined after core features are validated.

### Color palette

| Role | Color / description |
|------|---------------------|
| All | Deferred post-MVP |

### Typography

| Role | Choice |
|------|--------|
| All | Deferred post-MVP (Tailwind defaults at launch) |

### Layout

Mobile-first, player-centric single column. Full visual refinement deferred.

### Component library

Tailwind CSS utility classes from day one. Component library choice (shadcn/ui or similar) deferred to design phase.

### Dark mode

Yes — system default, user-toggleable, persisted per user. Tailwind `dark:` variant.

---

## Data Model

### releases

| Field | Type | Notes |
|-------|------|-------|
| discogs_id | string | Unique, indexed |
| title | string | |
| artist | string | |
| year | integer | |
| genre | string | First/primary Discogs genre |
| video_url | string | `videos[0]` YouTube URL; never blank |
| cover_url | string | Discogs cover image URL |
| discogs_url | string | Canonical Discogs release URL |

Discogs releases are cached on first fetch. A release is only persisted if it has a non-empty `video_url`.

---

### users

| Field | Type | Notes |
|-------|------|-------|
| email | string | Unique, from Google profile |
| name | string | Display name |
| avatar_url | string | Google profile picture |
| google_uid | string | Unique, from OmniAuth |
| history_paused | boolean | Default false |

---

### history_entries

| Field | Type | Notes |
|-------|------|-------|
| user_id | references users | |
| release_id | references releases | |
| played_at | datetime | |

Ordered by `played_at DESC`. No uniqueness constraint — the same release can appear multiple times.

---

### likes

| Field | Type | Notes |
|-------|------|-------|
| user_id | references users | |
| release_id | references releases | |

Unique constraint on `(user_id, release_id)`.

---

### mixtapes

| Field | Type | Notes |
|-------|------|-------|
| user_id | references users | |
| title | string | |
| share_token | string | Unique, URL-safe, generated on create |

Mixtape is shareable at creation time via `/m/:share_token`. Deleting the mixtape invalidates the link.

---

### mixtape_tracks

| Field | Type | Notes |
|-------|------|-------|
| mixtape_id | references mixtapes | |
| release_id | references releases | |
| position | integer | 1-based, maintained on reorder |

Unique constraint on `(mixtape_id, release_id)` — a release can only appear once per mixtape. Max 50 tracks enforced at the model level.

---

## Out of scope (MVP)

- Public user profiles
- Playing the liked collection as a queue
- Shareable collections
- Multi-video selection per release
- Discogs styles (sub-genre) filtering
- Saved filter presets
- Collaborative mixtape editing
- Mixtape comments or reactions
- History analytics
- Email notifications
- Any social or discovery-of-users features
