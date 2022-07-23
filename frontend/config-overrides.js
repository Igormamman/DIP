/* config-overrides.js */

module.exports = function override(config, env) {
    //do stuff with the webpack config...
    console.log('override')
    let loaders = config.resolve
    loaders.fallback = {
       "stream": require.resolve("stream-browserify"),  
    }

    return config;
  }