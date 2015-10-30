var path           = require('path')
var fs             = require('fs')
var TerraformError = exports.TerraformError = require("../error").TerraformError


/**
 * This is our processor map for all the helpers.
 *
 * This is our map. This ideally this would be dynamically created by
 * some sort of plugin architecture but lets not let perfect be the
 * enemy of good :)
 *
 */

var processors = exports.processors = {
  "html": ["jade", "ejs", "md"],
  "css" : ["styl", "less", "scss", "sass"],
  "js"  : ["coffee"]
}


/**
 * Build Priority List
 *
 * returns a priority list of files to look for on a given request.
 *
 * `css` and `html` are special extensions that will add `less` and `jade`
 * to the priority list (respectively).
 *
 * e.g
 *
 *    buildPriorityList("foobar")
 *      => ["foobar", "foobar.html", "foobar.jade", "foobar.html.jade"]
 *
 *    buildPriorityList("foobar.css")
 *      => ["foobar.css", "foobar.less", "foobar.css.less"]
 *
 *    buildPriorityList("foobar.html")
 *      => ["foobar.html", "foobar.jade", "foobar.html.jade"]
 *
 *    buildPriorityList("foobar.jade")
 *      => ["foobar.jade"]
 *
 *    buildPriorityList("foobar.html.jade.html")
 *      => ["foobar.html.jade.html", "foobar.html.jade.jade", "foobar.html.jade.html.jade"]
 *
 *    buildPriorityList("hello/foobar")
 *      => ["hello/foobar", "hello/foobar.html", "hello/foobar.jade", "hello/foobar.html.jade"]
 *
 */

var buildPriorityList = exports.buildPriorityList = function(filePath){

  var list = []

  /**
   * get extension
   */

  var ext       = path.extname(filePath).replace(/^\./, '')
  var processor = processors[ext]

  if(processor){

    // foo.html => foo.jade
    processor.forEach(function(p){
      var regexp = new RegExp(ext + '$')
      list.push(filePath.replace(regexp, p))
    })

    // foo.html => foo.html.jade
    processor.forEach(function(p){
      list.push(filePath + '.' + p)
    })

  }else{
    // assume template when unknown processor
    if(processors['html'].indexOf(ext) !== -1){
      list.push(filePath)
    }else{
      // foo.xml => foo.xml.jade
      processors['html'].forEach(function(p){
        list.push(filePath + '.' + p)
      })
    }
  }

  // remove leading and trailing slashes
  var list = list.map(function(item){ return item.replace(/^\/|^\\|\/$/g, '') })

  return list
}


/**
 * Find First File
 *
 * Takes a directory and an array of files. Returns the first file in the list that exists.
 *
 *    findFirstFile("/path/to/public", ["foo.html", "foo.jade", "foo.html.jade"])
 *      => "foo.jade"
 *
 * returns null if no file is found.
 *
 */

var findFirstFile = exports.findFirstFile = function(dir, arr) {
  var dirPath   = path.dirname(path.join(dir, arr[0]))
  var fullPath  = path.resolve(dirPath)

  try{
    var list = fs.readdirSync(fullPath)
  }catch(e){
    var list = []
  }

  var first = null

  if(list){
    arr.reverse().map(function(item){
      var fileName = path.basename(item)
      if(list.indexOf(fileName) !== -1){
        first = item
      }
    })
  }

  return first
}


var layoutPriorityList = buildPriorityList("_layout.html")

/**
 * Find Nearest Layout
 *
 * Walks up the tree to find the nearest _layout file. Returns relative path to root
 *
 *    findNearestLayout("/path/to/public", "/path/to/public/nested/dir")
 *      => "_layout.jade"
 *
 * returns null if no layout is found.
 *
 */

var findNearestLayout = exports.findNearestLayout = function(rootPath, dirPath) {

  // lets make sure we are working with an absolute path
  var dirPath = path.resolve(rootPath, (dirPath || ""))
  var layout  = findFirstFile(dirPath, layoutPriorityList)

  // if we have a layout we return relative path
  if(layout !== null)
    return path.relative(rootPath, path.join(dirPath, layout))

  return path.relative(rootPath, dirPath) !== ''
    ? findNearestLayout(rootPath, path.resolve(dirPath, ".."))
    : null  // we reached the root
}


/**
 * Is Empty
 *
 * Checks if Object is empty & returns true or false.
 *
 */

var isEmpty = function(obj) {
  for(var prop in obj) {
    if(Object.prototype.hasOwnProperty.call(obj,prop))
      return false;
  }
  return true;
}


/**
 *
 * Walks directory and build the data object.
 *
 * If we call the dataTree on the public dir
 *
 *     public/
 *       |- _data.json
 *       |- articles/
 *       |   `- _data.json
 *       `- people
 *           `- _data.json
 *
 * We get the following...
 *
 *     {
 *       "_data": {...},
 *       "articles": {
 *         "_data": {...}
 *       },
 *       "people": {
 *         "_data": {...}
 *       }
 *     }
 *
 */

var dataTree = exports.dataTree = function (filename) {
  var dirPath   = path.resolve(filename)
  try{
    var list = fs.readdirSync(dirPath)
  }catch(e){
    e.source    = "Config"
    e.dest      = "Config"
    e.lineno    = -1
    e.filename  = filename
    e.stack     = null
    throw new TerraformError(e)
  }

  var obj       = {}
  obj._contents  = []

  try{
    var dataPath = path.resolve(dirPath, "_data.json")
    var fileData = fs.readFileSync(dataPath)
    obj._data     = JSON.parse(fileData)
  }catch(e){
    if(e.code && e.code === "ENOENT"){
      // data file failed or does not exist
    }else{
      e.source    = "Data"
      e.dest      = "Globals"
      e.lineno    = -1
      e.filename  = dataPath
      e.stack     = fileData.toString()
      throw new TerraformError(e)
    }
    //console.log(e.code)

  }

  list.forEach(function(file){
    var filePath = path.resolve(dirPath, file)
    var stat     = fs.statSync(filePath)

    if(stat.isDirectory()){
      if(['_', '.'].indexOf(file[0]) === -1){
        var d = dataTree(filePath)
        if(!isEmpty(d))
          obj[file] = d
      }
    }else{
      if(["_", "."].indexOf(file[0]) === -1 ) obj._contents.push(outputPath(file))
    }
  })

  if(obj._contents.length == 0)
    delete obj._contents

  return obj
}


/**
 *
 * Get Current
 *
 * returns current object based on the path of source file
 *
 *    getCurrent("foo/bar/baz.jade")
 *      => { path: ["foo", "bar", "baz"], source: "baz" }
 *
 *    getCurrent("index.html")
 *      => { path: ["index"], source: "index" }
 *
 */

exports.getCurrent = function(sourcePath){

  // windows
  sourcePath = sourcePath.replace(/\\/g,'/')

  // this could be a tad smarter
  var namespace = sourcePath.split(".")[0].split("/")

  return {
    source: namespace[namespace.length -1],
    path: namespace
  }
}


/**
 *
 * Source Type
 *
 * Returns processor based on file path.
 *
 *    sourceType("foobar.jade")  =>  "jade"
 *    sourceType("foobar.less")  =>  "less"
 *    sourceType("foobar")       =>  null
 *
 */

exports.sourceType = function(sourcePath){
  return path.extname(sourcePath).replace(/^\./, '')
}


/**
 *
 * Walk Data Tree
 *
 * Recursive function that returns the data object accociated with path.
 *
 *     var globals = {
 *       "public": {
 *         "articles": {
 *           "_data": {
 *             "hello-world": "You Found Me!"
 *           }
 *         }
 *       }
 *     }
 *
 *     walkData(["public", "articles", "hello-world"], globals) => "You Found Me!"
 */

var walkData = exports.walkData = function(tail, obj){
  var tail = tail.slice(0)  // clone array.
  var head = tail.shift()

  if(obj.hasOwnProperty(head)){
    return walkData(tail, obj[head])
  }else if(obj.hasOwnProperty("_data")){
    return obj["_data"][head]
      ? obj["_data"][head]
      : null

  }else{
    return null
  }
}


/**
 *
 * Output Path
 *
 * Returns output path output for given source file.
 *
 * allowAlternateExtensions: true. May be turned off for one less RegExp execution
 * (and faster performance)
 *
 * eg.
 *     foo.jade => foo.html
 *     foo.html.jade => foo.html
 */

var outputPath = exports.outputPath = function(source, allowAlternateExtensions){
  if (allowAlternateExtensions == undefined) {
    allowAlternateExtensions = true;
  }
  for(var targetExtension in processors){ // .html, .css, .js
    if (processors.hasOwnProperty(targetExtension)) {
      processors[targetExtension].forEach(function(sourceExtension){ // .jade, .ejs, .md
        if (allowAlternateExtensions) {
          // Search for a alternate extension before the known source extension e.g. foobar.bar.jade
          var alternateFileExtension = new RegExp("^.*\\.(\\w{3,4})\\." + sourceExtension + "$")
          var match = alternateFileExtension.exec(source)
          if (match) {
            // If found, replace the target extension with the alternate
            targetExtension = match[1]
          }
        }
        // Let's try to match known sourceExtensions. Optionally, capture the targetExtension BEFORE.
        var regexp = new RegExp('(\\.' + targetExtension + ')?\\.' + sourceExtension + '$')
        // If found, will replace the whole extension with target extension, e.g.:
        // .html.jade -> .html (matches)
        // .jade      -> .html (matches)
        // .html      -> .html (doesn't matches - doesn't matter)
        // .foo       -> .foo (doesn't matches - doesn't matter)
        source = source.replace(regexp, "." + targetExtension)
      })
    }
  }

  return source
}


/**
 *
 * Output Type
 *
 * Returns output type source file.
 *
 * eg.
 *     foo.jade       => html
 *     foo.html.jade  => html
 */

var outputType = exports.outputType = function(source){
  var op = outputPath(source)
  return path.extname(op).replace(/^\./, '')
}

/**
 *
 * Should Ignore
 *
 * Checks to see if path should be ignored.
 *
 * eg.
 *     shouldIgnore('_foo.html')         => true
 *     shouldIgnore('foo_.html')         => false
 *     shouldIgnore('_foo/bar.html')     => true
 *     shouldIgnore('foo/_bar.html')     => true
 *     shouldIgnore('foo/_bar/baz.html') => true
 */

exports.shouldIgnore = function(filePath){

  // remove starting and trailing slashed
  filePath = filePath.replace(/^\/|\/$/g, '')

  // create array out of path
  var arr = filePath.split(path.sep)

  // test for starting underscore, .git, .gitignore
  var map = arr.map(function(item){
    return item[0] === "_" || item.indexOf(".git") === 0
  })

  // return if any item starts with underscore
  return map.indexOf(true) !== -1
}


/**
 *
 * isTemplate
 *
 * returns true if file is a template file
 *
 * eg.
 *     isTemplate('foo.jade')         => true
 *     isTemplate('foo.md')           => true
 *     isTemplate('foo.html')         => false
 *     isTemplate('foo/bar.jade')     => true
 *     isTemplate('foo/bar.md')       => true
 *     isTemplate('foo/bar.html')     => false
 *     isTemplate('foo.less')         => false
 *     isTemplate('foo.css')          => false
 *     isTemplate('foo.bar.baz.jade') => true
 */

exports.isTemplate = function(filePath){
  var ext = path.extname(filePath).replace(/^\./, '')

  return processors["html"].indexOf(ext) !== -1
}

/**
 *
 * isStylesheet
 *
 * returns true if file is a pre-processor stylesheet file
 *
 * eg.
 *     isTemplate('foo.less')         => true
 *     isTemplate('foo.md')           => false
 *     isTemplate('foo.css')          => false
 *     isTemplate('foo.bar.baz.less') => true
 */

exports.isStylesheet = function(filePath){
  var ext = path.extname(filePath).replace(/^\./, '')

  return processors["css"].indexOf(ext) !== -1
}


/**
 * isJavaScript(filePath)
 *
 * returns true if file is a pre-processor javascript file
 *
 * eg.
 *     isJavaScript('foo.coffee')         => true
 *     isJavaScript('foo.md')             => false
 *     isJavaScript('foo.css')            => false
 *     isJavaScript('foo.bar.baz.coffee') => true
 */

exports.isJavaScript = function(filePath){
  var ext = path.extname(filePath).replace(/^\./, '')

  return processors["js"].indexOf(ext) !== -1
}
