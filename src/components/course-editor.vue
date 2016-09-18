<template>
  <div class="ui segment">
    <h3>
      <span v-if="isNew">New</span>
      <span v-else>Edit</span>
      Course
    </h3>
    <form class="ui form">
      <div class="field">
        <label>Title</label>
        <input v-model="course.title">
      </div>
      <div class="field">
        <label>Description</label>
        <textarea v-model="course.description"></textarea>
      </div>
      <div class="field">
        <label>Photo</label>
        <img v-show="course.photo" class="ui small image" :src="course.photo">
        <div class="ui green button">Select Photo</div>
      </div>
    </form>
  </div>
</template>

<script>
  import { User, Course } from '../services'

  export default {
    data () {
      return {
        isNew: false,
        course: {
          title: '',
          description: '',
          photo: ''
        }
      }
    },
    created () {
      if (this.$route.params.id === 'new') {
        this.isNew = true
      } else {
        // load course
        Course.get(this.$route.params.id)
          .subscribe(
            (course) => {
              console.log(course)
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
              this.user.photo = f.downloadURL
            },
            () => {
              this.uploading = false
            }
          )
      }
    }
  }
</script>
