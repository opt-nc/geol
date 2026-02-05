import React, {useEffect, useState} from 'react';
import Layout from '@theme/Layout';
import Link from '@docusaurus/Link';
import Translate, { translate } from '@docusaurus/Translate';
import Fuse from 'fuse.js';
import './search.css';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';

export default function SearchPage() {
  const [index, setIndex] = useState([]);
  const [query, setQuery] = useState('');
  const [results, setResults] = useState([]);
  const { siteConfig, i18n } = useDocusaurusContext();

  useEffect(() => {
    const baseUrl = (siteConfig && siteConfig.baseUrl) ? siteConfig.baseUrl.replace(/\/$/, '') : '';
    const currentLocale = i18n?.currentLocale || (typeof window !== 'undefined' && window.location && window.location.pathname.split('/')?.[1]) || 'en';
    const defaultLocale = i18n?.defaultLocale || 'en';

    const candidates = [];
    if (currentLocale && currentLocale !== defaultLocale) candidates.push(`${baseUrl}/${currentLocale}/search-index.json`);
    candidates.push(`${baseUrl}/search-index.json`);

    // try candidates in order until one succeeds
    (async () => {
      for (const p of candidates) {
        try {
          const res = await fetch(p);
          if (!res.ok) continue;
          const data = await res.json();
          if (Array.isArray(data) && data.length > 0) { setIndex(data); return; }
          // accept empty index but continue to try fallback if empty
          if (Array.isArray(data)) { setIndex(data); }
        } catch (e) {
          // try next
        }
      }
      // fallback to empty
      setIndex([]);
    })();
  }, [siteConfig, i18n]);

  useEffect(() => {
    if (!query || index.length === 0) { setResults([]); return; }

    // Exclude release version fragments (they look like `/releases/#v...`)
    const filteredIndex = index.filter(i => !((i.url || '').includes('/releases/#')));

    const fuse = new Fuse(filteredIndex, { keys: ['title', 'text'], threshold: 0.4, includeScore: true });
    const raw = fuse.search(query).slice(0, 400);
    const qlow = query.toLowerCase();
    // Also find direct substring matches in `text` to boost exact content matches
    const direct = filteredIndex.filter(i => (i.text || '').toLowerCase().includes(qlow));

    // Merge Fuse + direct into candidates list with scores
    const candidates = [];
    for (const r of raw) candidates.push({ item: r.item, score: typeof r.score === 'number' ? r.score : 1 });
    for (const it of direct) candidates.push({ item: it, score: 0 });

    // Group by base page URL (before '#')
    const groups = new Map();
    for (const c of candidates) {
      const url = (c.item.url || '');
      const base = (url.split('#')[0]) || '/';
      if (!groups.has(base)) groups.set(base, []);
      groups.get(base).push(c);
    }

    function escapeHtml(str) {
      return str.replace(/[&<>"']/g, s => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[s]));
    }

    function makeSnippet(text, q) {
      if (!text) return '';
      const qt = (q || '').toString();
      const qlow = qt.toLowerCase();
      const t = text.toString();
      // simple full excerpt when no query
      if (!qlow) {
        const s = t.slice(0, 140) + (t.length > 140 ? '…' : '');
        return escapeHtml(s);
      }
      const li = t.toLowerCase().indexOf(qlow);
      const start = li === -1 ? 0 : Math.max(0, li - 40);
      const end = li === -1 ? Math.min(t.length, 140) : Math.min(t.length, li + qlow.length + 40);
      const prefix = start > 0 ? '…' : '';
      const suffix = end < t.length ? '…' : '';
      // extract and escape
      const rawSlice = prefix + t.slice(start, end).trim() + suffix;
      const escaped = escapeHtml(rawSlice);
      // highlight query occurrences (case-insensitive)
      try {
        const qEsc = qt.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
        const re = new RegExp(qEsc, 'gi');
        return escaped.replace(re, m => `<span class="search-match">${m}</span>`);
      } catch (e) {
        return escaped;
      }
    }

    const resultsArr = [];
    for (const [base, items] of groups.entries()) {
      // For each page, pick the best candidate
      let best = null;
      for (const it of items) if (!best || it.score < best.score) best = it;
      const pageEntry = index.find(e => (e.url || '').split('#')[0] === base) || {};
      // prefer a fragment candidate if present
      const fragItem = items.find(i => (i.item.url || '').includes('#')) || null;
      const chosen = fragItem || best;
      if (!chosen) continue;
      const isFragment = (chosen.item.url || '').includes('#');
      let outUrl = isFragment ? chosen.item.url : (pageEntry.url || base);
      if (isFragment) {
        const parts = outUrl.split('#');
        const frag = parts[1] || '';
        const baseOnly = parts[0] || outUrl || '/';
        outUrl = `${baseOnly}${baseOnly.includes('?') ? '&' : '?'}anchorTitle=${encodeURIComponent(chosen.item.title)}#${frag}`;
      }
      const title = pageEntry.title || chosen.item.title || '';
      const subtitle = (isFragment && chosen.item.title && chosen.item.title !== title) ? chosen.item.title : '';
      const snippet = makeSnippet(chosen.item.text || pageEntry.text || '', query);
      resultsArr.push({ title, subtitle, url: outUrl, snippet, score: chosen.score });
    }

    // sort by score and limit
    resultsArr.sort((a, b) => (a.score || 1) - (b.score || 1));
    setResults(resultsArr.slice(0, 50));
  }, [query, index]);

  return (
    <Layout title={translate({id: 'theme.search.title', message: 'Search'})}>
      <main className="container margin-vert--lg">
        <h1><Translate id="theme.search.heading">Search</Translate></h1>
        <div className="search-box">
          <input
            value={query}
            onChange={e=>setQuery(e.target.value)}
            placeholder={translate({id: 'theme.search.placeholder', message: 'Search docs...'})}
          />
          {query && (
            <button
              className="search-clear"
              onClick={() => setQuery('')}
              aria-label={translate({id: 'theme.search.clearAria', message: 'Clear search'})}
            >
              ×
            </button>
          )}
        </div>
        <div className="search-results">
          {results.length === 0 && query && <p><Translate id="theme.search.noResults">No results</Translate></p>}
          {results.map((r, i) => (
            <div key={i} className="search-result-item">
              <h3><Link to={r.url}>{r.title}</Link></h3>
              {r.snippet && <p className="search-snippet" dangerouslySetInnerHTML={{ __html: r.snippet }} />}
            </div>
          ))}
        </div>
      </main>
    </Layout>
  );
}
