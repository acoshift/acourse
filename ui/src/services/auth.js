import Firebase from './firebase'

export default {
  signIn (email, password) {
    window.ga('send', 'event', 'auth', 'signin')
    return Firebase.signInWithEmailAndPassword(email, password)
  },
  signInWithFacebook () {
    window.ga('send', 'event', 'auth', 'facebook')
    return Firebase.signInWithProvider(Firebase.provider.facebook)
  },
  signInWithGoogle () {
    window.ga('send', 'event', 'auth', 'google')
    return Firebase.signInWithProvider(Firebase.provider.google)
  },
  signInWithGithub () {
    window.ga('send', 'event', 'auth', 'github')
    return Firebase.signInWithProvider(Firebase.provider.github)
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
  },
  requireUser () {
    return Firebase.currentUser.filter((x) => !!x)
  },
  linkFacebook () {
    window.ga('send', 'event', 'auth', 'link', 'facebook')
    return Firebase.linkProvider(Firebase.provider.facebook)
  },
  linkGoogle () {
    window.ga('send', 'event', 'auth', 'link', 'google')
    return Firebase.linkProvider(Firebase.provider.google)
  },
  linkGithub () {
    window.ga('send', 'event', 'auth', 'link', 'github')
    return Firebase.linkProvider(Firebase.provider.github)
  },
  unlinkFacebook () {
    window.ga('send', 'event', 'auth', 'unlink', 'facebook')
    return Firebase.unlinkProvider(Firebase.provider.facebook)
  },
  unlinkGoogle () {
    window.ga('send', 'event', 'auth', 'unlink', 'google')
    return Firebase.unlinkProvider(Firebase.provider.google)
  },
  unlinkGithub () {
    window.ga('send', 'event', 'auth', 'unlink', 'github')
    return Firebase.unlinkProvider(Firebase.provider.github)
  }
}
