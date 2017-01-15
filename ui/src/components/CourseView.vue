<template lang="pug">
  div
    CourseHeader(:course='course', v-if='course')
    CourseEnrollPanel(v-if='course && !course.enrolled && !course.owned && course.options.enroll && !course.$preload', :course='course')
    CourseOwnerPanel(v-if='course && course.owned', :course='course')
    // <course-student-panel v-if="course && course.enrolled" :course="course"></course-student-panel>
    CourseVideo(v-if='course && course.video', :src='course.video')
    CourseDetail(:course='course', v-if='course')
    CourseContent(:contents='course.contents', v-if='course && course.contents')
</template>

<script>
import { Course, Auth } from 'services'
import CourseHeader from './CourseHeader'
import CourseVideo from './CourseVideo'
import CourseDetail from './CourseDetail'
import CourseContent from './CourseContent'
import CourseOwnerPanel from './CourseOwnerPanel'
import CourseEnrollPanel from './CourseEnrollPanel'
import CourseStudentPanel from './courseStudentPanel'

export default {
  components: {
    CourseHeader,
    CourseVideo,
    CourseDetail,
    CourseContent,
    CourseOwnerPanel,
    CourseEnrollPanel,
    CourseStudentPanel
  },
  subscriptions () {
    return {
      course: this.$$route
        .flatMap((route) => Course.get(route.params.id)),
      currentUser: Auth.currentUser()
        .flatMap(() => this.$$route.first())
        .do((route) => Course.fetch(route.params.id))
    }
  }
}
</script>
