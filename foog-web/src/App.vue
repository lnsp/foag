<template>
  <div id="app" class="container">
    <div id="nav" class="row align-items-center">
      <div class="col-auto">
        <h3 class="brand">foog</h3>
      </div>
      <div class="col-auto">
        <router-link to="/">Home</router-link>
        <router-link to="/deploy">Deploy</router-link>
      </div>
    </div>
    <router-view/>
  </div>
</template>

<style>
@import '../node_modules/bootstrap/dist/css/bootstrap.css';
@import '../node_modules/bootstrap-vue/dist/bootstrap-vue.css';
@import '../node_modules/codemirror/lib/codemirror.css';

#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto,
    Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #3b3b3b;
}

.brand {
  font-family: Futura;
}

#nav {
  margin-top: 2em;
}

#nav a {
  font-weight: bold;
  margin-right: 1em;
}
</style>
