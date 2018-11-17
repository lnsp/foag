<template>
  <div>
    <codemirror v-model="code" :options="{ theme: 'base16-dark' }"/>
      <b-dropdown text="Deploy" variant="primary">
        <b-dropdown-item @click="deploy('go')">as Go</b-dropdown-item>
        <b-dropdown-item @click="deploy('swift')">as Swift</b-dropdown-item>
        <b-dropdown-item @click="deploy('js')">as NodeJS</b-dropdown-item>
        <b-dropdown-item @click="deploy('c')">as C</b-dropdown-item>
      </b-dropdown>
  </div>
</template>

<script>
export default {
  name: 'deployFunction',
  data() {
    return {
      code: '// just some examples here',
    };
  },
  methods: {
    deploy(lang) {
      this.axios.post(`http://localhost:8080/deploy?lang=${lang}`, this.code)
        .then(() => {
          this.$router.push('/');
        });
    },
  },
};
</script>
