interface SearchResult {
  subject: string;
  from: string;
  to: string;
  body: string;
}

interface SearchResponse {
  results: SearchResult[];
  total: number;
  currentPage: number;
  maxResults: number;
}