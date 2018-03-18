import Vue from 'vue';
import General from '../../components/engines/General.vue';

const state = {
    engines: {
        'default': General,
        'boltdb': General,
        'badger': General
    },
    selectedEngine: 'default',
};

const mutations = {
    'SET_SELECTED_ENGINE'(state, engine) {
        state.selectedEngine = engine;
    },
};

const actions = {
    init: ({commit}) => {
        const url = `init`;

        Vue.http.get(url)
            .then(response => response.json())
            .then(data => {
                if (data) {
                    if (data.engine) {
                        if (data.engine in state.engines) {
                            commit('SET_SELECTED_ENGINE', data.engine);
                        }
                    }
                }
            })
    },
};

const getters = {
    selectedEngine: state => {
        return state.selectedEngine;
    },
    engines: state => {
        return state.engines;
    },
    selectedEngineComponent: state => {
        if (state.selectedEngine in state.engines) {
            return state.engines[state.selectedEngine];
        }

        return state.engines['default'];
    },
};

export default {
    state,
    mutations,
    actions,
    getters
}
