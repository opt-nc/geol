import React from 'react';
import ComponentCreator from '@docusaurus/ComponentCreator';

export default [
  {
    path: '/blog',
    component: ComponentCreator('/blog', '303'),
    exact: true
  },
  {
    path: '/blog/archive',
    component: ComponentCreator('/blog/archive', '182'),
    exact: true
  },
  {
    path: '/blog/authors',
    component: ComponentCreator('/blog/authors', '0b7'),
    exact: true
  },
  {
    path: '/blog/authors/adriens',
    component: ComponentCreator('/blog/authors/adriens', '9ab'),
    exact: true
  },
  {
    path: '/blog/dev-challenge-hacktoberfest',
    component: ComponentCreator('/blog/dev-challenge-hacktoberfest', 'd53'),
    exact: true
  },
  {
    path: '/blog/end-of-life-date-first-article',
    component: ComponentCreator('/blog/end-of-life-date-first-article', '4b8'),
    exact: true
  },
  {
    path: '/blog/first-hackathon',
    component: ComponentCreator('/blog/first-hackathon', '031'),
    exact: true
  },
  {
    path: '/blog/mdx-blog-post',
    component: ComponentCreator('/blog/mdx-blog-post', 'e9f'),
    exact: true
  },
  {
    path: '/blog/tags',
    component: ComponentCreator('/blog/tags', '287'),
    exact: true
  },
  {
    path: '/blog/tags/docusaurus',
    component: ComponentCreator('/blog/tags/docusaurus', 'b6f'),
    exact: true
  },
  {
    path: '/blog/tags/facebook',
    component: ComponentCreator('/blog/tags/facebook', '858'),
    exact: true
  },
  {
    path: '/blog/tags/hello',
    component: ComponentCreator('/blog/tags/hello', '4ad'),
    exact: true
  },
  {
    path: '/blog/tags/hola',
    component: ComponentCreator('/blog/tags/hola', '439'),
    exact: true
  },
  {
    path: '/blog/welcome',
    component: ComponentCreator('/blog/welcome', 'd2b'),
    exact: true
  },
  {
    path: '/markdown-page',
    component: ComponentCreator('/markdown-page', '3d7'),
    exact: true
  },
  {
    path: '/docs',
    component: ComponentCreator('/docs', 'be0'),
    routes: [
      {
        path: '/docs',
        component: ComponentCreator('/docs', '295'),
        routes: [
          {
            path: '/docs',
            component: ComponentCreator('/docs', '95c'),
            routes: [
              {
                path: '/docs/category/tutorial---advanced',
                component: ComponentCreator('/docs/category/tutorial---advanced', 'b61'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/docs/category/tutorial---basics',
                component: ComponentCreator('/docs/category/tutorial---basics', '20e'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/docs/intro',
                component: ComponentCreator('/docs/intro', '61d'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/docs/tutorial-basics/check-command',
                component: ComponentCreator('/docs/tutorial-basics/check-command', '0aa'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/docs/tutorial-basics/congratulations',
                component: ComponentCreator('/docs/tutorial-basics/congratulations', '458'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/docs/tutorial-basics/deal-with-keys',
                component: ComponentCreator('/docs/tutorial-basics/deal-with-keys', 'b1c'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/docs/tutorial-basics/markdown-features',
                component: ComponentCreator('/docs/tutorial-basics/markdown-features', 'b05'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/docs/tutorial-basics/reporting-in-termial',
                component: ComponentCreator('/docs/tutorial-basics/reporting-in-termial', 'a2d'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/docs/tutorial-extras/asciidoc-toolchain',
                component: ComponentCreator('/docs/tutorial-extras/asciidoc-toolchain', '578'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/docs/tutorial-extras/dataengineering',
                component: ComponentCreator('/docs/tutorial-extras/dataengineering', '93b'),
                exact: true,
                sidebar: "tutorialSidebar"
              }
            ]
          }
        ]
      }
    ]
  },
  {
    path: '/',
    component: ComponentCreator('/', '2e1'),
    exact: true
  },
  {
    path: '*',
    component: ComponentCreator('*'),
  },
];
