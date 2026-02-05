import React, {useEffect, useRef} from 'react';

export default function AnchorFallback() {
  const lastKeyRef = useRef('');

  useEffect(() => {
    if (typeof window === 'undefined') return;

    const getAnchorTitle = () => new URLSearchParams(window.location.search || '').get('anchorTitle');
    const anchorTitle = getAnchorTitle();
    if (!anchorTitle) return;

    const normalize = (s) => String(s).trim();
    const makeCandidates = (title) => {
      const t = normalize(title).toLowerCase();
      const r1 = t.replace(/[\u2013\u2014]/g, '-').replace(/[^a-z0-9\s-]/g, '').replace(/\s+/g, '-');
      const r2 = r1.replace(/-+/g, '-');
      const r3 = r1.replace(/-+/g, '--');
      return Array.from(new Set([r1, r2, r3, encodeURIComponent(title)]));
    };

    const tryScrollOnce = (title) => {
      // First try the explicit hash if present
      if (window.location.hash) {
        const h = window.location.hash.slice(1);
        const tryIds = new Set();
        tryIds.add(h);
        // common variations: triple vs double hyphens and numeric suffixes that Docusaurus may add
        tryIds.add(h.replace(/---/g, '--'));
        tryIds.add(h.replace(/--/g, '---'));
        tryIds.add(h + '-1');
        tryIds.add(h.replace(/---/g, '--') + '-1');
        tryIds.add(h.replace(/--/g, '---') + '-1');
        for (const id of tryIds) {
          const el = document.getElementById(id);
          if (el) { el.scrollIntoView({behavior: 'smooth', block: 'start'}); return true; }
          const q = document.querySelector(`[id='${id}']`);
          if (q) { q.scrollIntoView({behavior: 'smooth', block: 'start'}); return true; }
        }
      }
      // Otherwise try candidates derived from anchorTitle
      const candidates = makeCandidates(title || getAnchorTitle());
      for (const c of candidates) {
        // try candidate and common variants
        const vars = [c, c.replace(/---/g, '--'), c.replace(/--/g, '---'), c + '-1', c.replace(/---/g, '--') + '-1'];
        for (const v of vars) {
          const el = document.getElementById(v) || document.querySelector(`[id='${v}']`);
          if (el) { el.scrollIntoView({behavior: 'smooth', block: 'start'}); return true; }
        }
      }
      return false;
    };

    // Run immediately and also whenever the location changes (client-side navigation)
    const runAttempts = (title) => {
      let attempts = 0;
      const maxAttempts = 8;
      const intervals = [0, 50, 150, 300, 600, 1000, 1200, 1500];
      const runner = () => {
        attempts++;
        if (tryScrollOnce(title) || attempts >= maxAttempts) return;
        setTimeout(runner, intervals[Math.min(attempts, intervals.length - 1)]);
      };
      runner();
    };

    // Provide a 'locationchange' event by monkey-patching history methods
    const origPush = history.pushState;
    const origReplace = history.replaceState;
    const dispatchLocationChange = () => window.dispatchEvent(new Event('locationchange'));
    history.pushState = function () { const r = origPush.apply(this, arguments); dispatchLocationChange(); return r; };
    history.replaceState = function () { const r = origReplace.apply(this, arguments); dispatchLocationChange(); return r; };

    const onLocationChange = () => {
      const key = `${window.location.pathname}${window.location.search}${window.location.hash}`;
      if (key === lastKeyRef.current) return;
      lastKeyRef.current = key;
      runAttempts(getAnchorTitle());
    };

    window.addEventListener('locationchange', onLocationChange);
    window.addEventListener('hashchange', onLocationChange);

    // initial attempt
    onLocationChange();

    return () => {
      window.removeEventListener('locationchange', onLocationChange);
      window.removeEventListener('hashchange', onLocationChange);
      history.pushState = origPush;
      history.replaceState = origReplace;
    };
  }, []);

  return null;
}
