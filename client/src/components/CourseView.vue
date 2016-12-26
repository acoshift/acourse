<template>
  <div>
    <course-header :course="course" v-if="course"></course-header>
    <course-apply-panel v-if="course && !course.enrolled && !course.owned && course.enroll" :course="course" @refresh="refresh++"></course-apply-panel>
    <course-owner-panel v-if="course && course.owned" :course="course"></course-owner-panel>
    <!-- <course-student-panel v-if="course && course.enrolled" :course="course"></course-student-panel> -->
    <course-video v-if="(course && course.video) && (course.enrolled || course.owned)" :src="course.video"></course-video>
    <course-detail :course="course" v-if="course"></course-detail>
    <course-content :contents="course.contents" v-if="course && course.contents"></course-content>
  </div>
</template>

<script>
import { Course, Loader } from 'services'
import CourseHeader from './CourseHeader'
import CourseVideo from './CourseVideo'
import CourseDetail from './CourseDetail'
import CourseContent from './CourseContent'
import CourseOwnerPanel from './CourseOwnerPanel'
import CourseApplyPanel from './CourseApplyPanel'
import CourseStudentPanel from './courseStudentPanel'

export default {
  components: {
    CourseHeader,
    CourseVideo,
    CourseDetail,
    CourseContent,
    CourseOwnerPanel,
    CourseApplyPanel,
    CourseStudentPanel
  },
  data () {
    return {
      refresh: 0
    }
  },
  subscriptions () {
    return {
      course: this.$watchAsObservable('refresh')
        .do(() => { Loader.start('course') })
        .flatMap(() => Course.get(this.$route.params.id))
        .do(() => { Loader.stop('course') })
    }
  },
  created () {
    this.refresh++
  }
}
</script>
