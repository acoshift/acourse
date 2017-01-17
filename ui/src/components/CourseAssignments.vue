<template lang="pug">
  .ui.segment(v-if="course")
    .ui.segment(v-for="(x, i) in course.assignments")
      h4.ui.header {{ x.title }}
      div(v-html="marked(x.description)")
      // div(v-if="userAssignments")
      //   div(v-for="(y, i) in userAssignments[x.id]")
      //     a(target="_bank", :href="y.url") {{ i }}
      //     br
      .ui.basic.segment
        .ui.green.button(v-if="x.open", @click="selectFile(i)") Upload
</template>

<script>
import { Course, Me, Loader, Document } from 'services'
import marked from 'marked'

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
    },
    marked (data) {
      if (!data) return ''
      return marked(data)
    }
  }
}
</script>
