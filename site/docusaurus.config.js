// @ts-check
// `@type` JSDoc annotations allow editor autocompletion and type checking
// (when paired with `@ts-check`).
// There are various equivalent ways to declare your Docusaurus config.
// See: https://docusaurus.io/docs/api/docusaurus-config

import {themes as prismThemes} from 'prism-react-renderer';

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'geol',
  tagline: 'Software End Of Life management is too important to be  boring ',
  // Favicon (restored to original image)
  favicon: 'img/logo-no-name-gradient.png',

  future: {
    v4: true,
  },

  url: 'https://your-docusaurus-site.example.com',
  baseUrl: '/',

  organizationName: 'opt-nc',
  projectName: 'geol',

  onBrokenLinks: 'throw',

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */ ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
        },
        blog: {
          showReadingTime: true,
          feedOptions: { type: ['rss', 'atom'], xslt: true },
          editUrl: 'https://github.com/opt-nc/geol',
          onInlineTags: 'warn',
          onInlineAuthors: 'warn',
          onUntruncatedBlogPosts: 'warn',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  clientModules: [
    require.resolve('./src/js/no-image-drag.js'),
    require.resolve('./src/js/lazy-images-client.js'),
    require.resolve('./src/js/left-panel-toggle.js'),
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */ ({
      metadata: [
        {name: 'description', content: 'geol — manage software end-of-life and track supported versions.'},
        {name: 'og:site_name', content: 'geol'},
        {name: 'og:type', content: 'website'},
        {name: 'twitter:card', content: 'summary_large_image'},
      ],
      image: 'img/docusaurus-social-card.jpg',
      colorMode: { respectPrefersColorScheme: true },
      navbar: {
        title: 'geol',
        logo: { alt: 'geol Logo', src: 'img/logo-no-name-gradient.png' },
        items: [
          { to: '/search', label: 'Search', position: 'right' },
          {
            label: 'Menu',
            position: 'left',
            className: 'left-panel-dropdown',
            items: [
                { type: 'docSidebar', sidebarId: 'tutorialSidebar', label: 'Tutorial' },
                { to: '/blog', label: 'Blog' },
              ],
          },
          { href: 'https://github.com/opt-nc/geol', label: 'GitHub', position: 'right' },
          // locale dropdown removed (only English site)
        ],
      },
        footer: {
        style: 'dark',
        links: [
          { title: 'Docs', items: [{ label: 'Tutorial', to: '/docs/intro' }] },
          {
            title: 'Community',
            items: [
              { label: 'YouTube Playlist', href: 'https://www.youtube.com/playlist?list=PL7GdrgVAWcDhTeW8rjUCZVp2Mic1fx-j9' },
              { label: 'geol DEV.to series', href: 'https://dev.to/adriens/series/34740' },
              { label: 'endoflife.date DEV.to Series', href: 'https://dev.to/adriens/series/21232' },
              { label: 'Project kanban', href: 'https://github.com/orgs/opt-nc/projects/28' },
            ],
          },
          { title: 'More', items: [{ label: 'Blog', to: '/blog' }, { label: 'GitHub', href: 'https://github.com/opt-nc/geol' }] },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} geol contributors. Built with Docusaurus.`,
      },
      prism: { theme: prismThemes.github, darkTheme: prismThemes.dracula },
    }),
};

export default config;
