<template>
  <div class="ui segment">
    <h3 class="ui header">Edit Profile</h3>
    <form class="ui form" @submit.prevent="submit">
      <div class="field">
        <label>Photo</label>
        <avatar v-show="user.photo" :src="user.photo" size="small"></avatar>
        <div class="ui green button" :class="{loading: uploading}" @click="$refs.photo.click()">Select Photo</div>
        <input ref="photo" type="file" class="hidden" @change="uploadPhoto">
      </div>
      <div class="field">
        <label>Name</label>
        <input v-model="user.name">
      </div>
      <div class="field">
        <label>About me</label>
        <input v-model="user.aboutMe" maxlength="40">
      </div>
      <button class="ui blue submit button" :class="{loading: saving}">Save</button>
      <router-link to="/profile" class="ui red button">Cancel</router-link>
    </form>
  </div>
</template>

<style>
  img.circular.image {
    margin: 10px;
  }
</style>

<script>
  import { User } from '../services'
  import Avatar from './avatar'
  import _ from 'lodash'

  export default {
    components: {
      Avatar
    },
    data () {
      return {
        user: {
          photo: '',
          name: '',
          aboutMe: ''
        },
        uploading: false,
        saving: false
      }
    },
    created () {
      User.me()
        .first()
        .subscribe(
          (user) => {
            this.user = _.pick(user, _.keys(this.user))
          }
        )
    },
    methods: {
      submit () {
        if (this.saving) return
        this.saving = true
        User.updateMe(this.user)
          .subscribe(
            () => {
              this.saving = false
              this.$router.push('/profile')
            },
            () => {
              this.saving = false
            }
          )
      },
      uploadPhoto () {
        if (this.uploading) return
        const file = this.$refs.photo.files[0]
        if (!file) return
        this.uploading = true
        User.uploadMePhoto(file)
          .subscribe(
            (f) => {
              this.uploading = false
              this.user.photo = f.downloadURL
            },
            () => {
              this.uploading = false
            }
          )
      }
    }
  }
</script>
