const { src, dest, series } = require('gulp')
const concat = require('gulp-concat')
const sass = require('gulp-sass')
const freeze = require('gulp-freeze')
const filenames = require('gulp-filenames')
const purgecss = require('gulp-purgecss')
const noop = require('gulp-noop')
const fs = require('fs-extra')

const prod = process.env.NODE_ENV === 'production'

const sassOption = {
	outputStyle: 'compressed',
	includePaths: 'node_modules'
}

function style () {
	return src('./style/main.scss')
		.pipe(sass(sassOption).on('error', sass.logError))
		.pipe(concat('style.css'))
		.pipe(prod ? purgecss({ content: ['template/**/*.tmpl'] }) : noop())
		.pipe(prod ? freeze() : noop())
		.pipe(filenames('style'))
		.pipe(dest('./assets'))
}

async function static () {
	await fs.mkdirp('.build')
	await fs.remove('.build/static.yaml')
	await fs.appendFile('.build/static.yaml', 'style.css: ' + filenames.get('style')[0] + '\n')
}

exports.style = style
exports.static = static
exports.default = series(style, static)
