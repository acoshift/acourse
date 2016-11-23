<template>
  <div class="ui segment">
    <div class="ui form">
      <textarea readonly>{{ messages }}</textarea>
    </div>
  </div>
</template>

<style scoped>
  .form {
    height: calc(100vh - 17rem);
  }

  textarea {
    height: 100%;
    min-height: 100% !important;
    max-height: 100% !important;
  }
</style>

<script>
  import { Course, Loader, User } from '../services'
  import { Observable } from 'rxjs/Observable'
  import reduce from 'lodash/fp/reduce'
  import orderBy from 'lodash/fp/orderBy'

  export default {
    data () {
      return {
        messages: ''
      }
    },
    created () {
      Loader.start('messages')
      Course.allMessages(this.$route.params.id)
        .flatMap(Observable.from)
        .flatMap((x) => User.get(x.u).first(), (x, u) => ({ ...x, user: u }))
        .toArray()
        .map(orderBy(['t'], ['asc']))
        .map(reduce((p, v) => `${p}${v.user.name || 'Anonymous'}: ${v.m}\n`, ''))
        .finally(() => { Loader.stop('messages') })
        .subscribe(
          (messages) => {
            this.messages = messages
          }
        )
    }
  }
</script>
