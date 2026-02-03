import React, {useState, useRef} from 'react';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import styles from './CLIPlayground.module.css';

function defaultSimulator(cmd){
  const c = (cmd||'').trim();
  if(!c) return '';
  if(c.startsWith('geol check')) return `12:14 WRN > Cache last updated XXXX-XX-XX 09:22:13, older than 24 hours. Updating the cache...
12:14 INF > Number of products=431 elapsed time (ms)=1158
12:14 INF > Number of tags=83 elapsed time (ms)=401
12:14 INF > Number of categories=9 elapsed time (ms)=418
12:14 INF > You have the latest geol version !
12:14 ERR > quarkus 3.12 (quarkus) is 1y 6m 1d past EOL (EOL: 2024-07-31)
12:14 ERR > maven 3.6 (maven) is 4y 10m 5d past EOL (EOL: 2021-03-30)
12:14 WRN > traefik 2.11 (traefik) is nearing EOL in 3d (EOL: 2026-02-01)
12:14 INF > Found skip:true for non-existent-product 1.0, product will be skipped
## MySuperApp
 Software     │ Version │ EOL Date   │ Status │ Days  │ Is Latest │ Latest 
──────────────┼─────────┼────────────┼────────┼───────┼───────────┼────────
 maven        │ 3.6     │ 2021-03-30 │ EOL    │ -1765 │ false     │ 3.9    
 quarkus      │ 3.12    │ 2024-07-31 │ EOL    │ -546  │ false     │ 3.30   
 traefik      │ 2.11    │ 2026-02-01 │ WARN   │ 3     │ false     │ 3.6    
 ubuntu       │ 25.10   │ 2026-07-01 │ OK     │ 153   │ true      │ 25.10  
 postgresql   │ 14      │ 2026-11-12 │ OK     │ 287   │ false     │ 18     
 java temurin │ 21      │ 2029-12-31 │ OK     │ 1432  │ false     │ 25     
 opensearch   │ 2       │            │ OK     │ -     │ false     │ 3      
.`;

  if(c === 'geol version') return `2.5.1\n12:15 INF > Checking the latest geol version...\n12:15 INF > You have the latest geol version !`;

  if(c === 'geol about') return `
  ____     ______     ______    __ 
/  _____| |   ____|  /  __  \\  |  |     
|  |  __  |  |__    |  |  |  | |  |     
|  | |_ | |   __|   |  |  |  | |  |     
|  |__| | |  |____  |  '--'  | |  '---.
 \\______| |_______|  \\______/  |_______|
⏳ Tech doesn’t last forever. Awareness does.

--- Build Info ---
GitVersion:          2.5.1
Git Commit:          1c6b0eb252ede344a5a703ac52c065262603a150
BuildDate:           2026-01-07T03:56:22Z
BuiltBy:             GoReleaser
GoVersion:           go1.25.5
Compiler:            gc
Platform:            linux/amd64

--- Ressources ---
Licence:             Apache-2.0
Code:                https://github.com/opt-nc/geol
Roadmap:             https://github.com/orgs/opt-nc/projects/28
API:                 https://endoflife.date
`;

  if(c === 'geol list products') return `+ adonisjs
+ akeneo-pim
+ alibaba-ack
+ alibaba-dragonwell
+ almalinux
+ alpine-linux
+ amazon-aurora-postgresql
+ amazon-cdk
+ amazon-corretto
+ amazon-documentdb
+ amazon-eks
+ amazon-glue
+ amazon-linux
+ amazon-msk
+ amazon-neptune
+ amazon-rds-mariadb
+ amazon-rds-mysql
+ amazon-rds-postgresql
+ android
+ angular
+ angularjs
+ ansible
+ ansible-core
+ ant
+ antix
+ apache-activemq
+ apache-activemq-artemis
+ apache-airflow
+ apache-apisix
+ apache-camel
+ apache-cassandra
+ apache-couchdb
+ apache-flink
+ apache-groovy
+ apache-hadoop
+ apache-hop
+ apache-http-server
+ apache-kafka
+ apache-lucene
+ apache-maven
+ apache-nifi
+ apache-pulsar
+ apache-spark
+ apache-struts
+ apache-subversion
+ api-platform
+ apple-watch
+ arangodb
+ argo-cd
+ argo-workflows
+ artifactory
+ authentik
+ aws-lambda
+ azul-zulu
+ azure-devops-server
+ azure-kubernetes-service
+ backdrop
+ bamboo
+ bazel
+ beats
+ behat
+ bellsoft-liberica
+ big-ip
+ bigbluebutton
+ bitbucket
+ bitcoin-core
+ blender
+ bootstrap
+ boundary
+ bun
+ cachet
+ caddy
+ cakephp
+ calico
+ centos
+ centos-stream
+ centreon
+ cert-manager
+ cfengine
+ chef-infra-client
+ chef-infra-server
+ chef-inspec
+ chef-supermarket
+ chef-workstation
+ chrome
+ cilium
+ cisco-ios-xe
+ citrix-vad
+ ckeditor
+ clamav
+ clear-linux
+ cloud-sql-auth-proxy
+ cnspec
+ cockroachdb
+ coder
+ coldfusion
+ commvault
+ composer
+ confluence
+ consul
+ containerd
+ contao
+ contour
+ controlm
+ cortex-xdr
+ cos
+ couchbase-server
+ craft-cms
+ dbt-core
+ dce
+ debian
+ deno
+ dependency-track
+ devuan
+ django
+ docker-engine
+ dotnet
+ dotnetfx
+ dovecot
+ drupal
+ drush
+ duckdb
+ eclipse-jetty
+ eclipse-temurin
+ elasticsearch
+ electron
+ elixir
+ emberjs
+ envoy
+ erlang
+ eslint
+ esxi
+ etcd
+ eurolinux
+ exim
+ express
+ fairphone
+ fedora
+ ffmpeg
+ filemaker
+ firefox
+ fluent-bit
+ flux
+ font-awesome
+ foreman
+ forgejo
+ fortios
+ freebsd
+ gatekeeper
+ gerrit
+ ghc
+ gitlab
+ gleam
+ go
+ goaccess
+ godot
+ google-kubernetes-engine
+ google-nexus
+ gorilla
+ graalvm-ce
+ gradle
+ grafana
+ grafana-loki
+ grails
+ graylog
+ greenlight
+ grumphp
+ grunt
+ gstreamer
+ guzzle
+ haproxy
+ harbor
+ hashicorp-packer
+ hashicorp-vault
+ hbase
+ hibernate-orm
+ horizon
+ ibm-aix
+ ibm-db2
+ ibm-i
+ ibm-mq
+ ibm-semeru-runtime
+ icinga
+ icinga-web
+ idl
+ influxdb
+ intel-processors
+ internet-explorer
+ ionic
+ ios
+ ipad
+ ipados
+ iphone
+ isc-dhcp
+ istio
+ jaeger
+ jekyll
+ jenkins
+ jhipster
+ jira-software
+ joomla
+ jquery
+ jquery-ui
+ jreleaser
+ julia
+ kde-plasma
+ keda
+ keycloak
+ kibana
+ kindle
+ kirby
+ kong-gateway
+ kotlin
+ kubernetes
+ kubernetes-csi-node-driver-registrar
+ kubernetes-node-feature-discovery
+ kuma
+ kyverno
+ laravel
+ ldap-account-manager
+ libreoffice
+ lineageos
+ linux
+ linuxmint
+ liquibase
+ log4j
+ logstash
+ looker
+ lua
+ macos
+ mageia
+ magento
+ mandrel
+ mariadb
+ mastodon
+ matomo
+ mattermost
+ mautic
+ mediawiki
+ meilisearch
+ memcached
+ micronaut
+ microsoft-build-of-openjdk
+ mongodb
+ moodle
+ motorola-mobility
+ msexchange
+ mssqlserver
+ mulesoft-runtime
+ mxlinux
+ mysql
+ neo4j
+ neos
+ netapp-ontap
+ netbackup-appliance-os
+ netbsd
+ nextcloud
+ nextjs
+ nexus
+ nginx
+ nix
+ nixos
+ nodejs
+ nokia
+ nomad
+ notepad-plus-plus
+ numpy
+ nutanix-aos
+ nutanix-files
+ nutanix-prism
+ nuxt
+ nvidia
+ nvidia-gpu
+ nvm
+ office
+ oneplus
+ openbao
+ openbsd
+ openjdk-builds-from-oracle
+ opensearch
+ openssl
+ opensuse
+ opentofu
+ openvpn
+ openwrt
+ openzfs
+ opnsense
+ oracle-apex
+ oracle-database
+ oracle-graalvm
+ oracle-jdk
+ oracle-linux
+ oracle-solaris
+ ovirt
+ pangp
+ panos
+ pci-dss
+ perl
+ phoenix-framework
+ photon
+ php
+ phpbb
+ phpmyadmin
+ pigeonhole
+ pixel
+ pixel-watch
+ plesk
+ plone
+ pnpm
+ podman
+ pop-os
+ postfix
+ postgresql
+ postmarketos
+ powershell
+ privatebin
+ proftpd
+ prometheus
+ protractor
+ proxmox-ve
+ puppet
+ python
+ qt
+ quarkus-framework
+ quasar
+ rabbitmq
+ rails
+ rancher
+ raspberry-pi
+ react
+ react-native
+ red-hat-openshift
+ redhat-build-of-openjdk
+ redhat-jboss-eap
+ redhat-satellite
+ redis
+ redmine
+ rhel
+ robo
+ rocket-chat
+ rocky-linux
+ ros
+ ros-2
+ roundcube
+ rtpengine
+ ruby
+ rust
+ salt
+ samsung-galaxy-tab
+ samsung-galaxy-watch
+ samsung-mobile
+ sapmachine
+ scala
+ sharepoint
+ shopware
+ silverstripe
+ slackware
+ sles
+ sns-firmware
+ sns-hardware
+ sns-smc
+ solr
+ sonarqube-community
+ sonarqube-server
+ sony-xperia
+ sourcegraph
+ splunk
+ spring-boot
+ spring-framework
+ sqlite
+ squid
+ statamic
+ steamos
+ surface
+ suse-linux-micro
+ suse-manager
+ svelte
+ symfony
+ tails
+ tailwind-css
+ tarantool
+ tarteaucitron
+ telegraf
+ teleport
+ terraform
+ thumbor
+ tls
+ tomcat
+ traefik
+ tvos
+ twig
+ typo3
+ ubuntu
+ umbraco
+ unity
+ unrealircd
+ valkey
+ varnish
+ vcenter
+ veeam-backup-and-replication
+ veeam-backup-for-microsoft-365
+ veeam-one
+ virtualbox
+ visionos
+ visual-cobol
+ visual-studio
+ vitess
+ vmware-cloud-foundation
+ vmware-harbor-registry
+ vmware-srm
+ vue
+ vuetify
+ wagtail
+ watchos
+ weakforced
+ weechat
+ windows
+ windows-embedded
+ windows-nano-server
+ windows-powershell
+ windows-server
+ windows-server-core
+ wireshark
+ wordpress
+ xcp-ng
+ yarn
+ yocto
+ zabbix
+ zentyal
+ zerto
+ zookeeper

431 products listed`;

  if(c.startsWith('geol product extended') || c.startsWith('geol extended product')) return `# Products
          
## windows
          
 Cycle               │ Release    │ Latest     │ Support    │ EOL        
─────────────────────┼────────────┼────────────┼────────────┼────────────
 11-25h2-e           │ 2025-09-30 │ 10.0.26200 │ 2028-10-10 │ 2028-10-10 
 11-25h2-w           │ 2025-09-30 │ 10.0.26200 │ 2027-10-12 │ 2027-10-12 
 11-24h2-iot-lts LTS │ 2024-10-01 │ 10.0.26100 │ 2029-10-09 │ 2034-10-10 
 11-24h2-e-lts LTS   │ 2024-10-01 │ 10.0.26100 │ 2029-10-09 │ 2029-10-09 
 11-24h2-e           │ 2024-10-01 │ 10.0.26100 │ 2027-10-12 │ 2027-10-12 
 ...                 │ ...        │ ...        │ ...        │ ...        
─────────────────────┴────────────┴────────────┴────────────┴────────────
45 rows (5 shown)`;

  return `${c}: command not found`;
}

export default function CLIPlayground({examples}){
  const {i18n} = useDocusaurusContext();
  const locale = (i18n && i18n.currentLocale) ? i18n.currentLocale : 'en';
  const [history, setHistory] = useState([]);
  const [input, setInput] = useState('');
  const [copied, setCopied] = useState(null);
  const endRef = useRef(null);
  const isFr = locale && locale.toLowerCase().startsWith('fr');
  const t = {
    run: isFr ? 'Exécuter' : 'Run',
    use: isFr ? 'Utiliser' : 'Use',
    copy: isFr ? 'Copier' : 'Copy',
    copied: isFr ? 'Copié' : 'Copied',
    examples: isFr ? 'Exemples' : 'Examples',
    available: isFr ? 'Commandes disponibles' : 'Available commands',
    note: isFr ? "Seules les commandes affichées ci-dessous peuvent être utilisées dans ce terminal simulé. Les autres commandes ne sont pas disponibles ici — installez geol sur votre machine pour exécuter toutes les commandes localement." : "Only the commands displayed below can be used in this simulated terminal. Other commands are not available here — please install geol on your machine to run all commands locally.",
    cliInput: isFr ? 'Entrée du terminal' : 'CLI input',
    cliRegion: isFr ? 'Terminal interactif' : 'CLI playground',
  };

  function runCommand(cmd){
    const out = defaultSimulator(cmd);
    setHistory(h => [...h, {cmd, output: out}]);
    setInput('');
    setTimeout(()=> endRef.current && endRef.current.scrollIntoView({behavior:'smooth'}),50);
  }

  async function copyText(text){
    try{
      await navigator.clipboard.writeText(text);
      setCopied(text);
      setTimeout(()=> setCopied(null), 1200);
    }catch(e){
      // fallback
    }
  }

  const availableCommands = [
    'geol check',
    'geol version',
    'geol product extended windows -n5',
    'geol list products',
    'geol about',
  ];

  return (
    <div className={styles.wrapper}>
      <div className={styles.terminal} role="region" aria-label={t.cliRegion}>
        <div className={styles.screen}>
          {history.map((h, i)=> (
            <div key={i} className={styles.entry}>
              <div className={styles.prompt}>&gt; {h.cmd}</div>
              <pre className={styles.output}>{h.output}</pre>
            </div>
          ))}
          <div ref={endRef} />
        </div>
        <form className={styles.inputRow} onSubmit={(e)=>{e.preventDefault(); runCommand(input)}}>
          <input className={styles.input} value={input} onChange={(e)=>setInput(e.target.value)} placeholder={isFr? 'Tapez une commande (ex. geol check)': 'Type a command (eg. geol check)'} aria-label={t.cliInput} />
          <button className={styles.run} type="submit">{t.run}</button>
        </form>
      </div>
      {/* Examples panel rendered below the terminal */}
      <div className={styles.examples} aria-label={t.examples}>
        {(() => {
          const hidden = new Set(['geol check', 'geol version', 'geol product extended windows -n5']);
          const visible = (examples || []).filter(e => !hidden.has(e));
          if (!visible.length) return null;
          return visible.map((ex, idx) => (
            <div key={idx} className={styles.exampleItem}>
              <div className={styles.code}>{ex}</div>
              <div className={styles.buttons}>
                <button type="button" className={styles.btnSmall} aria-label={`${t.use} ${ex}`} title={`${t.use} ${ex}`} onClick={()=>{ setInput(ex); }}>{t.use}</button>
                <button type="button" className={styles.btnSmall} aria-label={`${t.copy} ${ex}`} title={`${t.copy} ${ex}`} onClick={()=>copyText(ex)}>{copied===ex? t.copied : t.copy}</button>
                <button type="button" className={styles.btnSmall} aria-label={`${t.run} ${ex}`} title={`${t.run} ${ex}`} onClick={()=>runCommand(ex)}>{t.run}</button>
              </div>
            </div>
          ));
        })()}
      </div>
      <div className={styles.commands} aria-label={t.available}>
        <div className={styles.commandsTitle}>{t.available}</div>
        <div className={styles.commandGrid}>
          {availableCommands.map((cmd, i)=> (
            <div key={i} className={`${styles.commandItem} ${cmd && cmd.trim() === 'geol about' ? styles.aboutCmd : ''}`}>
              <div className={styles.code}>{cmd}</div>
              <div className={styles.buttons}>
                <button type="button" className={styles.btnSmall} aria-label={`${t.use} ${cmd}`} title={`${t.use} ${cmd}`} onClick={()=>{ setInput(cmd); }}>{t.use}</button>
                <button type="button" className={styles.btnSmall} aria-label={`${t.copy} ${cmd}`} title={`${t.copy} ${cmd}`} onClick={()=>copyText(cmd)}>{copied===cmd? t.copied : t.copy}</button>
                <button type="button" className={styles.btnSmall} aria-label={`${t.run} ${cmd}`} title={`${t.run} ${cmd}`} onClick={()=>runCommand(cmd)}>{t.run}</button>
              </div>
            </div>
          ))}
        </div>
        <div className={styles.note}>
          <>{t.note}</>
        </div>
      </div>
    </div>
  );
}
