<template>
  <div class="ui borderless top fixed menu">
    <router-link class="item" to="/home">
      <img src="../assets/acourse.svg">
    </router-link>
    <div class="right menu">
      <div class="ui dropdown item" ref="dropdownAdmin" v-if="currentUser && currentUser.role && currentUser.role.admin">
        Admin <i class="dropdown icon"></i>
        <div class="menu">
          <router-link class="item" to="/admin/payment">Payment</router-link>
          <router-link class="item" to="/admin/payment/history">Payment History</router-link>
        </div>
      </div>
      <div ref="dropdownUser" v-if="currentUser" class="ui dropdown item" style="padding-top: 0.5rem; padding-bottom: 0.5rem;">
        <user-avatar :user="currentUser"></user-avatar>
        <i class="dropdown icon"></i>
        <div class="menu">
          <router-link class="item" to="/profile">Profile</router-link>
          <a class="item" @click="signOut">Sign Out</a>
        </div>
      </div>
      <div v-if="currentUser === null" style="padding-top: 0.5rem; padding-bottom: 0.5rem;">
        <div class="item">
          <div class="ui blue button" @click="openAuth">Sign In</div>
        </div>
      </div>
    </div>
    <AuthModal ref="auth"/>
  </div>
</template>

<script>
import { mapState } from 'vuex'
import { Auth } from 'services'
import UserAvatar from './UserAvatar'
import AuthModal from './AuthModal'

export default {
  components: {
    UserAvatar,
    AuthModal
  },
  computed: {
    ...mapState(['currentUser'])
  },
  updated () {
    $(this.$refs.dropdownUser).dropdown({ action: 'hide' })
    $(this.$refs.dropdownAdmin).dropdown({ action: 'hide' })
  },
  methods: {
    signOut () {
      Auth.signOut().subscribe(() => {
        this.$nextTick(() => {
          this.$router.push('/')
        })
      })
    },
    openAuth () {
      this.$refs.auth.open()
    }
  }
}
</script>
