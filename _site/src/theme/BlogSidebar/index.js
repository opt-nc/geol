import React from 'react';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';

export default function BlogSidebar() {
  const {i18n} = useDocusaurusContext();
  const locale = (i18n && i18n.currentLocale) || 'en';
  const title = locale.startsWith('fr') ? 'Articles récents' : 'Recent posts';

  return (
    <aside className="col col--3">
      <nav className="sidebar_re4s thin-scrollbar blog-sidebar-no-select" aria-label="Blog recent posts navigation">
        <div className="sidebarItemTitle_pO2u margin-bottom--md">{title}</div>

        <div role="group">
          <h3 className="yearGroupHeading_rMGB">2026</h3>
          <ul className="sidebarItemList_Yudw clean-list">
            <li className="sidebarItem__DBe"><a className="sidebarItemLink_mo7H" href="/blog/geol-v2-3-0-export-endoflife-to-duckdb">geol v2.3.0 - export endoflife.date to duckdb</a></li>
          </ul>
        </div>

        <div role="group">
          <h3 className="yearGroupHeading_rMGB">2025</h3>
          <ul className="sidebarItemList_Yudw clean-list">
            <li className="sidebarItem__DBe"><a className="sidebarItemLink_mo7H" href="/blog/geol-fin-de-vie-ci-cd-distrobox-v1-3-0-full-edition">⏳ geol — Fin de vie logicielle ci cd distrobox (v1.3.0)</a></li>
            <li className="sidebarItem__DBe"><a className="sidebarItemLink_mo7H" href="/blog/geol-1-3-0-unboxing-the-check-command">📢 geol 1.3.0 — unboxing: the check command</a></li>
            <li className="sidebarItem__DBe"><a className="sidebarItemLink_mo7H" href="/blog/mvp-unboxing-geol-devops-secops-eol-cli">📢 MVP Unboxing geol — a devops/secops CLI to manage EOLs</a></li>
            <li className="sidebarItem__DBe"><a className="sidebarItemLink_mo7H" href="/blog/first-hackathon">Geol Go inception during Hackathon</a></li>
          </ul>
        </div>

        <div role="group">
          <h3 className="yearGroupHeading_rMGB">2022</h3>
          <ul className="sidebarItemList_Yudw clean-list">
            <li className="sidebarItem__DBe"><a className="sidebarItemLink_mo7H" href="/blog/end-of-life-date-first-article">Manage EoLs like a boss with endoflife.date</a></li>
          </ul>
        </div>

        
      </nav>
    </aside>
  );
}
