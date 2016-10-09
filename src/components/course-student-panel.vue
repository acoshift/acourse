<template>
  <div class="ui segment">
    <div :class="{disabled: isAttended || !course.attend, loading: attending}" class="ui blue button" @click="attend">Attend</div>
    <router-link class="ui yellow button" :to="`/course/${course.id}/chat`">Chat room</router-link>
    <router-link class="ui teal button" :to="`/course/${course.id}/assignment`">Assignments</router-link>
  </div>
</template>

<script>
  import { Course, Document } from '../services'

  export default {
    props: ['course'],
    data () {
      return {
        attending: false,
        isAttended: false,
        $isAttended: null
      }
    },
    methods: {
      created () {
        this.$isAttended = Course.isAttended(this.course.id)
          .subscribe(
            (isAttended) => {
              this.isAttended = isAttended
            }
          )
      },
      destroyed () {
        this.$isAttended.unsubscribe()
      },
      attend () {
        this.attending = true
        Course.attend(this.course.id, this.course.attend)
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
