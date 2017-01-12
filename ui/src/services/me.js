import { Observable } from 'rxjs/Observable'
import { BehaviorSubject } from 'rxjs/BehaviorSubject'
import RPC from './rpc'
import Auth from './auth'
import User from './user'
import orderBy from 'lodash/fp/orderBy'
import map from 'lodash/map'

const $me = new BehaviorSubject(false)

const fetch = () => {
  Auth.currentUser()
    .first()
    .flatMap((user) => user ? RPC.get('/acourse.UserService/GetMe', true) : Observable.of(null))
    .subscribe((reply) => {
      if (reply && reply.user) {
        $me.next({
          ...reply.user,
          role: reply.role
        })
      } else {
        $me.next(null)
      }
    })
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

const get = () => $me.asObservable().filter((x) => x !== false)
const update = (data) => RPC.post('/acourse.UserService/UpdateMe', { user: data })
const ownCourses = () => Auth.requireUser()
  .flatMap(({ uid }) => RPC.post('/acourse.CourseService/ListCourses', { owner: uid, enrollCount: true }, true))
  .map(mapCourseReply)
  .map(orderBy(['createdAt'], ['desc']))
const courses = () => Auth.requireUser()
  .flatMap(({ uid }) => RPC.get(`/acourse.CourseService/ListEnrolledCourses`, true))
  .map(mapCourseReply)
  .map(orderBy(['createdAt'], ['desc']))

Auth.currentUser().subscribe(fetch)

export default {
  fetch,
  get,
  ownCourses,
  courses,
  upload (file) {
    return Auth.currentUser()
      .first()
      .flatMap(({ uid }) => User.upload(uid, file))
  },
  update
}
