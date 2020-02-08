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

const baseCodeExample = `${pre}go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/vardius/gorouter/v4"
    "github.com/vardius/gorouter/v4/context"
)

func index(w http.ResponseWriter, _ *http.Request) {
    fmt.Fprint(w, "Welcome!\\n")
}

func main() {
    router := gorouter.New()
    router.GET("/", http.HandlerFunc(index))

    log.Fatal(http.ListenAndServe(":8080", router))
}
${pre}`;

const description = `
Pick router that **does not** slow down with response size and maintains high performance for large and deep route tree.

Most of the router benchmarks out there test only against root route, which does not give a great feedback.

**gorouter** is designed to scale well against deep route tree.
Its architecture allows to keep high performance with low memory usage no matter how deep and big route tree is.
Built-in middleware system enables you to build complex solutions keeping performance at its best!

- extensive set of features
- compatible with multiple http packages
- use of native context

allows you to keep your business logic decoupled from external dependencies.
`;

class HomeSplash extends React.Component {
  render() {
    const {siteConfig, language = ''} = this.props;
    const {baseUrl, docsUrl} = siteConfig;
    const docsPart = `${docsUrl ? `${docsUrl}/` : ''}`;
    const langPart = `${language ? `${language}/` : ''}`;
    const docUrl = doc => `${baseUrl}${docsPart}${langPart}${doc}`;
    const pageUrl = page => `${baseUrl}${langPart}${page}`;

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

    return (
      <SplashContainer>
        <Logo img_src={`${baseUrl}img/logo.png`} />
        <div className="inner">
          <ProjectTitle tagline={siteConfig.tagline} title={siteConfig.title} />
          <Container className="codeSample">
            <MarkdownBlock>{baseCodeExample}</MarkdownBlock>
          </Container>
          <PromoSection>
            <Button href={docUrl('installation.html')}>Documentation</Button>
            <Button href={pageUrl('help')}>Help</Button>
            <Button href="https://github.com/vardius/gorouter">GitHub</Button>
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
          align={props.align || 'center'}
          contents={props.children}
          layout={props.layout}
        />
      </Container>
    );

    const Description = () => (
      <Block background="light" align="left">
          {[
            {
              title: 'Are you looking for a router that can handle deep route trees and large response sizes ?',
              content: description,
              image: baseUrl + 'img/logo.png',
              imageAlign: 'left',
            }
        ]}
      </Block>
    );

    const Features = () => (
      <div>
        <h1 className="paddingBottom" style={{textAlign: 'center'}} >Features</h1>
        <Block layout="fourColumn" background="dark">
          {[
            {
              key:"routing",
              title: 'Routing System',
              content: 'Parameters with flexible patterns including regexp wildcards.',
            },
            {
              key:"middleware",
              title: 'Middleware System',
              content: 'Build-in middleware system with order by priority.',
            },
            {
              key:"authentication",
              title: 'Authentication',
              content: 'Easy authentication.',
            },
            {
              key:"fasthttp",
              title: 'Fast HTTP',
              content: 'Multiple implementations. Support for native net/http or valyala/fasthttp.',
            },
            {
              key:"files",
              title: 'Serving Files',
              content: 'Out of box static files serving.',
            },
            {
              key:"multidomain",
              title: 'Multidomain',
              content: 'Easy multidomain setup.',
            },
            {
              key:"http2",
              title: 'HTTP2 Support',
              content: 'Support for HTTP2.',
            },
            {
              key:"memory",
              title: 'Low memory usage',
              content: 'Efficient and low memory usage, performent and flexible for any response size, no matter depth of the route tree.',
            },
          ]}
        </Block>
      </div>
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
          <Description />
          <Features />
          <Showcase />
        </div>
      </div>
    );
  }
}

module.exports = Index;
