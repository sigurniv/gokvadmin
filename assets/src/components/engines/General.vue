<template>
    <div>
        <div>
            <div class="panel panel-default">
                <div class="panel-heading">Search</div>

                <div class="form-buttons">
                    <div class="row">
                        <div class="col-lg-4">
                            <div class="input-group">
                                <label>Bucket / Namespace:</label>
                                <input
                                        class="form-control"
                                        placeholder="Bucket"
                                        v-model="bucket"
                                >
                            </div>
                        </div>
                    </div>

                    <div class="row">
                        <div class="col-lg-4">
                            <div class="input-group">
                                <label>Key:</label>
                                <input
                                        class="form-control"
                                        placeholder="Search keys.."
                                        v-model="key"
                                >
                            </div>
                        </div>
                    </div>

                    <div class="row">
                        <div class="col-lg-4">
                            <button
                                    class="btn btn-success"
                                    type="button"
                                    @click="searchKey"
                            >Search
                            </button>
                            <button
                                    class="btn btn-success"
                                    type="button"
                                    @click="searchKeyByPrefix"
                            >Search by prefix
                            </button>
                        </div>

                    </div>
                </div>
            </div>
        </div>

        <div class="alert alert-success" role="alert" v-if="success.length > 0">
            <a class="alert-link">{{ success }}</a>
        </div>

        <div class="alert alert-danger" role="alert" v-if="error.length > 0">
            <a class="alert-link">{{ error }}</a>
        </div>

        <div class="panel panel-default" style="margin-top: 20px;">
            <div class="panel-heading clearfix">
                <span>{{searchTitle}}</span>
                <span class="pull-right">
                    <button
                            class="btn btn-success"
                            type="button"
                            @click="setAddingNewKey(true)"
                    >Set key</button>
                </span>
            </div>
        </div>

        <table class="table table-bordered table-hover">
            <thead>
            <tr>
                <th>Key</th>
                <th>Value</th>
                <th></th>
            </tr>
            </thead>
            <tbody>
            <tr v-if="addingNewKey">
                <td>
                    <textarea
                            class="form-control"
                            placeholder="Key"
                            v-model="newKey.key"
                    ></textarea>
                </td>
                <td>
                    <textarea
                            class="form-control"
                            placeholder="Value"
                            v-model="newKey.value"
                    ></textarea>
                </td>
                <td>
                    <button
                            class="btn btn-success"
                            type="button"
                            @click="saveKey"
                    >Save
                    </button>

                    <button
                            class="btn btn-warning"
                            type="button"
                            @click="cancelAddKey"
                    >Cancel
                    </button>
                </td>
            </tr>

            <tr v-for="(value, index) in values" v-on:dblclick="editValue(value)">
                <td>{{ value.key }}</td>
                <td>
                    <div>{{ value.value }}</div>
                </td>
                <td>
                    <button
                            class="btn btn-danger"
                            type="button"
                            @click="deleteKey(value)"
                    >Delete
                    </button>
                </td>
            </tr>
            </tbody>
        </table>

        <ul class="pagination" v-if="searchType == 'prefix'">
            <li
                    class="paginate_button previous"
                    :class="{disabled: !pagination.hasPrevious}"
                    @click="searchKeyByPrefix('previous')"

            ><a href="#">Previous</a></li>
            <li
                    class="paginate_button next"
                    :class="{disabled: !pagination.hasNext}"
                    @click="searchKeyByPrefix('next')"
            ><a>Next</a></li>
        </ul>
    </div>
</template>

<script>
    import {mapActions} from 'vuex';
    import {mapGetters} from 'vuex';


    export default {
        data() {
            return {
                key: '',
                bucket: 'MyBucket',
                searchTitle: 'Search results',
                searchType: "",
            }
        },
        computed: {
            ...mapGetters([
                'values',
                'addingNewKey',
                'success',
                'error',
                'newKey',
                'pagination'
            ]),
        },
        methods: {
            ...mapActions({
                fetchKey: 'searchKey',
                fetchKeyByPrefix: 'searchKeyByPrefix',
                setAddingNewKey: 'setAddingNewKey',
                storeKey: 'addKey',
                removeKey: 'deleteKey',
                clearAlerts: 'removeAlerts',
            }),
            editValue(value){

            },
            searchKey(){
                if (this.key.length > 0) {
                    this.searchType = "key";
                    this.fetchKey({key: this.key, bucket: this.bucket});
                    this.searchTitle = `Search results for key : "${this.key}"`;
                    this.clearAlerts();
                }
            },
            searchKeyByPrefix(type) {
                this.searchType = "prefix";

                this.fetchKeyByPrefix({
                    key: this.key,
                    bucket: this.bucket,
                    type: type,
                });
                this.searchTitle = `Search results for prefix : "${this.key}"`;
                this.clearAlerts();
            },
            cancelAddKey(){
                this.setAddingNewKey(false);
                this.newKey.key = "";
                this.newKey.value = "";
            },
            saveKey(){
                if (this.newKey.key.length > 0) {
                    this.storeKey({
                        key: this.newKey.key,
                        value: this.newKey.value,
                        bucket: this.bucket,
                    })
                }
            },
            deleteKey(value){
                this.removeKey({...value, bucket: this.bucket});
            }
        }
    }
</script>

<style scoped>
    .form-buttons {
        padding: 20px;
    }

    .form-buttons .row {
        margin-bottom: 10px;
    }

    .panel {
        margin-bottom: 0px;
        border-bottom: 0px;
    }

    .panel-heading {
        font-size: 15px;
        font-weight: bold;
    }

    .search-btn {
        border-radius: 0px;
        border-left: 0px;
        border-right: 0px;
    }

    .alert {
        margin: 30px;
    }

    .paginate_button {
        cursor: pointer;
    }

    textarea {
        resize: vertical;
    }
</style>