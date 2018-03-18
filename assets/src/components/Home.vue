<template>
    <div>
        <app-header></app-header>
        <div class="container-fluid">
            <div class="row">
                <app-sidebar></app-sidebar>
                <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
                    <component v-bind:is="selectedEngineComponent">

                    </component>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    import Header from './Header.vue';
    import Sidebar from './Sidebar.vue';

    import {mapGetters} from 'vuex';

    export default {
        computed: {
            ...mapGetters([
                'selectedEngineComponent',
            ]),
        },
        components: {
            appHeader: Header,
            appSidebar: Sidebar,
        }
    }
</script>

<style scoped>
    .main {
        padding: 20px;
        padding-top: 70px;
    }
    @media (min-width: 768px) {
        .main {
            padding-right: 40px;
            padding-left: 40px;
        }
    }
    .main .page-header {
        margin-top: 0;
    }
</style>