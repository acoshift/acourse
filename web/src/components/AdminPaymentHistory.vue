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
        .flatMap(() => Payment.history())
        .do(() => { Loader.stop('payment') })
    }
  },
  created () {
    this.refresh++
  }
}
</script>
