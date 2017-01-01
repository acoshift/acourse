// https://github.com/shelljs/shelljs
/* globals env mkdir cp */
require('shelljs/global')
env.NODE_ENV = 'production'

const path = require('path')
const config = require('../config')
const ora = require('ora')
const webpack = require('webpack')
const webpackConfig = require('./webpack.local.conf')

const spinner = ora('building for local...')
spinner.start()

const assetsPath = path.join(config.local.assetsRoot, config.local.assetsSubDirectory)
// rm('-rf', assetsPath)
mkdir('-p', assetsPath)
cp('-R', 'static/*', assetsPath)

webpack(webpackConfig, (err, stats) => {
  spinner.stop()
  if (err) throw err
  process.stdout.write(stats.toString({
    colors: true,
    modules: false,
    children: false,
    chunks: false,
    chunkModules: false
  }) + '\n')
})
