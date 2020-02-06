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

const baseCodeExample = `<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
${pre}go
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
${pre}
<!--valyala/fasthttp-->
${pre}go
package main

import (
    "fmt"
    "log"

    "github.com/valyala/fasthttp"
    "github.com/vardius/gorouter/v4"
)

func index(_ *fasthttp.RequestCtx) {
    fmt.Print("Welcome!\\n")
}

func main() {
    router := gorouter.NewFastHTTPRouter()
    router.GET("/", index)

    log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleFastHTTP))
}
${pre}
<!--END_DOCUSAURUS_CODE_TABS-->`;

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

    return (
      <SplashContainer>
        <Logo img_src={`${baseUrl}img/logo.png`} />
        <div className="inner">
          <ProjectTitle tagline={siteConfig.tagline} title={siteConfig.title} />
          <Container className="codeSample">
            <MarkdownBlock>{baseCodeExample}</MarkdownBlock>
          </Container>
          <PromoSection>
            <Button href={docUrl('installation.html')}>Get Started</Button>
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
      <div>
        <h1 className="paddingBottom" style={{textAlign: 'center'}} >Features</h1>
        <Block layout="fourColumn" background="dark">
          {[
            {
              key:"routing",
              title: 'Routing System',
              content: 'Routing with static and named parameters, easy setup for wildcards and regexp wildcards',
            },
            {
              key:"middleware",
              title: 'Middleware System',
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
              title: 'HTTP2 Support',
              content: 'Support for HTTP2, use it if you need it.',
            },
            {
              key:"memory",
              title: 'Low memory usage',
              content: 'Efficient and low memory usage, router implementation keeps allocations at 0!',
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
          <Features />
          <Showcase />
        </div>
      </div>
    );
  }
}

module.exports = Index;
