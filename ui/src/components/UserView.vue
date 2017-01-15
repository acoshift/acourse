<template lang="pug">
  div
    .ui.segment
      UserProfile(:user='user', v-show='user')
    .ui.segment(v-if='ownCourses')
      h3.ui.header Courses own by {{ user && user.name || 'Anonymous' }}
      .ui.four.stackable.cards(v-if='ownCourses')
        CourseCard(v-for='x in ownCourses', :course='x')
    .ui.segment(v-if='courses')
      h3.ui.header {{ user && user.name || 'Anonymous' }}'s Courses
      .ui.four.stackable.cards
        CourseCard(v-for='x in courses', :course='x')
</template>

<script>
import { User, Course, Loader } from 'services'
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
      user: User.get(id).do(() => { Loader.stop('user') }).catch(() => { this.$router.replace('/') }),
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
