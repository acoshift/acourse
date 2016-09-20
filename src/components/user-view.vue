<template>
  <div class="ui segment" :class="{loading}">
    <user-profile :user="user" v-show="!loading"></user-profile>
  </div>
</template>

<script>
  import { User } from '../services'
  import UserProfile from './user-profile'
  import _ from 'lodash'

  export default {
    components: {
      UserProfile
    },
    data () {
      return {
        user: null,
        loading: false
      }
    },
    created () {
      this.init()
    },
    watch: {
      $route () {
        this.init()
      }
    },
    methods: {
      init () {
        this.loading = true
        User.get(this.$route.params.id)
          .subscribe(
            (user) => {
              this.loading = false
              if (_.isEmpty(user)) {
                this.user = null
                return
              }
              this.user = user
            },
            () => {
              this.loading = false
            }
          )
      }
    }
  }
</script>
