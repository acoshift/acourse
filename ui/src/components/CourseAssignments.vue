<template lang="pug">
  .ui.segment(v-if="course")
    .ui.grid
      .three.column.row(v-for="x in course.assignments")
        .five.wide.column {{ x.title }}
        .column
          // div(v-if="userAssignments")
          //   div(v-for="(y, i) in userAssignments[x.id]")
          //     a(target="_bank", :href="y.url") {{ i }}
          //     br
        .two.wide.column
          .ui.green.button(v-if="x.open", @click="selectFile(x.id)") Upload
</template>

<script>
import { Course, Me, Loader, Document } from 'services'

export default {
  data () {
    return {
      courseId: this.$route.params.id,
      select: 0
    }
  },
  subscriptions () {
    Loader
    // Loader.start('assignments')
    return {
      course: this.$$route
        .flatMap((route) => Course.get(route.params.id))
      // userAssignments: Me.getCourseAssignments(this.courseId).do(() => Loader.stop('assignments'))
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
