<template lang="pug">
  div
    .ui.segment(style="padding-bottom: 2rem;")
      .ui.teal.button(@click="createAssignment") Add Assignment
      .ui.divider
      .ui.stackable.grid
        .two.column.row(v-for="(x, i) in assignments")
          .column
            h3 {{i + 1}}. {{ x.title }}
            div(v-html="marked(x.description)")
          .action.column
            .ui.red.button(v-if="x.open", @click="closeAssignment(x)") Close
            .ui.blue.button(v-else, @click="openAssignment(x)") Open
            .ui.green.button(@click="editAssignment(x)") Edit
    .ui.small.modal(ref="assignmentModal")
      .header
        span(v-if="assignment.id") Edit
        span(v-else) Add
        | &nbsp;Assignment
      .content
        .ui.form
          .field
            label Title
            input(v-model="assignment.title")
          .field
            label Description (markdown)
            textarea(v-model="assignment.description")
          .ui.fluid.blue.button(@click="submitAssignment")
            span(v-if="assignment.id") Edit
            span(v-else) Create
    // .ui.segment
      .ui.styled.fluid.accordion
        div(v-for="x in assignments")
          .title
            .dropdown.icon
            | {{ x.title }}
            span ({{ x.users.length }})
          .ui.content
            .ui.stackable.grid
              .two.column.row(v-for="u in x.users")
                .column
                  UserAvatar(:user="u")
                .column
                  div(v-for="(a, f) in u.files")
                    span {{ a.timestamp | date('YYYY/MM/DD HH:mm') }}
                    a(:href="a.url", target="_blank") {{ f }}
</template>

<style scoped>
  .action.column > .button {
    width: 120px;
  }
</style>

<script>
import { Loader, User, Assignment } from 'services'
import UserAvatar from './UserAvatar'

User
export default {
  components: {
    UserAvatar
  },
  data () {
    return {
      courseId: this.$route.params.id,
      assignment: {
        id: '',
        title: '',
        description: ''
      },
      refresh: 0
    }
  },
  subscriptions () {
    return {
      assignments: this.$watchAsObservable('refresh')
        .do(() => Loader.start('assignment'))
        .flatMap(() => Assignment.list(this.courseId))
        .do(() => Loader.stop('assignment'))
    }
  },
  created () {
    this.refresh++
    // Loader.start('assignment')
    // this.$subscribeTo(Assignment.get(this.courseId),
    //   (assignments) => {
    //     Loader.stop('assignment')
    //     this.assignments = flow(
    //       keys,
    //       map((id) => ({
    //         id,
    //         ...assignments.code[id],
    //         users: flow(
    //           keys,
    //           map((x) => ({
    //             id: x,
    //             files: assignments.user[x][id]
    //           })),
    //           filter((x) => !isEmpty(x.files)),
    //           forEach((x) => User.inject(x))
    //         )(assignments.user)
    //       }))
    //     )(assignments.code)
    //     setTimeout(() => {
    //       this.$nextTick(() => {
    //         $('.accordion').accordion()
    //       })
    //     }, 500)
    //   })
  },
  updated () {
    this.$nextTick(() => {
      $('.accordion').accordion()
    })
  },
  methods: {
    openAssignmentModal () {
      $(this.$refs.assignmentModal).modal('show')
    },
    submitAssignment () {
      ;(this.assignment.id
        ? Assignment.update(this.assignment)
        : Assignment.create({ ...this.assignment, courseId: this.courseId }))
        .subscribe(
          () => {
            this.assignment = {
              id: '',
              title: '',
              description: ''
            }
            $(this.$refs.assignmentModal).modal('hide')
            this.refresh++
          }
        )
    },
    openAssignment (x) {
      Assignment.open(x.id)
        .subscribe(
          () => this.refresh++
        )
    },
    closeAssignment (x) {
      Assignment.close(x.id)
        .subscribe(
          () => this.refresh++
        )
    },
    createAssignment () {
      this.assignment = {
        id: '',
        title: '',
        description: ''
      }
      this.openAssignmentModal()
    },
    editAssignment (x) {
      this.assignment = {
        id: x.id,
        title: x.title,
        description: x.description
      }
      this.openAssignmentModal()
    }
  }
}
</script>
