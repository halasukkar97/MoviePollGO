import { useEffect, useMemo, useState } from 'react';
import type { ChangeEvent, FormEvent } from 'react';
import { useLocation, useParams } from 'react-router-dom';
import { apiClient } from '../../api/client';
import type { ExternalMovie } from '../../api/client';
import type {
  LoadedPollState,
  MovieDraftValues,
  MovieSearchState,
  PollPageProps,
  PollRouteState,
  ToastState,
} from './interfaces';
import './Poll.scss';

const initialMovieDraft: MovieDraftValues = {
  title: '',
};

const initialMovieSearch: MovieSearchState = {
  suggestions: [],
  selectedMovie: null,
  isSearching: false,
  searchError: '',
  hasSearched: false,
};

const toastDurationMs = 30000;

function hasVotingEnded(deadline: string, isClosed: boolean) {
  return isClosed || (deadline ? new Date(deadline).getTime() < Date.now() : false);
}

function formatDate(deadline: string) {
  if (!deadline) {
    return '';
  }

  return new Intl.DateTimeFormat(undefined, { dateStyle: 'medium' }).format(new Date(deadline));
}

function getReleaseYear(movie: ExternalMovie) {
  return movie.release_date ? new Date(movie.release_date).getFullYear() : 0;
}

function getExternalPosterURL(movie: ExternalMovie) {
  return movie.poster_url ?? movie.posterUrl ?? '';
}

// PollPage shows one public poll by pollCode and will become the voting workspace.
export function PollPage({ t }: PollPageProps) {
  const { pollCode = '' } = useParams();
  const location = useLocation();
  const routeState = location.state as PollRouteState | null;
  // pollState keeps the loaded poll and request status together.
  const [pollState, setPollState] = useState<LoadedPollState>({
    poll: null,
    isLoading: true,
    errorMessage: '',
  });
  const [movieDraft, setMovieDraft] = useState<MovieDraftValues>(initialMovieDraft);
  const [movieSearch, setMovieSearch] = useState<MovieSearchState>(initialMovieSearch);
  const [isAddingMovie, setIsAddingMovie] = useState(false);
  const [toast, setToast] = useState<ToastState | null>(
    routeState?.createdPollCode
      ? {
          id: Date.now(),
          type: 'success',
          message: t('poll.created'),
          detail: routeState.createdPollCode,
        }
      : null,
  );

  const votingEnded = useMemo(
    () => hasVotingEnded(pollState.poll?.deadline ?? '', pollState.poll?.isClosed ?? false),
    [pollState.poll],
  );

  // Toasts stay visible long enough to copy details, then close themselves.
  useEffect(() => {
    if (!toast) {
      return;
    }

    const timeoutID = window.setTimeout(() => setToast(null), toastDurationMs);
    return () => window.clearTimeout(timeoutID);
  }, [toast]);

  // loadPoll fetches the public poll by pollCode from the route.
  useEffect(() => {
    async function loadPoll() {
      setPollState({ poll: null, isLoading: true, errorMessage: '' });

      try {
        const poll = await apiClient.getPoll(pollCode);
        setPollState({ poll, isLoading: false, errorMessage: '' });
      } catch (error) {
        setPollState({
          poll: null,
          isLoading: false,
          errorMessage: error instanceof Error ? error.message : t('poll.notFound'),
        });
      }
    }

    if (pollCode) {
      loadPoll();
    }
  }, [pollCode, t]);

  // Search TMDB shortly after the user stops typing.
  useEffect(() => {
    const query = movieDraft.title.trim();

    if (query.length < 2 || movieSearch.selectedMovie) {
      setMovieSearch((currentSearch) => ({
        ...currentSearch,
        suggestions: [],
        isSearching: false,
        searchError: '',
        hasSearched: false,
      }));
      return;
    }

    const timeoutID = window.setTimeout(async () => {
      setMovieSearch((currentSearch) => ({
        ...currentSearch,
        isSearching: true,
        searchError: '',
        hasSearched: true,
      }));

      try {
        const suggestions = await apiClient.searchMovies(query);
        setMovieSearch((currentSearch) => ({
          ...currentSearch,
          suggestions,
          isSearching: false,
          searchError: '',
          hasSearched: true,
        }));
      } catch {
        setMovieSearch((currentSearch) => ({
          ...currentSearch,
          suggestions: [],
          isSearching: false,
          searchError: t('poll.searchError'),
          hasSearched: true,
        }));
      }
    }, 350);

    return () => window.clearTimeout(timeoutID);
  }, [movieDraft.title, movieSearch.selectedMovie, t]);

  async function refreshPoll() {
    const poll = await apiClient.getPoll(pollCode);
    setPollState({ poll, isLoading: false, errorMessage: '' });
  }

  function showToast(nextToast: Omit<ToastState, 'id'>) {
    setToast({ ...nextToast, id: Date.now() });
  }

  // handleMovieDraftChange keeps the add-movie starter form connected to state.
  function handleMovieDraftChange(event: ChangeEvent<HTMLInputElement>) {
    const { value } = event.target;

    setMovieDraft({ title: value });
    setMovieSearch((currentSearch) => ({
      ...currentSearch,
      selectedMovie: null,
    }));
  }

  function handleSelectMovie(movie: ExternalMovie) {
    setMovieDraft({ title: movie.title });
    setMovieSearch({
      suggestions: [],
      selectedMovie: movie,
      isSearching: false,
      searchError: '',
      hasSearched: true,
    });
  }

  // handleAddMovie sends the selected TMDB movie to the backend create movie endpoint.
  async function handleAddMovie(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!pollState.poll) {
      return;
    }

    if (!movieSearch.selectedMovie) {
      showToast({ type: 'error', message: t('poll.chooseMovieFirst') });
      return;
    }

    setIsAddingMovie(true);

    try {
      await apiClient.createMovie({
        title: movieSearch.selectedMovie.title,
        pollId: pollState.poll.id,
        releaseYear: getReleaseYear(movieSearch.selectedMovie),
        description: movieSearch.selectedMovie.overview,
        posterUrl: getExternalPosterURL(movieSearch.selectedMovie),
      });

      setMovieDraft(initialMovieDraft);
      setMovieSearch(initialMovieSearch);
      await refreshPoll();
      showToast({ type: 'success', message: t('poll.addMovieSuccess') });
    } catch (error) {
      showToast({
        type: 'error',
        message: error instanceof Error ? error.message : t('poll.addMovieError'),
      });
    } finally {
      setIsAddingMovie(false);
    }
  }

  if (pollState.isLoading) {
    return <section className="page poll-page"><p>{t('poll.loading')}</p></section>;
  }

  if (!pollState.poll) {
    return <section className="page poll-page"><p>{pollState.errorMessage || t('poll.notFound')}</p></section>;
  }

  const selectedMovie = movieSearch.selectedMovie;
  const shouldShowNoMoviesFound =
    movieSearch.hasSearched &&
    !movieSearch.isSearching &&
    !movieSearch.searchError &&
    movieDraft.title.trim().length >= 2 &&
    movieSearch.suggestions.length === 0 &&
    !selectedMovie;

  return (
    <section className="page poll-page">
      {toast ? (
        <div className={'toast toast--' + toast.type} role={toast.type === 'error' ? 'alert' : 'status'}>
          <button className="toast-close" type="button" onClick={() => setToast(null)} aria-label={t('toast.close')}>
            x
          </button>
          <strong>{toast.message}</strong>
          {toast.detail ? <p>{t('poll.code')}: <code>{toast.detail}</code></p> : null}
        </div>
      ) : null}

      <div className="poll-heading">
        <div className="poll-title-block">
          <h1>{pollState.poll.name}</h1>
        </div>
        <div className="poll-meta-card">
          <div className="poll-meta-item">
            <span>{t('poll.code')}</span>
            <strong><code>{pollState.poll.pollCode}</code></strong>
          </div>
          <div className="poll-meta-divider" />
          <div className="poll-meta-item">
            <span>{t('poll.endVotingOn')}</span>
            <strong>{formatDate(pollState.poll.deadline)}</strong>
          </div>
        </div>
      </div>

      {votingEnded ? (
        <div className="feedback expired-message" role="status">
          {t('poll.votingEnded')}
        </div>
      ) : null}

      <section className="poll-workspace">
        <form className="form movie-search-form" onSubmit={handleAddMovie}>
          <h2>{t('poll.addMovies')}</h2>
          <label>
            {t('poll.movieTitle')}
            <input
              name="title"
              type="text"
              placeholder={t('poll.movieTitlePlaceholder')}
              value={movieDraft.title}
              onChange={handleMovieDraftChange}
              disabled={votingEnded}
              autoComplete="off"
            />
          </label>

          <div className="movie-search-status" aria-live="polite">
            {movieDraft.title.trim().length < 2 && !selectedMovie ? t('poll.searchHint') : null}
            {movieSearch.isSearching ? t('poll.searchLoading') : null}
            {movieSearch.searchError ? movieSearch.searchError : null}
            {shouldShowNoMoviesFound ? t('poll.noMoviesFound') : null}
          </div>

          {movieSearch.suggestions.length > 0 ? (
            <div className="suggestions-list">
              {movieSearch.suggestions.map((movie) => (
                <button key={movie.id} type="button" onClick={() => handleSelectMovie(movie)}>
                  <span>{movie.title}</span>
                  <span>{getReleaseYear(movie) || ''}</span>
                </button>
              ))}
            </div>
          ) : null}

          {selectedMovie ? (
            <div className="selected-movie-preview">
              {getExternalPosterURL(selectedMovie) ? (
                <img src={getExternalPosterURL(selectedMovie)} alt={t('poll.posterAlt')} />
              ) : null}
              <div>
                <strong>{t('poll.selectedMovie')}</strong>
                <h3>{selectedMovie.title}</h3>
                <p>{getReleaseYear(selectedMovie) || ''}</p>
                {selectedMovie.overview ? <p>{selectedMovie.overview}</p> : null}
              </div>
            </div>
          ) : null}

          <button type="submit" disabled={votingEnded || isAddingMovie}>
            {isAddingMovie ? t('poll.addingMovie') : t('poll.addMovieButton')}
          </button>
        </form>

        <section className="movie-grid-section">
          <h2>{t('poll.moviesInPoll')}</h2>
          <div className="movie-grid">
            {pollState.poll.movies.length > 0 ? (
              pollState.poll.movies.map((movie) => {
                const poster = movie.posterUrl;
                const releaseYear = movie.releaseYear;
                const description = movie.description;

                return (
                  <article className="movie-card" key={movie.id}>
                    <div className="movie-card-poster">
                      {poster ? <img src={poster} alt={t('poll.posterAlt')} /> : <span>{t('poll.noPoster')}</span>}
                    </div>
                    <div className="movie-card-body">
                      <div className="movie-card-heading">
                        <h3>{movie.title}</h3>
                        {releaseYear ? <span>{releaseYear}</span> : null}
                      </div>
                      {description ? <p>{description}</p> : null}
                      {!votingEnded ? <button type="button">{t('poll.voteButton')}</button> : null}
                    </div>
                  </article>
                );
              })
            ) : (
              <p>{t('poll.noMovies')}</p>
            )}
          </div>
        </section>
      </section>
    </section>
  );
}
