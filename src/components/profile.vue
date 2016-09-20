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
    <div class="ui right aligned basic segment">
      <router-link class="ui green button" to="/profile/edit">Edit</router-link>
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
      this.loading = true
      User.me()
        .subscribe(
          (user) => {
            this.loading = false
            this.user = !_.isEmpty(user) ? user : null
          },
          () => {
            this.loading = false
          }
        )
    },
    methods: {
    }
  }
</script>
