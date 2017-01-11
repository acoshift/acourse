import RPC from './rpc'
import reverse from 'lodash/fp/reverse'
import map from 'lodash/map'
import find from 'lodash/find'

const mapPaymentsReply = (res) => map(res.payments, (payment) => ({
  ...payment,
  user: find(res.users, { id: payment.userId }),
  course: find(res.courses, { id: payment.courseId })
}))

const list = () => RPC.post('/acourse.PaymentService/ListPayments', {}, true)
  .map(mapPaymentsReply)

const history = () => RPC.post('/acourse.PaymentService/ListPayments', { history: true }, true)
  .map(mapPaymentsReply)
  .map(reverse)

const approve = (ids) => RPC.post('/acourse.PaymentService/ApprovePayments', { ids }, true)
const reject = (ids) => RPC.post('/acourse.PaymentService/RejectPayments', { ids }, true)

export default {
  list,
  history,
  approve,
  reject
}
