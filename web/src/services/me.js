import { Observable } from 'rxjs/Observable'
import { BehaviorSubject } from 'rxjs/BehaviorSubject'
import API from './api'
import Auth from './auth'
import User from './user'
import orderBy from 'lodash/fp/orderBy'

const $me = new BehaviorSubject(false)

Auth.currentUser()
  .flatMap((user) => user ? API.get(`/user/${user.uid}`, true) : Observable.of(null))
  .subscribe((user) => {
    $me.next(user)
  })

export default {
  fetch () {
    Auth.requireUser()
      .first()
      .flatMap(({ uid }) => API.get(`/user/${uid}`, true))
      .subscribe((user) => {
        $me.next(user)
      })
  },
  get () {
    return $me.asObservable()
      .filter((x) => x !== false)
  },
  ownCourses () {
    return Auth.requireUser()
      .flatMap(({ uid }) => API.get(`/course?owner=${uid}`, true))
      .map(orderBy(['createdAt'], ['desc']))
  },
  courses () {
    return Auth.requireUser()
      .flatMap(({ uid }) => API.get(`/course?student=${uid}`, true))
      .map(orderBy(['createdAt'], ['desc']))
  },
  upload (file) {
    return Auth.currentUser()
      .first()
      .flatMap(({ uid }) => User.upload(uid, file))
  },
  update (data) {
    return Auth.requireUser()
      .flatMap(({ uid }) => API.patch(`/user/${uid}`, data))
  }
}
