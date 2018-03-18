import Vue from 'vue'
import VueRouter from 'vue-router';
import VueResource from 'vue-resource';

import App from './App.vue'
import {routes} from './routes';
import store from './store/store';
import * as Cookies from 'js-cookie'

Vue.use(VueRouter);
Vue.use(VueResource);

Vue.http.options.root = process.env.API_URL;

Vue.http.interceptors.push((request, next) => {
    request.headers.set('Authorization', store.state.auth.token);
    next((response) => {
        if (response.status == 401) {
            store.commit('SET_TOKEN', "");
            if (router.history.current.name != 'login') {
                router.go('/login');
            }
        }
    });
});

const router = new VueRouter({
    mode: 'history',
    routes
});


router.beforeEach((to, from, next) => {
    const token = Cookies.get("token") || "";
    store.commit('SET_TOKEN', token);

    if (to.name != 'login' && !store.state.auth.token.length) {
        next("/login")
    } else {
        next();
    }
});

new Vue({
    el: '#app',
    render: h => h(App),
    store,
    router,
});
