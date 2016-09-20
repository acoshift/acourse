import Firebase from './firebase'
import User from './user'
import _ from 'lodash'
import { Observable } from 'rxjs'

export default {
  signIn (email, password) {
    return Firebase.signInWithEmailAndPassword(email, password)
  },
  signInWithFacebook () {
    return Firebase.signInWithFacebook()
      .flatMap((auth) => this.saveUserProfile(auth), (x) => x)
  },
  signInWithGoogle () {
    return Firebase.signInWithGoogle()
      .flatMap((auth) => this.saveUserProfile(auth), (x) => x)
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
  get currentUser () {
    return Firebase.currentUser.filter((x) => x !== undefined)
  },
  saveUserProfile (auth) {
    return User.get(auth.uid)
      .flatMap((user) => {
        const data = _.pick(user, ['name', 'photo'])
        if (!data.name || !data.photo) {
          if (!data.name) {
            data.name = auth.displayName
          }
          if (!data.photo) {
            data.photo = auth.photoURL
          }
          return User.update(auth.uid, data)
        }
        return Observable.of({})
      })
  }
}
