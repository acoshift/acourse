import Firebase from './firebase'

export default {
  signIn (email, password) {
    return Firebase.signInWithEmailAndPassword(email, password)
  },
  signUp (email, password) {
    return Firebase.createUserWithEmailAndPassword(email, password)
  },
  signOut () {
    return Firebase.signOut()
  },
  get currentUser () {
    return Firebase.currentUser.filter((x) => x !== undefined)
  }
}
