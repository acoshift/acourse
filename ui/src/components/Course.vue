<template lang="pug">
  .ui.basic.segment(:class='{loading: !course}')
    .ui.huge.breadcrumb(style='padding-bottom: 1.5rem;')
      router-link.section(to='/home') Courses
      i.right.chevron.icon.divider
      router-link.section(:to='`/course/${courseId}`', :tag="$route.name === 'courseView' && 'div' || 'a'", active-class='active', exact='') {{ course && course.title || courseId }}
      i.right.chevron.icon.divider(v-show="$route.name !== 'courseView'")
      .active.section(v-show="$route.name === 'courseEdit'") Edit
      .active.section(v-show="$route.name === 'courseNew'") New
      .active.section(v-show="$route.name === 'courseAssignment'") Assignments
      .active.section(v-show="$route.name === 'courseAttend'") Attendants
      .active.section(v-show="$route.name === 'courseAssignmentEdit'") Edit Assignment
    router-view
</template>

<style scoped>
  @media only screen and (max-width: 500px) {
    .breadcrumb {
      font-size: 1.05rem !important;
    }
  }
</style>

<script>
import { Course, Document } from 'services'

export default {
  data () {
    return {
      courseId: this.$route.params.id
    }
  },
  subscriptions () {
    return {
      course: Course.get(this.courseId)
        .do((course) => Document.setCourse(course))
    }
  },
  created () {
    Course.fetch(this.courseId)
      .subscribe(null, () => {
        this.$router.replace('/')
      })
  }
}
</script>
