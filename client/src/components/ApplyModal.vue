<template lang="pug">
  .ui.small.modal
    .header Enroll
    .content
      div
        h4 Upload
        p.description
          | วิธีการลงทะเบียน Course English for Programmer
          br
          br
          | 1. โอนเงินจำนวน #[b {{ price }}] บาท ไปที่
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
        .ui.fluid.green.button(@click="upload") Upload
</template>

<script>
import { Document, Course } from '../services'

export default {
  props: {
    course: {
      type: Object,
      required: true
    }
  },
  computed: {
    price () {
      if (this.course.discount) {
        return this.course.discountedPrice
      }
      return this.course.price
    }
  },
  methods: {
    show () {
      $(this.$el).modal('show')
    },
    upload () {
      Document.uploadModal.open('image/*')
        .flatMap((file) => Course.enroll(this.course.id, { url: file.downloadURL }))
        .subscribe(
          () => {
            Document.openSuccessModal('Success', 'Your enroll request success!.')
            this.$emit('refresh')
          },
          (err) => {
            Document.openErrorModal('Upload Error', err && err.message || err)
          }
        )
    }
  }
}
</script>
