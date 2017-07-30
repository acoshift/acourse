const md5 = require('md5')
const sass = require('node-sass')
const fs = require('fs')
const babel = require('babel-core')

let config = ''

let buf = sass.renderSync({
  file: 'style/main.scss',
  outputStyle: 'compressed'
}).css

let outFile = 'style.' + md5(buf) + '.css'
fs.writeFileSync('static/' + outFile, buf)
config += 'style.css: ' + outFile + '\n'

buf = babel.transformFileSync('script/main.js', {
  presets: [ 'env' ],
  minified: true,
  compact: true
}).code
outFile = 'script.' + md5(buf) + '.js'
fs.writeFileSync('static/' + outFile, buf)
config += 'script.js: ' + outFile + '\n'

fs.writeFileSync('static.yaml', config)
