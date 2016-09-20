import Firebase from './firebase'
import Auth from './auth'

export default {
  get (id) {
    return Firebase.onValue(`user/${id}`)
  },
  me () {
    return Auth.currentUser
      .flatMap((user) => this.get(user.uid))
      .map((x) => x || {})
  },
  uploadMePhoto (file) {
    return Auth.currentUser
      .first()
      .flatMap((user) => Firebase.upload(`user/${user.uid}-${Date.now()}`, file))
  },
  update (id, data) {
    return Firebase.update(`user/${id}`, data)
  },
  updateMe (data) {
    return Auth.currentUser
      .first()
      .flatMap((user) => this.update(user.uid, data))
  },
  addCourseMe (courseId) {
    return Auth.currentUser
      .first()
      .flatMap((user) => Firebase.set(`user/${user.uid}/course/${courseId}`, true))
  }
}
