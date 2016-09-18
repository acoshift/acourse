<template>
  <div id="auth" class="ui fluid container">
    <div class="ui center aligned grid">
      <div class="row">
        <div class="column">
          Logo
        </div>
      </div>
      <div class="row">
        <h1>Acourse</h1>
      </div>
      <div class="row">
        <div class="ui center aligned stacked segment">
          <div v-show="state === 0">
            <h3>Sign In</h3>
            <div class="ui fluid left aligned container">
              <form class="ui form" :class="{error}" @submit.prevent="signIn">
                <div class="field">
                  <label>Email</label>
                  <input v-model="email" @focus="resetError">
                </div>
                <div class="field">
                  <label>Password</label>
                  <input v-model="password" type="password" @focus="resetError">
                </div>
                <div class="ui error message">
                  {{ error }}
                </div>
                <div class="ui right aligned basic segment">
                  <a href="#" @click="state = 1">Forget password ?</a>
                  /
                  <a href="#" @click="state = 2">Register</a>
                </div>
                <button class="ui blue fluid submit button" :class="{loading}">Sign In</button>
              </form>
            </div>
          </div>
          <div v-show="state === 1">
            <h3>Forgot Password</h3>
            <div class="ui fluid left aligned container">
              <form class="ui form" :class="{error}" @submit.prevent="forgot">
                <div class="field">
                  <label>Email</label>
                  <input v-model="email" @focus="resetError">
                </div>
                <div class="ui fluid basic segment">
                  <a href="#" @click="state = 0"><i class="left chevron icon"></i> Sign In</a>
                </div>
                <button class="ui blue fluid submit button" :class="{loading}">Reset password</button>
              </form>
            </div>
          </div>
          <div v-show="state === 2">
            <h3>Register</h3>
            <div class="ui fuild left aligned container">
              <form class="ui form" :class="{error}" @submit.prevent="signUp">
                <div class="field">
                  <label>Email</label>
                  <input v-model="email" @focus="resetError">
                </div>
                <div class="field">
                  <label>Password</label>
                  <input v-model="password" type="password" @focus="resetError">
                </div>
                <div class="ui error message">
                  {{ error }}
                </div>
                <div class="ui fluid basic segment">
                  <a href="#" @click="state = 0"><i class="left chevron icon"></i> Sign In</a>
                </div>
                <button class="ui blue fluid submit button" :class="{loading}">Sign Up</button>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  #auth {
    padding-top: 100px;
  }

  .segment {
    width: 500px;
  }

  .basic.segment {
    padding-top: 0;
    padding-bottom: 0;
    padding-left: 0;
    padding-right: 35px;
  }
</style>

<script>
  import { Auth } from '../services'
  import Vue from 'vue'

  export default {
    data () {
      return {
        email: '',
        password: '',
        state: 0,
        error: '',
        loading: false
      }
    },
    methods: {
      signIn () {
        if (this.loading) return
        this.loading = true
        this.error = ''
        Auth.signIn(this.email, this.password)
          .subscribe(
            () => {
              this.loading = false
              this.gotoHome()
            },
            () => {
              this.error = 'Email or password wrong'
              this.loading = false
            }
          )
      },
      forgot () {
        if (this.loading) return
        this.loading = true
        this.error = ''
      },
      signUp () {
        if (this.loading) return
        this.loading = true
        this.error = ''
        Auth.signUp(this.email, this.password)
          .subscribe(
            (res) => {
              this.loading = false
              this.gotoHome()
            },
            (err) => {
              this.loading = false
              this.error = err.message
            }
          )
      },
      resetError () {
        this.error = ''
      },
      gotoHome () {
        Vue.nextTick(() => {
          this.$router.push('/home')
        })
      }
    },
    watch: {
      state () {
        this.resetError()
      }
    }
  }
</script>
