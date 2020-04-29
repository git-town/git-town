
var sqwish = require('sqwish');
var detect = require('language-classifier');
var html = require('html-minifier');
var uglify = require('uglify-js');

/**
 * Minify a `string` of an unknown language.
 *
 * @param {String} string
 * @param {Object} opts
 * @return {String}
 */

module.exports = exports = function minify (string, opts) {
  var lang = detect(string);
  var fn = exports[lang];
  if (!fn) throw new Error('Unsupported language: ' + lang);
  return fn(string, opts);
};

/**
 * Minify a Javascript `string` with optional `opts`.
 *
 * @param {String} string
 * @param {Object} opts
 * @return {String}
 */

exports.js =
exports.javascript = function (string, opts) {
  opts = opts || {};
  opts.fromString = true;
  return uglify.minify(string, opts).code;
};

/**
 * Minify a CSS `string`.
 *
 * @param {String} string
 * @return {String}
 */

exports.css = function (string) {
  return sqwish.minify(string);
};

/**
 * Minify an HTML `string`.
 *
 * @param {String} string
 * @return {String}
 */

exports.html = function (string) {
  return html.minify(string);
};
