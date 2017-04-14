require.config({
    baseUrl: '',
    paths: {
        'vue': 'https://cdn.bootcss.com/vue/2.2.6/vue.min',
        'vue_resource': 'https://cdn.bootcss.com/vue-resource/1.2.1/vue-resource.min',
        'vue_router': 'https://cdn.bootcss.com/vue-router/2.3.1/vue-router.min',
        'routes': '/assets/routes'
    },
    shim: {
        vue: {
            exports: 'Vue'
        },
        vue_resource: {
            exports: 'VueResource'
        },
        vue_router: {
            exports: 'VueRouter'
        }
    }
});

require([
    'vue',
    'vue_resource',
    'vue_router',
    'routes'
], function (Vue, VueResource, VueRouter, AppRoutes) {
    Vue.use(VueResource);
    Vue.use(VueRouter);

    let router = new VueRouter({
        mode: 'history',
        routes: AppRoutes
    });

    new Vue({
        el: '#app',
        router: router
    });
});