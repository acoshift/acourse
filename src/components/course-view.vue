<template>
  <div>
    <course-header :course="course" v-if="course"></course-header>
    <course-apply-panel v-if="course && !isApply && !isOwn" :course="course"></course-apply-panel>
    <course-owner-panel v-if="course && isOwn" :course="course"></course-owner-panel>
    <course-student-panel v-if="course && isApply" :course="course"></course-student-panel>
    <course-video v-if="course && course.video" :src="course.video"></course-video>
    <course-detail :course="course" v-if="course"></course-detail>
    <course-content :contents="contents" v-if="contents"></course-content>
    <students :users="students" v-if="students"></students>
  </div>
</template>

<script>
  import { Auth, User, Course, Loader } from '../services'
  import { Observable } from 'rxjs'
  import get from 'lodash/fp/get'
  import keys from 'lodash/fp/keys'
  import isEmpty from 'lodash/fp/isEmpty'
  import orderBy from 'lodash/fp/orderBy'
  import CourseHeader from './course-header'
  import CourseVideo from './course-video'
  import CourseDetail from './course-detail'
  import CourseContent from './course-content'
  import CourseOwnerPanel from './course-owner-panel'
  import CourseApplyPanel from './course-apply-panel'
  import CourseStudentPanel from './course-student-panel'
  import Students from './students'

  export default {
    components: {
      CourseHeader,
      CourseVideo,
      CourseDetail,
      CourseContent,
      CourseOwnerPanel,
      CourseApplyPanel,
      CourseStudentPanel,
      Students
    },
    data () {
      return {
        courseId: '',
        course: null,
        contents: null,
        isOwn: false,
        isApply: false,
        applying: false,
        students: null,
        attending: false,
        $course: null
      }
    },
    beforeCreate () {
      Loader.start('course')
    },
    mounted () {
      this.$nextTick(() => {
        this.courseId = this.$route.params.id

        this.$course = Observable.combineLatest(
          Auth.currentUser().first(),
          Course.get(this.courseId)
            .map((course) => ({ ...course, owner: { id: course.owner } }))
            .do((course) => User.inject(course.owner))
        )
          .flatMap(([user, course]) =>
            Course.content(this.courseId)
              .catch(() => Observable.of(null)),
            (p, contents) => [...p, contents]
          )
          .flatMap(([user, course, contents]) =>
            Observable.of(course.student)
              .map(keys)
              .flatMap((users) => Observable.from(users))
              .flatMap((id) => User.get(id).first())
              .toArray()
              .map(orderBy(['name'], ['asc'])),
            (p, students) => [...p, students]
          )
          .do(() => Loader.stop('course'))
          .subscribe(
            ([user, course, contents, students]) => {
              this.course = course
              this.contents = !isEmpty(contents) && contents || null
              if (course.owner.id === user.uid) this.isOwn = true
              this.isApply = !!get(user.uid)(course.student)
              this.students = students
            }
          )
      })
    },
    destroyed () {
      this.$course.unsubscribe()
    }
  }
</script>
