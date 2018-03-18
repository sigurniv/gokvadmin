<template>
    <div>
        <app-header></app-header>
        <div class="container login-from">
            <div class="row">
                <div class="col-md-4 col-md-offset-4">
                    <div class="login-panel panel panel-default">
                        <div class="panel-heading">
                            <h3 class="panel-title">Sign In</h3>
                        </div>
                        <div class="panel-body">
                            <form role="form">
                                <fieldset>
                                    <div class="form-group">
                                        <input v-model="login"
                                               class="form-control"
                                               placeholder="Login"
                                               type="text"
                                               autofocus>
                                    </div>
                                    <div class="form-group">
                                        <input v-model="password"
                                               class="form-control"
                                               placeholder="Password"
                                               type="password">
                                    </div>
                                    <a class="btn btn-lg btn-success btn-block" @click="auth">Login</a>

                                    <div class="alert alert-danger form-error" role="alert" v-if="loginError.length > 0">
                                        <a class="alert-link">{{ loginError }}</a>
                                    </div>
                                </fieldset>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    import Header from './Header.vue';
    import Sidebar from './Sidebar.vue';
    import {mapActions} from 'vuex';
    import {mapGetters} from 'vuex';

    export default {
        data(){
            return {
                login: "",
                password: "",
            }
        },
        computed: {
            ...mapGetters([
                'loginError'
            ]),
        },
        components: {
            appHeader: Header,
            appSidebar: Sidebar,
        },
        methods: {
            ...mapActions({
                appLogin: 'login',
            }),
            auth(){
                this.appLogin({"login": this.login, "password": this.password, "router": this.$router});
            }
        },
    }
</script>

<style scoped>
    .login-from {
        padding-top: 120px;
    }
    .form-error{
        margin-top: 10px;
        text-align: center;
    }
</style>