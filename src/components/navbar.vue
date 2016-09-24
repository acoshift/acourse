<template>
  <div class="ui borderless top fixed menu">
    <router-link class="item" to="/home">
      <img src="../assets/acourse.svg">
    </router-link>
    <div class="right menu">
      <div class="item">
        <router-link to="/home">Home</router-link>
      </div>
      <div class="item" style="padding: 0 0.5rem;">
        <router-link to="/profile">
          <avatar :src="user && user.photo" size="tiny"></avatar>
          {{ user && user.name || 'Anonymous' }}
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
