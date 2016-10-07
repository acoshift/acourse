import Firebase from './firebase'

export default {
  signIn (email, password) {
    window.ga('send', 'event', 'auth', 'signin')
    return Firebase.signInWithEmailAndPassword(email, password)
  },
  signInWithFacebook () {
    window.ga('send', 'event', 'auth', 'facebook')
    return Firebase.signInWithFacebook()
  },
  signInWithGoogle () {
    window.ga('send', 'event', 'auth', 'google')
    return Firebase.signInWithGoogle()
  },
  signUp (email, password) {
    window.ga('send', 'event', 'auth', 'signup')
    return Firebase.createUserWithEmailAndPassword(email, password)
  },
  signOut () {
    window.ga('send', 'event', 'auth', 'signout')
    return Firebase.signOut()
  },
  resetPassword (email) {
    window.ga('send', 'event', 'auth', 'resetPassword')
    return Firebase.sendPasswordResetEmail(email)
  },
  currentUser () {
    return Firebase.currentUser.filter((x) => x !== undefined)
  }
}
