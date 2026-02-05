import React from 'react';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import useBaseUrl from '@docusaurus/useBaseUrl';

export default function LocaleDropdown({ mobile }) {
  const { i18n } = useDocusaurusContext();
  const { currentLocale, locales } = i18n || { currentLocale: 'en', locales: ['en', 'fr'] };
  const baseUrl = useBaseUrl('/');

  const labelFor = (loc) => (loc === 'fr' ? 'Français' : loc === 'en' ? 'English' : loc);

  const icon = (
    <svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true" className="iconLanguage_nlXk"><path fill="currentColor" d="M12.87 15.07l-2.54-2.51.03-.03c1.74-1.94 2.98-4.17 3.71-6.53H17V4h-7V2H8v2H1v1.99h11.17C11.5 7.92 10.44 9.75 9 11.35 8.07 10.32 7.3 9.19 6.69 8h-2c.73 1.63 1.73 3.17 2.98 4.56l-5.09 5.02L4 19l5-5 3.11 3.11.76-2.04zM18.5 10h-2L12 22h2l1.12-3h4.75L21 22h2l-4.5-12zm-2.62 7l1.62-4.33L19.12 17h-3.24z"/></svg>
  );

  const triggerAria = currentLocale === 'fr' ? 'Sélecteur de langue, langue actuelle: Français' : `Language selector, current: ${labelFor(currentLocale)}`;

  return (
    <div className={`navbar__item dropdown ${mobile ? 'menu__link--mobile' : ''} dropdown--hoverable dropdown--right`}>
      <a href="#" role="button" aria-haspopup="true" aria-expanded="false" className={`navbar__link ${!mobile ? 'only-icon-locale' : ''}`} title={triggerAria} aria-label={triggerAria}>
        {icon}
        {mobile ? <span style={{marginLeft:8}}>{labelFor(currentLocale)}</span> : null}
      </a>
      <ul className="dropdown__menu">
        {locales.map((loc) => {
          const to = loc === (i18n?.defaultLocale || 'en') ? `${baseUrl}` : `${baseUrl}${loc}/`;
          const isActive = loc === currentLocale;
          return (
            <li key={loc}>
              <Link to={to} className={`dropdown__link ${isActive ? (mobile ? 'menu__link--active' : 'dropdown__link--active') : ''}`} aria-current={isActive ? 'page' : undefined} lang={loc}>
                {labelFor(loc)}
              </Link>
            </li>
          );
        })}
      </ul>
    </div>
  );
}
