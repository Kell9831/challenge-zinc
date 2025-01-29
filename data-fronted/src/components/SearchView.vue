<template>
  <div class="flex">

    <!-- Tabla de resultados -->
    <div class="w-2/3">
      <input type="text" v-model="query" placeholder="Escribe algo para buscar" class="border p-2 rounded w-full mb-4"
        @input="debouncedSearch" />
      <!-- <button @click="search" class="bg-blue-500 text-white px-4 py-2 rounded mb-4">Buscar</button> -->

      <div v-if="isLoading" class="text-gray-500">Cargando...</div>
      <div v-if="error" class="text-red-500">Error: {{ error }}</div>

      <table class="table-auto border-collapse w-full">
        <thead>
          <tr>
            <th class="border px-4 py-2">Asunto</th>
            <th class="border px-4 py-2">De</th>
            <th class="border px-4 py-2">Para</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(email, index) in paddedResults" :key="index"
            :class="{ 'bg-blue-100': selectedRow === index, 'hover:bg-gray-100': true }" @click="selectRow(index)"
            class="border cursor-pointer transition">

            <td class="border px-4 py-2">{{ email?.subject || '' }}</td>
            <td class="border px-4 py-2">{{ email?.from || '' }}</td>
            <td class="border px-4 py-2">{{ email?.to || '' }}</td>
          </tr>
        </tbody>
      </table>

      <div v-if="totalPages > 1" class="mt-4 flex justify-center space-x-2">
        <button :disabled="currentPage === 1" @click="prevPage"
          class="px-4 py-2 bg-blue-500 text-white rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-blue-600">
          Anterior
        </button>

        <span>Página {{ currentPage }} de {{ totalPages }}</span>

        <button :disabled="currentPage >= totalPages" @click="nextPage"
          class="px-4 py-2 bg-blue-500 text-white rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-blue-600">
          Siguiente
        </button>

      </div>

      <div v-if="totalResults > 0" class="mt-2">
        Mostrando {{ results.length }} de {{ totalResults }} resultados.
      </div>

      <div class="w-2/3 mb-4">
        <label for="size" class="mr-2">Mostrar</label>
        <select id="size" v-model="resultsPerPage" @change="fetchData" class="border p-2 rounded">
          <option value="5">5</option>
          <option value="10">10</option>
          <option value="15">15</option>
          <option value="20">20</option>
        </select>
        <span class="ml-2">resultados por página</span>
      </div>
    </div>

    <!-- Panel de detalles -->
    <div class="w-1/3 border-l p-4 flex flex-col">
      <h2 class="text-lg font-bold mb-2">Detalle del mensaje</h2>
      <div v-if="selectedEmail" class="flex-1 overflow-y-auto max-h-[calc(100vh-200px)]">
        <p><strong>Asunto:</strong> {{ selectedEmail.subject }}</p>
        <p><strong>De:</strong> {{ selectedEmail.from }}</p>
        <p><strong>Para:</strong> {{ selectedEmail.to }}</p>
        <p><strong>Mensaje:</strong></p>
        <p class="whitespace-pre-wrap border rounded p-2 bg-gray-50" v-html="highlightBody(selectedEmail.body)"></p>
      </div>
      <div v-else>
        <p>Seleccione un mensaje para ver los detalles.</p>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { ref, computed } from 'vue';
import { useSearchStore } from '@/stores/search';

export interface Email {
  subject: string;
  from: string;
  to: string;
  body: string;
}

export default {
  setup() {
    const store = useSearchStore();

    const query = ref<string>('');
    const selectedRow = ref<number | null>(null);

    let debounceTimeout: any = null;

    const debouncedSearch = () => {
      clearTimeout(debounceTimeout);
      debounceTimeout = setTimeout(() => {
        store.query = query.value; // Actualiza el query en el store
        store.fetchResults(store.query, 1, store.resultsPerPage);
      }, 500);
    };
    // Función para manejar el cambio de tamaño
    const fetchData = () => {
      store.fetchResults(query.value, 1, store.resultsPerPage);
    };

    // Sincronizar `resultsPerPage` entre el componente y la tienda
    const resultsPerPage = computed({
      get: () => store.resultsPerPage,
      set: (value: number) => {
        store.resultsPerPage = value; 
        fetchData(); 
      },
    });


    const selectRow = (index: number) => {
      selectedRow.value = index;
    };

    const nextPage = () => {
      store.nextPage();
    };

    const prevPage = () => {
      store.prevPage();
    };

    const paddedResults = computed<Email[]>(() => {
      const padding = Math.max(0, 10 - store.results.length);
      return [...store.results, ...Array(padding).fill(null)];
    });

    const selectedEmail = computed<Email | null>(() =>
      selectedRow.value !== null ? store.results[selectedRow.value] : null
    );

    const highlightBody = (body: string) => {
      if (!query.value) return body;
      const regex = new RegExp(`(${query.value})`, 'gi');
      return body.replace(regex, '<strong>$1</strong>');
    };

    return {
      query,
      debouncedSearch,
      resultsPerPage,
      selectRow,
      nextPage,
      prevPage,
      paddedResults,
      selectedRow,
      selectedEmail,
      highlightBody,
      fetchData,
      results: computed(() => store.results),
      totalResults: computed(() => store.totalResults),
      totalPages: computed(() => store.totalPages),
      currentPage: computed(() => store.currentPage),
      isLoading: computed(() => store.isLoading),
      error: computed(() => store.error),
    };
  },
};
</script>

