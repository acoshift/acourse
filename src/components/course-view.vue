<template>
  <div class="ui basic segment" :class="{loading}">
    <div class="ui massive breadcrumb">
      <router-link class="section" to="/home">Courses</router-link>
      <i class="right chevron icon divider"></i>
      <div class="active section">{{ course && course.title || courseId }}</div>
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
    <div class="ui segment" v-if="course">
      <div class="ui center aligned grid">
        <div class="row">
          <div class="column">
            <img :src="course.photo" class="ui centered big image">
          </div>
        </div>
        <div class="row" style="padding-top: 0;">
          <div class="column">
            <h1>{{ course.title }}</h1>
          </div>
        </div>
        <div class="row" style="margin-top: -2rem; margin-bottom: 1rem;">
          <div class="column">
            <i>{{ course.start | date('DD/MM/YYYY') }}</i>
          </div>
        </div>
        <div class="two column middle aligned row" style="margin-top: -30px !important;">
          <div class="right aligned column" style="padding-right: 2px;">
            <router-link :to="`/user/${course.owner.id}`">
              <avatar :src="course.owner.photo" size="mini"></avatar>
            </router-link>
          </div>
          <div class="left aligned column" style="padding-left: 2px;">
            <router-link :to="`/user/${course.owner.id}`">
              <h3>{{ course.owner.name || 'Anonymous' }}</h3>
            </router-link>
          </div>
        </div>
        <div class="row" v-if="!isApply && !isOwn">
          <div class="column">
            <div class="ui green join button" :class="{loading: applying}" @click="apply">Apply</div>
          </div>
        </div>
        <div class="row">
          <div class="column">
            <p class="description">{{ course.description }}</p>
          </div>
        </div>
      </div>
    </div>
    <div id="student" class="ui segment">
      <h3 class="ui header">Students <span v-if="students">({{ students.length }})</span></h3>
      <div class="ui stackable three column grid">
        <div class="column" v-for="x in students">
          <router-link :to="`/user/${x.id}`">
            <avatar :src="x.photo" size="tiny"></avatar>
            {{ x.name || 'Anonymous' }}
          </router-link>
        </div>
      </div>
    </div>
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

<style lang="scss" scoped>
  p.description {
    text-align: left;
    white-space: pre-wrap;
    word-break: break-word;
  }

  .join.button {
    width: 180px;
  }

  #student.segment {
    .ui.grid {
      .column {
        padding: .6rem;
      }
    }
  }
</style>

<script>
  import { Auth, User, Course } from '../services'
  import { Observable } from 'rxjs'
  import get from 'lodash/fp/get'
  import keys from 'lodash/fp/keys'
  import Avatar from './avatar'

  export default {
    components: {
      Avatar
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
            .flatMap((course) => User.getOnce(course.owner), (course, owner) => ({...course, owner}))
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
            },
            () => {
              // not found
              this.$router.replace('/home')
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
