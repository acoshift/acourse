import Firebase from './firebase'

export default {
  list () {
    return Firebase.onArrayValue('course')
  },
  get (id) {
    return Firebase.onValue(`course/${id}`)
  }
}
