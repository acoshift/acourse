document.addEventListener('DOMContentLoaded', () => {
  const navbar = document.querySelector('#navbar')
  const delta = 5
  const navbarHeight = navbar.offsetHeight || navbar.clientHeight
  let didScroll = false
  let lastScrollTop = 0

  window.onscroll = () => {
    didScroll = true
  }

  const hasScrolled = () => {
    const st = window.scrollY

    if (Math.abs(lastScrollTop - st) <= delta) return

    if (st > lastScrollTop && st > navbarHeight) {
      navbar.classList.remove('nav-down')
      navbar.classList.add('nav-up')
    } else {
      if (st + window.innerHeight < document.documentElement.scrollHeight) {
        navbar.classList.remove('nav-up')
        navbar.classList.add('nav-down')
      }
    }
    lastScrollTop = st
  }

  setInterval(() => {
    if (didScroll) {
      hasScrolled()
      didScroll = false
    }
  }, 500)
})

/* eslint-disable no-unused-vars */
function bindFileInputImage (input, img) {
  if (!input || !img) return
  input.onchange = () => {
    const fp = input.files && input.files[0]
    if (fp) {
      const reader = new window.FileReader()
      reader.onload = (e) => {
        img.src = e.target.result
      }
      reader.readAsDataURL(fp)
    }
  }
}
