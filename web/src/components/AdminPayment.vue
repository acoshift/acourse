<template lang="pug">
  .ui.segment
    h3.ui.header Payments
    table.ui.compact.selectable.celled.table
      thead
        tr
          th ID
          th Course
          th User
          th Image
          th Original Price
          th Price
          th Created At
          th.four.wide Actions
      tbody
        tr(v-for='x in list')
          td {{ x.id }}
          td {{ x.course.title }}
          td
            UserAvatar(:user='x.user')
          td
            a(:href='x.url')
              img.ui.tiny.image(:src='x.url')
          td {{ x.originalPrice }}
          td {{ x.price }}
          td {{ x.createdAt | date('YYYY/MM/DD, HH:mm') }}
          td
            .ui.green.button(@click='approve(x)') Approve
            .ui.red.button(@click='reject(x)') Reject
</template>

<script>
import { Payment, Loader } from 'services'
import UserAvatar from './UserAvatar'

export default {
  components: {
    UserAvatar
  },
  data () {
    return {
      refresh: 0
    }
  },
  subscriptions () {
    return {
      list: this.$watchAsObservable('refresh')
        .do(() => { Loader.start('payment') })
        .flatMap(() => Payment.list())
        .do(() => { Loader.stop('payment') })
    }
  },
  created () {
    this.refresh++
  },
  methods: {
    approve (x) {
      if (!window.confirm(`Approve ${x.user.name} to ${x.course.title} ?`)) return
      Loader.start('approve')
      Payment.approve(x.id)
        .finally(() => { Loader.stop('approve') })
        .subscribe(
          () => {
            this.refresh++
          }
        )
    },
    reject (x) {
      if (!window.confirm(`Reject ${x.user.name} from ${x.course.title} ?`)) return
      Loader.start('reject')
      Payment.reject(x.id)
        .finally(() => { Loader.stop('reject') })
        .subscribe(
          () => {
            this.refresh++
          }
        )
    }
  }
}
</script>
