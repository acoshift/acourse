import Layout from './layout'

const Home = () => System.import('./home.vue')
const Auth = () => System.import('./auth.vue')
const Profile = () => System.import('./profile.vue')
const ProfileEdit = () => System.import('./profile-edit.vue')
const UserView = () => System.import('./user-view.vue')
const Course = () => System.import('./course.vue')
const CourseView = () => System.import('./course-view.vue')
const CourseChat = () => System.import('./course-chat.vue')
const CourseEditor = () => System.import('./course-editor.vue')
const CourseChatHistory = () => System.import('./course-chat-history.vue')
const CourseAttend = () => System.import('./course-attend.vue')
const CourseAssignment = () => System.import('./course-assignment.vue')
const CourseAssignmentEdit = () => System.import('./course-assignment-edit.vue')

export {
  Auth,
  Layout,
  Home,
  Profile,
  ProfileEdit,
  Course,
  CourseEditor,
  CourseView,
  UserView,
  CourseChat,
  CourseChatHistory,
  CourseAttend,
  CourseAssignment,
  CourseAssignmentEdit
}
