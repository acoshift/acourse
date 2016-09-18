<template>
  <div class="ui segment">
    <h3 class="ui header">Profile</h3>
    <div v-if="user" class="ui grid">
      <div class="two column middle aligned row">
        <div class="three wide column">
          <avatar :src="user.photo" size="small"></avatar>
        </div>
        <div class="column">
          <h3>{{ user.name }} <span v-show="user.email">({{ user.email }})</span></h3>
        </div>
      </div>
    </div>
    <div v-else>
      <div class="ui yellow message">No Profile Data</div>
    </div>
    <div class="ui basic fluid segment">
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
