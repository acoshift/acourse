const md5 = require('md5')
const sass = require('node-sass')
const fs = require('fs')

const buf = sass.renderSync({
  file: 'style/main.scss',
  outputStyle: 'compressed'
}).css

const hash = md5(buf)
const outFile = 'style.'+hash+'.css'

fs.writeFileSync('static/' + outFile, buf)
fs.writeFileSync('static.yaml', 'style.css: ' + outFile)
