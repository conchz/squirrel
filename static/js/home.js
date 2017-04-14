define(['vue'], function (Vue) {
    return Vue.extend({
        template: `
            <div>
                Hi, Welcome to use Squirrel!
                <ul>
                    <li v-for='link in links'>
                        <a :href='link.url' target="_blank">{{link.title}}</a>
                    </li>
                </ul>
            </div>
        `,
        data: function () {
            return {
                links: [
                    {
                        'url': 'https://golang.org/',
                        'title': 'Golang'
                    },
                    {
                        'url': 'https://echo.labstack.com/',
                        'title': 'Echo'
                    },
                    {
                        'url': 'https://github.com/vuejs/vue',
                        'title': 'Vue.js'
                    }
                ]
            }
        }
    });
});