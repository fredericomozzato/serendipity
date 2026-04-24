import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  static values = {
    currentKey: { type: String, default: "player:current_track" },
    historyKey: { type: String, default: "player:history" },
    filtersKey: { type: String, default: "player:active_genres" }
  }

  connect() {
    this._boundHandleFilterChanged = this.handleFilterChanged.bind(this)
    this._boundHandleTrackLoaded = this.handleTrackLoaded.bind(this)
    this._boundHandleHistoryPush = this.handleHistoryPush.bind(this)
    this._boundHandleHistoryPop = this.handleHistoryPop.bind(this)

    this.element.addEventListener("genre-filter:changed", this._boundHandleFilterChanged)
    this.element.addEventListener("player:track-loaded", this._boundHandleTrackLoaded)
    this.element.addEventListener("player:history-push", this._boundHandleHistoryPush)
    this.element.addEventListener("player:history-pop", this._boundHandleHistoryPop)

    const track = this.getCurrentTrack()
    const history = this.getHistory()
    const activeGenres = this.getActiveGenres()

    requestAnimationFrame(() => {
      this.element.dispatchEvent(new CustomEvent("history:restored", {
        bubbles: true,
        detail: { track, history, activeGenres }
      }))
    })
  }

  disconnect() {
    this.element.removeEventListener("genre-filter:changed", this._boundHandleFilterChanged)
    this.element.removeEventListener("player:track-loaded", this._boundHandleTrackLoaded)
    this.element.removeEventListener("player:history-push", this._boundHandleHistoryPush)
    this.element.removeEventListener("player:history-pop", this._boundHandleHistoryPop)
  }

  handleFilterChanged(event) {
    this.setActiveGenres(event.detail.activeGenres || [])
  }

  handleTrackLoaded(event) {
    this.setCurrentTrack(event.detail.track)
  }

  handleHistoryPush(event) {
    const history = this.getHistory()
    history.push(event.detail.track)
    if (history.length > 50) history.shift()
    this.setHistory(history)
    this.dispatchHistoryUpdated(history)
  }

  handleHistoryPop(event) {
    const history = this.getHistory()
    if (history.length === 0) return

    const prevTrack = history.pop()
    this.setHistory(history)
    this.dispatchHistoryUpdated(history)

    this.element.dispatchEvent(new CustomEvent("player:load-track", {
      bubbles: true,
      detail: { track: prevTrack }
    }))
  }

  dispatchHistoryUpdated(history) {
    this.element.dispatchEvent(new CustomEvent("history:updated", {
      bubbles: true,
      detail: { history }
    }))
  }

  getCurrentTrack() {
    try {
      const raw = localStorage.getItem(this.currentKeyValue)
      return raw ? JSON.parse(raw) : null
    } catch {
      return null
    }
  }

  setCurrentTrack(track) {
    localStorage.setItem(this.currentKeyValue, JSON.stringify(track))
  }

  getHistory() {
    try {
      const raw = localStorage.getItem(this.historyKeyValue)
      return raw ? JSON.parse(raw) : []
    } catch {
      return []
    }
  }

  setHistory(history) {
    localStorage.setItem(this.historyKeyValue, JSON.stringify(history))
  }

  getActiveGenres() {
    try {
      const raw = localStorage.getItem(this.filtersKeyValue)
      return raw ? JSON.parse(raw) : []
    } catch {
      return []
    }
  }

  setActiveGenres(genres) {
    localStorage.setItem(this.filtersKeyValue, JSON.stringify(genres))
  }
}
