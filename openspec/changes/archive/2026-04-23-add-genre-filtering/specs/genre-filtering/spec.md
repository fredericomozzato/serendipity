## ADDED Requirements

### Requirement: System displays available genre filter tags
The system SHALL render a row of clickable genre tags above the player, sourced from the canonical Discogs genre list.

#### Scenario: User visits the player page
- **WHEN** a user navigates to the player route
- **THEN** the system renders the genre filter tag bar with all available genres as small rounded-rectangles in an inactive (off-white background with black text) state

### Requirement: User can toggle genre filter tags on and off
The system SHALL allow users to activate and deactivate individual genre filter tags by clicking them.

#### Scenario: User activates a genre filter
- **WHEN** a user clicks an inactive genre tag
- **THEN** the tag renders as a small rounded-rectangle with a blue background and white text, and the system includes that genre in subsequent track fetches

#### Scenario: User deactivates a genre filter
- **WHEN** a user clicks an active genre tag
- **THEN** the tag returns to an inactive (off-white background with black text) state and the system stops including that genre in track fetches

#### Scenario: User selects multiple genre filters
- **WHEN** a user activates two or more genre tags
- **THEN** all selected genres are considered active and the system uses additive (OR) logic for track selection

### Requirement: System persists active genre filters across sessions
The system SHALL store the set of active genre filters in browser `localStorage` so it survives page refreshes, tab closes, and browser restarts.

#### Scenario: User returns after closing the browser
- **WHEN** the user reopens the player page after a previous session with active filters
- **THEN** the system restores the active genre filters from `localStorage` and renders the corresponding tags as active

#### Scenario: Page refresh during a session
- **WHEN** the user refreshes the player page with active filters
- **THEN** the system restores the active genre filters from `localStorage`

### Requirement: Server applies active genre filters when fetching random tracks
The system SHALL pass active genre filters to the server when requesting a random track, and the server SHALL use them to bias the Discogs API query.

#### Scenario: Requesting a track with no active filters
- **WHEN** the user requests a new random song with no genre filters active
- **THEN** the server fetches a completely random release from Discogs without genre constraints

#### Scenario: Requesting a track with active filters
- **WHEN** the user requests a new random song with one or more active genre filters
- **THEN** the server selects one genre at random from the active set and queries the Discogs API with that genre parameter
