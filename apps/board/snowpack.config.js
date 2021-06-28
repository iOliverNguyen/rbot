/** @type {import("snowpack").SnowpackUserConfig } */
export default {
  exclude: [
    '**/*_*',           // nodejs files
    '**/__*',           // ignored files
    '**/testdata/**/*', // test data, not code
    '**/node_modules/**/*',
  ],
  routes: [
    {
      match: 'routes',
      src: '.*',
      dest: '/index.html',
    },
    {
      src: '/api/.*',
      dest: (req, res) => {
        return proxy.web(req, res, {
          hostname: 'localhost',
          port: 8000,
        });
      },
    },
  ],
  mount: {
    // the main source, served at /_/
    'src': '/_/',

    // to serve alias as module, the polyfill directory must be included here
    '../_polyfill': '/_/polyfill',

    // static files, served at /_/
    'static': {
      url: '/_/',
      static: true,
    },

    // the index.html file, it's served at root (/)
    'zentry': {
      url: '/',
      static: true,
    },
  },
  alias: {
    'path': '../_polyfill/path.js',
    'perf_hooks': '../_polyfill/perf-hooks.js',
  },
  plugins: [
    '@snowpack/plugin-sass',
    '@snowpack/plugin-svelte',
    ["snowpack-plugin-inliner", {
      "exts": ["jpg", "png", "svg"],
      "limit": 66666,
    }],
  ],
  packageOptions: {
    knownEntrypoints: [],
    // source: 'remote',
    // types: true,
  },
  devOptions: {},
  buildOptions: {},
  optimize: {
    bundle: true,
    minify: process.env.PROJECT_MINIFY === '1',
    target: 'es2018',
  },
};
