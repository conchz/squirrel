'use strict';

const path = require('path'),
  exec = require('child_process').exec,
  gulp = require('gulp'),
  gutil = require('gulp-util'),
  rimraf = require('rimraf');

gulp.task('default', ['dev-server']);

gulp.task('clean', () =>
  rimraf('dist/*', err => {
    if (err) throw new gutil.PluginError('clean', err)
  })
);

gulp.task('build-prod', ['clean'], cb =>
  exec('node build/build.js', function (err, stdout, stderr) {
    gutil.log(stdout);
    gutil.log(stderr);
    cb(err);
  })
);

gulp.task('dev-server', () => {
  exec('node build/dev-server.js', function (err) {
    if (err) {
      throw new gutil.PluginError('dev-server', err)
    }
  });
});
