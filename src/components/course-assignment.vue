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
  import { Observable } from 'rxjs'

  export default {
    data () {
      return {
        courseId: '',
        course: null,
        select: 0,
        assignments: null,
        userAssignments: null
      }
    },
    created () {
      Loader.start('course')
      this.courseId = this.$route.params.id
      Course.get(this.courseId)
        .subscribe(
          (course) => {
            Loader.stop('course')
            this.course = course
          },
          () => {
            this.$router.replace('/home')
          }
        )

      Loader.start('assignment')
      Observable.combineLatest(
        Assignment.getCode(this.courseId),
        Me.getCourseAssignments(this.courseId)
      )
        .subscribe(
          ([assignments, userAssignments]) => {
            if (Loader.has('assignment')) Loader.stop('assignment')
            this.assignments = assignments
            this.userAssignments = userAssignments
          }
        )
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
