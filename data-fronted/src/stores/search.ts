import { defineStore } from 'pinia';

interface Email {
  subject: string;
  from: string;
  to: string;
  body: string;
}

interface ApiResponse {
  results: Email[];
  total_results: number;
  page: number;
  total_pages: number;
  results_per_page: number;
}

export const useSearchStore = defineStore('search', {
  state: () => ({
    query: '',
    results: [] as Email[], 
    totalResults: 0, 
    currentPage: 1, 
    totalPages: 0, 
    resultsPerPage: 5, 
    isLoading: false, 
    error: null as string | null, 
  }),
  
  actions: {
    // Acción para realizar la búsqueda
    async fetchResults(query: string, page: number = 1, size: number = 5) {
      if (this.isLoading) return;
      this.isLoading = true;
      this.error = null;

      try {
        // Realizando la solicitud a la API
        console.log(`Fetching results: Query=${query}, Page=${page}, Size=${size}`);
        const response = await fetch(`http://localhost:8080/search?page=${page}&size=${size}`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            term: query,
            // field: '_all',
          }),
        });

        if (!response.ok) {
          throw new Error('Error al buscar los datos');
        }

        // Analizando la respuesta
        const data: ApiResponse = await response.json();
        
        // Actualizando el estado con los resultados de la API
        this.results = data.results;
        this.totalResults = data.total_results;
        this.currentPage = data.page;
        this.totalPages = Math.ceil(data.total_results / size);
        this.resultsPerPage = data.results_per_page;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Error desconocido';
      } finally {
        this.isLoading = false;
      }
    },

    // Acción para avanzar a la siguiente página
  nextPage() {
    if (this.currentPage < this.totalPages) {
    this.fetchResults(this.query, this.currentPage + 1, this.resultsPerPage);
   }
  },

// Acción para retroceder a la página anterior
  prevPage() {
    if (this.currentPage > 1) {
      this.fetchResults(this.query, this.currentPage - 1, this.resultsPerPage);
   }
  },
  },
  });
