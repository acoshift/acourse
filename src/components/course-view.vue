<template>
  <div>
    <course-header :course="course" v-if="course"></course-header>
    <div v-if="!isApply && !isOwn" class="ui segment">
      <div class="ui blue button" style="width: 180px;" :class="{loading: applying}" @click="apply">Apply</div>
    </div>
    <course-owner-panel v-if="course && isOwn" :course="course"></course-owner-panel>
    <div v-if="isApply" class="ui segment">
      <div :class="{disabled: isAttended || !course.attend, loading: attending}" class="ui blue button" @click="attend">Attend</div>
      <router-link class="ui yellow button" :to="`/course/${courseId}/chat`">Chat room</router-link>
      <router-link class="ui teal button" :to="`/course/${courseId}/assignment`">Assignments</router-link>
    </div>
    <course-detail :course="course" v-if="course"></course-detail>
    <course-content :contents="contents" v-if="contents"></course-content>
    <students :users="students" v-if="students"></students>
    <success-modal ref="successModal"></success-modal>
  </div>
</template>

<script>
  import { Auth, User, Course, Loader } from '../services'
  import { Observable } from 'rxjs'
  import get from 'lodash/fp/get'
  import keys from 'lodash/fp/keys'
  import CourseHeader from './course-header'
  import CourseDetail from './course-detail'
  import CourseContent from './course-content'
  import CourseOwnerPanel from './course-owner-panel'
  import Students from './students'
  import SuccessModal from './success-modal'

  export default {
    components: {
      CourseHeader,
      CourseDetail,
      CourseContent,
      CourseOwnerPanel,
      Students,
      SuccessModal
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
        ob: [],
        isAttended: true
      }
    },
    beforeCreate () {
      Loader.start('course')
    },
    mounted () {
      this.$nextTick(() => {
        this.courseId = this.$route.params.id

        this.ob.push(Observable.combineLatest(
          Auth.currentUser().first(),
          Course.get(this.courseId)
            .map((course) => ({ ...course, owner: { id: course.owner } }))
            .do((course) => User.inject(course.owner))
        )
          .flatMap(([user, course]) =>
            Course.content(this.courseId)
              .catch(() => Observable.of(null)),
            ([user, course], contents) => [user, course, contents]
          )
          .do(() => Loader.stop('course'))
          .subscribe(
            ([user, course, contents]) => {
              this.course = course
              this.contents = contents
              if (course.owner.id === user.uid) this.isOwn = true
              this.isApply = !!get(user.uid)(course.student)

              this.ob.push(Observable.of(course.student)
                .map(keys)
                .flatMap((users) => Observable.from(users))
                .flatMap((id) => User.get(id).first())
                .toArray()
                .subscribe(
                  (students) => {
                    this.students = students
                  },
                  () => {
                    this.students = null
                  }
                )
              )

              Course.isAttended(this.courseId)
                .subscribe(
                  (isAttended) => {
                    this.isAttended = isAttended
                  }
                )
            }
          )
        )
      })
    },
    destroyed () {
      this.ob.forEach((x) => x.unsubscribe())
    },
    methods: {
      apply () {
        if (this.applying) return
        this.applying = true
        Course.join(this.courseId)
          .finally(() => { this.applying = false })
          .subscribe(
            () => {
              this.$refs.successModal.show('Success', 'You have applied to this course.')
            }
          )
      },
      attend () {
        this.attending = true
        Course.attend(this.courseId, this.course.attend)
          .finally(() => { this.attending = false })
          .subscribe(
            () => {
              this.$refs.successModal.show()
              this.$refs.successModal.show('Success', 'You have attended to this section.')
            },
            () => {
              window.alert('Error')
            }
          )
      }
    }
  }
</script>
