<template>
  <div>
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
          <input ref="photo" type="file" class="hidden" @change="uploadPhoto" accept="image/*">
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
  import flow from 'lodash/fp/flow'
  import defaults from 'lodash/fp/defaults'
  import pick from 'lodash/fp/pick'
  import keys from 'lodash/fp/keys'

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
          .finally(() => { this.loading = false })
          .subscribe(
            ([user, course]) => {
              if (course.owner !== user.uid) return this.$router.replace(`/course/${this.courseId}`)
              this.course = flow(
                pick(keys(this.course)),
                defaults(this.course)
              )(course)
            }
          )
      }
    },
    mounted () {
      $('.checkbox').checkbox()
    },
    methods: {
      uploadPhoto () {
        if (this.uploading) return
        const file = this.$refs.photo.files[0]
        if (!file) return
        this.uploading = true
        User.uploadMePhoto(file)
          .finally(() => { this.uploading = false })
          .subscribe(
            (f) => {
              this.course.photo = f.downloadURL
            }
          )
      },
      submit () {
        if (this.saving) return
        this.saving = true
        if (this.isNew) {
          Course.create(this.course)
            .finally(() => { this.saving = false })
            .subscribe(
              (courseId) => {
                this.$router.push(`/course/${courseId}`)
              }
            )
        } else {
          Course.save(this.courseId, this.course)
            .finally(() => { this.saving = false })
            .subscribe(
              () => {
                this.$router.push(`/course/${this.courseId}`)
              }
            )
        }
      }
    }
  }
</script>
