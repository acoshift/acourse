<template lang="pug">
  .ui.segment
    .ui.stackable.equal.width.grid
      .column
        router-link.ui.green.fluid.button(:to='`/course/${course.id}/edit`') Edit
      .column(v-if='!course.options.attend')
        .ui.teal.fluid.button(@click='openAttend', :class='{loading: attending}') Open Attend
      .column(v-else)
        .ui.red.fluid.button(@click='closeAttend', :class='{loading: attending}') Close Attend
      .column(v-if='course.hasAssignment')
        router-link.ui.blue.fluid.button(:to='`/course/${course.id}/assignment/edit`') Assignments
      .column(v-if='course.canAttend')
        router-link.ui.blue.fluid.button(:to='`/course/${course.id}/attend`') Attendants
</template>

<script>
import { Course } from 'services'

export default {
  props: {
    course: {
      type: Object,
      required: true
    }
  },
  data () {
    return {
      attending: false
    }
  },
  methods: {
    openAttend () {
      if (this.attending) return
      this.attending = true
      Course.openAttend(this.course.id)
        .finally(() => { this.attending = false })
        .subscribe()
    },
    closeAttend () {
      if (this.attending) return
      this.attending = true
      Course.closeAttend(this.course.id)
        .finally(() => { this.attending = false })
        .subscribe()
    }
  }
}
</script>
