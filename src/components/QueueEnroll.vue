<template lang="pug">
  .ui.segment
    h3.ui.header Queue Enroll
    .ui.segment(v-for="x in list")
      h4.ui.header {{ x.course.title }}
      table.ui.compact.selectable.celled.table
        thead
          tr
            th User
            th Slip
            th.four.wide Actions
        tbody
          tr(v-for="y in x.users")
            td
              user-avatar(:user="y")
            td
              a(:href="y.detail.url")
                img.ui.tiny.image(:src="y.detail.url")
            td
              .ui.green.button(@click="approve(x, y)") Approve
              .ui.red.button(@click="reject(x, y)") Reject
</template>

<script>
import { Course, User, Loader } from '../services'
import forEach from 'lodash/fp/forEach'
import UserAvatar from './UserAvatar'

export default {
  components: {
    UserAvatar
  },
  subscriptions: {
    list: Course.queueEnroll()
      .do(forEach((x) => forEach((u) => User.inject(u))(x.users)))
  },
  methods: {
    approve (x, y, btn) {
      if (!window.confirm(`Approve ${y.name} to ${x.course.title} ?`)) return
      Loader.start('approve')
      Course.addUser(x.course.id, y.id)
        .flatMap(() => User.addCourse(y.id, x.course.id))
        .flatMap(() => Course.removeQueueEnroll(x.course.id, y.id))
        .finally(() => { Loader.stop('approve') })
        .subscribe()
    },
    reject (x, y) {
      if (!window.confirm(`Reject ${y.name} from ${x.course.title} ?`)) return
      Loader.start('reject')
      Course.removeQueueEnroll(x.course.id, y.id)
        .finally(() => { Loader.stop('reject') })
        .subscribe()
    }
  }
}
</script>
