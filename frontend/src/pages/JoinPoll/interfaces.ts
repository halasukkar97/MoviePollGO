import type { TranslationKey } from '../../i18n/useTranslations';

export interface JoinPollPageProps {
  savedName: string;
  t: (key: TranslationKey) => string;
}
