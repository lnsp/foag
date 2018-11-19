<template>
  <div>
    <codemirror v-model="code" :options="{ theme: 'base16-dark', lineNumbers: true, mode: 'javascript' }" class="mb-3"/>
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
  name: 'build',
  data() {
    return {
      code: '// You can use JavaScript, Go, Swift and C.',
    };
  },
  methods: {
    deploy(lang) {
      this.axios.post(`${this.$endpoint}/deploy?lang=${lang}`, this.code)
        .then(() => {
          this.$router.push('/');
        });
    },
  },
};
</script>
