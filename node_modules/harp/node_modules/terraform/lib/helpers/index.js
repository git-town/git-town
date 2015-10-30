
/**
 * If evironment is production we load the memozied
 * file instead of the raw file so that the functions
 * are cached. In development we can't have memoization
 * because the developer may have change some files
 * and they shouldn't have to restart the server.
 *
 */

module.exports = (process.env.NODE_ENV == "production")
  ? require('./memoized')
  : require('./raw')
