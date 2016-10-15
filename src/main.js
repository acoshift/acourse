import Vue from 'vue'
import VueRx from 'vue-rx'
import Raven from 'raven-js'
import App from './app'
import './filters'

import 'semantic-ui-css/components/accordion.min.js'
import 'semantic-ui-css/components/checkbox.min.js'
import 'semantic-ui-css/components/dimmer.min.js'
import 'semantic-ui-css/components/dropdown.min.js'
import 'semantic-ui-css/components/embed.min.js'
import 'semantic-ui-css/components/modal.min.js'
import 'semantic-ui-css/components/transition.min.js'
import 'semantic-ui-css/components/progress.min.js'

import { Firebase } from './services'
import router from './router'

Raven
  .config('https://fda9f1b21cd04a4585b9f9051b37a466@sentry.io/103020')
  .install()

Vue.use(VueRx)

Firebase.init()

/* eslint-disable no-new */
new Vue({
  router,
  render: (h) => h(App)
}).$mount('app')
