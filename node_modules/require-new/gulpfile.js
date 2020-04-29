var gulp = require('gulp'),
    mocha = require('gulp-mocha'),
    eslint = require('gulp-eslint'),
    GitDown = require('gitdown');

gulp.task('lint', function () {
    return gulp
        .src(['./src/**/*.js', './src/tests/**/*.js'])
        .pipe(eslint())
        .pipe(eslint.format())
        .pipe(eslint.failOnError());
});

gulp.task('test', ['lint'], function () {
    return gulp
        .src(['./tests/*.js'], {read: false})
        .pipe(mocha());
});

gulp.task('gitdown', function () {
    return GitDown
        .read('.gitdown/README.md')
        .write('README.md');
});

gulp.task('watch', function () {
    gulp.watch(['./src/*', './tests/*'], ['default']);
    gulp.watch(['./.gitdown/*'], ['gitdown']);
});

gulp.task('default', ['test']);
