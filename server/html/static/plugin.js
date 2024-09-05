function compose(plugins) {
    return function (msg, next) {
      let index = -1
      return dispatch(0)
      function dispatch (i) {
        if (i <= index) return 'next() called multiple times'
        index = i
        let fn = plugins[i]
        if (i === plugins.length) fn = next
        if (!fn) return true
        if(typeof(fn) !== 'function') return 'plugin ' + (i+1) + ' must be composed of functions!'
        try {
          return fn(msg, dispatch.bind(null, i + 1))
        } catch (err) {
          return err
        }
      }
    }
}

const plugins = [];
// plugin_list

self.onmessage = function(e) {
    const msg = e.data;
    try {
        compose(plugins)(e.data);
    } catch(err) {}
    self.postMessage(msg);
}