import Firebase from './firebase'
import API from './api'

export default {
  get (id) {
    return API.get(`/user/${id}`)
  },
  upload (id, file) {
    return Firebase.upload(`user/${id}/${Date.now()}`, file)
  },
  update (id, data) {
    return Firebase.update(`user/${id}`, data)
  }
}
