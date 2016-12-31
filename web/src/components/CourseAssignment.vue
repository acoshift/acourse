<template lang="pug">
  .ui.segment(v-if="course")
    .ui.grid
      .three.column.row(v-for="x in assignments")
        .five.wide.column {{ x.title }}
        .column
          div(v-if="userAssignments")
            div(v-for="(y, i) in userAssignments[x.id]")
              a(target="_bank", :href="y.url") {{ i }}
              br
        two.wide.column
          .ui.green.button(v-if="x.open", @click="selectFile(x.id)") Upload
</template>

<script>
import { Course, Assignment, Me, Loader, Document } from 'services'

export default {
  data () {
    return {
      courseId: this.$route.params.id,
      select: 0
    }
  },
  subscriptions () {
    Loader.start('course')
    Loader.start('assignment')
    Loader.start('userAssignments')
    return {
      course: Course.get(this.courseId).do(() => Loader.stop('course')).catch(() => { this.$router.replace('/home') }),
      assignments: Assignment.getCode(this.courseId).do(() => Loader.stop('assignment')),
      userAssignments: Me.getCourseAssignments(this.courseId).do(() => Loader.stop('userAssignments'))
    }
  },
  methods: {
    selectFile (select) {
      this.select = select
      Document.uploadModal.open('*/*')
        .flatMap((file) => Me.submitCourseAssignment(this.courseId, this.select, file.downloadURL))
        .subscribe(
          null,
          (err) => {
            Document.openErrorModal('Upload Error', err && err.message || err)
          }
        )
    }
  }
}
</script>
