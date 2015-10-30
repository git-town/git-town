var classifier = require("./lib/bayesian");

var bayes = new classifier.Bayesian({
    backend : { type: 'Redis' }
  });

      bayes.trainAll([], function() {
          bayes.classify("replica watches", function(cat) {
            console.log("goes here");
          });
        })