import Firebase from './firebase'

export default {
  signIn (email, password) {
    return Firebase.signInWithEmailAndPassword(email, password)
  },
  signInWithFacebook () {
    return Firebase.signInWithFacebook()
  },
  signInWithGoogle () {
    return Firebase.signInWithGoogle()
  },
  signUp (email, password) {
    return Firebase.createUserWithEmailAndPassword(email, password)
  },
  signOut () {
    return Firebase.signOut()
  },
  resetPassword (email) {
    return Firebase.sendPasswordResetEmail(email)
  },
  currentUser () {
    return Firebase.currentUser.filter((x) => x !== undefined)
  }
}
