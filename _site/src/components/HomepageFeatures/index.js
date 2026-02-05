import clsx from 'clsx';
import Heading from '@theme/Heading';
import Translate, {translate} from '@docusaurus/Translate';
import styles from './styles.module.css';

const FeatureList = [
  {
    title: translate({id: 'homepage.feature.easy.title', message: 'Easy to use'}),
    Svg: (props) => (
      <img
        {...props}
        src={require('@site/static/img/mouse-cursor-click-icon-png-ojcVeq.png').default}
        alt={translate({id: 'homepage.feature.easy.alt', message: 'Click mouse'})}
      />
    ),
    description: (
      <>
        <code>geol</code>{' '}
        <Translate id="homepage.feature.easy.description">
          was designed be easily installed and used to get you ready and efficient with software end-of-life management.
        </Translate>
      </>
    ),
  },
  {
    title: translate({id: 'homepage.feature.focus.title', message: 'Focus on What Matters'}),
    Svg: (props) => (
      <img
        {...props}
        src={require('@site/static/img/Img_central_sans_dino.png').default}
        alt={translate({id: 'homepage.feature.focus.alt', message: 'Central image'})}
      />
    ),
    description: (
      <>
        <code>geol</code>{' '}
        <Translate id="homepage.feature.focus.description">
          lets you focus on software end-of-life management, exporting useful reports while we handle the data gathering and processing.
        </Translate>
      </>
    ),
  },
  {
    title: translate({id: 'homepage.feature.go.title', message: 'Powered by Go'}),
    Svg: (props) => (
      <img
        {...props}
        src={require('@site/static/img/go_transparent.png').default}
        alt={translate({id: 'homepage.feature.go.alt', message: 'Go language'})}
      />
    ),
    description: (
      <>
        <Translate id="homepage.feature.go.description">
          Built with the Go programming language,
        </Translate>{' '}
        <code>geol</code>{' '}
        <Translate id="homepage.feature.go.description_tail">
          is a fast, efficient, and reliable tool for managing software end-of-life information on any platform.
        </Translate>
      </>
    ),
  },
];

function Feature({Svg, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className={clsx('text--center padding-horiz--md', !title && styles.noTitle)}>
        <Heading as="h2" style={{visibility: title ? 'visible' : 'hidden'}} aria-hidden={!title}>
          {title || ' '}
        </Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
