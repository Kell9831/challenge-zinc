export interface Email {
    subject: string;
    from: string;
    to: string;
    body: string;
  }

 export interface ApiResponse {
    results: Email[];
    total_results: number;
    page: number;
    total_pages: number;
    results_per_page: number;
  }