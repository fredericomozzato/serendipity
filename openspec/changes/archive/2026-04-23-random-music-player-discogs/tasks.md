## 1. Setup

- [x] 1.1 Add `player` route mapped to `PlayerController#index`
- [x] 1.2 Create `PlayerController` with `index` action
- [x] 1.3 Create empty `app/views/player/index.html.erb` view
- [x] 1.4 Verify the `/player` page loads successfully

## 2. Discogs API Integration

- [x] 2.1 Create a `DiscogsService` class to wrap Discogs API calls
- [x] 2.2 Implement method to fetch a random release using Discogs search or database endpoints
- [x] 2.3 Implement method to extract track and artist info from a release
- [x] 2.4 Implement method to parse YouTube video ID from Discogs `videos` URIs (handle `youtube.com/watch?v=` and `youtu.be/` formats)
- [x] 2.5 Implement retry loop that skips releases without YouTube links, up to 5 attempts
- [x] 2.6 Create `PlayerController#random_track` JSON endpoint that returns `{ title, artist, video_id }`
- [x] 2.7 Add error handling when no playable track is found after retries

## 3. Player UI

- [x] 3.1 Add YouTube iframe embed container to `index.html.erb`
- [x] 3.2 Add song info display elements (track title and artist name)
- [x] 3.3 Add "Next Random Song" button
- [x] 3.4 Add "Previous Song" button
- [x] 3.5 Style the player page with minimal, clean CSS (centered layout, responsive iframe)

## 4. Frontend Interactivity

- [x] 4.1 Create `app/javascript/player.js` to handle player logic
- [x] 4.2 Implement `loadTrack(track)` to update iframe src and song info text
- [x] 4.3 Implement `fetchRandomTrack()` to call the JSON endpoint and invoke `loadTrack`
- [x] 4.4 Implement `nextTrack()` that fetches a new track, pushes current track to `localStorage` history, and loads the new one
- [x] 4.5 Implement `previousTrack()` that pops from `localStorage` history and loads the prior track
- [x] 4.6 Wire up button click events to `nextTrack` and `previousTrack`
- [x] 4.7 On page load, restore current track from `localStorage` if present; otherwise fetch a new random track
- [x] 4.8 Disable or guard the "Previous" button when history is empty

## 5. Polish & Error Handling

- [x] 5.1 Display a user-friendly message when no playable track can be found
- [x] 5.2 Add basic loading state while fetching a new track
- [x] 5.3 Handle Discogs API failures gracefully (network errors, rate limits)
- [x] 5.4 Ensure the iframe does not break layout on mobile screens
- [x] 5.5 Do a manual end-to-end test: load page, skip songs, go back, refresh page, close and reopen browser, verify localStorage restore
