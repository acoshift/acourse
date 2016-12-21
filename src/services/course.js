import Firebase from './firebase'
import { Observable } from 'rxjs/Observable'

import flow from 'lodash/fp/flow'
import map from 'lodash/fp/map'
import values from 'lodash/fp/values'
import identity from 'lodash/fp/identity'
import keys from 'lodash/fp/keys'
import flatMap from 'lodash/fp/flatMap'
import countBy from 'lodash/fp/countBy'
import toPairs from 'lodash/fp/toPairs'
import orderBy from 'lodash/fp/orderBy'
import reverse from 'lodash/fp/reverse'
import reduce from 'lodash/fp/reduce'
import filter from 'lodash/fp/filter'
import isEmpty from 'lodash/fp/isEmpty'

const reduceNoCap = reduce.convert({ cap: false })

export default {
  list () {
    const ref = Firebase.ref('course').orderByChild('public').equalTo(true)
    return Firebase.onArrayValue(ref)
      .map(reverse)
  },
  get (id) {
    return Firebase.onValue(`course/${id}`)
      .flatMap((course) => course ? Observable.of(course) : Observable.throw(new Error('Course not found')))
      .map((course) => ({id, ...course}))
  },
  create (data) {
    data.timestamp = Firebase.timestamp
    return Firebase.push('course', data)
      .map((snapshot) => snapshot.key)
  },
  save (id, data) {
    return Firebase.update(`course/${id}`, data)
  },
  content (id) {
    return Firebase.onArrayValue(`content/${id}`)
  },
  saveContent (id, data) {
    return Firebase.set(`content/${id}`, data)
  },
  favorite (id, userId) {
    return Firebase.set(`course/${id}/favorite/${userId}`, true)
  },
  unfavorite (id, userId) {
    return Firebase.remove(`course/${id}/favorite/${userId}`)
  },
  enroll (id, userId, code) {
    return Firebase.set(`course-private/${id}/enroll/${userId}`, code)
      .flatMap(() => Firebase.set(`course/${id}/student/${userId}`, true))
  },
  ownBy (userId) {
    const ref = Firebase.ref('course').orderByChild('owner').equalTo(userId)
    return Firebase.onArrayValue(ref)
      .map(reverse)
  },
  // sendMessage (id, userId, text) {
  //   return Firebase.push(`chat/${id}`, {
  //     u: userId,
  //     m: text,
  //     t: Firebase.timestamp
  //   })
  // },
  // messages (id, limit) {
  //   let ref = Firebase.ref(`chat/${id}`).orderByKey()
  //   if (limit) {
  //     ref = ref.limitToLast(limit)
  //   }
  //   return Firebase.onChildAdded(ref)
  // },
  // allMessages (id) {
  //   const ref = Firebase.ref(`chat/${id}`).orderByKey()
  //   return Firebase.onArrayValue(ref)
  //     .first()
  // },
  // lastMessage (id) {
  //   const ref = Firebase.ref(`chat/${id}`).orderByKey().limitToLast(1)
  //   return Firebase.onValue(ref)
  // },
  attend (id, userId) {
    return this.get(id)
      .first()
      .map((course) => course.attend)
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
  },
  codes (id) {
    return Firebase.onValue(`course-private/${id}/code`)
      .map(keys)
  },
  saveCodes (id, codes) {
    return Firebase.set(`course-private/${id}/code`, flow(
      filter((x) => !!x),
      reduce((p, v) => { p[v] = true; return p }, {})
    )(codes))
  },
  setQueueEnroll (id, userId, url) {
    return Firebase.set(`queue-enroll/${id}/${userId}`, {
      url,
      timestamp: Firebase.timestamp
    })
  },
  queueEnroll () {
    return Firebase.onValue(`queue-enroll`)
      .map((x) => isEmpty(x) ? [] : flow(
        reduceNoCap((p, v, k) => { p.push({ course: k, users: v }); return p }, []),
        map((x) => ({
          ...x,
          users: reduceNoCap((p, v, k) => { p.push({ id: k, detail: v }); return p }, [])(x.users)
        }))
      )(x))
      .flatMap((xs) =>
        isEmpty(xs)
          ? Observable.of([])
          : Observable.combineLatest(...xs.map((x) =>
              this.get(x.course).first().map((course) => ({ ...x, course })))))
  },
  removeQueueEnroll (id, userId) {
    return Firebase.remove(`queue-enroll/${id}/${userId}`)
  },
  addUser (id, userId) {
    return Firebase.set(`course/${id}/student/${userId}`, true)
  }
}
