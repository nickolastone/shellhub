<template>
  <v-form>
    <v-row>
      <v-col class="mb-6">
        <h3 class="mb-5">
          Security
        </h3>

        <div
          class="ml-3"
        >
          <v-checkbox
            v-model="sessionRecord"
            :disabled="!hasAuthorization"
            label="Enable session record"
          />

          <p>
            Session record is a feature that allows you to check logged activity when
            connecting to a device.
          </p>
        </div>
      </v-col>
    </v-row>
  </v-form>
</template>

<script>

import hasPermission from '@/components/filter/permission';

export default {
  name: 'SettingSecurityComponent',

  filters: { hasPermission },

  props: {
    hasTenant: {
      type: Boolean,
      default: false,
    },
  },

  data() {
    return {
      action: 'enableSessionRecord',
    };
  },

  computed: {
    sessionRecord: {
      get() {
        return this.$store.getters['security/get'];
      },

      async set(value) {
        const data = {
          id: localStorage.getItem('tenant'),
          status: value,
        };
        try {
          await this.$store.dispatch('security/set', data);
          this.$store.dispatch('snackbar/showSnackbarSuccessAction', this.$success.namespaceEdit);
        } catch {
          this.$store.dispatch('snackbar/showSnackbarErrorDefault');
        }
      },
    },

    hasAuthorization() {
      const role = this.$store.getters['auth/role'];
      if (role !== '') {
        return hasPermission(
          this.$authorizer.role[role],
          this.$actions.namespace[this.action],
        );
      }

      return false;
    },
  },

  async created() {
    try {
      if (this.hasTenant) {
        await this.$store.dispatch('security/get');
      }
    } catch {
      this.$store.dispatch('snackbar/showSnackbarErrorDefault');
    }
  },
};
</script>
