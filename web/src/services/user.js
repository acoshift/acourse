import Firebase from './firebase'
import API from './api'
import orderBy from 'lodash/fp/orderBy'

export default {
  get (id) {
    return API.get(`/user/${id}`)
  },
  ownCourses (id) {
    return API.get(`/course?owner=${id}`)
      .map(orderBy(['createdAt'], ['desc']))
  },
  courses (id) {
    return API.get(`/course?student=${id}`)
      .map(orderBy(['createdAt'], ['desc']))
  },
  upload (id, file) {
    return Firebase.upload(`user/${id}/${Date.now()}`, file)
  },
  update (id, data) {
    return Firebase.update(`user/${id}`, data)
  }
}
