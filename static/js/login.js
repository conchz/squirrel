define(['vue'], function (Vue) {
    return Vue.extend({
        template: `
            <div class="col-sm-4 col-sm-offset-4">
                <h2>Log In</h2>
                <p>Log in to your account to get some great quotes.</p>
                <div class="alert alert-danger" v-if="error">
                    <p>{{ error }}</p>
                </div>
                <div class="form-group">
                    <input type="text" class="form-control"
                           placeholder="Enter your username" v-model="credentials.username">
                </div>
                <div class="form-group">
                    <input type="password" class="form-control"
                           placeholder="Enter your password" v-model="credentials.password">
                </div>
                <button class="btn btn-primary" @click="submit()">Access</button>
            </div>
        `,
        data() {
            return {
                credentials: {
                    username: '',
                    password: ''
                },
                error: ''
            }
        },
        methods: {
            submit() {
                let credentials = {
                    username: this.credentials.username,
                    password: this.credentials.password
                };
                console.log(credentials);
            }
        }
    });
});