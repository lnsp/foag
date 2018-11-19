import Vue from 'vue';
import Router from 'vue-router';
import Home from './views/Home.vue';

Vue.use(Router);

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
    },
    {
      path: '/build',
      name: 'build',
      component: () => import(/* webpackChunkName: "about" */ './views/Build.vue'),
    },
    {
      path: '/details/:id',
      name: 'details',
      component: () => import('./views/Details.vue'),
    },
    {
      path: '/route/:id',
      name: 'addRoute',
      component: () => import('./views/Route.vue'),
    },
    {
      path: '/route',
      name: 'route',
      component: () => import('./views/Route.vue'),
    },
  ],
});
