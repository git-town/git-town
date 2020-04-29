# classifier

`classifier` is a JavaScript naive [Bayesian classifier](http://en.wikipedia.org/wiki/Bayesian_spam_filtering) with backends for Redis and localStorage:

```javascript
var bayes = new classifier.Bayesian();

bayes.train("cheap replica watches", 'spam');
bayes.train("I don't know if this works on windows", 'not');

var category = bayes.classify("free watches");   // "spam"
```

# using in node
If you have [node](http://nodejs.org/) you can install with [npm](http://github.com/isaacs/npm):

	npm install classifier

# using in the browser
Download the latest [classifier.js](http://github.com/harthur/classifier/downloads). In the browser you can only use the `localStorage` and (default) memory backends.

# Redis backend
You can store the classifier state in [Redis](http://redis.io/) for persisting and training from multiple sources:

```javascript
var bayes = new classifier.Bayesian({
  backend: {
    type: 'Redis',
    options: {
      hostname: 'localhost', // default
      port: 6379,            // default
      name: 'emailspam'      // namespace for persisting
    }
  }
});

bayes.train("cheap replica watches", "spam", function() {
  console.log("trained");
});

bayes.classify("free watches", function(category) {
  console.log("classified in: " + category);
});
```

# JSON
You can serialize and load in the classifier's state with JSON:

```javascript
var json = bayes.toJSON();

bayes.fromJSON(json);
```

# Other options
`Bayesian()` takes an options hash that you can define these properties in:

### backend
The backend property takes a `type` which is one of `'Redis'`, `'localStorage'`, or `'memory'`(default). The backend also has an `options` hash. The Redis backend takes `hostname`, `port`, `password`, `name`, `db`, and `error` (an error callback) in its options. The localStorage backend takes `name` for namespacing.

### thresholds
Specify the classification thresholds for each category. To classify an item in a category with a threshold of `x` the probably that item is in the category has to be more than `x` times the probability that it's in any other category. Default value is `1`. A common threshold setting for spam is:

```
thresholds: {
  spam: 3,
  not: 1
}
```

### default
The default category to throw an item in if it can't be classified in any of the categories. The default value of `default` is `"unclassified"`.


