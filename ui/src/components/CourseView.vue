<template lang="pug">
  div
    CourseHeader(v-if='course', :course='course')
    CourseEnrollPanel(v-if='course && !course.enrolled && !course.owned && course.options.enroll && !course.$preload', :course='course')
    CourseOwnerPanel(v-if='course && course.owned', :course='course')
    CourseStudentPanel(v-if='course && course.enrolled', :course='course')
    CourseVideo(v-if='course && course.video', :src='course.video')
    CourseDetail(v-if='course', :course='course')
    CourseContent(v-if='course && course.contents', :contents='course.contents')
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
        .skip(1)
        .flatMap(() => this.$$route.first())
        .do((route) => Course.fetch(route.params.id))
    }
  }
}
</script>
