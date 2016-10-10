<template>
  <div class="ui borderless top fixed menu">
    <router-link class="item" to="/home">
      <img src="../assets/acourse.svg">
    </router-link>
    <div class="right menu">
      <div ref="dropdownUser" v-if="user" class="ui dropdown item" style="padding-top: 0.5rem; padding-bottom: 0.5rem;">
        <user-avatar :user="user"></user-avatar>
        <i class="dropdown icon"></i>
        <div class="menu">
          <router-link class="item" to="/profile">Profile</router-link>
          <a class="item" @click="signOut">Sign Out</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import { Auth, User } from '../services'
  import UserAvatar from './user-avatar'

  export default {
    components: {
      UserAvatar
    },
    data () {
      return {
        user: null,
        $user: null
      }
    },
    mounted () {
      this.$user = Auth.currentUser()
        .flatMap(({ uid }) => User.getProfile(uid))
        .subscribe(
          (user) => {
            this.user = user
            this.$nextTick(() => {
              $(this.$refs.dropdownUser).dropdown({ action: 'hide' })
            })
          }
        )
    },
    destroyed () {
      this.$user.unsubscribe()
    },
    methods: {
      signOut () {
        Auth.signOut()
          .subscribe(
            () => {
              this.$nextTick(() => {
                this.$router.push('/')
              })
            }
          )
      }
    }
  }
</script>
