import Firebase from './firebase'
import Auth from './auth'
import { Observable } from 'rxjs'
import _ from 'lodash'

export default {
  get (id) {
    return Firebase.onValue(`user/${id}`)
      .map((user) => ({id, ...user}))
  },
  getOnce (id) {
    return Firebase.onceValue(`user/${id}`)
      .map((user) => ({id, ...user}))
  },
  isInstructor (id) {
    return Firebase.onValue(`instructor/${id}`)
      .map((x) => !!x)
  },
  me () {
    return Auth.currentUser()
      .flatMap((auth) =>
        Observable.combineLatest(
          this.get(auth.uid),
          this.isInstructor(auth.uid)
        ), (auth, [user, instructor]) => ({id: auth.uid, ...user, instructor}))
  },
  getProfileAndInstructor (id) {
    return Observable.combineLatest(
      this.get(id),
      this.isInstructor(id)
    )
      .map(([user, instructor]) => ({id, ...user, instructor}))
  },
  uploadMePhoto (file) {
    return Auth.currentUser()
      .first()
      .flatMap(({ uid }) => Firebase.upload(`user/${uid}/${Date.now()}`, file))
  },
  upload (file) {
    return Auth.currentUser()
      .first()
      .flatMap(({ uid }) => Firebase.upload(`user/${uid}/${Date.now()}`, file))
  },
  update (id, data) {
    return Firebase.update(`user/${id}`, data)
  },
  updateMe (data) {
    return Auth.currentUser()
      .first()
      .flatMap((user) => this.update(user.uid, data))
  },
  addCourseMe (courseId) {
    return Auth.currentUser()
      .first()
      .flatMap((user) => Firebase.set(`user/${user.uid}/course/${courseId}`, true))
  },
  saveAuthProfile ({ uid, displayName, photoURL }) {
    return this.get(uid)
      .flatMap((user) => {
        const data = _.pick(user, ['name', 'photo'])
        if (!data.name || !data.photo) {
          if (!data.name) {
            data.name = displayName
          }
          if (!data.photo) {
            data.photo = photoURL
          }
          return this.update(uid, data)
        }
        return Observable.of({})
      })
  }
}
