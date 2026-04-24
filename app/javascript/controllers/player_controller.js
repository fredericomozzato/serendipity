import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  static targets = [
    "iframe", "title", "artist",
    "nextButton", "prevButton",
    "loading", "error"
  ]

  connect() {
    this.activeGenres = []
    this._currentTrack = null

    this._boundHandleHistoryRestored = this.handleHistoryRestored.bind(this)
    this._boundHandleFilterChanged = this.handleFilterChanged.bind(this)
    this._boundHandleHistoryUpdated = this.handleHistoryUpdated.bind(this)
    this._boundHandleLoadTrack = this.handleLoadTrack.bind(this)

    this.element.addEventListener("history:restored", this._boundHandleHistoryRestored)
    this.element.addEventListener("genre-filter:changed", this._boundHandleFilterChanged)
    this.element.addEventListener("history:updated", this._boundHandleHistoryUpdated)
    this.element.addEventListener("player:load-track", this._boundHandleLoadTrack)
  }

  disconnect() {
    this.element.removeEventListener("history:restored", this._boundHandleHistoryRestored)
    this.element.removeEventListener("genre-filter:changed", this._boundHandleFilterChanged)
    this.element.removeEventListener("history:updated", this._boundHandleHistoryUpdated)
    this.element.removeEventListener("player:load-track", this._boundHandleLoadTrack)
  }

  handleHistoryRestored(event) {
    const { track, history, activeGenres } = event.detail
    this.activeGenres = activeGenres || []
    this.prevButtonTarget.disabled = !history || history.length === 0
    if (track) {
      this.loadTrack(track)
    } else {
      this.nextTrack()
    }
  }

  handleFilterChanged(event) {
    this.activeGenres = event.detail.activeGenres || []
  }

  handleHistoryUpdated(event) {
    this.prevButtonTarget.disabled = event.detail.history.length === 0
  }

  handleLoadTrack(event) {
    this.loadTrack(event.detail.track)
  }

  async nextTrack() {
    const currentTrack = this._currentTrack
    const newTrack = await this.fetchRandomTrack()

    if (newTrack) {
      if (currentTrack) {
        this.element.dispatchEvent(new CustomEvent("player:history-push", {
          bubbles: true,
          detail: { track: currentTrack }
        }))
      }
      this.loadTrack(newTrack)
    }
  }

  previousTrack() {
    this.element.dispatchEvent(new CustomEvent("player:history-pop", {
      bubbles: true
    }))
  }

  async fetchRandomTrack() {
    this.showLoading()

    try {
      const url = new URL("/player/random_track", window.location.origin)
      this.activeGenres.forEach(genre => url.searchParams.append("genres[]", genre))

      const response = await fetch(url)

      if (!response.ok) {
        if (response.status === 404) {
          throw new Error("No playable track found. Please try again.")
        }
        throw new Error("Failed to fetch a random track.")
      }

      return await response.json()
    } catch (error) {
      this.showError(error.message)
      return null
    } finally {
      this.hideLoading()
    }
  }

  loadTrack(track) {
    if (!track) return

    this.iframeTarget.src = `https://www.youtube.com/embed/${track.video_id}?autoplay=1`
    this.titleTarget.textContent = track.title
    this.artistTarget.textContent = track.artist
    this._currentTrack = track

    this.element.dispatchEvent(new CustomEvent("player:track-loaded", {
      bubbles: true,
      detail: { track }
    }))
  }

  showLoading() {
    this.loadingTarget.classList.remove("hidden")
    this.errorTarget.classList.add("hidden")
    this.nextButtonTarget.disabled = true
  }

  hideLoading() {
    this.loadingTarget.classList.add("hidden")
    this.nextButtonTarget.disabled = false
  }

  showError(message) {
    this.hideLoading()
    this.errorTarget.textContent = message
    this.errorTarget.classList.remove("hidden")
  }
}
