import firebase from 'firebase'
import { Observable, BehaviorSubject } from 'rxjs'
import _ from 'lodash'

export default {
  init () {
    firebase.initializeApp({
      apiKey: 'AIzaSyBa6jPVcSPbLzjuei9-d0u9q7M-9rGjmb8',
      authDomain: 'acourse-d9d0a.firebaseapp.com',
      databaseURL: 'https://acourse-d9d0a.firebaseio.com',
      storageBucket: 'acourse-d9d0a.appspot.com',
      messagingSenderId: '582047384847'
    })
    this.currentUser = new BehaviorSubject()

    firebase.auth().onAuthStateChanged((user) => {
      this.currentUser.next(user)
      if (user) {
        window.ga('set', 'userId', user.uid)
        this.ref('user').on('value', () => {})
      }
    })
  },
  signInWithEmailAndPassword (email, password) {
    return Observable.fromPromise(firebase.auth().signInWithEmailAndPassword(email, password))
  },
  signInWithFacebook () {
    return Observable.fromPromise(firebase.auth().signInWithPopup(new firebase.auth.FacebookAuthProvider()))
  },
  signInWithGoogle () {
    return Observable.fromPromise(firebase.auth().signInWithPopup(new firebase.auth.GoogleAuthProvider()))
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
      ref = _.isString(ref) ? this.ref(ref) : ref
      ref.on('value', (snapshot) => {
        o.next(snapshot.val())
      })
    })
  },
  onceValue (ref) {
    return Observable.create((o) => {
      ref = _.isString(ref) ? this.ref(ref) : ref
      ref.once('value', (snapshot) => {
        o.next(snapshot.val())
        o.complete()
      })
    })
  },
  onChildAdded (ref) {
    return Observable.create((o) => {
      ref = _.isString(ref) ? this.ref(ref) : ref
      ref.on('child_added', (snapshot) => {
        o.next(snapshot.val())
      })
    })
  },
  onArrayValue (ref) {
    return Observable.create((o) => {
      ref = _.isString(ref) ? this.ref(ref) : ref
      ref.on('value', (snapshots) => {
        const result = []
        snapshots.forEach((snapshot) => {
          result.push({
            id: snapshot.key,
            ...snapshot.val()
          })
        })
        o.next(result)
      })
    })
  },
  upload (path, file) {
    const ref = firebase.storage().ref(path)
    return Observable.fromPromise(ref.put(file))
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
