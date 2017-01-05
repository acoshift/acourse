import Layout from './Layout'
const Home = () => System.import('./Home')
const Profile = () => System.import('./Profile')
const ProfileEdit = () => System.import('./ProfileEdit')
const CourseEditor = () => System.import('./CourseEditor')
const CourseView = () => System.import('./CourseView')
// import CourseView from './CourseView'
// import UserView from './UserView'
const Course = () => System.import('./Course')
// import CourseAttend from './CourseAttend'
// import CourseAssignment from './CourseAssignment'
// import CourseAssignmentEdit from './CourseAssignmentEdit'
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
  // CourseAssignment,
  // CourseAssignmentEdit,
  AdminPayment,
  AdminPaymentHistory,
  Privacy
}
