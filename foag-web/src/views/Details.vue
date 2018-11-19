<template>
    <b-card v-if="deployment" class="mb-3" title="Deployment" :subTitle="deployment.ID"
      :border-variant="deployment.Ready ? 'success' : 'danger'">
      <b-list-group class="mt-3">
        <b-list-group-item>
          <div class="subheader">URL</div>
          <a :href="$endpoint+deployment.URL">{{ deployment.URL }}</a>
        </b-list-group-item>
        <b-list-group-item>
          <div class="subheader">Details</div>
          Deployed {{ deployment.Date | moment('from') }}
        </b-list-group-item>
        <b-list-group-item v-if="logs">
          <div class="subheader">Logs</div>
          <pre>{{ logs }}</pre>
        </b-list-group-item>
      </b-list-group>
    </b-card>
</template>

<script>
export default {
  name: 'details',
  data() {
    return {
      deployment: null,
      logs: null,
      refreshTimeout: null,
    };
  },
  methods: {
    refresh() {
      this.axios.get(`${this.$endpoint}/describe/${this.$route.params.id}`).then((resp) => {
        this.deployment = resp.data;
        if (this.deployment && !this.deployment.Ready) {
          window.setTimeout(this.refresh, 1000);
        }
      });
      this.axios.get(`${this.$endpoint}/logs/${this.$route.params.id}`).then((resp) => {
        this.logs = resp.data;
      });
    },
  },
  mounted() {
    this.refresh();
  },
  beforeDestroy() {
    window.clearTimeout(this.refreshTimeout);
  },
};
</script>

<style scoped>
.subheader {
  font-size: 0.8em;
  font-weight: 700;
  text-transform: uppercase;
  color: #636363;
  margin-bottom: 0.5em;
}
</style>
