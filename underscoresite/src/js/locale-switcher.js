// Client-side locale switcher used by navbar HTML buttons
// Exposes window.switchLocale(locale) to switch between 'en' and 'fr'.
(function () {
  if (typeof window === 'undefined') return;

  function joinPath(prefix, path) {
    if (!path) return prefix || '/';
    if (prefix.endsWith('/') && path.startsWith('/')) return prefix + path.slice(1);
    return prefix + path;
  }

  window.switchLocale = function (locale) {
    try {
      var path = window.location.pathname || '/';
      if (locale === 'fr') {
        if (path.startsWith('/fr/')) return; // already french
        if (path === '/') {
          window.location.pathname = '/fr/';
        } else {
          window.location.pathname = joinPath('/fr', path);
        }
      } else {
        // switch to English: remove leading /fr
        if (path === '/fr/' || path === '/fr') {
          window.location.pathname = '/';
        } else if (path.startsWith('/fr/')) {
          window.location.pathname = path.replace(/^\/fr/, '') || '/';
        }
      }
    } catch (e) {
      console.error('locale switch failed', e);
      window.location.reload();
    }
  };
})();
