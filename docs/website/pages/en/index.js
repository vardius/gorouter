/**
 * Copyright (c) 2017-present, Facebook, Inc.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

const React = require('react');

const CompLibrary = require('../../core/CompLibrary.js');

const MarkdownBlock = CompLibrary.MarkdownBlock; /* Used to read markdown */
const Container = CompLibrary.Container;
const GridBlock = CompLibrary.GridBlock;

const pre = "```";

const baseCodeExample = `${pre}
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vardius/gorouter/v4"
)

func index(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Welcome!")
}

func main() {
	router := gorouter.New()
	router.GET("/", http.HandlerFunc(index))

	log.Fatal(http.ListenAndServe(":8080", router))
}
${pre}`;

class HomeSplash extends React.Component {
  render() {
    const {siteConfig, language = ''} = this.props;
    const {baseUrl, docsUrl} = siteConfig;
    const docsPart = `${docsUrl ? `${docsUrl}/` : ''}`;
    const langPart = `${language ? `${language}/` : ''}`;
    const docUrl = doc => `${baseUrl}${docsPart}${langPart}${doc}`;

    const SplashContainer = props => (
      <div className="homeContainer">
        <div className="homeSplashFade">
          <div className="wrapper homeWrapper">{props.children}</div>
        </div>
      </div>
    );

    const Logo = props => (
      <div className="projectLogo">
        <img src={props.img_src} alt="Project Logo" />
      </div>
    );

    const ProjectTitle = props => (
      <h2 className="projectTitle">
        {props.title}
        <small>{props.tagline}</small>
      </h2>
    );

    const PromoSection = props => (
      <div className="section promoSection">
        <div className="promoRow">
          <div className="pluginRowBlock">{props.children}</div>
        </div>
      </div>
    );

    const Button = props => (
      <div className="pluginWrapper buttonWrapper">
        <a className="button" href={props.href} target={props.target}>
          {props.children}
        </a>
      </div>
    );

    const Block = props => (
      <GridBlock
        background={props.background}
        className="codeSample"
        contents={props.children}
        layout={props.layout}
      />
    );

    const BasicExample = () => (
      <Block background="light">
        {[
          {
            key: "example",
            content: baseCodeExample,
            image: `${baseUrl}img/logo.png`,
            imageAlign: 'right',
          },
        ]}
      </Block>
    );

    return (
      <SplashContainer>
        {/* <Logo img_src={`${baseUrl}img/logo.png`} /> */}
        <div className="inner">
          <ProjectTitle tagline={siteConfig.tagline} title={siteConfig.title} />

          <BasicExample />

          <PromoSection>
            <Button href={docUrl('intro.html')}>Get Started</Button>
          </PromoSection>
        </div>
      </SplashContainer>
    );
  }
}

class Index extends React.Component {
  render() {
    const {config: siteConfig, language = ''} = this.props;
    const {baseUrl} = siteConfig;

    const Block = props => (
      <Container
        padding={['bottom', 'top']}
        id={props.id}
        background={props.background}>
        <GridBlock
          align="center"
          contents={props.children}
          layout={props.layout}
        />
      </Container>
    );

    const Features = () => (
      <Block layout="fourColumn" background="dark">
        {[
          {
            key:"middleware",
            title: 'Middleware',
            content: 'Build-in middleware system with order by priority',
          },
          {
            key:"authentication",
            title: 'Authentication',
            content: 'Easy authentication',
          },
          {
            key:"fasthttp",
            title: 'Fast HTTP',
            content: 'Multiple implementations. Support for native net/http or valyala/fasthttp.',
          },
          {
            key:"files",
            title: 'Serving Files',
            content: 'Out fof box static files serving',
          },
          {
            key:"multidomain",
            title: 'Multidomain',
            content: 'Easy multidomain setup',
          },
          {
            key:"http2",
            title: 'HTTP2',
            content: 'Support for HTTP2',
          },
        ]}
      </Block>
    );

    const Showcase = () => {
      if ((siteConfig.users || []).length === 0) {
        return null;
      }

      const showcase = siteConfig.users
        .filter(user => user.pinned)
        .map(user => (
          <a href={user.infoLink} key={user.infoLink}>
            <img src={user.image} alt={user.caption} title={user.caption} />
          </a>
        ));

      const pageUrl = page => baseUrl + (language ? `${language}/` : '') + page;

      return (
        <div className="productShowcaseSection paddingBottom">
          <h2>Who is Using This?</h2>
          <p>This project is used by all these people</p>
          <div className="logos">{showcase}</div>
          <div className="more-users">
            <a className="button" href={pageUrl('users.html')}>
              More {siteConfig.title} Users
            </a>
          </div>
        </div>
      );
    };

    return (
      <div>
        <HomeSplash siteConfig={siteConfig} language={language} />
        <div className="mainContainer">
          <Features />
          <Showcase />
        </div>
      </div>
    );
  }
}

module.exports = Index;
