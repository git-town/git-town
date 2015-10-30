var assert = require('should'),
    classifier = require("../../lib/bayesian");

describe('synchronous backends', function() {
  it('classify with in-memory backend', function() {
    testBasic(new classifier.Bayesian());
  })

  it('classify with localStorage backend', function() {
    testBasic(new classifier.Bayesian({
      backend : {
        type: 'localStorage',
        options: {
          name: 'testnamespace',
          testing: true
        }
      }
    }))
  })
})

describe('synchronous backends', function() {
  var expected = {"cats":{"spam":3,"notspam":2},"words":{"vicodin":{"spam":1},"pharmacy":{"spam":1},"on":{"spam":1,"notspam":1},"cheap":{"spam":1},"replica":{"spam":1},"watches":{"spam":1},"viagra":{"spam":1},"pills":{"spam":1},"unknown":{"notspam":1},"command":{"notspam":1},"line":{"notspam":1},"parameters":{"notspam":1},"I":{"notspam":1},"don":{"notspam":1},"t":{"notspam":1},"know":{"notspam":1},"if":{"notspam":1},"this":{"notspam":1},"works":{"notspam":1},"Windows":{"notspam":1}}};

  describe('toJSON', function() {
    var data = [{input: "vicodin pharmacy on", output: "spam"},
            {input: "cheap replica watches", output: "spam"},
            {input: "viagra pills", output: "spam"},
            {input: "unknown command line parameters", output: "notspam"},
            {input: "I don't know if this works on Windows", output: "notspam"}];
    it('toJSON() with memory backend', function() {
      var bayes = new classifier.Bayesian();
      bayes.trainAll(data);
      assert.deepEqual(bayes.toJSON(), expected);
    });

    it('toJSON() with localStorage backend', function() {
      var bayes = new classifier.Bayesian({
        backend : {
          type: 'localStorage',
          options: {
            name: 'testnamespace',
            testing: true
          }
        }
      });
      bayes.trainAll(data);
      assert.deepEqual(bayes.toJSON(), expected);
    })
  });

  describe('fromJSON', function() {
    it('toJSON() with memory backend', function() {
      var bayes = new classifier.Bayesian();
      bayes = bayes.fromJSON(expected);

      assert.deepEqual(bayes.toJSON(), expected);
    });

    it('toJSON() with localStorage backend', function() {
      var bayes = new classifier.Bayesian({
        backend : {
          type: 'localStorage',
          options: {
            name: 'testnamespace',
            testing: true
          }
        }
      });
      bayes = bayes.fromJSON(expected);

      assert.deepEqual(bayes.toJSON(), expected);
    })
  })
})

function testBasic(bayes) {
  var spam = ["vicodin pharmacy",
              "all quality replica watches marked down",
              "cheap replica watches",
              "receive more traffic by gaining a higher ranking in search engines",
              "viagra pills",
              "watches chanel tag heuer",
              "watches at low prices"];
  spam.forEach(function(text) {
    bayes.train(text, 'spam');
  });

  var not = ["unknown command line parameters",
             "I don't know if this works on Windows",
             "recently made changed to terms of service agreement",
             "does anyone know about this",
             "this is a bit out of date",
             "the startup options need linking"];
  not.forEach(function(text) {
    bayes.train(text, 'notspam');
  });

  assert.equal(bayes.classify("replica watches"),"spam");
  assert.equal(bayes.classify("check out the docs"), "notspam");
  assert.equal(bayes.classify("recently, I've been thinking that I should"), "notspam");
  assert.equal(bayes.classify("come buy these cheap pills"), "spam");
}
