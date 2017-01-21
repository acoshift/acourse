import RPC from './rpc'
import response from './response'

const list = (courseId) => RPC.invoke('/acourse.AssignmentService/ListAssignments', { courseId })
  .map(response.assignments)

const create = ({ courseId, title, description }) => RPC.invoke('/acourse.AssignmentService/CreateAssignment', {
  courseId,
  title,
  description
})

const open = (assignmentId) => RPC.invoke('/acourse.AssignmentService/OpenAssignment', { assignmentId })
const close = (assignmentId) => RPC.invoke('/acourse.AssignmentService/CloseAssignment', { assignmentId })

export default {
  list,
  create,
  open,
  close
}

