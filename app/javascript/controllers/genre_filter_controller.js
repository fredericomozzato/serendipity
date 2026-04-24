import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  static targets = ["tag"]

  connect() {
    this._boundHandleHistoryRestored = this.handleHistoryRestored.bind(this)
    this.element.addEventListener("history:restored", this._boundHandleHistoryRestored)
  }

  disconnect() {
    this.element.removeEventListener("history:restored", this._boundHandleHistoryRestored)
  }

  handleHistoryRestored(event) {
    const { activeGenres } = event.detail
    this.tagTargets.forEach(tag => {
      if (activeGenres?.includes(tag.dataset.genre)) {
        tag.classList.add("active")
      } else {
        tag.classList.remove("active")
      }
    })
    this.notifyFilterChange()
  }

  toggle(event) {
    event.currentTarget.classList.toggle("active")
    this.notifyFilterChange()
  }

  notifyFilterChange() {
    const activeGenres = this.tagTargets
      .filter(tag => tag.classList.contains("active"))
      .map(tag => tag.dataset.genre)

    this.element.dispatchEvent(new CustomEvent("genre-filter:changed", {
      bubbles: true,
      detail: { activeGenres }
    }))
  }
}
