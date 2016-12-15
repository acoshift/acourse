import Auth from './auth'
import User from './user'
import Course from './course'
import Assignment from './assignment'
import { Observable } from 'rxjs/Observable'
import orderBy from 'lodash/fp/orderBy'

export default {
  get () {
    return Auth.currentUser()
      .flatMap(({ uid }) =>
        Observable.combineLatest(
          User.get(uid),
          User.isInstructor(uid)
        ), ({ uid }, [user, instructor]) => ({id: uid, ...user, instructor}))
  },
  getProfile () {
    return Auth.currentUser()
      .flatMap(({ uid }) => User.getProfile(uid))
  },
  ownCourses () {
    return Auth.currentUser()
      .flatMap(({ uid }) => Course.ownBy(uid))
  },
  courses () {
    return Auth.currentUser()
      .flatMap(({ uid }) => User.courses(uid))
      .flatMap((courseIds) => Observable.combineLatest(...courseIds.map((id) => Course.get(id))))
      .map(orderBy(['timestamp'], ['desc']))
  },
  upload (file) {
    return Auth.currentUser()
      .first()
      .flatMap(({ uid }) => User.upload(uid, file))
  },
  update (data) {
    return Auth.currentUser()
      .first()
      .flatMap(({ uid }) => User.update(uid, data))
  },
  addCourse (id) {
    return Auth.currentUser()
      .first()
      .flatMap(({ uid }) => User.addCourse(uid, id))
  },
  applyCourse (id, code) {
    code = code || true
    return Auth.currentUser()
      .first()
      .flatMap(({ uid }) => Course.enroll(id, uid, code))
      .flatMap(() => this.addCourse(id))
  },
  favoriteCourse (id) {
    return Auth.currentUser()
      .first()
      .flatMap(({ uid }) => Course.favorite(id, uid))
  },
  unfavoriteCourse (id) {
    return Auth.currentUser()
      .first()
      .flatMap(({ uid }) => Course.unfavorite(id, uid))
  },
  isAttendedCourse (id) {
    return Auth.currentUser()
      .first()
      .flatMap(({ uid }) => Course.isAttended(id, uid))
  },
  attendCourse (id) {
    return Auth.currentUser()
      .first()
      .flatMap(({ uid }) => Course.attend(id, uid))
  },
  sendMessage (id, text) {
    return Auth.currentUser()
      .first()
      .flatMap(({ uid }) => Course.sendMessage(id, uid, text))
  },
  submitCourseAssignment (id, assignmentId, url) {
    return Auth.currentUser()
      .first()
      .flatMap(({ uid }) => Assignment.submit(id, uid, assignmentId, url))
  },
  getCourseAssignments (id) {
    return Auth.currentUser()
      .first()
      .flatMap(({ uid }) => Assignment.getUser(id, uid))
  },
  submitCourseQueueEnroll (id, url) {
    return Auth.currentUser()
      .first()
      .flatMap(({ uid }) => Course.setQueueEnroll(id, uid, url))
  }
}
