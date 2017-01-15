import RPC from './rpc'
import reverse from 'lodash/fp/reverse'
import response from './response'

const list = () => RPC.invoke('/acourse.PaymentService/ListWaitingPayments', null, true)
  .map(response.payments)

const history = () => RPC.invoke('/acourse.PaymentService/ListHistoryPayments', null, true)
  .map(response.payments)
  .map(reverse)

const approve = (paymentIds) => RPC.invoke('/acourse.PaymentService/ApprovePayments', { paymentIds }, true)
const reject = (paymentIds) => RPC.invoke('/acourse.PaymentService/RejectPayments', { paymentIds }, true)
const updatePrice = (paymentId, price) => RPC.invoke('/acourse.PaymentService/UpdatePrice', { paymentId, price }, true)

export default {
  list,
  history,
  approve,
  reject,
  updatePrice
}
