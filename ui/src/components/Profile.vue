<template lang="pug">
  div
    .ui.segment
      UserProfile(:user='currentUser', v-show='currentUser')
      .ui.right.aligned.basic.segment
        router-link.ui.green.edit.button(to='/profile/edit') Edit
      .ui.basic.segment
        h4.ui.header User ID
        p(v-if="currentUser") {{ currentUser.id }}
      .ui.basic.segment
        h4.ui.header Link Accounts
        .ui.google.plus.button(v-if='!isLinkedGoogle', @click='linkGoogle')
          i.google.plus.icon
          |  Link Google
        .ui.red.basic.button(v-else='', @click='unlinkGoogle', :class='{disabled: !canUnlink}')
          i.google.plus.icon
          |  Unlink Google
        .ui.facebook.button(v-if='!isLinkedFacebook', @click='linkFacebook')
          i.facebook.icon
          |  Link Facebook
        .ui.blue.basic.button(v-else='', @click='unlinkFacebook', :class='{disabled: !canUnlink}')
          i.facebook.icon
          |  Unlink Facebook
        .ui.black.button(v-if='!isLinkedGithub', @click='linkGithub')
          i.github.icon
          |  Link Github
        .ui.black.basic.button(v-else='', @click='unlinkGithub', :class='{disabled: !canUnlink}')
          i.github.icon
          |  Unlink Github
    .ui.segment(v-if='currentUser && currentUser.role && currentUser.role.instructor', :class='{loading: !ownCourses}')
      h3.ui.header My Own Courses
      router-link.ui.blue.button(to='/course/new') Create new course
      .ui.four.stackable.cards(v-if='ownCourses')
        CourseCard(v-for='x in ownCourses', :course='x', :hideprice='true')
    .ui.segment(:class='{loading: !myCourses}')
      h3.ui.header My Courses
      .ui.four.stackable.cards
        CourseCard(v-for='x in myCourses', :course='x', :hideprice='true')
</template>

<style lang="scss" scoped>
  .cards {
    padding-top: 30px;
  }

  .edit.button {
    width: 140px;
  }

  .basic.segment {
    & > .ui.button {
      margin: .3rem;
    }
  }
</style>

<script>
import { Auth, Me, Firebase } from 'services'
import UserProfile from './UserProfile'
import CourseCard from './CourseCard'
import some from 'lodash/fp/some'

export default {
  components: {
    UserProfile,
    CourseCard
  },
  subscriptions () {
    return {
      providerData: Auth.requireUser()
        .map((user) => user.providerData),
      currentUser: Me.get(),
      ownCourses: Me.ownCourses(),
      myCourses: Me.courses()
    }
  },
  computed: {
    isLinkedGoogle () {
      return some((x) => x.providerId === Firebase.provider.google.providerId)(this.providerData)
    },
    isLinkedFacebook () {
      return some((x) => x.providerId === Firebase.provider.facebook.providerId)(this.providerData)
    },
    isLinkedGithub () {
      return some((x) => x.providerId === Firebase.provider.github.providerId)(this.providerData)
    },
    canUnlink () {
      return this.providerData && this.providerData.length > 1
    }
  },
  methods: {
    linkGoogle () {
      Auth.linkGoogle().subscribe()
    },
    linkFacebook () {
      Auth.linkFacebook().subscribe()
    },
    linkGithub () {
      Auth.linkGithub().subscribe()
    },
    unlinkGoogle () {
      if (!this.canUnlink) return
      Auth.unlinkGoogle().subscribe()
    },
    unlinkFacebook () {
      if (!this.canUnlink) return
      Auth.unlinkFacebook().subscribe()
    },
    unlinkGithub () {
      if (!this.canUnlink) return
      Auth.unlinkGithub().subscribe()
    }
  }
}
</script>
