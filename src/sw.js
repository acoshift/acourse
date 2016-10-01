/* eslint-env serviceworker */

self.addEventListener('install', (event) => {
  self.skipWaiting()
})

self.addEventListener('push', (event, a) => {
  console.log(a)
  event.waitUntil(
    self.registration.showNotification('Acourse')
  )
})
