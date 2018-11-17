import Vue from 'vue';
import BootstrapVue from 'bootstrap-vue';
import axios from 'axios';
import VueAxios from 'vue-axios';
import VueCodemirror from 'vue-codemirror';
import App from './App.vue';
import router from './router';


Vue.config.productionTip = false;
Vue.use(BootstrapVue);
Vue.use(VueAxios, axios);
Vue.use(VueCodemirror, { theme: 'base16-dark' });
Vue.use(require('vue-moment'));

new Vue({
  router,
  render: h => h(App),
}).$mount('#app');
