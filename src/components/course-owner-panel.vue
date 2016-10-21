<template>
  <div class="ui segment">
    <div class="ui stackable equal width grid">
      <div class="column">
        <router-link class="ui green fluid button" :to="`/course/${course.id}/edit`">Edit</router-link>
      </div>
      <div class="column">
        <router-link class="ui yellow fluid button" :to="`/course/${course.id}/chat`">Chat room</router-link>
      </div>
      <div v-if="course.canAttend && !course.attend" class="column">
        <div class="ui teal fluid button" @click="openAttendModal">Open Attend</div>
      </div>
      <div v-if="course.canAttend && course.attend" class="column">
        <div class="ui red fluid button" @click="closeAttend" :class="{loading: removingCode}">Close Attend</div>
      </div>
      <div class="column">
        <router-link :to="`/course/${course.id}/assignment/edit`" class="ui blue fluid button">Assignments</router-link>
      </div>
      <div class="column">
        <router-link class="ui blue fluid button" :to="`/course/${course.id}/attend`">Attendants</router-link>
      </div>
    </div>
    <div class="ui small modal" ref="attendModal">
      <div class="header">
        Set Attend Code
      </div>
      <div class="content">
        <div class="ui form">
          <div class="field">
            <label>Enter Code</label>
            <input v-model="attendCode">
          </div>
          <div v-if="attendError" class="ui red message">{{ attendError }}</div>
          <div class="ui fluid blue button" @click="submitAttend" :class="{loading: submitingAttendCode}">OK</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import { Course } from '../services'
  import moment from 'moment'

  export default {
    props: ['course'],
    data () {
      return {
        attendCode: '',
        attendError: '',
        submitingAttendCode: false,
        removingCode: false
      }
    },
    methods: {
      openAttendModal () {
        this.attendError = ''
        this.attendCode = moment().format('DDMMYYYY')
        $(this.$refs.attendModal).modal('show')
      },
      submitAttend () {
        this.attendError = ''
        this.submitingAttendCode = true
        Course.setAttendCode(this.course.id, this.attendCode)
          .finally(() => { this.submitingAttendCode = false })
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
      closeAttend () {
        this.removingCode = true
        Course.removeAttendCode(this.course.id)
          .finally(() => { this.removingCode = false })
          .subscribe()
      }
    }
  }
</script>
