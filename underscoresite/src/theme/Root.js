import React, {useEffect} from 'react';
import OriginalRoot from '@theme-original/Root';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';

export default function Root(props){
  const {i18n} = useDocusaurusContext();

  useEffect(() => {
    if (typeof window === 'undefined') return;
    try {
      const skip = new URLSearchParams(window.location.search).get('no_locale_redirect');
      if (skip) return;
      const already = localStorage.getItem('geol_locale_redirect_done');
      if (already) return;
      const path = window.location.pathname || '/';
      // Only redirect the site root to the preferred locale once.
      if (path === '/' || path === '') {
        const nav = navigator.language || navigator.userLanguage || '';
        if (nav && nav.toLowerCase().startsWith('fr')) {
          const target = '/fr/';
          localStorage.setItem('geol_locale_redirect_done', '1');
          window.location.replace(target);
        }
      }
    } catch (e) {
      // silent
    }
  }, [i18n]);

  return (
    <>
      <OriginalRoot {...props} />
    </>
  );
}
