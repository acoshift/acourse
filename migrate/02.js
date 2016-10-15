// cleanup user's courses

const firebase = require('firebase')

firebase.initializeApp({
  serviceAccount: '../.private/acourse-e38f3d69ac48.json',
  databaseURL: 'https://acourse-d9d0a.firebaseio.com'
})

const userRef = firebase.database().ref('user')

userRef.on('child_added', (snapshot) => {
  snapshot.child('course').remove()
})
