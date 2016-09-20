<template>
  <div>
    <div class="ui segment" :class="{loading}">
      <user-profile :user="user" v-show="!loading"></user-profile>
    </div>
    <div class="ui segment" v-if="user && ownCourses">
      <h3 class="ui header">Courses own by {{ user.name }}</h3>
      <div class="four stackable cards" v-if="ownCourses">
        <course-card v-for="x in ownCourses" :course="x"></course-card>
      </div>
    </div>
    <div class="ui segment" v-if="user && courses">
      <h3 class="ui header">My Courses</h3>
      <div class="four stackable cards">
        <course-card v-for="x in courses" :course="x"></course-card>
      </div>
    </div>
  </div>
</template>

<script>
  import { User } from '../services'
  import UserProfile from './user-profile'
  import _ from 'lodash'

  export default {
    components: {
      UserProfile
    },
    data () {
      return {
        user: null,
        loading: false,
        ownCourses: null,
        courses: null
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
        User.get(this.$route.params.id)
          .subscribe(
            (user) => {
              this.loading = false
              if (_.isEmpty(user)) {
                this.user = null
                return
              }
              this.user = user
            },
            () => {
              this.loading = false
            }
          )
      }
    }
  }
</script>
