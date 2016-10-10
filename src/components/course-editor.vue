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
          <textarea v-model="course.description" rows="15"></textarea>
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
        <div class="field">
          <label>Video</label>
          <input v-model="course.video">
        </div>
        <div class="ui divider"></div>
        <div style="padding-bottom: 1rem;">
          <h3>Contents</h3>
          <div class="ui green button" @click="addContent">Add Content</div>
          <div class="ui segment" v-for="(x, i) in contents">
            <h4 class="ui header">Content {{ i + 1 }} <i class="red remove link icon" @click="removeContent(i)"></i></h4>
            <textarea v-model="x.content" rows="5"></textarea>
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
  import { Auth, User, Course, Loader } from '../services'
  import { Observable } from 'rxjs'
  import flow from 'lodash/fp/flow'
  import defaults from 'lodash/fp/defaults'
  import pick from 'lodash/fp/pick'
  import keys from 'lodash/fp/keys'
  import map from 'lodash/fp/map'

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
          open: false,
          video: ''
        },
        contents: [],
        courseId: '',
        uploading: false,
        saving: false
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
        Loader.start('course')
        this.courseId = this.$route.params.id
        Observable.forkJoin(
          Auth.currentUser().first(),
          Course.get(this.courseId).first(),
          Course.content(this.courseId).first()
        )
          .subscribe(
            ([user, course, contents]) => {
              Loader.stop('course')
              if (course.owner !== user.uid) return this.$router.replace(`/course/${this.courseId}`)
              this.course = flow(
                pick(keys(this.course)),
                defaults(this.course)
              )(course)
              this.contents = contents && map(pick('content'))(contents) || []
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
        User.upload(file)
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
            .flatMap((courseId) => Course.saveContent(courseId, this.contents), (courseId) => courseId)
            .finally(() => { this.saving = false })
            .subscribe(
              (courseId) => {
                this.$router.push(`/course/${courseId}`)
              }
            )
        } else {
          Course.save(this.courseId, this.course)
            .flatMap(() => Course.saveContent(this.courseId, this.contents))
            .finally(() => { this.saving = false })
            .subscribe(
              () => {
                this.$router.push(`/course/${this.courseId}`)
              }
            )
        }
      },
      addContent () {
        this.contents.push({
          content: ''
        })
      },
      removeContent (position) {
        this.contents.splice(position, 1)
      }
    }
  }
</script>
