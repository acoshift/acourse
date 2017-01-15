import { BehaviorSubject } from 'rxjs/BehaviorSubject'
import RPC from './rpc'
import find from 'lodash/find'
import response from './response'

const $courses = new BehaviorSubject(false)
const $course = new BehaviorSubject({})

if (window.$$state) {
  if (window.$$state.courses) {
    $courses.next(response.courses(window.$$state.courses))
  }
  if (window.$$state.course) {
    const course = response.course(window.$$state.course)
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

const fetchList = () => {
  RPC.invoke('/acourse.CourseService/ListPublicCourses')
    .map(response.courses)
    .subscribe((courses) => {
      $courses.next(courses)
    })
}

const list = () =>
  $courses.asObservable()

const listAll = () =>
  RPC.invoke('/acourse.CourseService/ListCourses', null, true)
    .map(response.courses)

const fetch = (url) => {
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
}

const get = (url) =>
  $course.asObservable()
    .map((course) => course[url] || find(course, { url }))
    .filter((x) => !!x)

const create = (data) =>
  RPC.invoke('/acourse.CourseService/CreateCourse', data, true)
    .map(({ id }) => id)

const save = (id, data) =>
  RPC.invoke('/acourse.CourseService/UpdateCourse', { id, ...data }, true)
    .do(() => fetch(id))

const enroll = (courseId, { code, url, price }) =>
  RPC.invoke('/acourse.CourseService/EnrollCourse', { courseId, code, url, price }, true)

const openAttend = (courseId) =>
  RPC.invoke('/acourse.CourseService/OpenAttend', { courseId }, true)
    .flatMap(() => fetch(courseId))

const closeAttend = (courseId) =>
  RPC.invoke('/acourse.CourseService/CloseAttend', { courseId }, true)
    .flatMap(() => fetch(courseId))

export default {
  fetchList,
  list,
  listAll,
  fetch,
  get,
  create,
  save,
  enroll,
  openAttend,
  closeAttend
}
