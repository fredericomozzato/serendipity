document.addEventListener("DOMContentLoaded", () => {
  const iframe = document.getElementById("youtube-player")
  const titleEl = document.getElementById("track-title")
  const artistEl = document.getElementById("track-artist")
  const nextButton = document.getElementById("next-button")
  const prevButton = document.getElementById("prev-button")
  const loadingEl = document.getElementById("loading")
  const errorEl = document.getElementById("error-message")

  // Only run on the player page
  if (!iframe || !nextButton) return

  const STORAGE_KEY_CURRENT = "player:current_track"
  const STORAGE_KEY_HISTORY = "player:history"

  function showLoading() {
    loadingEl.classList.remove("hidden")
    errorEl.classList.add("hidden")
    nextButton.disabled = true
  }

  function hideLoading() {
    loadingEl.classList.add("hidden")
    nextButton.disabled = false
  }

  function showError(message) {
    hideLoading()
    errorEl.textContent = message
    errorEl.classList.remove("hidden")
  }

  function loadTrack(track) {
    if (!track) return

    iframe.src = `https://www.youtube.com/embed/${track.video_id}`
    titleEl.textContent = track.title
    artistEl.textContent = track.artist

    localStorage.setItem(STORAGE_KEY_CURRENT, JSON.stringify(track))
    updatePrevButton()
  }

  async function fetchRandomTrack() {
    showLoading()

    try {
      const response = await fetch("/player/random_track")

      if (!response.ok) {
        if (response.status === 404) {
          throw new Error("No playable track found. Please try again.")
        }
        throw new Error("Failed to fetch a random track.")
      }

      const track = await response.json()
      return track
    } catch (error) {
      showError(error.message)
      return null
    } finally {
      hideLoading()
    }
  }

  async function nextTrack() {
    const currentTrack = getCurrentTrack()
    const newTrack = await fetchRandomTrack()

    if (newTrack) {
      if (currentTrack) {
        pushToHistory(currentTrack)
      }
      loadTrack(newTrack)
    }
  }

  function previousTrack() {
    const history = getHistory()
    if (history.length === 0) return

    const prevTrack = history.pop()
    setHistory(history)
    loadTrack(prevTrack)
  }

  function pushToHistory(track) {
    const history = getHistory()
    history.push(track)
    // Keep last 50 tracks
    if (history.length > 50) {
      history.shift()
    }
    setHistory(history)
  }

  function getHistory() {
    try {
      const raw = localStorage.getItem(STORAGE_KEY_HISTORY)
      return raw ? JSON.parse(raw) : []
    } catch {
      return []
    }
  }

  function setHistory(history) {
    localStorage.setItem(STORAGE_KEY_HISTORY, JSON.stringify(history))
  }

  function getCurrentTrack() {
    try {
      const raw = localStorage.getItem(STORAGE_KEY_CURRENT)
      return raw ? JSON.parse(raw) : null
    } catch {
      return null
    }
  }

  function updatePrevButton() {
    const history = getHistory()
    prevButton.disabled = history.length === 0
  }

  // Event listeners
  nextButton.addEventListener("click", nextTrack)
  prevButton.addEventListener("click", previousTrack)

  // Initialize
  const savedTrack = getCurrentTrack()
  if (savedTrack) {
    loadTrack(savedTrack)
  } else {
    nextTrack()
  }
})
