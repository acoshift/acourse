<template>
  <div>
    <course-header :course="course" v-if="course"></course-header>
    <course-enroll-panel v-if="course && !course.enrolled && !course.owned && course.enroll" :course="course"></course-enroll-panel>
    <course-owner-panel v-if="course && course.owned" :course="course"></course-owner-panel>
    <!-- <course-student-panel v-if="course && course.enrolled" :course="course"></course-student-panel> -->
    <course-video v-if="(course && course.video) && (course.enrolled || course.owned)" :src="course.video"></course-video>
    <course-detail :course="course" v-if="course"></course-detail>
    <course-content :contents="course.contents" v-if="course && course.contents"></course-content>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
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
  computed: {
    ...mapGetters({
      course: 'currentCourse'
    })
  }
}
</script>
