import { ui, defaultLang, showDefaultLang } from './ui';


export function getLangFromUrl(url: URL) {
  const [, lang] = url.pathname.split('/');
  if (lang || '' in ui) return lang as keyof typeof ui;
  return defaultLang;
}

export function useTranslatedPath(lang: keyof typeof ui) {
  return function translatePath(path: string, l: string = lang) {
    return !showDefaultLang && l === defaultLang ? path : `/${l}${path}`
  }
}

export function useTranslations(lang: keyof typeof ui): UseTranslations {
  return function t(key: keyof typeof ui[typeof defaultLang]) {
    return ui[lang][key] || ui[defaultLang][key];
  }
}

export type UseTranslations = (key: keyof typeof ui[typeof defaultLang]) => string; 