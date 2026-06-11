import type { Poll } from '../../api/client';
import type { TranslationKey } from '../../i18n/useTranslations';

export interface HistoryPageProps {
  t: (key: TranslationKey) => string;
}

export interface HistoryState {
  polls: Poll[];
  isLoading: boolean;
  errorMessage: string;
}
