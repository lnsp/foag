<template>
  <div id="app" class="container">
    <div id="nav" class="row align-items-center">
      <div class="col-auto">
        <h3 class="brand">foag</h3>
      </div>
      <div class="col-auto">
        <router-link to="/">Home</router-link>
        <router-link to="/build">Build</router-link>
        <router-link to="/route">Route</router-link>
      </div>
    </div>
    <router-view class="mt-3"/>
  </div>
</template>

<style>
@import '../node_modules/bootstrap/dist/css/bootstrap.css';
@import '../node_modules/bootstrap-vue/dist/bootstrap-vue.css';
@import '../node_modules/codemirror/lib/codemirror.css';
@import '../node_modules/codemirror/theme/base16-dark.css';

#app {
  font-family: 'Fira Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto,
    Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #3b3b3b;
}

#nav {
  margin-top: 2em;
}

#nav a:hover {
  text-decoration: none;
}

#nav .router-link-exact-active {
  color: #3b3b3b;
  padding-bottom: 5px;
  border-bottom: 2px solid #3b3b3b;
}

#nav a {
  font-weight: bold;
  margin-right: 1em;
}
</style>
