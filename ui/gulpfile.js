'use strict'

let exec = require('child_process').exec
let gulp = require('gulp')
let gutil = require('gulp-util')
let rimraf = require('rimraf')
let config = require('./config')

gulp.task('default', ['dev-server'])

gulp.task('clean', () =>
  rimraf('dist/*', err => {
    if (err) throw new gutil.PluginError('clean', err)
  })
)

gulp.task('build-prod', ['clean'], cb =>
  exec('node build/build.js', function (err, stdout, stderr) {
    gutil.log(stdout)
    gutil.log(stderr)
    cb(err)
  })
)

gulp.task('dev-server', () => {
  exec('node build/dev-server.js', function (err) {
    if (err) {
      throw new gutil.PluginError('dev-server', err)
    }
  })

  gutil.log('Listening at', gutil.colors.magenta('http://localhost:' + config.dev.port))
})
