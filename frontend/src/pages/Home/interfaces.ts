import type { FormEvent } from 'react';
import type { TranslationKey } from '../../i18n/useTranslations';

export interface HomePageProps {
  draftName: string;
  isEditingName: boolean;
  onDraftNameChange: (name: string) => void;
  onSaveName: (name: string) => void;
  t: (key: TranslationKey) => string;
}

export type NameFormSubmitEvent = FormEvent<HTMLFormElement>;
