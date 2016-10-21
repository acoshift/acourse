<template>
  <div class="ui card">
    <router-link class="image" :to="`/course/${course.id}`">
      <img :src="course.photo">
    </router-link>
    <div class="content">
      <router-link class="header" :to="`/course/${course.id}`">{{ course.title | trim(40) }}</router-link>
      <div class="meta">
        <span class="date">{{ course.start | date('DD/MM/YYYY') }}</span>
      </div>
      <div class="description">
        {{ course.shortDescription }}
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
  import { Auth, Me } from '../services'
  import keys from 'lodash/fp/keys'
  import get from 'lodash/fp/get'

  export default {
    props: {
      course: {
        type: Object,
        required: true
      }
    },
    data () {
      return {
        uid: Auth.currentUser().first().map(({ uid }) => uid)
      }
    },
    computed: {
      favorites () {
        return keys(this.course.favorite).length
      },
      students () {
        return keys(this.course.student).length
      },
      isFav () {
        return !!get(this.uid)(this.course.favorite)
      }
    },
    methods: {
      fav () {
        if (this.isFav) {
          Me.unfavoriteCourse(this.course.id).subscribe()
        } else {
          Me.favoriteCourse(this.course.id).subscribe()
        }
      }
    }
  }
</script>
