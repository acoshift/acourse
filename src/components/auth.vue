<template>
  <div id="auth" class="ui fluid container">
    <div class="ui center aligned grid">
      <div class="row">
        <div class="column">
          <img src="../assets/acourse.svg" class="ui centered image" style="width: 220px; height: 220px">
        </div>
      </div>
      <div class="row">
        <h1>Acourse</h1>
      </div>
      <div class="row">
        <div class="ui center aligned stacked segment">
          <div v-show="state === 0">
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
                <div class="ui right aligned basic fluid container" style="padding-bottom: 1rem;">
                  <a href="#" @click="state = 1">Forget password ?</a>
                  /
                  <a href="#" @click="state = 2">Register</a>
                </div>
                <button class="ui blue fluid submit button" :class="{loading}">Sign In</button>
              </form>
            </div>
          </div>
          <div v-show="state === 1">
            <h3>Reset Password</h3>
            <div class="ui fluid left aligned container">
              <form class="ui form" :class="{error}" @submit.prevent="forgot">
                <div class="field">
                  <label>Email</label>
                  <input v-model="email" @focus="resetError">
                </div>
                <div class="ui fluid basic segment">
                  <a href="#" @click="state = 0"><i class="left chevron icon"></i> Sign In</a>
                </div>
                <button class="ui blue fluid submit button" :class="{loading}">Reset</button>
              </form>
            </div>
          </div>
          <div v-show="state === 2">
            <h3>Register</h3>
            <div class="ui fluid left aligned container">
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
      <div class="ui horizontal divider">
        Or
      </div>
      <div class="row">
        <div class="ui center aligned segment">
          <div class="ui facebook fluid button" @click="facebookSignIn" :class="{loading: facebookLoading}"><i class="facebook f icon"></i>Sign In with Facebook</div>
          <br>
          <div class="ui google plus fluid button" :class="{loading: googleLoading}" @click="googleSignIn"><i class="google plus icon"></i>Sign In with Google+</div>
          <div class="ui error message" v-if="providerError">
            {{ providerError }}
          </div>
        </div>
      </div>
    </div>
    <div class="ui small modal" ref="successModal">
      <div class="image content">
        <div class="ui centered image">
          <i class="huge icons">
            <i class="green big thin circle icon"></i>
            <i class="green check icon"></i>
          </i>
        </div>
      </div>
      <div class="description" style="text-align: center;">
        <div class="ui header">Success</div>
        <p>Please check you email to reset password.</p>
        <div ref="closeButton" class="ui close button">OK</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  #auth {
    padding-top: 10px;
  }

  .segment {
    min-width: 320px;
  }

  @media only screen and (min-width: 500px) {
    #auth {
      padding-top: 80px;
    }

    .segment {
      min-width: 500px;
    }
  }

  .basic.segment {
    padding-top: 0;
    padding-bottom: 0;
    padding-left: 0;
    padding-right: 35px;
  }

  .modal {
    padding-bottom: 30px;
  }

  .modal .close.button {
    width: 180px;
    margin-top: 10px;
  }
</style>

<script>
  import { Auth, User } from '../services'
  import Vue from 'vue'

  export default {
    data () {
      return {
        email: '',
        password: '',
        state: 0,
        error: '',
        loading: false,
        providerError: '',
        facebookLoading: false,
        googleLoading: false
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
        Auth.resetPassword(this.email)
          .subscribe(
            () => {
              this.loading = false
              this.email = ''
              window.$(this.$refs.successModal)
                .modal('attach events', this.$refs.closeButton, 'hide')
                .modal('show')
            },
            (err) => {
              this.loading = false
              this.error = err.message
            }
          )
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
      facebookSignIn () {
        if (this.facebookLoading) return
        this.facebookLoading = true
        this.providerError = ''
        Auth.signInWithFacebook()
          .flatMap((res) => User.saveAuthProfile(res.user), (x) => x)
          .subscribe(
            () => {
              this.facebookLoading = false
              this.gotoHome()
            },
            (err) => {
              this.facebookLoading = false
              this.providerError = err.message
            }
          )
      },
      googleSignIn () {
        if (this.googleLoading) return
        this.googleLoading = true
        this.providerError = ''
        Auth.signInWithGoogle()
          .flatMap((res) => User.saveAuthProfile(res.user), (x) => x)
          .subscribe(
            () => {
              this.googleLoading = false
              this.gotoHome()
            },
            (err) => {
              this.googleLoading = false
              this.providerError = err.message
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
