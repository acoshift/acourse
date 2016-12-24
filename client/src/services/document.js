import { Observable } from 'rxjs/Observable'
import { Subject } from 'rxjs/Subject'

const getOgElement = (property) => $(`meta[property="og\\:${property}"]`)
const getContent = ($meta) => $meta.attr('content')
const setContent = ($meta, value) => $meta.attr('content', value)

let $description = $('meta[name="description"]')
let $ogTitle = getOgElement('title')
let $ogUrl = getOgElement('url')
let $ogDescription = getOgElement('description')
let $ogImage = getOgElement('image')

let _title = 'Acourse'
let _description = getContent($description)
let _ogTitle = getContent($ogTitle)
let _ogUrl = getContent($ogUrl)
let _ogDescription = getContent($ogDescription)
let _ogImage = getContent($ogImage)

const title = (value) => { document.title = value ? `${value} | ${_title}` : `${_title}` }
const description = (value) => { setContent($description, value || _description) }
const ogDescription = (value) => { setContent($ogDescription, value || _ogDescription) }
const ogImage = (value) => { setContent($ogImage, value || _ogImage) }
const ogTitle = (value) => { setContent($ogTitle, value || _ogTitle) }
const ogUrl = (value) => { setContent($ogUrl, value || _ogUrl) }

export default {
  $successModal: new Subject(),
  $errorModal: new Subject(),
  uploadModal: null,
  visibilityChanged () {
    return Observable.fromEvent(document, 'visibilitychange')
      .map(() => document.hidden)
  },
  reset () {
    title()
    description()
    ogTitle()
    ogDescription()
    ogImage()
    ogUrl()
  },
  setCourse (course) {
    if (!course) return
    title(course.title)
    ogTitle(course.title)
    description(course.shortDescription)
    ogDescription(course.shortDescription)
    ogImage(course.photo)
    ogUrl(`https://acourse.io/course/${course.url || course.id}`)
  },
  openSuccessModal (title, description) {
    this.$successModal.next({title, description})
  },
  openErrorModal (title, description) {
    this.$errorModal.next({title, description})
  }
}
