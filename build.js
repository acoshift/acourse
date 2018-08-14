const md5 = require('md5')
const sass = require('node-sass')
const fs = require('fs')

let config = ''

let buf = sass.renderSync({
  file: 'style/main.scss',
  outputStyle: 'compressed'
}).css

let outFile = 'style.' + md5(buf) + '.css'
fs.writeFileSync('static/' + outFile, buf)
config += 'style.css: ' + outFile + '\n'

fs.writeFileSync('static.yaml', config)
