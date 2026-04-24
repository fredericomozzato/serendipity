## Context

Serendipity is a Rails 8 + Tailwind CSS music discovery app that plays random YouTube videos sourced from Discogs. The player page (`app/views/player/index.html.erb`) currently uses unstyled HTML elements with minimal Tailwind utility classes (`container mx-auto mt-28 px-5 flex`). The app has three Stimulus controllers (`player`, `history`, `genre-filter`) but no dedicated stylesheet beyond the default Tailwind setup. The goal is to transform the UI into a dark, vintage cassette-era experience without altering any functional behavior.

## Goals / Non-Goals

**Goals:**
- Deliver a single, cohesive dark theme with earthy brown, green, purple, and crimson/orange accents
- Evoke analog cassette / boombox nostalgia through typography, color, and subtle textures
- Implement the design entirely within the existing Rails + Tailwind CSS stack
- Ensure the theme applies consistently across the player layout, genre tags, controls, and track info

**Non-Goals:**
- No dark/light mode toggle or theme switching
- No new JavaScript frameworks or CSS-in-JS libraries
- No changes to server-side logic, Discogs API integration, or Stimulus controller behavior
- No custom icon fonts or heavy image assets (keep it CSS-driven)

## Decisions

**Decision: Use Tailwind CSS configuration (`tailwind.config.js`) to define the custom theme rather than ad-hoc utility classes.**
- *Rationale*: Centralizing colors, fonts, and spacing in the Tailwind config ensures consistency and makes the design system maintainable. It also allows using semantic class names (e.g., `bg-cassette-brown`) instead of arbitrary hex values scattered through ERB templates.
- *Alternative considered*: Inline arbitrary values (e.g., `bg-[#3E2723]`). Rejected because it duplicates values and is harder to update globally.

**Decision: Use Google Fonts for retro typography (e.g., "VT323" for labels and "Space Mono" for track metadata).**
- *Rationale*: Free, web-optimized fonts that evoke 1980s computing and cassette labeling without requiring self-hosting or asset pipeline changes.
- *Alternative considered*: Self-hosting font files. Rejected to keep the change simple and avoid Propshaft font pipeline complexity.

**Decision: Implement subtle noise/grain texture via a tiny base64-encoded SVG or CSS gradient rather than an external image asset.**
- *Rationale*: Avoids an extra HTTP request and keeps the asset pipeline clean. A CSS-only approach is sufficient for the subtle vintage effect.
- *Alternative considered*: Loading a PNG texture image. Rejected to minimize asset dependencies.

**Decision: Restructure the player DOM in `index.html.erb` to support a boombox/cassette-deck layout while preserving Stimulus data attributes.**
- *Rationale*: The current flat DOM structure cannot achieve the desired visual hierarchy (e.g., a "deck" area for the video, a "control panel" for buttons). A shallow restructure is needed, but all `data-controller`, `data-target`, and `data-action` attributes will be preserved exactly.
- *Alternative considered*: Purely CSS-based layout changes. Rejected because the current markup lacks semantic containers for the cassette aesthetic.

## Risks / Trade-offs

- **[Risk] Tailwind CDN / build caching issues when adding custom config** → *Mitigation*: Run `bin/rails tailwindcss:build` after config changes and verify in development.
- **[Risk] Font loading flash or layout shift (FOUT/FOIT)** → *Mitigation*: Use `font-display: swap` and ensure fallback system fonts have similar metrics.
- **[Risk] Accessibility concerns with low-contrast dark earthy tones** → *Mitigation*: Verify WCAG AA contrast ratios for all text/background combinations; adjust purple/brown shades if needed.
- **[Trade-off] Single theme excludes users who prefer light mode** → *Accepted*: Explicitly out of scope per non-goals.
