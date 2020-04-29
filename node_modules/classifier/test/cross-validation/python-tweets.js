var url = require("url"),
    assert = require("should"),
    cradle = require("cradle"),
    _ = require("underscore"),
    crossValidate = require("./cross-validate"),
    classifier = require("../../lib/bayesian");

function getData(couchUrl, callback) {
  var parts = url.parse(couchUrl);

  var client = new cradle.Connection(parts.hostname, parts.port || 80);
  var db = client.database(parts.pathname.replace(/^\//, ""));


  db.all({include_docs: true}, function(err, res) {
    if (err) {
      console.log("error getting data from " + url + ": ");
      console.log(err);
    }
    else {
      var data = _(res.rows).pluck("doc");
      callback(data);
    }
  });
}

function runTest(url, callback) {
  getData(url, function(data) {
    var result = crossValidate(classifier.Bayesian, {}, data);
    callback(result);
  });
}

describe('Bayesian cross-validation', function() {
  it('classify Python snake vs. language tweets with error < 0.3', function(done) {
    var couchUrl = "http://harth.iriscouch.com/pythontweets";
    runTest(couchUrl, function(result) {
      console.log("Cross-validation of Python tweets:\n");
      console.log(result);

      assert.ok(result.error < .3);
      done();
    })
  })
})
