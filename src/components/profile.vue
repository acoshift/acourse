<template>
  <div class="ui segment">
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
    <div v-else>
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

  export default {
    components: {
      Avatar
    },
    data () {
      return {
        user: null
      }
    },
    created () {
      User.me()
        .subscribe(
          (user) => {
            this.user = user
          }
        )
    },
    methods: {
    }
  }
</script>
