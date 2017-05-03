<template>
  <div>
    <div class="vux-demo">
      <img class="logo" src="../assets/vux_logo.png" @click="greet">
      <h1>Squirrel</h1>
    </div>
    <group title="cell demo">
      <cell title="Vux" value="Cool" is-link>
        <router-link to="/home">Go to Home</router-link>
      </cell>
    </group>
  </div>
</template>

<script>
  import {Group, Cell} from 'vux'

  export default {
    components: {
      Group,
      Cell
    },
    data () {
      return {
        // note: changing this line won't causes changes
        // with hot-reload because the reloaded component
        // preserves its current state and we are modifying
        // its initial state.
        msg: 'Hello World!'
      }
    },
    methods: {
      greet: function () {
        this.$vux.toast.show({
          text: 'Loading',
          time: 1500,
          isShowMask: true
        })
      }
    }
  }
</script>

<style>
  .vux-demo {
    text-align: center;
  }

  .logo {
    width: 100px;
    height: 100px
  }
</style>
