import map from 'lodash/map'
import find from 'lodash/find'

const course = (res) => ({
  ...res.course,
  owner: res.user,
  owned: !!res.owned,
  purchase: !!res.purchase,
  enrolled: !!res.enrolled,
  attended: !!res.attended
})

const courses = (res) => map(res.courses, (course) => ({
  ...course,
  owner: find(res.users, { id: course.owner }),
  student: (() => {
    const r = find(res.enrollCounts, { courseId: course.id })
    return r && r.count || 0
  })()
}))

const payments = (res) => map(res.payments, (payment) => ({
  ...payment,
  user: find(res.users, { id: payment.userId }),
  course: find(res.courses, { id: payment.courseId })
}))

const assignments = (res) => res.assignments

export default {
  course,
  courses,
  payments,
  assignments
}
