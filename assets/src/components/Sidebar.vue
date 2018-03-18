<template>
    <div class="col-sm-2 col-md-2 sidebar">
        <ul class="nav nav-sidebar">
            <li :class="{active : componentName == selectedEngine, hidden : componentName != selectedEngine}" v-for="(component, componentName) in engines">
                <a href="#">{{ componentName }}</a>
            </li>
        </ul>
    </div>
</template>

<script>
    import {mapActions} from 'vuex';
    import {mapGetters} from 'vuex';

    export default {
        beforeMount(){
            this.initEngine();
        },
        methods: {
            ...mapActions({
                initEngine: 'init',
            }),
        },
        computed: {
            ...mapGetters([
                'engines',
                'selectedEngine'
            ]),
        },
    }
</script>


<style scoped>
    /* Hide for mobile, show later */
    .sidebar {
        display: none;
    }

    .sidebar a:first-letter {
        text-transform: capitalize;
    }

    @media (min-width: 768px) {
        .sidebar {
            position: fixed;
            top: 51px;
            bottom: 0;
            left: 0;
            z-index: 1000;
            display: block;
            padding: 20px;
            overflow-x: hidden;
            overflow-y: auto; /* Scrollable contents if viewport is shorter than content. */
            background-color: #f5f5f5;
            border-right: 1px solid #eee;
        }
    }

    /* Sidebar navigation */
    .nav-sidebar {
        margin-right: -21px; /* 20px padding + 1px border */
        margin-bottom: 20px;
        margin-left: -20px;
    }

    .nav-sidebar > li > a {
        padding-right: 20px;
        padding-left: 20px;
    }

    .nav-sidebar > .active > a,
    .nav-sidebar > .active > a:hover,
    .nav-sidebar > .active > a:focus {
        color: #fff;
        background-color: #428bca;
    }

</style>