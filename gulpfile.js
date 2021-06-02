const { src, dest, series } = require("gulp");
const eslint = require('gulp-eslint');
const del = require('del');
const webpackStream = require("webpack-stream");
const webpack = require("webpack");

const webpackConfig = require("./webpack.config");

function clean() {
    return del(["dist/**"]);
}

function lint() {
    return src("src/**/*.tsx")
        .pipe(eslint({ useEslintrc: true, fix:true }))
        .pipe(eslint.format())
        .pipe(eslint.failAfterError())
        .pipe(dest("temp"));
}

function for_lint_change() {
    return src("temp/*.tsx")
        .pipe(dest("src"));
}

function temp_clean() {
    return del(["temp/**"]);
}

function use_webpack() {
    return webpackStream(webpackConfig, webpack)
      .pipe(dest("dist"));
}

function copy() {
    return src(["src/public/**/*.html","src/public/**/*.css"])
        .pipe(dest("dist/public"));
}

exports.default = series(clean, lint, for_lint_change, temp_clean, use_webpack, copy);
