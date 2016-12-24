import Firebase from './firebase'
import API from './api'
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
import isEmpty from 'lodash/fp/isEmpty'

const reduceNoCap = reduce.convert({ cap: false })

export default {
  list () {
    return API.get('/course')
  },
  get (url) {
    return API.get(`/course/${url}`)
  },
  create (data) {
    return API.post('/course', data, true)
  },
  save (id, data) {
    return API.patch(`/course/${id}`, data)
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
