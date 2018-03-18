import Vue from 'vue';
import * as Cookies from 'js-cookie'

const state = {
    token: "",
    loginError: "",
};

const mutations = {
    'SET_LOGIN_ERROR'(state, error){
        state.loginError = error;
    },
    'SET_TOKEN'(state, token){
        state.token = token;
        Cookies.set('token', token);
    }
};

const actions = {
    login: ({commit}, {login, password, router}) => {
        const url = 'auth';

        Vue.http.post(url, {login, password})
            .then(response => response.json())
            .then(data => {
                if (data.token) {
                    commit('SET_TOKEN', data.token);
                    router.push({name: 'home'})
                }

                if (data.error) {
                    commit('SET_LOGIN_ERROR', data.error);
                }
            });
    },
};

const getters = {
    loginError: state => {
        return state.loginError;
    }
};

export default {
    state,
    mutations,
    actions,
    getters
}
