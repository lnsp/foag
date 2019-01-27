<template>
  <div id="app" class="container">
    <div id="nav" class="row align-items-center">
      <div class="col-auto">
        <h3 class="brand"><img src="@/assets/logo.png" alt="foag" style="height: 1em"></h3>
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

<style lang="scss">

@import "~bootstrap/scss/functions";
@import "assets/_custom.scss";
@import "~bootstrap/scss/bootstrap";
@import '~bootstrap-vue/dist/bootstrap-vue.css';
@import '../node_modules/codemirror/lib/codemirror.css';
@import '../node_modules/codemirror/theme/base16-dark.css';

body {
  margin: 0;
}

#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto,
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
  outline: none;
}

#nav .router-link-exact-active {
  color: #3b3b3b;
  padding-bottom: 5px;
  border-bottom: 2px solid #3b3b3b;
  outline: none;
}

#nav a {
  font-weight: bold;
  margin-right: 1em;
}
</style>

