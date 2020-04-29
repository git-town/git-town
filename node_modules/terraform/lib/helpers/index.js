
/**
 * If evironment is production we load the memozied
 * file instead of the raw file so that the functions
 * are cached. In development we can't have memoization
 * because the developer may have change some files
 * and they shouldn't have to restart the server.
 *
 */

var memoizationIsEnabled = (process.env.NODE_ENV == 'production' && !~process.argv.indexOf('--disable-harp-cache')) || ~process.argv.indexOf('--enable-harp-cache')

module.exports = memoizationIsEnabled
  ? require('./memoized')
  : require('./raw')
