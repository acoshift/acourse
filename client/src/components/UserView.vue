<template>
  <div>
    <div class="ui segment">
      <user-profile :user="user" v-show="user"></user-profile>
    </div>
    <div class="ui segment" v-if="ownCourses">
      <h3 class="ui header">Courses own by {{ user && user.name || 'Anonymous' }}</h3>
      <div class="ui four stackable cards" v-if="ownCourses">
        <course-card v-for="x in ownCourses" :course="x"></course-card>
      </div>
    </div>
    <div class="ui segment" v-if="courses">
      <h3 class="ui header">{{ user && user.name || 'Anonymous' }}'s Courses</h3>
      <div class="ui four stackable cards">
        <course-card v-for="x in courses" :course="x"></course-card>
      </div>
    </div>
  </div>
</template>

<script>
import { User, Course, Loader } from '../services'
import UserProfile from './UserProfile'
import CourseCard from './CourseCard'
import isEmpty from 'lodash/fp/isEmpty'
import filter from 'lodash/fp/filter'
import orderBy from 'lodash/fp/orderBy'
import { Observable } from 'rxjs/Observable'

export default {
  components: {
    UserProfile,
    CourseCard
  },
  subscriptions () {
    const id = this.$route.params.id
    Loader.start('user')
    return {
      user: User.get(id).do(() => { Loader.stop('user') }).catch(() => { this.$router.replace('/home') }),
      courses: User.courses(id)
        .flatMap((courseIds) =>
          isEmpty(courseIds)
            ? Observable.of([])
            : Observable.combineLatest(...courseIds.map((id) => Course.get(id)))
              .map(filter((course) => !!course.public))
              .map(orderBy(['timestamp'], ['desc'])))
        .map((courses) => isEmpty(courses) ? null : courses),
      ownCourses: User.ownCourses(id).map((courses) => isEmpty(courses) ? null : courses)
    }
  }
}
</script>
