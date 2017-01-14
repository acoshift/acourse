import { BehaviorSubject } from 'rxjs/BehaviorSubject'
import RPC from './rpc'
import find from 'lodash/find'
import response from './response'

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
    RPC.invoke('/acourse.CourseService/ListPublicCourses')
      .map(response.courses)
      .subscribe((courses) => {
        $courses.next(courses)
      })
  },
  list () {
    return $courses.asObservable()
  },
  listAll () {
    return RPC.invoke('/acourse.CourseService/ListCourses', null, true)
      .map(response.courses)
  },
  fetch (url) {
    const ob = RPC.invoke('/acourse.CourseService/GetCourse', { courseId: url })
      .map(response.course)
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
    return RPC.invoke('/acourse.CourseService/CreateCourse', data, true)
      .map(({ id }) => id)
      .do((id) => this.fetch(id))
  },
  save (id, data) {
    return RPC.invoke('/acourse.CourseService/UpdateCourse', { id, ...data }, true)
      .do(() => this.fetch(id))
  },
  enroll (courseId, { code, url, price }) {
    return RPC.invoke('/acourse.CourseService/EnrollCourse', { courseId, code, url, price }, true)
  }
}
