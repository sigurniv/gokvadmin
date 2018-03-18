import Home from './components/Home.vue';
import Login from './components/Login.vue';

export const routes = [
    {path : '/', component: Home, name: 'home'},
    {path: '/login', component: Login, name: 'login'},
];