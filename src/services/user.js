import Firebase from './firebase'
import { Observable } from 'rxjs'
import pick from 'lodash/fp/pick'
import assign from 'lodash/extend'
import keys from 'lodash/fp/keys'

export default {
  get (id) {
    return Firebase.onValue(`user/${id}`)
      .map((user) => ({id, ...user}))
  },
  inject (obj) {
    this.get(obj.id)
      .first()
      .subscribe((user) => assign(obj, user))
  },
  isInstructor (id) {
    return Firebase.onValue(`instructor/${id}`)
      .map((x) => !!x)
  },
  courses (id) {
    return Firebase.onValue(`user-course/${id}`)
      .map(keys)
  },
  getProfile (id) {
    return Observable.combineLatest(
      this.get(id),
      this.isInstructor(id)
    )
      .map(([user, instructor]) => ({id, ...user, instructor}))
  },
  upload (id, file) {
    return Firebase.upload(`user/${id}/${Date.now()}`, file)
  },
  update (id, data) {
    return Firebase.update(`user/${id}`, data)
  },
  addCourse (id, courseId) {
    return Firebase.set(`user-course/${id}/${courseId}`, true)
  },
  saveAuthProfile ({ uid, displayName, photoURL }) {
    return this.get(uid)
      .flatMap((user) => {
        const data = pick(['name', 'photo'])(user)
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
