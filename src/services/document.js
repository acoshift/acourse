import { Observable, Subject } from 'rxjs'

export default {
  $successModal: new Subject(),
  $errorModal: new Subject(),
  visibilityChanged () {
    return Observable.fromEvent(document, 'visibilitychange')
      .map(() => document.hidden)
  },
  setTitle (title) {
    document.title = 'Acourse' + (title ? ` ${title}` : '')
  },
  openSuccessModal (title, description) {
    this.$successModal.next({title, description})
  },
  openErrorModal (title, description) {
    this.$errorModal.next({title, description})
  }
}
