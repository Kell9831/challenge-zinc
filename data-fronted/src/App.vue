<template>
  <div id="app">
    <SearchView /> 
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import SearchView from './components/SearchView.vue'; 

export default defineComponent({
  name: 'App',
  components: {
    SearchView,
  },
});
</script>

