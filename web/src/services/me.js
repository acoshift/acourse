import API from './api'
import Auth from './auth'
import User from './user'
import orderBy from 'lodash/fp/orderBy'

export default {
  get () {
    return Auth.requireUser()
      .flatMap(({ uid }) => API.get(`/user/${uid}`, true))
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
