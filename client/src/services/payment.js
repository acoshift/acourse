import API from './api'

export default {
  list () {
    return API.get('/payment')
  },
  approve (id) {
    return API.put(`/payment/${id}/approve`)
  },
  reject (id) {
    return API.put(`/payment/${id}/reject`)
  }
}
