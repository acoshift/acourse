import { BehaviorSubject } from 'rxjs/BehaviorSubject'
import API from './api'
import RPC from './rpc'
import find from 'lodash/find'
import map from 'lodash/map'

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

const mapCourseReply = (reply) => map(reply.courses, (course) => {
  const res = {
    ...course,
    owner: reply.users[course.owner]
  }
  if (reply.enrollCount) {
    res.student = reply.enrollCount[course.id]
  }
  return res
})

export default {
  fetchList () {
    RPC.post('/acourse.CourseService/ListCourses', { enrollCount: true })
      .map(mapCourseReply)
      .subscribe((courses) => {
        $courses.next(courses)
      })
  },
  list () {
    return $courses.asObservable()
  },
  listAll () {
    return RPC.post('/acourse.CourseService/ListCourses', { public: false }, true)
      .map(mapCourseReply)
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
      .do((id) => this.fetch(id))
  },
  save (id, data) {
    return API.patch(`/course/${id}`, data)
      .do(() => this.fetch(id))
  },
  enroll (id, { code, url, price }) {
    return API.put(`/course/${id}/enroll`, { code, url, price }, true)
  }
}
