import { ui, defaultLang, showDefaultLang, languages } from './ui';

export function getJustPath(url: URL): string {
  const pathname = url.pathname;
  for (const language in languages) {
    if (pathname.includes(language)) return pathname.replace(`/${language}`, "");
  }
  return pathname;
}

export function useTranslatedPath(lang: keyof typeof ui) {
  return function translatePath(path: string, l: string = lang) {
    return !showDefaultLang && l === defaultLang ? path : `/${l}${path}`
  }
}

export function useTranslations(lang: Language): UseTranslations {
  return function t(key: keyof typeof ui[typeof defaultLang]) {
    return ui[lang][key] || ui[defaultLang][key];
  }
}

export type Language = keyof typeof ui

export type UseTranslations = (key: keyof typeof ui[typeof defaultLang]) => string; 