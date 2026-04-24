## Why

The Serendipity music player currently has a generic, unstyled appearance that fails to communicate the app's nostalgic, discovery-driven personality. A cohesive vintage cassette-era visual identity will create an immersive, memorable experience that matches the joy of stumbling upon unexpected music.

## What Changes

- Introduce a single, permanent visual theme (no dark/light mode toggle) with a vintage analog cassette aesthetic
- Define a dark, earthy color palette: deep browns, muted greens, desaturated purples, with crimson/orange accents
- Add retro typography, subtle noise/grain textures, and cassette-inspired UI chrome
- Redesign the player layout with a vintage hi-fi / boombox feel
- Update genre filter tags, control buttons, and track info display to match the new aesthetic
- **BREAKING**: Replace all existing Tailwind utility classes and custom CSS with the new design system

## Capabilities

### New Capabilities
- `vintage-visual-design`: Establish the complete visual design system including color palette, typography, spacing, textures, and component styles for the cassette-era theme

### Modified Capabilities
- `genre-filtering`: Update genre tag styling requirements to use the new vintage color palette instead of the current blue/off-white scheme

## Impact

- **Frontend**: All CSS/Tailwind classes in `app/assets/stylesheets/`, `app/views/player/index.html.erb`, and Stimulus controllers
- **Asset pipeline**: New image/textures may be added via Propshaft
- **No API changes**: Purely presentational; server-side logic and Discogs integration remain unchanged
- **No dependency changes**: Continues using Tailwind CSS via `tailwindcss-rails`
