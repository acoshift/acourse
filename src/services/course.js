import Firebase from './firebase'
import User from './user'

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
    return Firebase.set(`course/${id}`, data)
  }
}
