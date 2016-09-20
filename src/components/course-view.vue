<template>
  <div class="ui basic segment" :class="{loading}">
    <div class="ui massive breadcrumb">
      <router-link v-if="!isMyCourse" class="section" to="/home">Courses</router-link>
      <router-link v-else class="section" to="/course">My Courses</router-link>
      <i class="right chevron icon divider"></i>
      <div class="active section">{{ course && course.title || courseId }}</div>
    </div>
    <div class="ui segment" v-if="course">
      <div class="ui center aligned grid">
        <div class="row">
          <div class="column">
            <img :src="course.photo" class="ui centered big image">
          </div>
        </div>
        <div class="row">
          <div class="column">
            <h1>{{ course.title }}</h1>
          </div>
        </div>
        <div class="two column middle aligned row">
          <div class="right aligned column" style="padding-right: 2px;">
            <router-link :to="`/user/${course.owner.id}`">
              <avatar :src="course.owner.photo" size="mini"></avatar>
            </router-link>
          </div>
          <div class="left aligned column" style="padding-left: 2px;">
            <router-link :to="`/user/${course.owner.id}`">
              <h3>{{ course.owner.name }}</h3>
            </router-link>
          </div>
        </div>
        <div class="row">
          <div class="column">
            <p class="description">{{ course.description }}</p>
          </div>
        </div>
        <div v-if="isMyCourse" class="right aligned row">
          <div class="column">
            <router-link class="ui green button" :to="`/course/${courseId}/edit`">Edit</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
  p.description {
    text-align: left;
    white-space: pre-line;
  }
</style>

<script>
  import { User, Course } from '../services'
  import { Observable } from 'rxjs'
  import _ from 'lodash'
  import Avatar from './avatar'

  export default {
    components: {
      Avatar
    },
    data () {
      return {
        courseId: '',
        course: null,
        isMyCourse: false,
        loading: false
      }
    },
    created () {
      this.init()
    },
    watch: {
      $route () {
        this.init()
      }
    },
    methods: {
      init () {
        this.loading = true
        this.courseId = this.$route.params.id

        Observable.combineLatest(
          User.me().first(),
          Course.get(this.courseId)
            .flatMap((course) => User.get(course.owner), (course, owner) => ({...course, owner: {id: course.owner, ...owner}}))
        )
          .subscribe(
            ([user, course]) => {
              this.loading = false
              this.course = course
              if (_.get(user.course, this.courseId)) this.isMyCourse = true
            },
            () => {
              this.loading = false
              // not found
              this.$router.replace('/home')
            }
          )
      }
    }
  }
</script>
