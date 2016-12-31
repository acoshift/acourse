<template lang="pug">
  .ui.segment
    h2
      span(v-if='isNew') New
      span(v-else='') Edit
      | &nbsp;Course
    form.ui.form(@submit.prevent='submit')
      .field
        label Cover Photo
        img.ui.medium.image(v-show='course.photo', :src='course.photo')
        .ui.green.button(@click='uploadPhoto') Select Photo
      .field
        label Title
        input(v-model='course.title', maxlength='40')
      .field
        label Short Descrption
        input(v-model='course.shortDescription', maxlength='60')
      .field
        label Description
        textarea(v-model='course.description', rows='15')
      .field
        label Start Date
        input(type='date', v-model='course.start')
      .field
        .ui.toggle.checkbox
          input.hidden(type='checkbox', v-model='course.attend')
          label Can Attend
      .field
        .ui.toggle.checkbox
          input.hidden(type='checkbox', v-model='course.assignment')
          label Has Assignment
      .field
        label Video ID
        input(v-model='course.video')
      .ui.divider
      div(style='padding-bottom: 1rem;')
        h3 Contents
        .ui.green.button(@click='addContent') Add Content
        .ui.segment(v-for='(x, i) in course.contents')
          h4.ui.header
            | Content {{ i + 1 }}
            i.red.remove.link.icon(@click='removeContent(i)')
          .ui.form
            .field
              label Title
              input(v-model='x.title')
            .field
              label Description
              textarea(v-model='x.description', rows='5')
      button.ui.blue.save.button(:class='{loading: saving}')
        span(v-if='isNew') Create
        span(v-else='') Save
      router-link.ui.red.cancel.button(:to='`/course/${courseURL}`') Cancel
</template>

<style scoped>
  img.image {
    margin: 10px;
  }

  .save.button {
    width: 160px;
  }
</style>

<script>
import { Document, Course, Loader } from 'services'
import moment from 'moment'
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
        shortDescription: '',
        description: '',
        photo: '',
        start: '',
        video: '',
        contents: [],
        attend: false,
        assignment: false
      },
      courseId: '',
      courseURL: this.$route.params.id,
      saving: false
    }
  },
  created () {
    if (!this.$route.params.id) {
      this.isNew = true
    } else {
      Loader.start('course')
      Course.get(this.courseURL).first()
        .subscribe(
          (course) => {
            Loader.stop('course')
            this.courseId = course.id
            if (!course.owned) return this.$router.replace(`/course/${this.courseURL}`)
            this.course = flow(
              pick(keys(this.course)),
              defaults(this.course)
            )(course)
            this.course.start = moment(this.course.start).format('YYYY-MM-DD')
            if (this.course.start === '0001-01-01') {
              this.course.start = ''
            }
          },
          () => {
            this.$router.replace('/')
          }
        )
    }
  },
  mounted () {
    $('.checkbox').checkbox()
  },
  methods: {
    uploadPhoto () {
      Document.uploadModal.open('image/*')
        .subscribe(
          (f) => {
            this.course.photo = f.downloadURL
          },
          (err) => {
            Document.openErrorModal('Upload Error', err.message)
          }
        )
    },
    submit () {
      if (this.saving) return
      this.saving = true
      if (this.isNew) {
        Course.create({...this.course, start: new Date(this.course.start)})
          .finally(() => { this.saving = false })
          .subscribe(
            (courseId) => {
              this.$router.push(`/course/${courseId}`)
            }
          )
      } else {
        Course.save(this.courseId, {...this.course, start: new Date(this.course.start)})
          .finally(() => { this.saving = false })
          .subscribe(
            () => {
              this.$router.push(`/course/${this.courseURL}`)
            }
          )
      }
    },
    addContent () {
      this.course.contents.push({
        content: ''
      })
    },
    removeContent (position) {
      this.course.contents.splice(position, 1)
    }
  }
}
</script>
