<template>
  <div class="ui basic segment" :class="{loading}">
    <div class="ui massive breadcrumb">
      <router-link class="section" to="/home">Courses</router-link>
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
        <div class="row" style="padding-top: 0;">
          <div class="column">
            <h1>{{ course.title }}</h1>
          </div>
        </div>
        <div class="two column middle aligned row" style="margin-top: -30px !important;">
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
        <div class="row" v-if="!isApply && !isOwn">
          <div class="column">
            <div class="ui green join button" :class="{loading: applying}" @click="apply">Apply</div>
          </div>
        </div>
        <div class="row" v-if="isApply">
          <div class="ui green message">You already apply this course.</div>
        </div>
        <div v-if="isApply || isOwn">
          <router-link class="ui yellow button" :to="`/course/${courseId}/chat`">Chat room</router-link>
        </div>
        <div class="row">
          <div class="column">
            <p class="description">{{ course.description }}</p>
          </div>
        </div>
        <div v-if="isOwn" class="right aligned row">
          <div class="column">
            <router-link class="ui green edit button" :to="`/course/${courseId}/edit`">Edit</router-link>
          </div>
        </div>
      </div>
    </div>
    <div class="ui segment">
      <h3 class="ui header">Students</h3>
      <div v-for="x in students">
        <router-link :to="`/user/${x.id}`">
          <avatar :src="x.photo" size="tiny"></avatar>
          {{ x.name }}
        </router-link>
      </div>
    </div>
  </div>
</template>

<style>
  p.description {
    text-align: left;
    white-space: pre-line;
  }

  .join.button {
    width: 180px;
  }

  .edit.button {
    width: 150px;
  }
</style>

<script>
  import { Auth, User, Course } from '../services'
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
        isOwn: false,
        loading: false,
        isApply: false,
        applying: false,
        students: null
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
          Auth.currentUser.first(),
          Course.get(this.courseId)
            .flatMap((course) => User.get(course.owner), (course, owner) => ({...course, owner}))
        )
          .subscribe(
            ([user, course]) => {
              this.loading = false
              this.course = course
              if (course.owner.id === user.uid) this.isOwn = true
              this.isApply = !!_.get(course.student, user.uid)

              Observable.of(course.student)
                .map(_.keys)
                .flatMap(Observable.from)
                .flatMap((id) => User.get(id).first())
                .toArray()
                .subscribe(
                  (students) => {
                    this.students = students
                  },
                  () => {
                    this.students = null
                  }
                )
            },
            () => {
              this.loading = false
              // not found
              this.$router.replace('/home')
            }
          )
      },
      apply () {
        if (this.applying) return
        this.applying = true
        Course.join(this.courseId).subscribe(
          () => {
            this.applying = false
          },
          () => {
            this.applying = false
          }
        )
      }
    }
  }
</script>
