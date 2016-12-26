import Firebase from './firebase'

export default {
  submit (courseId, userId, assignmentId, url) {
    return Firebase.push(`assignment/${courseId}/user/${userId}/${assignmentId}`, {
      url,
      timestamp: Firebase.timestamp
    })
  },
  addCode (courseId, { title }) {
    return Firebase.push(`assignment/${courseId}/code`, { title, open: true })
  },
  get (courseId) {
    return Firebase.onValue(`assignment/${courseId}`)
  },
  getCode (courseId) {
    return Firebase.onArrayValue(`assignment/${courseId}/code`)
  },
  getUser (courseId, userId) {
    return Firebase.onValue(`assignment/${courseId}/user/${userId}`)
  },
  open (courseId, assignmentId, value) {
    return Firebase.set(`assignment/${courseId}/code/${assignmentId}/open`, value)
  }
}
