<template>
  <div>
    <course-header :course="course" v-if="course"></course-header>
    <course-apply-panel v-if="course && !isApply && !isOwn && course.open" :course="course"></course-apply-panel>
    <course-owner-panel v-if="course && isOwn" :course="course"></course-owner-panel>
    <course-student-panel v-if="course && isApply" :course="course"></course-student-panel>
    <course-video v-if="(course && course.video) && (isApply || isOwn)" :src="course.video"></course-video>
    <course-detail :course="course" v-if="course"></course-detail>
    <course-content :contents="contents" v-if="contents"></course-content>
  </div>
</template>

<script>
import { Auth, User, Course, Loader } from '../services'
import { Observable } from 'rxjs/Observable'
import get from 'lodash/fp/get'
import isEmpty from 'lodash/fp/isEmpty'
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
      isOwn: false,
      isApply: false,
      applying: false,
      students: null,
      attending: false
    }
  },
  created () {
    Loader.start('course')
    this.$subscribeTo(Observable.combineLatest(
      Auth.currentUser().first(),
      Course.get(this.courseId)
        .map((course) => ({ ...course, owner: { id: course.owner } }))
        .do((course) => User.inject(course.owner)))
      .flatMap(([user, course]) =>
        Course.content(this.courseId)
          .catch(() => Observable.of(null)),
        (p, contents) => [...p, contents]
      )
      .do(() => Loader.stop('course')),
      ([user, course, contents]) => {
        this.course = course
        this.contents = !isEmpty(contents) && contents || null
        if (course.owner.id === user.uid) this.isOwn = true
        this.isApply = !!get(user.uid)(course.student)
      }
    )
  }
}
</script>
