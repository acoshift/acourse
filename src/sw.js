/* eslint-env serviceworker */

console.log(self)

self.addEventListener('install', (event) => {
  self.skipWaiting()
})

self.addEventListener('push', (event) => {
  console.log(event)
})
