<template>
  <div class="ui segment">
    <h3 class="ui header">Edit Profile</h3>
    <form class="ui form" @submit.prevent="submit">
      <div class="ui red message" v-if="error">{{ error }}</div>
      <div class="field">
        <label>Photo</label>
        <avatar v-show="user.photo" :src="user.photo" size="small"></avatar>
        <div class="ui green button" @click="selectPhoto">Select Photo</div>
      </div>
      <div class="field">
        <label>Name</label>
        <input v-model="user.name" maxlength="45">
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
  import { Loader, Me, Document } from '../services'
  import Avatar from './avatar'
  import pick from 'lodash/fp/pick'
  import keys from 'lodash/fp/keys'

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
        saving: false,
        error: ''
      }
    },
    beforeCreate () {
      Loader.start('user')
    },
    created () {
      Me.get()
        .first()
        .finally(() => { Loader.stop('user') })
        .subscribe(
          (user) => {
            this.user = pick(keys(this.user))(user)
          }
        )
    },
    methods: {
      submit () {
        this.error = ''
        if (this.saving) return
        this.saving = true
        Me.update(this.user)
          .finally(() => { this.saving = false })
          .subscribe(
            () => {
              this.$router.push('/profile')
            }
          )
      },
      selectPhoto () {
        Document.uploadModal.open('image/*')
          .subscribe(
            (f) => {
              this.user.photo = f.downloadURL
            },
            (err) => {
              Document.openErrorModal('Upload Error', err.message)
            }
          )
      }
    }
  }
</script>
