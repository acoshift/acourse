import firebase from 'firebase'
import { Observable } from 'rxjs/Observable'
import { Subject } from 'rxjs/Subject'
import { BehaviorSubject } from 'rxjs/BehaviorSubject'
import isString from 'lodash/fp/isString'
import Raven from 'raven-js'

const off = (ref, type, fn) => {
  // delay for reuse cached data
  setTimeout(() => {
    ref.off(type, fn)
  }, 10000)
}

const $currentUser = new BehaviorSubject(undefined)

export default {
  currentUser: $currentUser.filter((x) => x !== undefined).asObservable(),
  notificationToken: new BehaviorSubject(),
  notification: new Subject(),
  provider: {
    google: new firebase.auth.GoogleAuthProvider(),
    facebook: new firebase.auth.FacebookAuthProvider(),
    github: new firebase.auth.GithubAuthProvider()
  },

  init () {
    firebase.initializeApp(process.env.FIREBASE)

    const messaging = firebase.messaging()

    if (window.Notification && window.Notification.requestPermission) {
      messaging.requestPermission()
    }

    // messaging.getToken().then((token) => {
    //   this.notificationToken.next(token)
    // })

    // messaging.onTokenRefresh(() => {
    //   messaging.getToken().then((token) => {
    //     this.notificationToken.next(token)
    //   })
    // })

    // messaging.onMessage((payload) => {
    //   this.notification.next(payload)
    // })

    firebase.auth().onAuthStateChanged((user) => {
      $currentUser.next(user)
      if (user) {
        window.ga('set', 'userId', user.uid)
        Raven.setUserContext({
          id: user.uid,
          email: user.email
        })
      } else {
        window.ga('set', 'userId', null)
        Raven.setUserContext(null)
      }
    })

    // Observable.combineLatest(
    //   $currentUser,
    //   this.notificationToken
    // )
    //   .subscribe(
    //     ([user, token]) => {
    //       if (token) {
    //         if (user) {
    //           firebase.database().ref(`notification/${token}`).set({ user: user.uid })
    //         } else {
    //           firebase.database().ref(`notification/${token}`).remove()
    //         }
    //       }
    //     }
    //   )
  },
  getToken () {
    return this.currentUser
      .flatMap((user) => user ? Observable.fromPromise(user.getToken()) : Observable.of(''))
  },
  signInWithEmailAndPassword (email, password) {
    return Observable.fromPromise(firebase.auth().signInWithEmailAndPassword(email, password))
  },
  signInWithProvider (provider) {
    return Observable.fromPromise(firebase.auth().signInWithPopup(provider))
  },
  linkProvider (provider) {
    return Observable.fromPromise(firebase.auth().currentUser.linkWithPopup(provider))
      .do(() => {
        $currentUser.next(firebase.auth().currentUser)
      })
  },
  unlinkProvider (provider) {
    return Observable.fromPromise(firebase.auth().currentUser.unlink(provider.providerId))
      .do(() => {
        $currentUser.next(firebase.auth().currentUser)
      })
  },
  createUserWithEmailAndPassword (email, password) {
    return Observable.fromPromise(firebase.auth().createUserWithEmailAndPassword(email, password))
  },
  sendPasswordResetEmail (email) {
    return Observable.fromPromise(firebase.auth().sendPasswordResetEmail(email))
  },
  signOut () {
    return Observable.fromPromise(firebase.auth().signOut())
  },
  onValue (ref) {
    return Observable.create((o) => {
      ref = isString(ref) ? this.ref(ref) : ref
      const fn = ref.on('value', (snapshot) => {
        o.next(snapshot.val())
      }, (err) => { o.error(err) })
      return () => off(ref, 'value', fn)
    })
  },
  onceValue (ref) {
    return Observable.create((o) => {
      ref = isString(ref) ? this.ref(ref) : ref
      ref.once('value', (snapshot) => {
        o.next(snapshot.val())
        o.complete()
      }, (err) => { o.error(err) })
    })
  },
  onChildAdded (ref) {
    return Observable.create((o) => {
      ref = isString(ref) ? this.ref(ref) : ref
      const fn = ref.on('child_added', (snapshot) => {
        o.next(snapshot.val())
      }, (err) => { o.error(err) })
      return () => off(ref, 'child_added', fn)
    })
  },
  onArrayValue (ref) {
    return Observable.create((o) => {
      ref = isString(ref) ? this.ref(ref) : ref
      const fn = ref.on('value', (snapshots) => {
        const result = []
        snapshots.forEach((snapshot) => {
          result.push({
            id: snapshot.key,
            ...snapshot.val()
          })
        })
        o.next(result)
      }, (err) => { o.error(err) })
      return () => off(ref, 'value', fn)
    })
  },
  upload (path, file) {
    const ref = firebase.storage().ref(path)
    return Observable.of(ref.put(file))
  },
  set (path, data) {
    const ref = firebase.database().ref(path)
    return Observable.fromPromise(ref.set(data))
  },
  update (path, data) {
    const ref = firebase.database().ref(path)
    return Observable.fromPromise(ref.update(data))
  },
  push (path, data) {
    const ref = firebase.database().ref(path)
    return Observable.fromPromise(ref.push(data))
  },
  remove (path) {
    const ref = firebase.database().ref(path)
    return Observable.fromPromise(ref.remove())
  },
  get timestamp () {
    return firebase.database.ServerValue.TIMESTAMP
  },
  ref (path) {
    return firebase.database().ref(path)
  }
}
