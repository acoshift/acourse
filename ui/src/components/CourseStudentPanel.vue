<template lang="pug">
  .ui.segment(v-if='course.options.attend || course.options.assignment')
    .ui.stackable.equal.width.grid
      .column(v-if='course.options.attend')
        .ui.blue.fluid.button(:class='{disabled: !course.options.attend || course.attended, loading: attending}', @click='attend') Attend
      .column(v-if='course.options.assignment')
        router-link.ui.teal.fluid.button(:to='`/course/${course.id}/assignments`') Assignments
</template>

<script>
import { Course, Document } from 'services'

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
    attend () {
      if (this.attending) return
      this.attending = true
      Course.attend(this.course.id)
        .flatMap(() => Course.fetch(this.course.id))
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
