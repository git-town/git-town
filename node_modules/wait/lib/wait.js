(function() {
  var root;
  root = typeof window != "undefined" && window !== null ? window : exports;
  root.wait = function(delay, func) {
    return setTimeout(func, delay);
  };
  root.repeat = function(delay, func) {
    return setInterval(func, delay);
  };
  root.doAndRepeat = function(delay, func) {
    func();
    return setInterval(func, delay);
  };
  root.waitUntil = function(condition, delay, func) {
    var g, h;
    if (!func) {
      func = delay;
      delay = 100;
    }
    g = function() {
      if (condition()) {
        func();
        return clearInterval(h);
      }
    };
    return h = setInterval(g, delay);
  };
}).call(this);
