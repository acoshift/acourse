<template lang="pug">
  .ui.small.modal
    .header Enroll
    .content
      div
        h4 Upload
        p.description
          | วิธีการลงทะเบียน {{ course.title }}
          br
          br
          | 1. โอนเงินจำนวน #[b {{ calcPrice }}] บาท ไปที่
          br
          br
          | กฤษฎา เฉลิมสุข (Krissada Chalermsook)
          | 470-2-46894-4
          | ธนาคาร กสิกรไทย
          | สาขา ถนนแจ้งวัฒนะ
          br
          br
          | 2. ถ่ายรูป หรือ Capture screen หลักฐานการโอนเงินไว้
          | 3. Upload มาที่ระบบ
          | 4. รอการ approve จากเราภายในไม่เกิน 1 วันทำการ แล้วจากนั้นจะสามารถใช้งานระบบได้เลย
          br
          b *** สามารถติดต่อ admin ได้ที่ line : hideoaki กรณีที่ท่านไม่ได้รับการยืนยันภายใน 1 วันทำการ
          br
          b *** สมัครจากที่อื่นสามารถ Upload หลักฐานมาได้เหมือนกัน
        .ui.form
          .field
            label จำนวนเงิน (บาท)
            input(type="number", v-model.number="price")
        br
        .ui.fluid.green.button(@click="upload") Upload
        br
        img.ui.tiny.image(v-if="url", :url="url")
        br
        .ui.fluid.blue.button(@click="enroll", :class="{'loading disabled': loading}") Enroll
</template>

<script>
import { mapActions } from 'vuex'
import { Document, Course } from 'services'
import { Observable } from 'rxjs/Observable'

export default {
  props: {
    course: {
      type: Object,
      required: true
    }
  },
  data () {
    return {
      price: 0,
      url: '',
      code: '',
      loading: false
    }
  },
  computed: {
    calcPrice () {
      return this.course.discount ? this.course.discountedPrice : this.course.price
    }
  },
  methods: {
    ...mapActions(['fetchCurrentCourse']),
    show () {
      this.price = this.calcPrice
      $(this.$el).modal('show')
    },
    upload () {
      Document.uploadModal.open('image/*')
        .subscribe(
          (file) => {
            this.url = file.downloadURL
            Document.openSuccessModal('Success', 'Your enroll request success!.')
          },
          (err) => {
            Document.openErrorModal('Upload Error', err && err.message || err)
          }
        )
    },
    enroll () {
      Observable.of({})
        .do(() => { this.loading = true })
        .finally(() => { this.loading = false })
        .flatMap(() => Course.enroll(this.course.id, { url: this.url, price: this.price, code: this.code }))
        .subscribe(
          () => {
            Document.openSuccessModal('Success', 'Your enroll request success!.')
            this.fetchCurrentCourse()
          },
          (err) => {
            Document.openErrorModal('Upload Error', err && err.message || err)
          }
        )
    }
  }
}
</script>
