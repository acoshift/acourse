import { BehaviorSubject } from 'rxjs/BehaviorSubject'
import API from './api'

const $courses = new BehaviorSubject(false)

if (window.$$state) {
  if (window.$$state.courses) {
    $courses.next(window.$$state.courses)
  }
}

export default {
  fetchList () {
    API.get('/course').subscribe((courses) => {
      $courses.next(courses)
    })
  },
  list () {
    this.fetchList()
    return $courses.asObservable()
  },
  get (url) {
    return API.get(`/course/${url}`)
  },
  create (data) {
    return API.post('/course', data, true)
      .map(({ id }) => id)
  },
  save (id, data) {
    return API.patch(`/course/${id}`, data)
  },
  enroll (id, { code, url }) {
    return API.put(`/course/${id}/enroll`, { code, url }, true)
  }
}
