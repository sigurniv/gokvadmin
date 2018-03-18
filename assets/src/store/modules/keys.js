import Vue from 'vue';

const state = {
    searchKey: '',
    values: [],
    addingKey: false,
    newKey: {
        key: "",
        value: "",
    },
    success: "",
    error: "",
    pagination: {
        hasNext: true,
        hasPrevious: false,
        limit: 10,
        offset: 0,
    }
};

const mutations = {
        'SET_SEARCH_KEY'(state, key) {
            state.searchKey = key;
        },
        'SET_VALUES'(state, values){
            state.values = values;
        },
        'ADD_VALUE'(state, value) {
            if (value.exists) {
                state.values.unshift({
                    'value': value.value,
                    'key': value.key,
                    'bucket': value.bucket,
                });

                state.success = `Key "${value.key}" was successfully added`;
            }

            if (value.error) {
                state.error = value.error;
            }
        },
        'SET_ADDING_NEW_KEY'(state, isAdding){
            state.addingKey = isAdding;
        },
        'SET_NEW_KEY'(state, {key, value}){
            state.newKey.key = key;
            state.newKey.value = value;
        },
        'DELETE_VALUE'(state, {key, bucket}) {
            state.values = state.values.filter(item => {
                return item.key != key
            });
            state.success = `Key "${key}" was successfully deleted`;
        },
        'SET_ERROR'(state, value){
            state.error = value;
        },
        'SET_SUCCESS'(state, value){
            state.success = value;
        }
    }
    ;

const actions = {
    searchKey: ({commit}, {key, bucket}) => {
        commit('SET_VALUES', []);

        const url = `key/${key}?bucket=${bucket}`;
        Vue.http.get(url)
            .then(response => response.json())
            .then(data => {
                if (data) {
                    if (data.exists) {
                        commit('SET_VALUES', [
                            {
                                'value': data.value,
                                'key': data.key
                            }
                        ]);
                    }
                }
            })
    },
    searchKeyByPrefix: ({commit}, {key, bucket, type}) => {
        commit('SET_VALUES', []);

        const limit = state.pagination.limit;
        switch (type) {
            case 'next':
                state.pagination.hasPrevious = true;
                state.pagination.offset += limit;
                break;
            case 'previous':
                const nextOffset = state.pagination.offset - limit;
                state.pagination.offset = nextOffset;

                if (nextOffset < 0) {
                    state.pagination.offset = 0;
                    state.pagination.hasPrevious = false;
                }
                break;
            default:
                state.pagination.offset = 0;
        }

        const url = `prefix/key/${key}?bucket=${bucket}&limit=${limit}&offset=${state.pagination.offset}`;
        Vue.http.get(url)
            .then(response => response.json())
            .then(data => {
                if (data) {
                    if (data.length > 0) {
                        commit('SET_VALUES', data);

                    }
                }
                state.pagination.hasNext = data && data.length > 0;
            })

    },
    setAddingNewKey: ({commit}, isAdding) => {
        commit('SET_ADDING_NEW_KEY', isAdding);
    },
    addKey: ({commit}, {key, bucket, value}) => {
        const url = `key/${key}?bucket=${bucket}`;

        Vue.http.post(url, {value})
            .then(response => response.json())
            .then(data => {
                if (data) {
                    commit('ADD_VALUE', data);
                }

                commit('SET_ADDING_NEW_KEY', false);
                commit('SET_NEW_KEY', {key: "", value: ""})
            });
    },
    deleteKey: ({commit}, {key, bucket}) => {
        const url = `key/${key}?bucket=${bucket}`;

        Vue.http.delete(url)
            .then(response => response.json())
            .then(data => {
                if (data) {
                    if (data.success) {
                        commit('DELETE_VALUE', {key, bucket});
                    }
                }

                commit('SET_ADDING_NEW_KEY', false);
            });
    },
    removeAlerts: ({commit}) => {
        commit('SET_ERROR', "");
        commit('SET_SUCCESS', "");
    }
};

const getters = {
    searchKey: state => {
        return state.searchKey;
    },
    values: state => {
        return state.values;
    },
    addingNewKey: state => {
        return state.addingKey;
    },
    success: state => {
        return state.success;
    },
    error: state => {
        return state.error;
    },
    newKey: state => {
        return state.newKey;
    },
    pagination: state => {
        return state.pagination;
    }
};

export default {
    state,
    mutations,
    actions,
    getters
}
