root = window ? exports

root.wait = (delay, func) ->
  setTimeout func, delay

root.repeat = (delay, func) ->
  setInterval func, delay

root.doAndRepeat = (delay, func) ->
  func()
  setInterval func, delay

root.waitUntil = (condition, delay, func) ->
  unless func
    func = delay
    delay = 100
  g = ->
    if condition()
      func()
      clearInterval h
  h = setInterval g, delay