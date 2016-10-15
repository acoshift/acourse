<template>
  <div>
    <div class="ui segment">
      <user-profile :user="user" v-show="user"></user-profile>
      <div class="ui right aligned basic segment">
        <router-link class="ui green edit button" to="/profile/edit">Edit</router-link>
      </div>
    </div>
    <div class="ui segment" v-if="user && user.instructor">
      <h3 class="ui header">My Own Courses</h3>
      <router-link class="ui blue button" to="/course/new">Create new course</router-link>
      <div class="ui four stackable cards" v-if="ownCourses">
        <course-card v-for="x in ownCourses" :course="x"></course-card>
      </div>
    </div>
    <div class="ui segment" v-if="courses">
      <h3 class="ui header">My Courses</h3>
      <div class="ui four stackable cards">
        <course-card v-for="x in courses" :course="x"></course-card>
      </div>
    </div>
  </div>
</template>

<style>
  .cards {
    padding-top: 30px;
  }

  .edit.button {
    width: 140px;
  }
</style>

<script>
  import { Auth, User, Course, Loader } from '../services'
  import UserProfile from './user-profile'
  import CourseCard from './course-card'
  import { Observable } from 'rxjs'

  export default {
    components: {
      UserProfile,
      CourseCard
    },
    data () {
      return {
        user: Auth.currentUser()
          .flatMap(({ uid }) => User.getProfile(uid))
          .do(() => Loader.stop('user')),
        ownCourses: Auth.currentUser()
          .flatMap(({ uid }) => Course.ownBy(uid)),
        courses: Auth.currentUser()
          .flatMap(({ uid }) => User.courses(uid))
          .flatMap((courseIds) => Observable.combineLatest(...courseIds.map((id) => Course.get(id))))
      }
    },
    beforeCreate () {
      Loader.start('user')
    }
  }
</script>
