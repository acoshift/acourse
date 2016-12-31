<template lang="pug">
.ui.borderless.top.fixed.menu
  router-link.item(to='/home')
    img(src='../assets/acourse.svg')
  .right.menu
    .ui.dropdown.item(ref='dropdownAdmin', v-if='currentUser && currentUser.role && currentUser.role.admin')
      | Admin
      i.dropdown.icon
      .menu
        router-link.item(to='/admin/payment') Payment
        router-link.item(to='/admin/payment/history') Payment History
    .ui.dropdown.item(ref='dropdownUser', v-if='currentUser', style='padding-top: 0.5rem; padding-bottom: 0.5rem;')
      UserAvatar(:user='currentUser')
      i.dropdown.icon
      .menu
        router-link.item(to='/profile') Profile
        a.item(@click='signOut') Sign Out
    div(v-if='currentUser === null', style='padding-top: 0.5rem; padding-bottom: 0.5rem;')
      .item
        .ui.blue.button(@click='openAuth') Sign In
  AuthModal(ref='auth')
</template>

<script>
import { Auth, Me } from 'services'
import UserAvatar from './UserAvatar'
import AuthModal from './AuthModal'

export default {
  components: {
    UserAvatar,
    AuthModal
  },
  subscriptions () {
    return {
      currentUser: Me.get()
    }
  },
  mounted () {
    this.bindDOM()
  },
  updated () {
    this.bindDOM()
  },
  methods: {
    bindDOM () {
      $(this.$refs.dropdownUser).dropdown({ action: 'hide' })
      $(this.$refs.dropdownAdmin).dropdown({ action: 'hide' })
    },
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
