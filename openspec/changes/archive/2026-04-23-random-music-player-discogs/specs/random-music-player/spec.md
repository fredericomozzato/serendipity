## ADDED Requirements

### Requirement: Player page is accessible
The system SHALL provide a dedicated route that renders the music player page.

#### Scenario: User visits the player page
- **WHEN** a user navigates to the player route
- **THEN** the system renders the player page with a video container, song info area, and navigation controls

### Requirement: System fetches random releases from Discogs
The system SHALL query the Discogs API to retrieve random release data.

#### Scenario: Requesting a random track
- **WHEN** the user requests a new random song
- **THEN** the system fetches release data from the Discogs API

### Requirement: System filters releases for playable YouTube tracks
The system SHALL ignore releases and tracks that do not contain a valid YouTube video link and SHALL continue fetching until a playable track is found or a retry limit is reached.

#### Scenario: Release has a YouTube video link
- **WHEN** the system retrieves a release that includes a `videos` entry with a YouTube URL
- **THEN** the system extracts the video ID and prepares it for embedding

#### Scenario: Release lacks a YouTube video link
- **WHEN** the system retrieves a release without any `videos` entries or with only non-YouTube links
- **THEN** the system discards that release and fetches another random release

#### Scenario: Retry limit exceeded
- **WHEN** the system has attempted to fetch a playable track five times without success
- **THEN** the system returns an error message indicating no playable track was found

### Requirement: System embeds YouTube video player
The system SHALL render an embedded YouTube video player on the page using the extracted video ID.

#### Scenario: Embedding a valid video
- **WHEN** a playable track with a valid YouTube video ID is selected
- **THEN** the system renders an iframe pointing to `https://www.youtube.com/embed/<video_id>`

### Requirement: System displays current track information
The system SHALL display the current song title and artist name on the player page.

#### Scenario: Displaying track metadata
- **WHEN** a track is selected and ready to play
- **THEN** the page shows the track title and the artist name

### Requirement: User can skip to next random song
The system SHALL allow the user to request a new random song, which triggers a fresh Discogs fetch and updates the embedded player and track info.

#### Scenario: User clicks next
- **WHEN** the user activates the "Next" control
- **THEN** the system fetches a new random playable track, updates the embedded video, updates the displayed track info, and adds the new track to the play history

### Requirement: User can return to previous song
The system SHALL allow the user to navigate back to a previously played song from their stored history.

#### Scenario: User clicks previous
- **WHEN** the user activates the "Previous" control
- **THEN** the system loads the prior track from stored history, updates the embedded video, and updates the displayed track info

#### Scenario: No previous track exists
- **WHEN** the user activates the "Previous" control and there is no prior track in stored history
- **THEN** the system does nothing or disables the control

### Requirement: Play history persists across sessions
The system SHALL store the recent song history in browser `localStorage` so it survives page refreshes, tab closes, and browser restarts.

#### Scenario: User returns after closing the browser
- **WHEN** the user reopens the player page after a previous session
- **THEN** the system restores the current track and history from `localStorage`

#### Scenario: Page refresh during a session
- **WHEN** the user refreshes the player page
- **THEN** the system restores the current track and history from `localStorage`
