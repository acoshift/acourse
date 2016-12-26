import API from './api'

export default {
  list () {
    return API.get('/course')
  },
  get (url) {
    return API.get(`/course/${url}`)
  },
  create (data) {
    return API.post('/course', data, true)
      .map(({ id }) => id)
  },
  save (id, data) {
    return API.patch(`/course/${id}`, data)
  },
  enroll (id, { code, url }) {
    return API.put(`/course/${id}/enroll`, { code, url }, true)
  }
}
