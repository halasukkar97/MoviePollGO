import { useState } from 'react';
import type { ChangeEvent, FormEvent } from 'react';
import type { TranslationKey } from '../../i18n/useTranslations';
import type { ResultsLookupFormValues } from './interfaces';
import './PollResults.scss';

const initialFormValues: ResultsLookupFormValues = {
  pollCode: '',
};

interface PollResultsPageProps {
  t: (key: TranslationKey) => string;
}

// PollResultsPage is responsible for looking up and showing a poll's results.
export function PollResultsPage({ t }: PollResultsPageProps) {
  // State stores the poll code that will be sent to the results API later.
  const [formValues, setFormValues] = useState<ResultsLookupFormValues>(initialFormValues);

  // handleChange keeps the poll code input connected to local state.
  function handleChange(event: ChangeEvent<HTMLInputElement>) {
    setFormValues({ pollCode: event.target.value });
  }

  // handleSubmit is where the poll-results API call will be added later.
  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
  }

  return (
    <section className="page poll-results-page">
      <h1>{t('results.title')}</h1>
      <form className="form" onSubmit={handleSubmit}>
        <label>
          {t('results.pollCode')}
          <input
            name="pollCode"
            type="text"
            placeholder={t('results.pollCodePlaceholder')}
            value={formValues.pollCode}
            onChange={handleChange}
          />
        </label>

        <button type="submit">{t('results.button')}</button>
      </form>
    </section>
  );
}
