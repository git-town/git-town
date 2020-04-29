var assert = require('assert')
var envy = require('..')

describe('envy-json', function(){

  it('accepts a JSON string as input', function(){
    var config = "{\"basicAuth\": \"$HARP_AUTH\"}"
    process.env.HARP_AUTH = "foo:bar"
    assert.equal(envy(config).basicAuth, "foo:bar")
  })

  it('accepts a Javascript object as input', function(){
    var config = {color: "$COLOR"}
    process.env.COLOR = "chartreuse"
    assert.equal(envy(config).color, "chartreuse")
  })

  it('accepts a filename as input', function(){
    var config = envy(__dirname+"/fixtures/harp.json")
    assert.equal(config.name, "My Harp App")
  })

  it('replaces values starting with $ with their process.env equivalent', function(){
    var config = {
      name: "My Harp App",
      basicAuth: "$HARP_AUTH",
      environment: "$NODE_ENV"
    }

    process.env.HARP_AUTH = "gibraltar:snoozleton"
    process.env.NODE_ENV = "production"

    assert.equal(envy(config).name, "My Harp App")
    assert.equal(envy(config).basicAuth, "gibraltar:snoozleton")
    assert.equal(envy(config).environment, "production")
  })

  it("ignores values with $ elsewhere in the string", function(){
    var config = {
      name: "My $Harp App$$",
    }

    assert.equal(envy(config).name, "My $Harp App$$")
  })

  it("supports nested objects", function(){
    var config = {
      global: {
        environment: "$NODE_ENV"
      }
    }

    process.env.NODE_ENV = "production"

    assert.equal(envy(config).global.environment, "production")
  })

})