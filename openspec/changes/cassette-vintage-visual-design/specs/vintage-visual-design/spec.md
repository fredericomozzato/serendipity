## ADDED Requirements

### Requirement: System uses a cohesive vintage cassette-era color palette
The system SHALL apply a single dark theme with earthy brown, muted green, desaturated purple, and crimson/orange accent colors defined in the Tailwind CSS configuration.

#### Scenario: User visits the player page
- **WHEN** a user navigates to the player route
- **THEN** the page background renders as a deep earthy brown, text renders in an off-white/cream color, and interactive elements use the crimson/orange accent

#### Scenario: User views the genre filter bar
- **WHEN** the genre filter bar is rendered
- **THEN** inactive tags use a muted green background with cream text, and active tags use a desaturated purple background with cream text

### Requirement: System uses retro typography
The system SHALL load and apply vintage-inspired web fonts for headings, labels, and track metadata.

#### Scenario: User reads track info
- **WHEN** track title and artist are displayed
- **THEN** the title uses a bold retro display font and the artist uses a monospaced font reminiscent of cassette liner notes

### Requirement: Player layout evokes analog cassette hardware
The system SHALL restructure the player page into a layout reminiscent of a vintage boombox or cassette deck, with distinct zones for the video "window", track info "display", and control "panel".

#### Scenario: User views the player
- **WHEN** the player page is rendered
- **THEN** the video wrapper has a dark recessed bezel appearance, the controls are grouped in a panel with tactile button styling, and the overall layout is centered with generous vintage-inspired spacing

### Requirement: Control buttons have tactile cassette-era styling
The system SHALL style the Previous and Next buttons with a tactile, push-button aesthetic consistent with analog hardware.

#### Scenario: User hovers over a control button
- **WHEN** the user hovers over a control button
- **THEN** the button provides visual feedback (e.g., slight glow or shift) using the crimson/orange accent color

#### Scenario: User clicks a control button
- **WHEN** the user clicks a control button
- **THEN** the button shows an active pressed state

## MODIFIED Requirements

### Requirement: System displays available genre filter tags
The system SHALL render a row of clickable genre tags above the player, sourced from the canonical Discogs genre list.

#### Scenario: User visits the player page
- **WHEN** a user navigates to the player route
- **THEN** the system renders the genre filter tag bar with all available genres as small rounded-rectangles in an inactive state with a muted green background and cream text

### Requirement: User can toggle genre filter tags on and off
The system SHALL allow users to activate and deactivate individual genre filter tags by clicking them.

#### Scenario: User activates a genre filter
- **WHEN** a user clicks an inactive genre tag
- **THEN** the tag renders as a small rounded-rectangle with a desaturated purple background and cream text, and the system includes that genre in subsequent track fetches

#### Scenario: User deactivates a genre filter
- **WHEN** a user clicks an active genre tag
- **THEN** the tag returns to an inactive state with a muted green background and cream text and the system stops including that genre in track fetches

#### Scenario: User selects multiple genre filters
- **WHEN** a user activates two or more genre tags
- **THEN** all selected genres are considered active and the system uses additive (OR) logic for track selection
