## MODIFIED Requirements

### Requirement: Player page is accessible
The system SHALL provide a dedicated route that renders the music player page, and the root path SHALL redirect to the player page.

#### Scenario: User visits the root path
- **WHEN** a user navigates to the root path
- **THEN** the system redirects the user to the player page

#### Scenario: User visits the player path
- **WHEN** a user navigates to the player route
- **THEN** the system renders the player page with a video container, song info area, and navigation controls

### Requirement: System fetches random releases from Discogs
The system SHALL query the Discogs API to retrieve random release data, optionally constrained by a genre filter parameter.

#### Scenario: Requesting a random track
- **WHEN** the user requests a new random song
- **THEN** the system fetches release data from the Discogs API, applying any active genre filters to the search query

### Requirement: User can skip to next random song
The system SHALL allow the user to request a new random song, which triggers a fresh Discogs fetch respecting the current active genre filters and updates the embedded player and track info.

#### Scenario: User clicks next
- **WHEN** the user activates the "Next" control
- **THEN** the system fetches a new random playable track using the active genre filters, updates the embedded video, updates the displayed track info, and adds the new track to the play history

### Requirement: Play history persists across sessions
The system SHALL store the recent song history and active genre filters in browser `localStorage` so it survives page refreshes, tab closes, and browser restarts.

#### Scenario: User returns after closing the browser
- **WHEN** the user reopens the player page after a previous session
- **THEN** the system restores the current track, history, and active genre filters from `localStorage`

#### Scenario: Page refresh during a session
- **WHEN** the user refreshes the player page
- **THEN** the system restores the current track, history, and active genre filters from `localStorage`
