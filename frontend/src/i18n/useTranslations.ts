import { useCallback, useEffect, useState } from 'react';
import translations from '../data/translations.json';

const LANGUAGE_STORAGE_KEY = 'votify:language';

export type Language = keyof typeof translations;
export type TranslationKey = keyof typeof translations.en;

const defaultLanguage: Language = 'en';

function readStoredLanguage(): Language {
  const storedLanguage = localStorage.getItem(LANGUAGE_STORAGE_KEY);
  return storedLanguage === 'de' || storedLanguage === 'en' ? storedLanguage : defaultLanguage;
}

export function useTranslations() {
  // language is saved locally so the same language is used next time the app opens.
  const [language, setLanguageState] = useState<Language>(readStoredLanguage);

  // setLanguage updates React state and localStorage together.
  const setLanguage = useCallback((nextLanguage: Language) => {
    localStorage.setItem(LANGUAGE_STORAGE_KEY, nextLanguage);
    setLanguageState(nextLanguage);
  }, []);

  // t returns translated UI text for the currently selected language.
  const t = useCallback((key: TranslationKey) => {
    return translations[language][key] ?? translations.en[key];
  }, [language]);

  useEffect(() => {
    document.documentElement.lang = language;
  }, [language]);

  return { language, languages: Object.keys(translations) as Language[], setLanguage, t };
}
