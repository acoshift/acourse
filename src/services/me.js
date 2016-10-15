import Auth from './auth'
import User from './user'
import Course from './course'
import { Observable } from 'rxjs'

export default {
  get () {
    return Auth.currentUser()
      .flatMap(({ uid }) =>
        Observable.combineLatest(
          User.get(uid),
          User.isInstructor(uid)
        ), ({ uid }, [user, instructor]) => ({id: uid, ...user, instructor}))
  },
  getProfile () {
    return Auth.currentUser()
      .flatMap(({ uid }) => User.getProfile(uid))
  },
  ownCourses () {
    return Auth.currentUser()
      .flatMap(({ uid }) => Course.ownBy(uid))
  },
  courses () {
    return Auth.currentUser()
      .flatMap(({ uid }) => User.courses(uid))
      .flatMap((courseIds) => Observable.combineLatest(...courseIds.map((id) => Course.get(id))))
  }
}
