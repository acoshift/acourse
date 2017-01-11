import { Observable } from 'rxjs/Observable'
import { BehaviorSubject } from 'rxjs/BehaviorSubject'
import API from './api'
import RPC from './rpc'
import Auth from './auth'
import User from './user'
import orderBy from 'lodash/fp/orderBy'

const $me = new BehaviorSubject(false)

const fetch = () => {
  Auth.currentUser()
    .first()
    .flatMap((user) => user ? RPC.get('/acourse.UserService/GetMe', true) : Observable.of(null))
    .subscribe((reply) => {
      $me.next(reply.user)
    })
}

const get = () => $me.asObservable().filter((x) => x !== false)
const update = (data) => RPC.post('/acourse.UserService/UpdateMe', { user: data })

Auth.currentUser().subscribe(fetch)

export default {
  fetch,
  get,
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
  update
}
