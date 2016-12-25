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
      <div v-if="user === null" style="padding-top: 0.5rem; padding-bottom: 0.5rem;">
        <div class="item">
          <div class="ui blue button" @click="openAuth">Sign In</div>
        </div>
      </div>
    </div>
    <AuthModal ref="auth"/>
  </div>
</template>

<script>
import { Auth, Me } from 'services'
import { Observable } from 'rxjs/Observable'
import UserAvatar from './UserAvatar'
import AuthModal from './AuthModal'

export default {
  components: {
    UserAvatar,
    AuthModal
  },
  subscriptions () {
    return {
      user: Auth.currentUser()
        .flatMap((user) => user ? Me.get() : Observable.of(null))
    }
  },
  updated () {
    $(this.$refs.dropdownUser).dropdown({ action: 'hide' })
  },
  methods: {
    openAuth () {
      this.$refs.auth.open()
    },
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
