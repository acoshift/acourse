<template>
  <div class="ui basic segment" :class="{loading}">
    <div class="ui massive breadcrumb">
      <router-link class="section" to="/home">Courses</router-link>
      <i class="right chevron icon divider"></i>
      <router-link class="section" :to="`/course/${courseId}`">{{ course && course.title || courseId }}</router-link>
      <i class="right chevron icon divider"></i>
      <div class="active section">Assignments</div>
    </div>
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
  import { Course, User } from '../services'
  import { Observable } from 'rxjs'

  export default {
    data () {
      return {
        courseId: '',
        loading: false,
        course: null,
        select: 0,
        assignments: null,
        userAssignments: null,
        uploading: false
      }
    },
    created () {
      this.courseId = this.$route.params.id
      this.loading = true
      Course.get(this.courseId)
        .subscribe(
          (course) => {
            this.loading = false
            this.course = course
          },
          () => {
            this.loading = false
            this.$router.replace('/home')
          }
        )

      Observable.combineLatest(
        Course.getAssignments(this.courseId),
        Course.getAssignmentUser(this.courseId)
      )
        .subscribe(
          ([assignments, userAssignments]) => {
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
          .subscribe(
            () => {
              this.uploading = false
            },
            () => {
              this.uploading = false
              window.alert('Error: Please check file size should less than 2MB')
            }
          )
      }
    }
  }
</script>
