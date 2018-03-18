import Vue from 'vue';
import Vuex from 'vuex';
import keys from './modules/keys';
import engines from './modules/engines';
import auth from './modules/auth';

Vue.use(Vuex);

export default new Vuex.Store({
    modules: {
        keys,
        engines,
        auth,
    }
})

