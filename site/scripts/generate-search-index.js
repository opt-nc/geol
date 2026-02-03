const fs = require('fs');
const path = require('path');

const docsDir = path.join(__dirname, '..', 'docs');
const pagesDir = path.join(__dirname, '..', 'src', 'pages');
const outDir = path.join(__dirname, '..', 'static');
const outFile = path.join(outDir, 'search-index.json');

function walk(dir) {
  const files = [];
  if (!fs.existsSync(dir)) return files;
  for (const name of fs.readdirSync(dir)) {
    const p = path.join(dir, name);
    const stat = fs.statSync(p);
    if (stat.isDirectory()) {
      files.push(...walk(p));
    } else if (/\.mdx?$/.test(name)) {
      files.push(p);
    }
  }
  return files;
}

function stripFrontMatter(content) {
  return content.replace(/^---[\s\S]*?---\s*/,'');
}

function frontMatterTitle(content) {
  const m = content.match(/^---\s*([\s\S]*?)\s*---/);
  if (!m) return null;
  const fm = m[1];
  const t = fm.match(/^[ \t]*title:\s*(?:"|')?(.+?)(?:"|')?\s*$/m);
  if (t) return t[1].trim();
  return null;
}

function frontMatterField(content, field) {
  const m = content.match(/^---\s*([\s\S]*?)\s*---/);
  if (!m) return null;
  const fm = m[1];
  const re = new RegExp('^[ \\t]*' + field + ':\s*(?:"|\')?(.+?)(?:"|\')?\s*$', 'm');
  const t = fm.match(re);
  return t ? t[1].trim() : null;
}

function frontMatterHasUnlisted(content) {
  const m = content.match(/^---\s*([\s\S]*?)\s*---/);
  if (!m) return false;
  const fm = m[1];
  return /(^|\n)[ \t]*unlisted:\s*true\b/.test(fm);
}

function stripMarkdown(md) {
  return md
    .replace(/```(?:\w*\n)?([\s\S]*?)```/g, '$1')
    .replace(/<[^>]+>/g, '')
    .replace(/\[(.*?)\]\([^)]*\)/g, '$1')
    .replace(/^#{1,6}\s+/gm, '')
    .replace(/^[>\-\*\+]\s+/gm, '')
    .replace(/`+/g, '')
    .replace(/~+/g, '')
    // remove emphasis markers **bold**, *italic*, __bold__, _italic_
    .replace(/\*\*(.*?)\*\*/g, '$1')
    .replace(/__(.*?)__/g, '$1')
    .replace(/\*(.*?)\*/g, '$1')
    .replace(/_(.*?)_/g, '$1')
    .replace(/\s+/g, ' ')
    .trim();
}

function titleFromContent(content, file) {
  const m = content.match(/^#\s+(.*)/m);
  if (m) return m[1].trim();
  return path.basename(file, path.extname(file));
}

function slugify(text) {
  return text
    .toString()
    .toLowerCase()
    .trim()
    .replace(/[\u2013\u2014]/g, '-')
    .replace(/[^a-z0-9\s-]/g, '')
    .replace(/\s+/g, '-');
}

function buildIndex() {
  // determine locales
  let locales = ['en'];
  let defaultLocale = 'en';
  try {
    const siteConfig = require(path.join(__dirname, '..', 'docusaurus.config.js'));
    if (siteConfig && siteConfig.i18n && Array.isArray(siteConfig.i18n.locales)) {
      locales = siteConfig.i18n.locales;
      defaultLocale = siteConfig.i18n.defaultLocale || locales[0];
    }
  } catch (e) {
    try {
      const i18nDirs = fs.readdirSync(path.join(__dirname, '..', 'i18n'));
      locales = Array.from(new Set(['en', ...i18nDirs]));
    } catch (e2) { locales = ['en']; }
  }

  const allLocaleItems = {};

  for (const locale of locales) {
    const items = [];
    const docsSource = (locale === defaultLocale) ? docsDir : path.join(__dirname, '..', 'i18n', locale, 'docusaurus-plugin-content-docs', 'current');
    const pagesSource = (locale === defaultLocale) ? pagesDir : path.join(__dirname, '..', 'i18n', locale, 'docusaurus-plugin-content-pages', 'default');

    const sources = [
      { dir: docsSource, prefix: '/docs/' },
      { dir: pagesSource, prefix: '/' },
    ];

    for (const src of sources) {
      // collect localized files first
      let files = walk(src.dir);
      // if locale has no translations for some docs/pages, include fallback default files
      if (locale !== defaultLocale) {
        const fallbackDir = (src.prefix === '/docs/') ? docsDir : pagesDir;
        const localizedRels = new Set(files.map(f => path.relative(src.dir, f).replace(/\\\\/g, '/')));
        const fallbackFiles = walk(fallbackDir);
        for (const ff of fallbackFiles) {
          const rel = path.relative(fallbackDir, ff).replace(/\\\\/g, '/');
          if (!localizedRels.has(rel)) files.push(ff);
        }
      }
      console.log('locale', locale, 'source', src.dir, 'found files', files.length);
      for (const f of files) {
        const raw = fs.readFileSync(f, 'utf8');
        // skip unlisted placeholders
        if (frontMatterHasUnlisted(raw)) continue;
        let nofm = stripFrontMatter(raw);
        // remove MDX/JS import lines to avoid indexing them (e.g. "import X from '...' ;")
        nofm = nofm.replace(/^[ \t]*import\s.+?;[ \t]*$/gm, '');
        if (/This page has been removed\.|Removed: the Markdown page/.test(raw)) continue;
        let rel = path.relative(src.dir, f).replace(/\\\\/g, '/').replace(/index\.mdx?$/,'').replace(/\.mdx?$/,'');
        let baseUrl = path.posix.join(src.prefix, rel || '/');
        if (!baseUrl.startsWith('/')) baseUrl = '/' + baseUrl;
        const fmTitle = frontMatterTitle(raw);
        const contentTitle = titleFromContent(nofm, f);
        const title = fmTitle || contentTitle || path.basename(f, path.extname(f));
        const text = stripMarkdown(nofm);
        let summary = text.slice(0, 300).trim();
        // Provide a better preview for the CLI playground page
        if (/cli-playground/.test(path.basename(f))) {
          if (locale && locale.toLowerCase().startsWith('fr')) {
            summary = "Voici un terminal vous permettant d'essayer quelques commandes et d'avoir le rendu exact de ce que vous pourriez avoir.";
          } else {
            summary = "Here is a terminal allowing you to try a few commands and see the kind of output you might get.";
          }
        }
        const fmDescription = frontMatterField(raw, 'description') || null;
        const fmTagsRaw = frontMatterField(raw, 'tags') || frontMatterField(raw, 'keywords') || null;
        // parse tags if present as YAML list or comma-separated
        let tags = null;
        if (fmTagsRaw) {
          try {
            // simple comma-separated fallback
            if (/\[.*\]/.test(fmTagsRaw)) {
              tags = fmTagsRaw.replace(/^[\s\[]+|[\]\s]+$/g, '').split(/,\s*/).map(s=>s.replace(/['"]/g,'').trim()).filter(Boolean);
            } else {
              tags = fmTagsRaw.split(/,\s*/).map(s=>s.replace(/['"]/g,'').trim()).filter(Boolean);
            }
          } catch (e) { tags = null; }
        }
        let pageUrl = baseUrl.endsWith('/') ? baseUrl : baseUrl + '/';
        // For non-default locales, ensure URLs are prefixed with the locale (e.g. /fr/docs/...)
        if (locale !== defaultLocale) {
          // use posix join to avoid duplicate slashes
          pageUrl = path.posix.join('/', locale, pageUrl);
          if (!pageUrl.endsWith('/')) pageUrl = pageUrl + '/';
        }

        // Normalize auto-generated page filenames created by Docusaurus i18n for pages
        // Example: /fr/site-src-pages-releases-md-7ef/ should map to the canonical /fr/releases/
        if (pageUrl.includes('site-src-pages-releases')) {
          pageUrl = pageUrl.replace(/\/site-src-pages-releases[^\/]*\//, '/releases/');
          if (locale !== defaultLocale) {
            // ensure locale prefix remains (e.g. /fr/releases/)
            pageUrl = path.posix.join('/', locale, 'releases/');
          } else {
            pageUrl = '/releases/';
          }
        }
        const excludePaths = ['/blog/2021-08-26-welcome', '/blog/welcome'];
        if (excludePaths.some(p => pageUrl.startsWith(p))) continue;
        items.push({ title, text, url: pageUrl, summary: fmDescription || summary, tags: tags || [], locale, type: src.prefix.replace(/\//g,'') || 'page' });

        const lines = nofm.split(/\r?\n/);
        for (let i = 0; i < lines.length; i++) {
          const line = lines[i];
          const m = line.match(/^(#{1,6})\s+(.*)$/);
          if (m) {
            const level = m[1].length;
            const heading = m[2].trim();
            if (level <= 1) continue;
            let j = i + 1;
            const collected = [];
            while (j < lines.length) {
              const next = lines[j];
              if (/^#{1,6}\s+/.test(next)) break;
              if (next.trim() === '' && collected.length > 0) break;
              collected.push(next);
              j++;
            }
            const block = stripMarkdown(collected.join('\n'));
            const slug = slugify(heading);
            const headingUrl = pageUrl.endsWith('/') ? pageUrl + '#' + slug : pageUrl + '/#' + slug;
            items.push({ title: heading, text: block, url: headingUrl, summary: block.slice(0,200), tags: [], locale, type: src.prefix.replace(/\//g,'') || 'page' });
          }
        }
      }
    }
    allLocaleItems[locale] = items;
    console.log('locale', locale, 'collected items', items.length);
  }

  // write default locale index
  // dedupe by normalized title, preferring docs > blog > pages
  const priority = (url) => {
    if (!url) return 3;
    if (url.startsWith('/docs/')) return 1;
    if (url.startsWith('/blog/')) return 2;
    return 3;
  };

  // default locale items
  const defaultItems = allLocaleItems[defaultLocale] || [];
  const map = new Map();
  for (const it of defaultItems) {
    const key = (it.title || '').toString().toLowerCase().trim();
    const existing = map.get(key);
    if (!existing) { map.set(key, it); continue; }
    const pNew = priority(it.url);
    const pOld = priority(existing.url);
    if (pNew < pOld) map.set(key, it);
  }
  let dedup = Array.from(map.values());
  // Remove any accidental aggregate page titled 'Versions' from the default index
  dedup = dedup.filter(it => (it.title || '').toString().trim().toLowerCase() !== 'versions');

  if (!fs.existsSync(outDir)) fs.mkdirSync(outDir);
  fs.writeFileSync(outFile, JSON.stringify(dedup, null, 2));
  console.log('Wrote', outFile, 'with', dedup.length, 'entries');

  // write other locales
  for (const locale of Object.keys(allLocaleItems)) {
    if (locale === defaultLocale) continue;
    const itemsForLocale = allLocaleItems[locale] || [];
    const mapLoc = new Map();
    for (const it of itemsForLocale) {
      const key = (it.title || '').toString().toLowerCase().trim();
      if (!mapLoc.has(key)) mapLoc.set(key, it);
    }
    let dedupLoc = Array.from(mapLoc.values());
    // Exclude the localized aggregate 'Versions' entry so it does not appear in search
    dedupLoc = dedupLoc.filter(it => (it.title || '').toString().trim().toLowerCase() !== 'versions');
    const localeOutDir = path.join(outDir, locale);
    if (!fs.existsSync(localeOutDir)) fs.mkdirSync(localeOutDir, { recursive: true });
    const localeFile = path.join(localeOutDir, 'search-index.json');
    fs.writeFileSync(localeFile, JSON.stringify(dedupLoc, null, 2));
    console.log('Wrote', localeFile, 'with', dedupLoc.length, 'entries');
  }
}

try { buildIndex(); } catch (e) { console.error(e); process.exit(1); }
