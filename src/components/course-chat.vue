<template>
  <div class="ui basic segment" :class="{loading}">
    <div class="ui massive breadcrumb">
      <router-link class="section" to="/home">Courses</router-link>
      <i class="right chevron icon divider"></i>
      <router-link class="section" :to="`/course/${courseId}`">{{ course && course.title || courseId }}</router-link>
      <i class="right chevron icon divider"></i>
      <div class="active section">Chat</div>
    </div>
    <div class="ui chat segment" ref="container">
      <div class="ui segment" id="chatBox" ref="chatBox">
        <div class="ui comments">
          <div v-for="x in messages" class="comment">
            <span class="avatar">
              <img :src="x.user.photo || '/static/icons/ic_face_black_48px.svg'">
            </span>
            <div class="content">
              <span class="author">{{ x.user.name || 'Anonymous' }}</span>
              <div class="metadata">
                <span class="date">{{ x.t | fromNow }}</span>
              </div>
              <div class="text">{{ x.m }}</div>
            </div>
          </div>
        </div>
      </div>
      <div class="ui segment">
        <div class="ui grid">
          <div class="row">
            <div class="column">
              <div class="ui fluid input">
                <input ref="input" v-model="input" @keyup.13="send"></input>
                <div class="ui basic icon button" @click="send"><i class="send icon"></i></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .chat.segment {
    padding: 0;
    overflow: hidden;
  }

  .chat.segment > .segment {
    padding: 0;
    margin: 0;
  }

  #chatBox {
    padding: 10px;
    overflow-y: scroll;
    overflow-x: hidden;
  }

  .chat.segment .button {
    margin: 0;
  }

  .avatar img {
    border-radius: 50% !important;
    height: 35px !important;
    object-fit: cover;
    object-position: top;
  }
</style>

<script>
  import { Course, User } from '../services'
  import Avatar from './avatar'
  import Vue from 'vue'
  import _ from 'lodash'

  export default {
    components: {
      Avatar
    },
    data () {
      return {
        loading: false,
        course: null,
        courseId: null,
        input: '',
        messages: [],
        ob: []
      }
    },
    created () {
      this.courseId = this.$route.params.id
      this.ob.push(Course.get(this.courseId)
        .subscribe(
          (course) => {
            this.loading = false
            this.course = course
          },
          () => {
            this.loading = false
            this.$router.replace('/home')
          }
        )
      )
      this.ob.push(Course.messages(this.courseId)
        .concatMap((message) => User.get(message.u).first(), (message, user) => ({...message, user}))
        .subscribe(
          (message) => {
            this.messages.push(message)
            Vue.nextTick(() => {
              window.$(this.$refs.chatBox).scrollTop(99999)
            })
          }
        )
      )
    },
    destroyed () {
      _.forEach(this.ob, (x) => x.unsubscribe())
    },
    mounted () {
      this.adjust()
      window.$(window).resize(this.adjust)
    },
    methods: {
      adjust () {
        let container = window.$(this.$refs.container)
        let box = window.$(this.$refs.chatBox)
        let h = window.innerHeight
        let input = window.$(this.$refs.input)
        container.height(() => h - container.offset().top - 20)
        box.height(() => h - box.offset().top - 80)
        box.scrollTop(99999)
        input.focus()
      },
      send () {
        const input = this.input
        if (!input || !input.trim()) return
        this.input = ''
        Course.sendMessage(this.courseId, input).subscribe()
      }
    }
  }
</script>
