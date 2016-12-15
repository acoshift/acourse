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
import { Course, User } from '../services'
import forEach from 'lodash/fp/forEach'
import UserAvatar from './UserAvatar'

export default {
  components: {
    UserAvatar
  },
  subscriptions: {
    list: Course.queueEnroll()
      .do(forEach((x) => forEach((u) => User.inject(u))(x.users)))
      .do(console.log)
  },
  methods: {
    approve (x) {

    },
    reject (x) {

    }
  }
}
</script>
