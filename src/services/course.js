import Firebase from './firebase'
import User from './user'
import Auth from './auth'

export default {
  list () {
    return Firebase.onArrayValue('course')
  },
  get (id) {
    return Firebase.onValue(`course/${id}`)
  },
  create (data) {
    data.timestamp = Firebase.timestamp
    return Firebase.push('course', data)
      .map((snapshot) => snapshot.key)
      .flatMap((key) =>
        User.me()
          .first()
          .flatMap((user) =>
            User.updateMe({
              ...user,
              course: {
                ...user.course,
                [key]: true
              }
            })
          )
          .map(() => key)
      )
  },
  save (id, data) {
    return Firebase.update(`course/${id}`, data)
  },
  favorite (id) {
    return Auth.currentUser
      .first()
      .flatMap((user) => Firebase.set(`course/${id}/favorite/${user.uid}`, true))
  },
  unfavorite (id) {
    return Auth.currentUser
      .first()
      .flatMap((user) => Firebase.remove(`course/${id}/favorite/${user.uid}`))
  }
}
