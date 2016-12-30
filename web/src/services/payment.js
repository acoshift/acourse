import API from './api'

export default {
  list () {
    return API.get('/payment')
  },
  history () {
    return API.get('/payment?history=true')
  },
  approve (id) {
    return API.put(`/payment/${id}/approve`)
  },
  reject (id) {
    return API.put(`/payment/${id}/reject`)
  }
}
