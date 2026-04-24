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
