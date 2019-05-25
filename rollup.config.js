import svelte from 'rollup-plugin-svelte'
import resolve from 'rollup-plugin-node-resolve'
import commonjs from 'rollup-plugin-commonjs'
import livereload from 'rollup-plugin-livereload'
import { terser } from 'rollup-plugin-terser'
import scss from 'rollup-plugin-scss'
import fs from 'fs-extra'
import crypto from 'crypto'
import purgecss from 'rollup-plugin-purgecss'
import { scss as svelteScss } from '@kazzkiq/svelte-preprocess-scss'

const production = !process.env.ROLLUP_WATCH

function hash (text) {
	if (!production) {
		return ''
	}
	return '_' + crypto
		.createHash('md5')
		.update(text, "utf8")
		.digest("hex")
}

export default {
	input: 'main.js',
	output: {
		sourcemap: true,
		format: 'iife',
		name: 'app',
		file: '.build/main.js'
	},
	plugins: [
		svelte({
			// enable run-time checks when not in production
			dev: !production,
			customElement: true,
			css: true,
			// we'll extract any component CSS out into
			// a separate file — better for performance
			// css: css => {
			// 	css.write('public/bundle.css');
			// }
			preprocess: {
				style: svelteScss()
			}
		}),
		{
			name: 'prepare',
			buildStart: async () => {
				await fs.remove('.build')
				await fs.mkdirp('.build')
			}
		},
		resolve(),
		commonjs(),
		purgecss({
			content: ['template/**/*.tmpl']
		}),
		scss({
			output: (styles) => {
				const fn = `style${hash(styles)}.css`
				fs.writeFileSync(fn, 'assets/' + styles, 'utf8')
				fs.appendFileSync('.build/static.yaml', `style.css: ${fn}` + '\n', 'utf8')
			},
			outputStyle: 'compressed',
			includePaths: [ 'node_modules' ]
		}),
		{
			name: 'output',
			onwrite: function(bundle, data) {
				const fn = `components${hash(data.code)}.js`

				fs.unlinkSync(bundle.file)
				fs.appendFileSync('.build/static.yaml', `components.js: ${fn}` + '\n', 'utf8')

				// let code = data.code;
				// if (bundle.sourcemap) {
				// 	const basename = path.basename(fileName);
				// 	data.map.file = basename;
				//
				// 	let url;
				// 	if (bundle.sourcemap === 'inline') {
				// 		url = data.map.toUrl();
				// 	} else {
				// 		url = basename + '.map';
				// 		fs.writeFileSync(fileName + '.map', data.map.toString());
				// 	}
				//
				// 	code += `\n//# sourceMappingURL=${url}`;
				// }

				fs.writeFileSync('assets/' + fn, data.code, 'utf8')
			}
		},

		// Watch the `public` directory and refresh the
		// browser on changes when not in production
		// !production && livereload('public'),

		// If we're building for production (npm run build
		// instead of npm run dev), minify
		production && terser()
	],
	watch: {
		clearScreen: false
	}
}
