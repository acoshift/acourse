import { Observable } from 'rxjs/Observable'
import { BehaviorSubject } from 'rxjs/BehaviorSubject'
import RPC from './rpc'
import Auth from './auth'
import User from './user'
import orderBy from 'lodash/fp/orderBy'
import response from './response'

const $me = new BehaviorSubject(false)

const fetch = () => {
  Auth.currentUser()
    .first()
    .flatMap((user) => user ? RPC.invoke('/acourse.UserService/GetMe', true) : Observable.of(null))
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

const get = () => $me.asObservable().filter((x) => x !== false)
const update = (data) => RPC.invoke('/acourse.UserService/UpdateMe', data)
const ownCourses = () => Auth.requireUser()
  .flatMap(({ uid }) => RPC.invoke('/acourse.CourseService/ListOwnCourses', { userId: uid }, true))
  .map(response.courses)
  .map(orderBy(['createdAt'], ['desc']))
const courses = () => Auth.requireUser()
  .flatMap(({ uid }) => RPC.invoke(`/acourse.CourseService/ListEnrolledCourses`, { userId: uid }, true))
  .map(response.courses)
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
