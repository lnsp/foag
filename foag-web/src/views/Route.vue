<template>
  <div>
    <b-card class="mb-3" title="Add an alias">
    <b-form inline>
      <b-input class="mr-2 mb-2 col-12 col-sm-auto flex-sm-fill"
        v-model="id" placeholder="Deployment ID" />
      <b-input class="mr-2 mb-2 col-12 col-sm-auto flex-sm-fill"
        v-model="alias" placeholder="Alias" />
      <b-button class="mb-2 col-12 col-sm-auto"
        variant="primary" @click="update">Add Alias</b-button>
    </b-form>
    </b-card>
    <b-table :fields="fields" :items="items">
      <template slot="For" slot-scope="data">
        <router-link :to="{ name: 'details', params: { id: data.value } }">
          {{ data.value | formatId }}
        </router-link>
      </template>
      <template slot="URL" slot-scope="data">
        <a :href="$endpoint+data.value">{{ data.value }}</a>
      </template>
    </b-table>
  </div>
</template>

<script>
export default {
  name: 'route',
  data() {
    return {
      id: '',
      alias: '',
      items: [],
      fields: [
        { key: 'Name' },
        { key: 'For' },
        { key: 'URL' },
      ],
    };
  },
  methods: {
    refresh() {
      this.axios.get(`${this.$endpoint}/listAlias`).then((resp) => { this.items = resp.data; });
    },
    update() {
      this.axios.post(`${this.$endpoint}/bind/${this.id}?to=${this.alias}`).then(this.refresh);
    },
  },
  mounted() {
    this.id = this.$route.params.id;
    this.refresh();
  },
};

</script>
