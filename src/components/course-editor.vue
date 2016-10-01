<template>
  <div class="ui basic segment" :class="{loading}">
    <div class="ui massive breadcrumb">
      <router-link class="section" to="/home">Courses</router-link>
      <i class="right chevron icon divider"></i>
      <div v-if="isNew" class="active section">Create New Course</div>
      <router-link v-if="courseId" class="section" :to="`/course/${courseId}`">{{ courseId }}</router-link>
      <i v-if="!isNew" class="right chevron icon divider"></i>
      <div v-if="!isNew" class="active section">Edit</div>
    </div>
    <div class="ui segment">
      <h2>
        <span v-if="isNew">New</span>
        <span v-else>Edit</span>
        Course
      </h2>
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
          <textarea v-model="course.description" rows="30"></textarea>
        </div>
        <div class="field">
          <label>Start Date</label>
          <input type="date" v-model="course.start">
        </div>
        <div class="field">
          <div class="ui toggle checkbox">
            <input type="checkbox" class="hidden" v-model="course.open">
            <label>Open public</label>
          </div>
        </div>
        <button class="ui blue save button" :class="{loading: saving}">
          <span v-if="isNew">Create</span>
          <span v-else>Save</span>
        </button>
        <router-link class="ui red cancel button" :to="`/course/${courseId}`">Cancel</router-link>
      </form>
    </div>
  </div>
</template>

<style>
  img.image {
    margin: 10px;
  }

  .save.button {
    width: 160px;
  }
</style>

<script>
  import { Auth, User, Course } from '../services'
  import { Observable } from 'rxjs'
  import _ from 'lodash'

  export default {
    data () {
      return {
        isNew: false,
        course: {
          title: '',
          description: '',
          photo: '',
          owner: '',
          start: '',
          open: false
        },
        courseId: '',
        uploading: false,
        saving: false,
        loading: false
      }
    },
    created () {
      if (!this.$route.params.id) {
        this.isNew = true
        Auth.currentUser()
          .first()
          .subscribe(
            (user) => {
              this.course.owner = user.uid
            }
          )
      } else {
        this.loading = true
        this.courseId = this.$route.params.id
        Observable.forkJoin(
          Auth.currentUser().first(),
          Course.get(this.$route.params.id).first()
        )
          .subscribe(
            ([user, course]) => {
              this.loading = false
              if (course.owner !== user.uid) return this.$router.replace(`/course/${this.courseId}`)
              this.course = _.defaults(_.pick(course, _.keys(this.course)), this.course)
            }
          )
      }
    },
    mounted () {
      window.$('.checkbox').checkbox()
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
        if (this.saving) return
        this.saving = true
        if (this.isNew) {
          Course.create(this.course)
            .subscribe(
              (courseId) => {
                this.saving = false
                this.$router.push(`/course/${courseId}`)
              },
              () => {
                this.saving = false
              }
            )
        } else {
          Course.save(this.courseId, this.course)
            .subscribe(
              () => {
                this.saving = false
                this.$router.push(`/course/${this.courseId}`)
              },
              () => {
                this.saving = false
              }
            )
        }
      }
    }
  }
</script>
