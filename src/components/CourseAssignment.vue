<template>
  <div>
    <div class="ui segment" v-if="course">
      <div class="ui grid">
        <div class="three column row" v-for="x in assignments">
          <div class="five wide column">
            {{ x.title }}
          </div>
          <div class="column">
            <div v-if="userAssignments">
              <div v-for="(y, i) in userAssignments[x.id]">
                <a target="_bank" :href="y.url">{{ i }}</a><br>
              </div>
            </div>
          </div>
          <div class="two wide column">
            <div class="ui green button" v-if="x.open" @click="selectFile(x.id)">Upload</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { Course, Assignment, Me, Loader, Document } from '../services'

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
      Document.uploadModal.open('image/*')
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
