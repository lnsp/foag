<template>
  <div class="home">
    <b-table :fields="fields" :items="items" sort-by="Date" :sort-desc="true">
      <template slot="ID" slot-scope="data">
        <router-link :to="{ name: 'details', params: { id: data.value } }">
          {{ data.value | formatId }}
        </router-link>
      </template>
      <template slot="Ready" slot-scope="data">
        <span :class="{ 'text-success': data.value, 'text-danger': !data.value }">&#x25a0;</span>
      </template>
      <template slot="Date" slot-scope="data">
        {{ data.value | moment('from') }}
      </template>
      <template slot="Actions" slot-scope="data">
        <router-link v-if="data.item.Ready"
          :to="{ name: 'addRoute', params: { id: data.item.ID } }">Bind</router-link>
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
      refreshTimer: null,
      items: [],
      fields: [
        { key: 'Ready', label: '' },
        { key: 'ID' },
        { key: 'Date' },
        { key: 'Language' },
        { key: 'Actions', label: 'Actions' },
      ],
    };
  },
  methods: {
    refresh() {
      this.axios.get(`${this.$endpoint}/list`).then((resp) => { this.items = resp.data; });
    },
  },
  created() {
    this.refresh();
    this.refreshTimer = window.setInterval(this.refresh, 10000);
  },
  beforeDestroy() {
    window.clearInterval(this.refreshTimer);
  },
};
</script>
