// API_BASE_URL points the React app at the Go backend.
// VITE_API_BASE_URL can override it for deployed or shared environments.
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080';

// request wraps fetch so every API call gets the same base URL, JSON headers, and error handling.
async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(API_BASE_URL + path, {
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
    ...options,
  });

  if (!response.ok) {
    throw new Error('Request failed with status ' + response.status);
  }

  return response.json() as Promise<T>;
}

export type CreatePollPayload = {
  name: string;
  maxVotesPerPerson: number;
  deadline: string;
};

export type CreateMoviePayload = {
  title: string;
  pollId: string;
  releaseYear: number;
  description: string;
  posterUrl: string;
};

export type Movie = {
  id: string;
  pollId: string;
  title: string;
  releaseYear: number;
  description: string;
  posterUrl: string;
};

export type ExternalMovie = {
  id: number;
  title: string;
  release_date: string;
  overview: string;
  poster_path: string;
  poster_url?: string;
  posterUrl?: string;
};

export type Vote = {
  id: string;
  pollId: string;
  userId: string;
  movieIds: string[];
};

export type Poll = {
  id: string;
  pollCode: string;
  name: string;
  maxVotesPerPerson: number;
  isClosed: boolean;
  deadline: string;
  movies: Movie[];
  votes: Vote[];
};

export type PollResults = Record<string, number>;

export const apiClient = {
  // listPolls gets every poll for the History page.
  listPolls: () => request<Poll[]>('/polls'),

  // createPoll sends the Create Poll form data and returns the new poll.
  createPoll: (payload: CreatePollPayload) =>
    request<Poll>('/polls', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),

  // getPoll opens one poll by its public poll code.
  getPoll: (pollCode: string) => request<Poll>('/polls/' + encodeURIComponent(pollCode)),

  // searchMovies asks the backend to search TMDB for movie suggestions.
  searchMovies: (query: string) =>
    request<ExternalMovie[]>('/movies/search?q=' + encodeURIComponent(query)),

  // createMovie adds one selected movie to the current poll.
  createMovie: (payload: CreateMoviePayload) =>
    request<Movie>('/movies', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),

  // getPollResults loads vote totals for one public poll code.
  getPollResults: (pollCode: string) =>
    request<PollResults>('/results?pollCode=' + encodeURIComponent(pollCode)),
};
