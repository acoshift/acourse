<template>
  <div>
    <div class="ui segment">
      <user-profile :user="currentUser" v-show="currentUser"></user-profile>
      <div class="ui right aligned basic segment">
        <router-link class="ui green edit button" to="/profile/edit">Edit</router-link>
      </div>
      <div class="ui basic segment">
        <h4 class="ui header">Link Accounts</h4>
        <div class="ui google plus button" v-if="!isLinkedGoogle" @click="linkGoogle"><i class="google plus icon"></i> Link Google</div>
        <div class="ui red basic button" v-else @click="unlinkGoogle" :class="{disabled: !canUnlink}"><i class="google plus icon"></i> Unlink Google</div>
        <div class="ui facebook button" v-if="!isLinkedFacebook" @click="linkFacebook"><i class="facebook icon"></i> Link Facebook</div>
        <div class="ui blue basic button" v-else @click="unlinkFacebook" :class="{disabled: !canUnlink}"><i class="facebook icon"></i> Unlink Facebook</div>
        <div class="ui black button" v-if="!isLinkedGithub" @click="linkGithub"><i class="github icon"></i> Link Github</div>
        <div class="ui black basic button" v-else @click="unlinkGithub" :class="{disabled: !canUnlink}"><i class="github icon"></i> Unlink Github</div>
      </div>
    </div>
    <div class="ui segment" v-if="currentUser && currentUser.role && currentUser.role.instructor" :class="{loading: !ownCourses}">
      <h3 class="ui header">My Own Courses</h3>
      <router-link class="ui blue button" to="/course/new">Create new course</router-link>
      <div class="ui four stackable cards" v-if="ownCourses">
        <course-card v-for="x in ownCourses" :course="x" :hidePrice="true"></course-card>
      </div>
    </div>
    <div class="ui segment" :class="{loading: !myCourses}">
      <h3 class="ui header">My Courses</h3>
      <div class="ui four stackable cards">
        <course-card v-for="x in myCourses" :course="x" :hidePrice="true"></course-card>
      </div>
    </div>
  </div>
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
