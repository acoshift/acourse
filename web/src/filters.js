import Vue from 'vue'
import moment from 'moment'

const time = new Vue({
  data () {
    return {
      now: null
    }
  },
  created () {
    setInterval(() => {
      this.now = Date.now()
    }, 60000)
  }
})

Vue.filter('date', (value, input) => {
  if (!value) return '-'
  const m = moment(value)
  // check isZero
  if (m.format('YYYY-MM-DD') === '0001-01-01') {
    return '-'
  }
  return moment(value).format(input)
})

Vue.filter('fromNow', (value) => {
  time.now
  if (!value) return '-'
  return moment(value).fromNow()
})

Vue.filter('trim', (value, input) => {
  value = value || ''
  if (value.length <= input) return value
  return value.substr(0, input) + '...'
})

Vue.filter('money', (value) => isFinite(value) ? value.toFixed(0) : '-')
