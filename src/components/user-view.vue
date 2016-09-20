<template>
  <div class="ui segment" :class="{loading}">
    <div v-if="user" class="ui center aligned grid" id="profile">
      <div class="row" style="padding-bottom: 0;">
        <avatar :src="user.photo" size="medium"></avatar>
      </div>
      <div class="row" style="padding-bottom: 0;">
        <h1>{{ user.name }}</h1>
      </div>
      <div class="row">
        <h3>{{ user.aboutMe }}</h3>
      </div>
    </div>
    <div v-if="!user && !loading">
      <div class="ui yellow message">No Profile Data</div>
    </div>
  </div>
</template>

<script>
  import { User } from '../services'
  import Avatar from './avatar'
  import _ from 'lodash'

  export default {
    components: {
      Avatar
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
