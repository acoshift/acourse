<template>
  <div class="ui card">
    <router-link class="image" :to="`/course/${course.id}`">
      <img :src="course.photo">
    </router-link>
    <div class="content">
      <router-link class="header" :to="`/course/${course.id}`">{{ course.title | trim(30) }}</router-link>
      <div class="meta">
        <span class="date">{{ course.start | date('DD/MM/YYYY') }}</span>
      </div>
      <div class="description">
        {{ course.description | trim(40) }}
      </div>
    </div>
    <div class="extra content">
      <div>
        <span class="right floated">
          <i class="user icon"></i>
          {{ students }}
        </span>
        <span>
          <i class="heart link icon" @click="fav" :class="{red: isFav, outline: !isFav}"></i>
          {{ favorites }}
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .card img {
    object-fit: cover;
    object-position: center center;
    height: 180px !important;
  }
</style>

<script>
  import { Auth, Course } from '../services'
  import _ from 'lodash'

  export default {
    props: ['course'],
    data () {
      return {
        isFav: false
      }
    },
    created () {
      this.init()
    },
    computed: {
      favorites () {
        return _.keys(this.course.favorite).length
      },
      students () {
        return _.keys(this.course.student).length
      }
    },
    watch: {
      course () {
        this.init()
      }
    },
    methods: {
      init () {
        Auth.currentUser()
          .first()
          .subscribe(
            (user) => {
              this.isFav = !!_.get(this.course.favorite, user.uid)
            }
          )
      },
      fav () {
        if (this.isFav) {
          Course.unfavorite(this.course.id).subscribe()
        } else {
          Course.favorite(this.course.id).subscribe()
        }
      }
    }
  }
</script>
