const firebase = require('firebase')
const _ = require('lodash')
const request = require('request')

firebase.initializeApp({
  serviceAccount: './.private/acourse-e38f3d69ac48.json',
  databaseURL: 'https://acourse-d9d0a.firebaseio.com'
})

const notiRef = firebase.database().ref('notification')

notiRef.once('value', (snapshot) => {
  const val = snapshot.val()
  const keys = _.keys(val)
  request.post('https://fcm.googleapis.com/fcm/send', {
    headers: {
      Authorization: 'key='
    },
    json: {
      notification: {
        title: 'Acourse',
        body: '',
        click_action: ''
      },
      data: {
        color: 'info',
        body: '',
        timeout: 10000
      },
      registration_ids: keys
    }
  }, (err, resp, body) => {
    if (err) return console.log(err)
    console.log(body)
    process.exit(0)
  })
})
