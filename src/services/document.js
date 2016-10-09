import { Observable } from 'rxjs'

export default {
  visibilityChanged () {
    return Observable.fromEvent(document, 'visibilitychange')
      .map(() => document.hidden)
  },
  setTitle (title) {
    document.title = 'Acourse' + (title ? ` ${title}` : '')
  }
}
