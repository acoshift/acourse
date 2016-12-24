/* eslint-env serviceworker */
/* globals firebase */

importScripts('https://www.gstatic.com/firebasejs/3.5.0/firebase-app.js')
importScripts('https://www.gstatic.com/firebasejs/3.5.0/firebase-messaging.js')

firebase.initializeApp({
  messagingSenderId: '582047384847'
})

const messaging = firebase.messaging()

messaging.setBackgroundMessageHandler(function (payload) {
  if (!payload.notification) return

  const notificationTitle = payload.notification.title || 'Acourse'
  const notificationOptions = {
    body: payload.notification.body,
    icon: '/static/acourse-120.png'
  }

  return self.registration.showNotification(notificationTitle, notificationOptions)
})
