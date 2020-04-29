var fs          = require('fs')
var path        = require('path')
var stylesheet  = require('./stylesheet')
var template    = require('./template')
var javascript  = require('./javascript')
var helpers     = require('./helpers')
var lodash = require('lodash')

/**
 * Expose Helpers
 *
 * We expose the helpers so that other libraries my use them.
 * Terraform is much more useful with these helpers.
 *
 */

exports.helpers = helpers


/**
 * Root
 *
 * This sets the home base directory for the app which affects
 * where we begin walking to build the home directory.
 *
 */

exports.root = function(root, globals){

  if(!globals){
    globals  = {}
  }

  var LayoutPriorityList = helpers.buildPriorityList("_layout.html")
  var layout  = helpers.findFirstFile(root, LayoutPriorityList)

  var data    = helpers.dataTree(root)

  var templateObject = { public: data }

  for(var key in globals){
    templateObject[key] = globals[key]
  }


  return {

    /**
     * Render
     *
     * This is the main method to to render a view. This function is
     * responsible to for figuring out the layout to use and sets the
     * `current` object.
     *
     */

    render: function(filePath, locals, callback){

      // get rid of leading slash (windows)
      filePath = filePath.replace(/^\\/g, '')

      // locals are optional
      if(!callback){
        callback = locals
        locals   = {}
      }


      /**
       * We ignore files that start with underscore
       */

      if(helpers.shouldIgnore(filePath)) return callback(null, null)


      /**
       * If template file we need to set current and other locals
       */

      if(helpers.isTemplate(filePath)) {

        /**
         * Current
         */
        locals._ = lodash
        locals.current = helpers.getCurrent(filePath)


        /**
         * Layout Priority:
         *
         *    1. passed into partial() function.
         *    2. in `_data.json` file.
         *    3. default layout.
         *    4. no layout
         */

        // 1. check for layout passed in
        if(!locals.hasOwnProperty('layout')){

          // 2. _data.json layout
          // TODO: Change this lookup relative to path.
          var templateLocals = helpers.walkData(locals.current.path, data)

          if(templateLocals && templateLocals.hasOwnProperty('layout')){
            if(templateLocals['layout'] === false){
              locals['layout'] = null
            } else if(templateLocals['layout'] !== true){

              // relative path
              var dirname = path.dirname(filePath)
              var layoutPriorityList = helpers.buildPriorityList(path.join(dirname, templateLocals['layout'] || ""))

              // absolute path (fallback)
              layoutPriorityList.push(templateLocals['layout'])

              // return first existing file
              // TODO: Throw error if null
              locals['layout'] = helpers.findFirstFile(root, layoutPriorityList)

            }
          }

          // 3. default _layout file
          if(!locals.hasOwnProperty('layout')){
            locals['layout'] = helpers.findDefaultLayout(root, filePath)
          }

          // 4. no layout (do nothing)
        }

        /**
         * TODO: understand again why we are doing this.
         */

        try{
          var error  = null
          var output = template(root, templateObject).partial(filePath, locals)
        }catch(e){
          var error  = e
          var output = null
        }finally{
          callback(error, output)
        }

      }else if(helpers.isStylesheet(filePath)){
        stylesheet(root, filePath, callback)
      }else if(helpers.isJavaScript(filePath)){
        javascript(root, filePath, callback)
      }else{
        callback(null, null)
      }


    }
  }

}
