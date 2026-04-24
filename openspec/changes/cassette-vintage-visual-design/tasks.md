## 1. Theme Foundation

- [x] 1.1 Add Google Fonts imports (`VT323`, `Space Mono`) to `app/views/layouts/application.html.erb`
- [x] 1.2 Define custom `@theme` in `app/assets/tailwind/application.css` with vintage color tokens (deep brown, muted green, desaturated purple, crimson/orange accent, cream text)
- [x] 1.3 Add base body styles in `app/assets/tailwind/application.css` (background color, text color, font family, subtle noise texture via CSS)
- [x] 1.4 Run `bin/rails tailwindcss:build` and verify the theme compiles without errors

## 2. Player Layout & Structure

- [x] 2.1 Restructure `app/views/player/index.html.erb` into a cassette-deck layout: wrap content in a `.player-container` with distinct zones for filter bar, video "window", track info "display", and control "panel"
- [x] 2.2 Ensure all `data-controller`, `data-target`, and `data-action` attributes are preserved exactly during DOM restructuring
- [x] 2.3 Apply Tailwind utility classes for layout (centering, spacing, max-width) using the custom theme tokens

## 3. Component Styling

- [x] 3.1 Rewrite `.video-wrapper` styles in `app/assets/stylesheets/application.css` for a dark recessed bezel look with inner shadows and rounded corners
- [x] 3.2 Rewrite `.track-info` styles with retro typography (VT323 for title, Space Mono for artist) and appropriate text colors from the theme
- [x] 3.3 Rewrite `.control-button` styles with tactile push-button aesthetic: dark background, cream text, crimson/orange hover glow, pressed active state
- [x] 3.4 Update disabled `.control-button` state to use muted theme colors instead of gray
- [x] 3.5 Style `.loading-state` and `.error-message` using theme colors (cream text on dark background, subtle orange accent for errors)

## 4. Genre Filter Theme Update

- [x] 4.1 Rewrite `.genre-tag` base styles: muted green background, cream text, subtle border
- [x] 4.2 Rewrite `.genre-tag.active` styles: desaturated purple background, cream text
- [x] 4.3 Add hover states for both active and inactive genre tags using theme colors
- [x] 4.4 Verify the Stimulus `genre-filter` controller correctly toggles the `active` class (no JS changes expected)

## 5. Polish & Verification

- [x] 5.1 Verify WCAG AA contrast ratios for all text/background combinations; adjust theme hex values if needed
- [x] 5.2 Test responsive behavior on mobile viewports; adjust spacing and font sizes
- [x] 5.3 Verify no visual regressions in player functionality (video loads, buttons work, genre filtering behaves correctly)
- [x] 5.4 Remove any now-unused Tailwind utility classes from the layout if applicable
- [x] 5.5 Run the full test suite (`bin/rspec`) to confirm no functional regressions
