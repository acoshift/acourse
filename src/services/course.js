import Firebase from './firebase'
import User from './user'
import Auth from './auth'
import { Observable } from 'rxjs'
import _ from 'lodash'

export default {
  list () {
    const ref = Firebase.ref('course').orderByChild('open').equalTo(true)
    return Firebase.onArrayValue(ref)
  },
  get (id) {
    return Firebase.onValue(`course/${id}`)
      .map((course) => ({id, ...course}))
  },
  create (data) {
    data.timestamp = Firebase.timestamp
    return Firebase.push('course', data)
      .map((snapshot) => snapshot.key)
      // .flatMap((key) =>
      //   User.me()
      //     .first()
      //     .flatMap((user) =>
      //       User.updateMe({
      //         ...user,
      //         course: {
      //           ...user.course,
      //           [key]: true
      //         }
      //       })
      //     )
      //     .map(() => key)
      // )
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
  },
  join (id) {
    return Auth.currentUser
      .first()
      .flatMap((user) =>
        Observable.forkJoin(
          Firebase.set(`course/${id}/student/${user.uid}`, true),
          User.addCourseMe(id)
        )
      )
  },
  ownBy (userId) {
    const ref = Firebase.ref('course').orderByChild('owner').equalTo(userId)
    return Observable.combineLatest(
      Auth.currentUser.first(),
      Firebase.onArrayValue(ref)
    )
      .map(([auth, courses]) => auth.uid === userId ? courses : _.filter(courses, (course) => course.open))
  },
  sendMessage (id, text) {
    return Auth.currentUser
      .first()
      .flatMap((auth) => Firebase.push(`chat/${id}`, {
        u: auth.uid,
        m: text,
        t: Firebase.timestamp
      }))
  },
  messages (id) {
    const ref = Firebase.ref(`chat/${id}`).orderByKey()
    return Firebase.onChildAdded(ref)
  },
  attend (id, code) {
    return Auth.currentUser
      .first()
      .flatMap((auth) => Firebase.set(`attend/${id}/${code}/${auth.uid}`, Firebase.timestamp))
  },
  setAttendCode (id, code) {
    return Firebase.set(`course/${id}/attend`, code)
  },
  removeAttendCode (id) {
    return Firebase.set(`course/${id}/attend`, null)
  },
  attendUsers (id) {
    return Firebase.onValue(`attend/${id}/user`)
      .map((users) => _.mapValues(users, (x, id) => ({id, count: _.keys(x).length})))
      .map(_.values)
      .flatMap((users) =>
        Observable.from(users)
          .flatMap((user) => User.get(user.id).first(), (user, data) => ({...user, ...data}))
          .toArray()
      )
  },
  isAttended (id) {
    return Observable.forkJoin(
      Auth.currentUser.first(),
      this.get(id).first().map((course) => course.attend)
    )
      .flatMap(([auth, code]) => Firebase.onValue(`attend/${id}/${code}/${auth.uid}`))
      .map((x) => !!x)
  }
}
