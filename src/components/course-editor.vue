<template>
  <div>
    <div class="ui massive breadcrumb">
      <router-link class="section" to="/course">My Courses</router-link>
      <i class="right chevron icon divider"></i>
      <div v-if="isNew" class="active section">Create New Course</div>
      <router-link v-if="courseId" class="section" :to="`/course/${courseId}`">{{ courseId }}</router-link>
      <i v-if="!isNew" class="right chevron icon divider"></i>
      <div v-if="!isNew" class="active section">Edit</div>
    </div>
    <div class="ui segment">
      <h3>
        <span v-if="isNew">New</span>
        <span v-else>Edit</span>
        Course
      </h3>
      <form class="ui form" @submit.prevent="submit">
        <div class="field">
          <label>Cover Photo</label>
          <img v-show="course.photo" class="ui medium image" :src="course.photo">
          <div class="ui green button" @click="$refs.photo.click()" :class="{loading: uploading}">Select Photo</div>
          <input ref="photo" type="file" class="hidden" @change="uploadPhoto">
        </div>
        <div class="field">
          <label>Title</label>
          <input v-model="course.title">
        </div>
        <div class="field">
          <label>Description</label>
          <textarea v-model="course.description"></textarea>
        </div>
        <button class="ui blue button" :class="{loading}">
          <span v-if="isNew">Create</span>
          <span v-else>Save</span>
        </button>
      </form>
    </div>
  </div>
</template>

<style>
  img.image {
    margin: 10px;
  }
</style>

<script>
  import { Auth, User, Course } from '../services'
  import { Observable } from 'rxjs'

  export default {
    data () {
      return {
        isNew: false,
        course: {
          title: '',
          description: '',
          photo: '',
          owner: ''
        },
        courseId: '',
        uploading: false,
        loading: false
      }
    },
    created () {
      if (!this.$route.params.id) {
        this.isNew = true
        Auth.currentUser
          .subscribe(
            (user) => {
              this.course.owner = user.uid
            }
          )
      } else {
        this.courseId = this.$route.params.id
        Observable.forkJoin(
          Auth.currentUser.first(),
          Course.get(this.$route.params.id)
        )
          .subscribe(
            ([user, course]) => {
              this.course = course
            }
          )
      }
    },
    methods: {
      uploadPhoto () {
        if (this.uploading) return
        const file = this.$refs.photo.files[0]
        if (!file) return
        this.uploading = true
        User.uploadMePhoto(file)
          .subscribe(
            (f) => {
              this.uploading = false
              this.course.photo = f.downloadURL
            },
            () => {
              this.uploading = false
            }
          )
      },
      submit () {
        if (this.loading) return
        this.loading = true
        if (this.isNew) {
          Course.create(this.course)
            .subscribe(
              (courseId) => {
                this.loading = false
                this.$router.push(`/course/${courseId}`)
              },
              () => {
                this.loading = false
              }
            )
        } else {
          Course.save(this.courseId, this.course)
            .subscribe(
              () => {
                this.loading = false
                this.$router.push(`/course/${this.courseId}`)
              },
              () => {
                this.loading = false
              }
            )
        }
      }
    }
  }
</script>
