<template>
  <div class="ui borderless top fixed menu">
    <router-link class="item" to="/home">
      <img src="../assets/acourse.svg">
    </router-link>
    <div class="right menu">
      <div class="item">
        <router-link to="/home">Home</router-link>
      </div>
      <div class="item">
        <router-link to="/course">My Courses</router-link>
      </div>
      <div class="item">
        <router-link to="/profile">Profile</router-link>
      </div>
      <div class="item">
        <a href="#" @click="signOut">Sign Out</a>
      </div>
    </div>
  </div>
</template>

<script>
  import { Auth } from '../services'
  import Vue from 'vue'

  export default {
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
