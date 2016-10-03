<template>
  <div>
    <div v-if="!isApply && !isOwn" class="ui segment">
      <div class="ui blue button" style="width: 180px;" :class="{loading: applying}" @click="apply">Apply</div>
    </div>
    <div v-if="isOwn" class="ui segment">
      <router-link class="ui green button" :to="`/course/${courseId}/edit`">Edit</router-link>
      <router-link class="ui yellow button" :to="`/course/${courseId}/chat`">Chat room</router-link>
      <div class="ui teal button" @click="openAttendModal">Open Attend</div>
      <div class="ui teal button" @click="openAssignmentModal">Add Assignment</div>
      <router-link class="ui blue button" :to="`/course/${courseId}/attend`">Attendants</router-link>
    </div>
    <div v-if="isApply" class="ui segment">
      <div :class="{disabled: isAttended || !course.attend}" class="ui blue button" @click="attend">Attend</div>
      <router-link class="ui yellow button" :to="`/course/${courseId}/chat`">Chat room</router-link>
      <router-link class="ui teal button" :to="`/course/${courseId}/assignment`">Assignments</router-link>
    </div>
    <course-detail :course="course" v-if="course"></course-detail>
    <students :users="students" v-if="students"></students>
    <div class="ui small modal" ref="attendModal">
      <div class="header">
        <span v-if="isOwn">Set Attend Code</span>
        <span v-else>Attend</span>
      </div>
      <div class="content">
        <div class="ui form">
          <div class="field">
            <label>Enter Code</label>
            <input v-model="attendCode">
          </div>
          <div v-if="attendError" class="ui red message">{{ attendError }}</div>
          <div class="ui fluid blue button" @click="submitAttend" :class="{loading: attending}">OK</div>
        </div>
      </div>
    </div>
    <div class="ui small modal" ref="assignmentModal">
      <div class="header">Add Assignment</div>
      <div class="content">
        <div class="ui form">
          <div class="field">
            <label>Title</label>
            <input v-model="assignmentCode">
          </div>
          <div class="ui fluid blue button" @click="submitAssignmentCode">OK</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import { Auth, User, Course } from '../services'
  import { Observable } from 'rxjs'
  import get from 'lodash/fp/get'
  import keys from 'lodash/fp/keys'
  import CourseDetail from './course-detail'
  import Students from './students'

  export default {
    components: {
      CourseDetail,
      Students
    },
    data () {
      return {
        courseId: '',
        course: null,
        isOwn: false,
        loading: false,
        isApply: false,
        applying: false,
        students: null,
        attendCode: '',
        attending: false,
        attendError: '',
        ob: [],
        isAttended: true,
        assignmentCode: ''
      }
    },
    mounted () {
      this.$nextTick(() => {
        this.loading = true
        this.courseId = this.$route.params.id

        this.ob.push(Observable.combineLatest(
          Auth.currentUser().first(),
          Course.get(this.courseId)
            .map((course) => ({ ...course, owner: { id: course.owner } }))
            .do((course) => User.inject(course.owner))
        )
          .finally(() => { this.loading = false })
          .subscribe(
            ([user, course]) => {
              this.loading = false
              this.course = course
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
          .subscribe()
      },
      attend () {
        Course.attend(this.courseId, this.course.attend)
          .subscribe(
            () => {
              window.alert('OK')
            },
            () => {
              window.alert('Error')
            }
          )
      },
      openAttendModal () {
        this.attendError = ''
        this.attendCode = ''
        $(this.$refs.attendModal).modal('show')
      },
      submitAttend () {
        this.attendError = ''
        this.attending = true
        Course.setAttendCode(this.courseId, this.attendCode)
          .finally(() => { this.attending = false })
          .subscribe(
            () => {
              this.attendCode = ''
              $(this.$refs.attendModal).modal('hide')
            },
            (err) => {
              this.attendError = err.message
            }
          )
      },
      openAssignmentModal () {
        $(this.$refs.assignmentModal).modal('show')
      },
      submitAssignmentCode () {
        Course.addAssignment(this.courseId, { title: this.assignmentCode })
          .subscribe(
            () => {
              window.alert('ok')
              $(this.$refs.assignmentModal).modal('hide')
            },
            () => {
              window.alert('error')
            }
          )
      }
    }
  }
</script>
