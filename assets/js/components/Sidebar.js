import BaseComponent from './BaseComponent'

export default class Sidebar extends BaseComponent {
  constructor({ el }) {
    super()

    this.el = document.querySelector('.sidebar')
    this.overlayEl = this.el.querySelector('.overlay')
    this.state = this.getCurrentState()

    this.overlayEl.addEventListener('click', () => this.close())
  }

  getCurrentState() {
    if (this.el.offsetWidth == 0) {
      return 'closed'
    } else {
      return 'open'
    }
  }

  toggle() {
    if (this.isOpen()) {
      this.close()
    } else {
      this.open()
    }
  }

  open() {
    this.transitionIn('opening', 'open').then(() => {
      this.overlayEl.classList.add('visible')
      this.setState('open')
    })
  }

  close() {
    this.transitionOut('opening', 'open').then(() => {
      this.overlayEl.classList.remove('visible')
      this.setState('closed')
    })
  }

  isOpen() {
    return this.state == 'open'
  }

  isClosed() {
    return this.state == 'closed'
  }

  setState(state) {
    this.state = state
  }
}
