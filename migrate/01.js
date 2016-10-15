// move user's courses from user/{id}/course to user-course

const firebase = require('firebase')
const _ = require('lodash')

firebase.initializeApp({
  serviceAccount: '../.private/acourse-e38f3d69ac48.json',
  databaseURL: 'https://acourse-d9d0a.firebaseio.com'
})

const userRef = firebase.database().ref('user')
const userCourseRef = firebase.database().ref('user-course')

userRef.on('child_added', (snapshot) => {
  const key = snapshot.key
  const user = snapshot.val()
  _.forEach(user.course, (data, x) => {
    userCourseRef.child(`${key}/${x}`).set(data)
  })
})
