var assert = require('should'),
    async = require('async'),
    classifier = require("../../lib/bayesian");

describe('Redis', function() {
  var data = [];

  var spam = ["vicodin pharmacy",
              "all quality replica watches marked down",
              "cheap replica watches",
              "receive more traffic by gaining a higher ranking in search engines",
              "viagra pills",
              "watches chanel tag heuer",
              "watches at low prices"];
  spam.forEach(function(text) {
    data.push({input: text, output: 'spam'});
  });

  var not = ["unknown command line parameters",
             "I don't know if this works on Windows",
             "recently made changed to terms of service agreement",
             "does anyone know about this",
             "this is a bit out of date",
             "the startup options need linbayesking"];
  not.forEach(function(text) {
    data.push({input: text, output: 'notspam'});
  });

  var bayes = new classifier.Bayesian({
    backend : { type: 'Redis' }
  });

  it('classify basic text', function(done) {
    bayes.trainAll(data, function() {
      async.series([
        function(done) {
          bayes.classify("replica watches", function(cat) {
            assert.equal(cat, "spam");
            done();
          });
        },
        function(done) {
          bayes.classify("check out the docs", function(cat) {
            assert.equal(cat, "notspam");
            done();
          });
        },
        function(done) {
          bayes.classify("recently, I've been thinking that we should", function(cat) {
            assert.equal(cat, "notspam");
            done();
          });
        },
        function(done) {
          bayes.classify("come buy these cheap pills", function(cat) {
            assert.equal(cat, "spam");
            done();
          });
        }
      ], done);
    });
  });

  it("train shouldn't require a callback", function() {
    assert.doesNotThrow(function() {
      bayes.train("cheap cialis", "spam");
      bayes.trainAll([{input: "hello dear sir", output: "notspam"}]);
    });
  })
})

describe('Redis JSON', function() {
  var expected = {"cats":{"spam":3,"notspam":2},"words":{"vicodin":{"spam":1},"pharmacy":{"spam":1},"on":{"spam":1,"notspam":1},"cheap":{"spam":1},"replica":{"spam":1},"watches":{"spam":1},"viagra":{"spam":1},"pills":{"spam":1},"unknown":{"notspam":1},"command":{"notspam":1},"line":{"notspam":1},"parameters":{"notspam":1},"I":{"notspam":1},"don":{"notspam":1},"t":{"notspam":1},"know":{"notspam":1},"if":{"notspam":1},"this":{"notspam":1},"works":{"notspam":1},"Windows":{"notspam":1}}};

  it('toJSON()', function(done) {
    var bayes = new classifier.Bayesian({
      backend : { type: 'Redis' }
    })

    var data = [{input: "vicodin pharmacy on", output: "spam"},
                {input: "cheap replica watches", output: "spam"},
                {input: "viagra pills", output: "spam"},
                {input: "unknown command line parameters", output: "notspam"},
                {input: "I don't know if this works on Windows", output: "notspam"}];

    bayes.trainAll(data, function() {
      bayes.toJSON(function(json) {
        assert.deepEqual(json, expected);
        done();
      });
    });
  });

  it('fromJSON()', function(done) {
    var bayes = new classifier.Bayesian({
      backend : { type: 'Redis' }
    });

    bayes.fromJSON(expected, function() {
      bayes.toJSON(function(json) {
        assert.deepEqual(json, expected);
        done();
      })
    });
  })
})
