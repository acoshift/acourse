<template>
  <div class="ui borderless top fixed menu">
    <router-link class="item" to="/home">
      <img src="../assets/acourse.svg">
    </router-link>
    <div class="right menu">
      <div class="item">
        <router-link to="/home">Home</router-link>
      </div>
      <div v-if="user" class="item" style="padding: 0 0.5rem;">
        <router-link to="/profile">
          <user-avatar :user="user"></user-avatar>
        </router-link>
      </div>
      <div class="item">
        <a href="#" @click="signOut">Sign Out</a>
      </div>
    </div>
  </div>
</template>

<script>
  import { Auth, User } from '../services'
  import Vue from 'vue'
  import UserAvatar from './user-avatar'

  export default {
    components: {
      UserAvatar
    },
    data () {
      return {
        user: Auth.currentUser()
          .flatMap(({ uid }) => User.getProfile(uid))
      }
    },
    methods: {
      signOut () {
        Auth.signOut()
          .subscribe(
            () => {
              Vue.nextTick(() => {
                this.$router.push('/')
              })
            }
          )
      }
    }
  }
</script>
