<template>
  <div>
    <div class="ui segment" v-if="course" :class="{loading: uploading}">
      <input type="file" style="display: none" ref="file" @change="selectedFile">
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
  import { Course, User, Loader } from '../services'
  import { Observable } from 'rxjs'

  export default {
    data () {
      return {
        courseId: '',
        course: null,
        select: 0,
        assignments: null,
        userAssignments: null,
        uploading: false
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
        Course.getAssignments(this.courseId),
        Course.getAssignmentUser(this.courseId)
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
        this.$refs.file.click()
      },
      selectedFile () {
        const file = this.$refs.file.files[0]
        if (!file) return
        this.uploading = true
        User.upload(file)
          .flatMap((file) => Course.addAssignmentFile(this.courseId, this.select, file.downloadURL))
          .finally(() => { this.uploading = false })
          .subscribe(
            null,
            () => {
              window.alert('Error: Please check file size should less than 5MB')
            }
          )
      }
    }
  }
</script>
