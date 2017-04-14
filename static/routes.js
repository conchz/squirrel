define([
    '/assets/js/home.js',
    '/assets/js/login.js'
], function (Home, Login) {
    return [
        {
            path: '/',
            component: Home
        },
        {
            path: '/login',
            component: Login
        }
    ]
});