<template>
  <div class="home">
    <b-table :fields="fields" :items="items">
      <template slot="ID" slot-scope="data">
        {{ data.value | formatId }}
      </template>
    </b-table>
  </div>
</template>

<script>
// @ is an alias to /src
export default {
  name: 'home',
  data() {
    return {
      items: [],
      fields: [ { key: 'ID' }, { key: 'Date' }, { key: 'Language' }, { key: 'Ready' } ]
    };
  },
  filters: {
    formatId: value => value.slice(0, 16),
  },
  mounted() {
    this.axios.get('http://localhost:8080/list').then((resp) => { this.items = resp.data; });
  },
};
</script>
