<template>
  <div class="ui segment">
    <div class="ui stackable equal width grid">
      <div class="column">
        <div :class="{disabled: isAttended || !course.attend, loading: attending}" class="ui blue fluid button" @click="attend">Attend</div>
      </div>
      <div class="column">
        <router-link class="ui yellow fluid button" :to="`/course/${course.id}/chat`">Chat room</router-link>
      </div>
      <div class="column">
        <router-link class="ui teal fluid button" :to="`/course/${course.id}/assignment`">Assignments</router-link>
      </div>
    </div>
  </div>
</template>

<script>
  import { Me, Document } from '../services'

  export default {
    props: ['course'],
    data () {
      return {
        attending: false,
        isAttended: false,
        $isAttended: null
      }
    },
    created () {
      this.$isAttended = Me.isAttendedCourse(this.course.id)
        .subscribe(
          (isAttended) => {
            this.isAttended = isAttended
          }
        )
    },
    destroyed () {
      this.$isAttended.unsubscribe()
    },
    methods: {
      attend () {
        this.attending = true
        Me.attendCourse(this.course.id)
          .finally(() => { this.attending = false })
          .subscribe(
            () => {
              Document.openSuccessModal('Success', 'You have attended to this section.')
            },
            (err) => {
              Document.openErrorModal('Attend Error', err && err.message || err)
            }
          )
      }
    }
  }
</script>
