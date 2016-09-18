import firebase from 'firebase'
import { Observable, BehaviorSubject } from 'rxjs'

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
    })
  },
  signInWithEmailAndPassword (email, password) {
    return Observable.fromPromise(firebase.auth().signInWithEmailAndPassword(email, password))
  },
  createUserWithEmailAndPassword (email, password) {
    return Observable.fromPromise(firebase.auth().createUserWithEmailAndPassword(email, password))
  },
  signOut () {
    return Observable.fromPromise(firebase.auth().signOut())
  },
  onValue (path) {
    const ref = firebase.database().ref(path)
    return Observable.bindCallback(ref.on.bind(ref))('value')
      .map((snapshot) => snapshot.val())
  },
  upload (path, file) {
    const ref = firebase.storage().ref(path)
    return Observable.fromPromise(ref.put(file))
  },
  set (path, data) {
    const ref = firebase.database().ref(path)
    return Observable.fromPromise(ref.set(data))
  }
}
