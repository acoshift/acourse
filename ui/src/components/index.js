import Layout from './Layout'
import Home from './Home'
import Profile from './Profile'
import ProfileEdit from './ProfileEdit'
import CourseEditor from './CourseEditor'
import CourseView from './CourseView'
// import UserView from './UserView'
import Course from './Course'
// import CourseAttend from './CourseAttend'
import CourseAssignments from './CourseAssignments'
import CourseAssignmentEdit from './CourseAssignmentEdit'
const AdminCourse = () => System.import('./AdminCourse')
const AdminPayment = () => System.import('./AdminPayment')
const AdminPaymentHistory = () => System.import('./AdminPaymentHistory')
const Privacy = () => System.import('./Privacy')

export {
  Layout,
  Home,
  Profile,
  ProfileEdit,
  Course,
  CourseEditor,
  CourseView,
  // UserView,
  // CourseAttend,
  CourseAssignments,
  CourseAssignmentEdit,
  AdminCourse,
  AdminPayment,
  AdminPaymentHistory,
  Privacy
}
