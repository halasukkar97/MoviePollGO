import { useState } from 'react';
import type { ChangeEvent, FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import type { JoinPollPageProps } from './interfaces';
import './JoinPoll.scss';

// JoinPollPage collects the poll ID needed to enter an existing poll.
export function JoinPollPage({ savedName, t }: JoinPollPageProps) {
  const navigate = useNavigate();
  // pollCode is the only value needed here; the voter name comes from localStorage later.
  const [pollCode, setPollCode] = useState('');

  // handlePollCodeChange keeps the poll code input connected to state.
  function handlePollCodeChange(event: ChangeEvent<HTMLInputElement>) {
    setPollCode(event.target.value);
  }

  // handleSubmit opens the public poll page by pollCode.
  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const trimmedPollCode = pollCode.trim();
    if (!trimmedPollCode) {
      return;
    }

    navigate('/polls/' + trimmedPollCode);
  }

  return (
    <section className="page join-poll-page">
      <h1>{t('join.title')}</h1>
      {savedName ? <p className="saved-name-note">{t('join.votingAs')} {savedName}</p> : null}

      <form className="form" onSubmit={handleSubmit}>
        <label>
          {t('join.pollCode')}
          <input
            name="pollCode"
            type="text"
            placeholder={t('join.pollCodePlaceholder')}
            value={pollCode}
            onChange={handlePollCodeChange}
          />
        </label>

        <button type="submit">{t('join.button')}</button>
      </form>
    </section>
  );
}
