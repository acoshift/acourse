<template>
  <div>
    <course-header :course="course" v-if="course"></course-header>
    <course-apply-panel v-if="course && !isApply && !isOwn && course.open" :course="course"></course-apply-panel>
    <course-owner-panel v-if="course && isOwn" :course="course"></course-owner-panel>
    <course-student-panel v-if="course && isApply" :course="course"></course-student-panel>
    <course-video v-if="(course && course.video) && (isApply || isOwn)" :src="course.video"></course-video>
    <course-detail :course="course" v-if="course"></course-detail>
    <course-content :contents="course.contents" v-if="course && course.contents"></course-content>
  </div>
</template>

<script>
import { Me, Course, Loader } from '../services'
import get from 'lodash/fp/get'
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
      courseId: this.$route.params.id,
      course: null,
      contents: null,
      isApply: false,
      applying: false,
      attending: false,
      user: null
    }
  },
  subscriptions () {
    Loader.start('course')
    return {
      user: Me.get(),
      course: Course.get(this.courseId)
        .finally(() => { Loader.stop('course') })
    }
  },
  computed: {
    isOwn () {
      if (!this.user || !this.course) {
        return false
      }
      return this.user.id === this.course.owner.id
    },
    isApply () {
      if (!this.user || !this.course) {
        return false
      }
      return !!get(this.user.id)(this.course.student)
    }
  }
}
</script>
