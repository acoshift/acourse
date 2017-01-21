<template lang="pug">
  .ui.segment(v-if="course")
    .ui.segment(v-for="(x, i) in assignments")
      h4.ui.header {{i+1}}. {{ x.title }}
      div(v-html="marked(x.description)")
      div(v-for="y in findUserAssignments(x.id)")
        a(target="_bank", :href="y.url") {{ y.id }} ({{ y.createdAt | date('YYYY/MM/DD HH:mm') }})
        br
      .ui.basic.segment
        .ui.green.button(v-if="x.open", @click="selectFile(x.id)") Upload
</template>

<script>
import { Course, Loader, Document, Assignment } from 'services'
import map from 'lodash/fp/map'
import filter from 'lodash/filter'

export default {
  data () {
    return {
      select: 0
    }
  },
  subscriptions () {
    return {
      course: this.$$route
        .flatMap((route) => Course.get(route.params.id)),
      assignments: this.$watchAsObservable('course')
        .pluck('newValue')
        .filter((x) => !!x)
        .flatMap((course) => Assignment.list(course.id)),
      userAssignments: this.$watchAsObservable('assignments')
        .pluck('newValue')
        .do(() => Loader.start('assignment'))
        .map(map((x) => x.id))
        .flatMap((ids) => Assignment.getUserAssignments(ids))
        .do(() => Loader.stop('assignment'))
        .do(console.log)
    }
  },
  methods: {
    selectFile (select) {
      this.select = select
      Document.uploadModal.open('*/*')
        .flatMap((file) => Assignment.submitUserAssignment(this.select, file.downloadURL))
        .subscribe(
          () => {
            this.assignments = { ...this.assignments }
          },
          (err) => {
            Document.openErrorModal('Upload Error', err && err.message || err)
          }
        )
    },
    findUserAssignments (assignmentId) {
      if (!this.userAssignments) return null
      return filter(this.userAssignments, (x) => x.assignmentId === assignmentId)
    }
  }
}
</script>
