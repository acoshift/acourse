import { Observable } from 'rxjs/Observable'
import { Subject } from 'rxjs/Subject'

export default {
  $successModal: new Subject(),
  $errorModal: new Subject(),
  uploadModal: null,
  visibilityChanged () {
    return Observable.fromEvent(document, 'visibilitychange')
      .map(() => document.hidden)
  },
  setTitle (after, before) {
    before = before ? `${before} ` : ''
    after = after ? ` ${after}` : ''
    document.title = `${before}Acourse${after}`
  },
  openSuccessModal (title, description) {
    this.$successModal.next({title, description})
  },
  openErrorModal (title, description) {
    this.$errorModal.next({title, description})
  }
}
