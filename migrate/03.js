// migrate course's user code to course-private to protect peek code

const admin = require('firebase-admin')
const _ = require('lodash')

admin.initializeApp({
  credential: admin.credential.applicationDefault(),
  databaseURL: 'https://acourse-d9d0a.firebaseio.com'
})

const database = admin.database()

const courseRef = database.ref('course')
const coursePrivateRef = database.ref('course-private')

courseRef.once('value', (snapshots) => {
  snapshots.forEach((snapshot) => {
    const course = snapshot.val()
    if (course.student) {
      _.forEach(course.student, (v, k) => {
        console.log(k + ' - ' + v)
        coursePrivateRef.child(snapshot.key).child('enroll').child(k).set(v)
        snapshot.ref.child('student').child(k).set(true)
      })
    }
  })
})
