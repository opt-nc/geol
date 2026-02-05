// Client-only script: safely set native lazy-loading on images and iframes.
// Guard against SSR: do not reference `window` at module evaluation time.
if (typeof window !== 'undefined') {
  window.addEventListener('load', () => {
    try {
      document.querySelectorAll('img').forEach(img => {
        if (!img.hasAttribute('loading')) img.setAttribute('loading', 'lazy');
      });
      document.querySelectorAll('iframe').forEach(fr => {
        if (!fr.hasAttribute('loading')) fr.setAttribute('loading', 'lazy');
      });
    } catch (e) {
      // ignore
    }
  });
}
