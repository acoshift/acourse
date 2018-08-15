const gulp = require('gulp')
const concat = require('gulp-concat')
const sass = require('gulp-sass')
const freeze = require('gulp-freeze')
const filenames = require('gulp-filenames')
const purgecss = require('gulp-purgecss')
const noop = require('gulp-noop')
const sync = require('gulp-sync')(gulp).sync
const fs = require('fs-extra')

const prod = process.env.NODE_ENV === 'production'

const sassOption = {
	outputStyle: 'compressed',
	includePaths: 'node_modules'
}

gulp.task('default', sync(['style', 'static']))

gulp.task('style', () => gulp
	.src('./style/main.scss')
	.pipe(sass(sassOption).on('error', sass.logError))
	.pipe(concat('style.css'))
	// .pipe(prod ? purgecss({ content: ['template/**/*.tmpl'] }) : noop())
	.pipe(prod ? freeze() : noop())
	.pipe(filenames('style'))
	.pipe(gulp.dest('./assets'))
)

gulp.task('static', async () => {
	await fs.mkdirp('.build')
	await fs.remove('.build/static.yaml')
	await fs.appendFile('.build/static.yaml', 'style.css: ' + filenames.get('style')[0] + '\n')
})
