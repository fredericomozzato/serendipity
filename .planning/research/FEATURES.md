# Feature Landscape

**Domain:** Random music discovery / web jukebox (no-auth, Discogs + YouTube)
**Researched:** 2026-03-31

---

## Table Stakes

Features users expect from any music playback product. Missing = product feels broken or
users bounce within seconds.

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Immediate playback on load | Core value prop; every music app starts playing on open | Low | Discogs Search API random-page + YouTube embed; silent retry on no video |
| Visible play / pause control | Players without a pause feel broken | Low | YouTube iframe API exposes `pauseVideo` / `playVideo`; keep button always in thumb reach |
| Skip (next track) | Radio apps always have a skip; Pandora, last.fm, Spotify all surface it prominently | Low | Triggers new Discogs fetch + embed swap |
| Back / history navigation | "What was that song?" is the #1 post-discovery frustration | Medium | Session array of 10-20 items; back button steps backward; no persistence across reloads is acceptable for V1 |
| Track metadata display | Artist, title, year, label — without this it's a black box | Low | Discogs release object contains all needed fields |
| Album art / cover display | Visual identity of what's playing; absent = app feels empty | Low | Discogs `thumb` or `cover_image` field |
| Genre filter | Users will immediately ask "can I hear only jazz?" — table stakes for a catalog-backed product | Medium | Discogs ~13 top-level genres; filter appended to Search API call |
| Loading indicator during fetch | API round-trip takes 0.5–2s; no feedback = users think app is broken | Low | Spinner or skeleton state between skip and video ready |
| Error state for unavailable video | YouTube embeds can become unavailable (geo-block, takedown, private); silent infinite retry is not acceptable | Medium | Show "searching..." feedback; auto-skip to next; cap retries visibly |
| Mobile-first touch targets | Majority of personal-use sessions on mobile; controls must be reachable with one thumb | Low | 48x48 px minimum targets; skip/pause in bottom half of viewport |

---

## Differentiators

Features that make Serendipity memorable or defensible. Not expected, but high value if done well.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Retro / neon-on-dark aesthetic | Personality and nostalgia; most music apps look the same (Spotify clone) | Low | Pure CSS; typography + color palette; spinning vinyl or VU meter animation optional |
| "Serendipity" framing in copy | Names the feeling the app creates; makes skipping feel like adventure, not frustration | Low | UX copy only; no code cost |
| Discogs catalog depth | 17M+ releases including deep cuts, obscure labels, world music; far broader than Spotify/YouTube Music | None (platform advantage) | Surface this: "Pulling from 17 million releases" in empty/loading state |
| Genre + decade filter combination | Filtering by genre is table stakes; adding era filtering (70s, 80s) is a meaningful discovery tool with zero extra API complexity | Medium | Discogs Search API supports `year` range; could add decade selector in V1.5 |
| "Now playing" link to Discogs release | Power users want to dig deeper — "who released this?", "is this on vinyl?" | Low | Discogs release URL is in the API response; one anchor tag |
| Animated loading state that fits the retro theme | Delightful; turns a 1-2s API wait into a moment of anticipation | Low | CSS animation only; "flipping through records" or "searching crates" copy |
| Session history visible as a track list | Show the last N songs as a scrollable mini-history | Medium | Adds visible context to back-navigation; state already exists in session array |

---

## Anti-Features

Things to deliberately NOT build in V1. Including these would increase complexity, delay
shipping, and dilute focus without proportional user value.

| Anti-Feature | Why Avoid | What to Do Instead |
|--------------|-----------|-------------------|
| User accounts / authentication | Requires auth flow, session management, email/password or OAuth — V1 has zero persistence need | Defer to V2; Docker Compose is structured to add PostgreSQL trivially |
| Persistent favorites / likes | Requires DB; without auth it's anonymous state that can't meaningfully survive reloads | Defer; when DB + auth land in V2 this becomes high-value |
| Playlists or queue management | Queue implies intent-driven listening; this app is serendipity-driven | Skip is the queue; don't fight the product's identity |
| Recommendation engine / ML | Complexity far exceeds V1 value; Discogs catalog + genre filter IS the discovery mechanism | The random-page strategy IS the algorithm |
| Social features (share, vote, collaborative jukebox) | Requires real-time infra (Action Cable/Redis) ruled out in PROJECT.md | Defer; shareable links to a specific release are low-cost if needed later |
| Volume control | YouTube embed handles volume at OS/browser level; in-app volume control adds zero value for web | Let the browser/OS do it |
| Equalizer / audio effects | Out of scope for a web jukebox; adds no value when audio passes through YouTube | Never; this is not a DAW |
| Offline support / PWA caching | Relies on YouTube streams; can't cache; PWA shell adds infra for no gain | Not applicable |
| Search / browse by artist or album | Browsing inverts the value prop; Serendipity's identity IS that you don't choose | Use genre filter for coarse control; let randomness do the rest |
| Podcast / non-music content | Discogs is a music catalog; mixing in podcasts erodes identity | Out of scope permanently |
| Ads or monetization layer | V1 is personal use + portfolio; premature monetization adds legal, UX, and infra complexity | Irrelevant for V1 |

---

## Feature Dependencies

```
Immediate playback on load
  → Track metadata display        (need Discogs response before rendering)
  → Album art display             (same Discogs response)
  → Loading indicator             (required to mask API latency)
  → Error state for unavailable video  (YouTube embed may fail)

Skip (next track)
  → Back / history navigation     (history requires knowing what played)
  → Loading indicator             (same latency masking need)
  → Error state                   (same retry logic)

Genre filter
  → Skip (next track)             (filter only activates on next fetch, not mid-track)
  → Session history               (filtered history should be tracked consistently)

Session history (back navigation)
  → Track metadata display        (metadata must be stored per entry to be useful)

"Now playing" Discogs link
  → Track metadata display        (release URL comes from same Discogs object)
```

---

## MVP Recommendation

### Build in V1

1. **Immediate autoplay on load** — the entire value prop lives here
2. **Loading indicator** — covers the 0.5-2s API fetch gap; without it the app feels frozen
3. **Error state + silent retry** — YouTube video unavailability is common (~60-80% of Discogs releases); must handle gracefully
4. **Track metadata + album art** — artist, title, year, label, cover; makes every random track feel like a discovery
5. **Skip button** — table stakes; always visible; in thumb zone
6. **Back / session history** — 10-20 track window; no DB needed
7. **Genre filter** — 13 Discogs top-level genres; single-select; clears on "All"
8. **Retro aesthetic** — neon on dark, vintage typography; this is a differentiator with zero code cost beyond CSS
9. **"Now playing" Discogs link** — single anchor tag; high value for curious users

### Defer to V2

- Favorites / likes (requires DB + auth)
- Session history as visible track list (nice-to-have; back button covers the core need)
- Decade / era filter combination (low cost but not blocking V1)
- Social / shareable links

---

## Sources

- [The Jukebox App — features overview](https://www.thejkbxapp.com/)
- [Pandora features and Music Genome Project (Quora comparison)](https://www.quora.com/What-are-the-differences-between-Last-fm-Pandora-Radio-and-Spotify-What-advantages-does-each-provide)
- [19 Must-Have Features for a Market-Leading Music App — SolGuruz](https://solguruz.com/blog/19-must-have-features-for-a-market-leading-music-app/)
- [Enhancing UX for Music Streaming Apps — Onething Design](https://www.onething.design/post/tuning-ux-for-music-streaming-apps)
- [UX Design Patterns for Loading — Pencil & Paper](https://www.pencilandpaper.io/articles/ux-pattern-analysis-loading-feedback)
- [Error Handling UX Design Patterns — Medium / Design Bootcamp](https://medium.com/design-bootcamp/error-handling-ux-design-patterns-c2a5bbae5f8d)
- [Mobile App UX: Thumb Zones and Gestures — Elaris Software](https://elaris.software/blog/mobile-ux-thumb-zones-2025/)
- [Algorithm-free music discovery: Lazyrecords launch — DJ Mag](https://djmag.com/news/algorithm-free-music-discovery-app-djs-lazyrecords-launched)
- [How to break free of Spotify's algorithm — MIT Technology Review](https://www.technologyreview.com/2024/08/16/1096276/spotify-algorithms-music-discovery-ux/)
- [Discogs iOS/Android app feature requests forum](https://www.discogs.com/forum/thread/707219)
- [Retro Music Player UI Design — Figma Community](https://www.figma.com/community/file/1058350090477787969/retro-music-player-ui-design)
