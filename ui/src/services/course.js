import { BehaviorSubject } from 'rxjs/BehaviorSubject'
import API from './api'
import find from 'lodash/find'

const $courses = new BehaviorSubject(false)
const $course = new BehaviorSubject({})

if (window.$$state) {
  if (window.$$state.courses) {
    $courses.next(window.$$state.courses)
  }
  if (window.$$state.course) {
    const course = window.$$state.course
    $course.first()
      .subscribe(($$course) => {
        $course.next({
          ...$$course,
          [course.id]: {
            ...course,
            $preload: true
          }
        })
      })
  }
}

export default {
  fetchList () {
    API.get('/course').subscribe((courses) => {
      $courses.next(courses)
    })
  },
  list () {
    return $courses.asObservable()
  },
  fetch (url) {
    const ob = API.get(`/course/${url}`)
      .share()
    ob
      .flatMap(() => $course.first(), (x, y) => [x, y])
      .subscribe(([course, $$course]) => {
        $course.next({
          ...$$course,
          [course.id]: course
        })
      }, () => {})
    return ob
  },
  get (url) {
    return $course.asObservable()
      .map((course) => course[url] || find(course, { url }))
      .filter((x) => !!x)
  },
  create (data) {
    return API.post('/course', data, true)
      .map(({ id }) => id)
  },
  save (id, data) {
    return API.patch(`/course/${id}`, data)
  },
  enroll (id, { code, url, price }) {
    return API.put(`/course/${id}/enroll`, { code, url, price }, true)
  }
}
