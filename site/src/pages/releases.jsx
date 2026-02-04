import React from 'react';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import En from '@site/src/components/Releases/en.mdx';
import ReleasesLayout from '@site/src/components/Releases/ReleasesLayout';

export default function ReleasesPage() {
  // English-only site: always render English releases content
  const Content = En;
  const title = 'Releases';
  return (
    <Layout title={title} description="Releases and changelog">
      <ReleasesLayout>
        <article>
          <div className="markdown">
            <Content />
          </div>
        </article>
      </ReleasesLayout>
    </Layout>
  );
}
