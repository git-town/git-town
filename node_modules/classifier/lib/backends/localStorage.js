var _ = require("underscore")._;

var LocalStorageBackend = function(options) {
  var options = options || {};
  var name = options.name || Math.floor(Math.random() * 100000);

  this.prefix = 'classifier.bayesian.' + name;

  if (options.testing) {
    this.storage = {};
  }
  else {
    this.storage = localStorage;
  }

  this.storage[this.prefix + '.cats'] = '{}';
}

LocalStorageBackend.prototype = {
  async : false,

  getCats : function() {
    return JSON.parse(this.storage[this.prefix + '.cats']);
  },

  setCats : function(cats) {
    this.storage[this.prefix + '.cats'] = JSON.stringify(cats);
  },

  getWordCount : function(word) {
    return JSON.parse(this.storage[this.prefix + '.words.' + word] || '{}');
  },

  setWordCount : function(word, counts) {
    this.storage[this.prefix + '.words.' + word] = JSON.stringify(counts);
  },

  getWordCounts : function(words) {
    var counts = {};
    words.forEach(function(word) {
      counts[word] = this.getWordCount(word);
    }, this);
    return counts;
  },

  incCounts : function(catIncs, wordIncs) {
    var cats = this.getCats();
    _(catIncs).each(function(inc, cat) {
      cats[cat] = cats[cat] + inc || inc;
    }, this);
    this.setCats(cats);

    _(wordIncs).each(function(incs, word) {
      var wordCounts = this.getWordCount(word);
      _(incs).each(function(inc, cat) {
        wordCounts[cat] = wordCounts[cat] + inc || inc;
      }, this);
      this.setWordCount(word, wordCounts);
    }, this);
  },

  toJSON : function() {
    var words = {};
    var regex = new RegExp("^" + this.prefix + "\.words\.(.+)$")
    for (var item in this.storage) {
      var match = regex.exec(item);
      if (match) {
        words[match[1]] = JSON.parse(this.storage[item]);
      }
    }
    return {
      cats: JSON.parse(this.storage[this.prefix + '.cats']),
      words: words
    };
  },

  fromJSON : function(json) {
    this.incCounts(json.cats, json.words);
  }
}

exports.LocalStorageBackend = LocalStorageBackend;