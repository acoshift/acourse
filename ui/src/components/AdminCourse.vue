<template lang="pug">
  .ui.segment
    h3.ui.header Course
    table.ui.compact.selectable.celled.table
      thead
        tr
          th ID
          th Title
          th Owner
          th Image
          th Type
          th Price
          th Discounted Price
          th Created At
          th Updated At
          th.four.wide Actions
      tbody
        tr(v-for='x in list')
          td {{ x.id }}
          td {{ x.title }}
          td
            UserAvatar(:user='x.owner')
          td
            a(:href='x.photo')
              img.ui.tiny.image(:src='x.photo')
          td {{ x.type }}
          td {{ x.price }}
          td {{ x.discountedPrice }}
          td {{ x.createdAt | date('YYYY/MM/DD, HH:mm') }}
          td {{ x.updatedAt | date('YYYY/MM/DD, HH:mm') }}
          td
            // .ui.green.button(@click='approve(x)') Approve
            // .ui.red.button(@click='reject(x)') Reject
</template>

<script>
import { Course, Loader } from 'services'
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
        .do(() => { Loader.start('course') })
        .flatMap(() => Course.listAll())
        .do(() => { Loader.stop('course') })
    }
  },
  created () {
    this.refresh++
  },
  methods: {
    // approve (x) {
    //   if (!window.confirm(`Approve ${x.user.name} to ${x.course.title} ?`)) return
    //   Loader.start('approve')
    //   Payment.approve(x.id)
    //     .finally(() => { Loader.stop('approve') })
    //     .subscribe(
    //       () => {
    //         this.refresh++
    //       }
    //     )
    // },
    // reject (x) {
    //   if (!window.confirm(`Reject ${x.user.name} from ${x.course.title} ?`)) return
    //   Loader.start('reject')
    //   Payment.reject(x.id)
    //     .finally(() => { Loader.stop('reject') })
    //     .subscribe(
    //       () => {
    //         this.refresh++
    //       }
    //     )
    // }
  }
}
</script>
