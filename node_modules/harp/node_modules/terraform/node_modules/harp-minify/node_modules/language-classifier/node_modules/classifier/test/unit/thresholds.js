var assert = require('should'),
    classifier = require("../../lib/bayesian");

describe('thresholds', function() {
  var bayes = new classifier.Bayesian({
    thresholds: {
      spam: 3,
      notspam: 1
    }
  });

  var spam = ["a c", "b a", "c e"];
  spam.forEach(function(text) {
    bayes.train(text, 'spam');
  });

  var not = ["d e", "e f", "f b"];
  not.forEach(function(text) {
    bayes.train(text, 'notspam');
  });

  it('categorize with default thresholds', function() {
    assert.equal(bayes.classify("a"), "spam");
    assert.equal(bayes.classify("b"), "notspam");
    assert.equal(bayes.classify("c"), "spam");
    assert.equal(bayes.classify("d"), "notspam");
    assert.equal(bayes.classify("e"), "notspam");
  })

  it('categorize with really high thresholds', function() {
    bayes.setThresholds({spam: 4, notspam: 4});

    assert.equal(bayes.classify("a"), "unclassified");
    assert.equal(bayes.classify("b"), "unclassified");
    assert.equal(bayes.classify("c"), "unclassified");
    assert.equal(bayes.classify("d"), "unclassified");
    assert.equal(bayes.classify("e"), "unclassified");
  })
})
