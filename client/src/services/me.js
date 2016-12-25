import API from './api'
import Auth from './auth'
import User from './user'
import Course from './course'
import Assignment from './assignment'
import orderBy from 'lodash/fp/orderBy'

export default {
  get () {
    return Auth.requireUser()
      .flatMap(({ uid }) => API.get(`/user/${uid}`, true))
  },
  ownCourses () {
    return Auth.requireUser()
      .flatMap(({ uid }) => API.get(`/course?owner=${uid}`, true))
      .map(orderBy(['createdAt'], ['desc']))
  },
  courses () {
    return Auth.requireUser()
      .flatMap(({ uid }) => API.get(`/course?student=${uid}`, true))
      .map(orderBy(['createdAt'], ['desc']))
  },
  upload (file) {
    return Auth.currentUser()
      .first()
      .flatMap(({ uid }) => User.upload(uid, file))
  },
  update (data) {
    return Auth.requireUser()
      .flatMap(({ uid }) => API.patch(`/user/${uid}`, data))
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
  }
}
