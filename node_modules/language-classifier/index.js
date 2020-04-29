
/**
 * Module dependencies.
 */

var Bayesian = require('classifier').Bayesian;

/**
 * Memory from classifier training.
 */

var memory = require('./memory.json');

/**
 * Classifier singleton.
 */

var classifier = new Bayesian;

// input training data

classifier.fromJSON(memory);

/**
 * Expose the classifier.
 */

module.exports = classifier.classify.bind(classifier);
