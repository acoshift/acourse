import Firebase from './firebase'

import flow from 'lodash/fp/flow'
import map from 'lodash/fp/map'
import values from 'lodash/fp/values'
import identity from 'lodash/fp/identity'
import keys from 'lodash/fp/keys'
import flatMap from 'lodash/fp/flatMap'
import countBy from 'lodash/fp/countBy'
import toPairs from 'lodash/fp/toPairs'
import orderBy from 'lodash/fp/orderBy'

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
    window.ga('send', 'event', 'course', 'create')
    data.timestamp = Firebase.timestamp
    return Firebase.push('course', data)
      .map((snapshot) => snapshot.key)
  },
  save (id, data) {
    window.ga('send', 'event', 'course', 'save', id)
    return Firebase.update(`course/${id}`, data)
  },
  content (id) {
    return Firebase.onArrayValue(`content/${id}`)
  },
  saveContent (id, data) {
    return Firebase.set(`content/${id}`, data)
  },
  favorite (id, userId) {
    window.ga('send', 'event', 'course', 'favorite', id)
    return Firebase.set(`course/${id}/favorite/${userId}`, true)
  },
  unfavorite (id, userId) {
    window.ga('send', 'event', 'course', 'unfavorite', id)
    return Firebase.remove(`course/${id}/favorite/${userId}`)
  },
  addStudent (id, userId) {
    window.ga('send', 'event', 'course', 'apply', id)
    return Firebase.set(`course/${id}/student/${userId}`, true)
  },
  ownBy (userId) {
    const ref = Firebase.ref('course').orderByChild('owner').equalTo(userId)
    return Firebase.onArrayValue(ref)
  },
  sendMessage (id, userId, text) {
    window.ga('send', 'event', 'course', 'sendMessage', id)
    return Firebase.push(`chat/${id}`, {
      u: userId,
      m: text,
      t: Firebase.timestamp
    })
  },
  messages (id, limit) {
    let ref = Firebase.ref(`chat/${id}`).orderByKey()
    if (limit) {
      ref = ref.limitToLast(limit)
    }
    return Firebase.onChildAdded(ref)
  },
  lastMessage (id) {
    const ref = Firebase.ref(`chat/${id}`).orderByKey().limitToLast(1)
    return Firebase.onValue(ref)
  },
  attend (id, userId) {
    return this.get(id)
      .first()
      .map((course) => course.attend)
      .do((code) => window.ga('send', 'event', 'course', 'attend', id, code))
      .flatMap((code) => Firebase.set(`attend/${id}/${code}/${userId}`, Firebase.timestamp))
  },
  setAttendCode (id, code) {
    return Firebase.set(`course/${id}/attend`, code)
  },
  removeAttendCode (id) {
    return Firebase.set(`course/${id}/attend`, null)
  },
  attendUsers (id) {
    return Firebase.onValue(`attend/${id}`)
      .map(
        flow(
          values,
          flatMap(keys),
          countBy(identity),
          toPairs,
          map((x) => ({ id: x[0], count: x[1] })),
          orderBy('count', 'desc')
        )
      )
  },
  isAttended (id, userId) {
    return this.get(id)
      .map((course) => course.attend)
      .flatMap((code) => Firebase.onValue(`attend/${id}/${code}/${userId}`))
      .map((x) => !!x)
  }
}
