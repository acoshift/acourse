import API from './api'
import reverse from 'lodash/fp/reverse'

export default {
  list () {
    return API.get('/payment')
  },
  history () {
    return API.get('/payment?history=true')
      .map(reverse)
  },
  approve (id) {
    return API.put(`/payment/${id}/approve`)
  },
  reject (id) {
    return API.put(`/payment/${id}/reject`)
  }
}
