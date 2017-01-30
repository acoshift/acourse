import Firebase from './firebase'
import RPC from './rpc'
import orderBy from 'lodash/fp/orderBy'
import response from './response'

export default {
  get (id) {
    return RPC.invoke('/acourse.UserService/GetUser', { userId: id })
  },
  ownCourses (id) {
    return RPC.invoke(`/course?owner=${id}`)
      .map(orderBy(['createdAt'], ['desc']))
  },
  courses (id) {
    return RPC.invoke('/acourse.CourseService/ListOwnCourses', { userId: id })
      .map(response.courses)
      .map(orderBy(['createdAt'], ['desc']))
  },
  upload (id, file) {
    return Firebase.upload(`user/${id}/${Date.now()}-${file.name}`, file)
  }
}
