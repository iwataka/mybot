import React from "react";
import Container from "react-bootstrap/Container";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import "bootstrap/dist/css/bootstrap.css";
import Navbar from "react-bootstrap/Navbar";
import Nav from "react-bootstrap/Nav";
import { LinkContainer } from "react-router-bootstrap";
import { FaGithub, FaTwitter, FaSlack } from "react-icons/fa";
import Table from "react-bootstrap/Table";
import Badge from "react-bootstrap/Badge";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import Figure from "react-bootstrap/Figure";
import Alert from "react-bootstrap/Alert";
import Accordion from "react-bootstrap/Accordion";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import update from "immutability-helper";
import "./App.css";

const httpStatusNotAuthenticated = 498;
const httpStatusNotSetup = 499;

class App extends React.Component<{}, any> {
  constructor(props: {}) {
    super(props);
    this.state = {
      auth: {
        status: 0,
      },
    };
    this.requireAuth = this.requireAuth.bind(this);
    this.setAuthStatus = this.setAuthStatus.bind(this);
    this.resetAuthStatus = this.resetAuthStatus.bind(this);
  }

  requireAuth(children: JSX.Element) {
    let auth = this.state.auth;

    if (auth.status === 0) {
      fetch("/api/auth/status", {
        credentials: "same-origin",
      }).then((res) => {
        this.setAuthStatus(res.status);
      });
    }

    if (auth.status === httpStatusNotSetup) {
      return <Navigate to="/web/setup" replace />;
    }
    if (auth.status === httpStatusNotAuthenticated) {
      return <Navigate to="/web/login" replace />;
    }
    if (200 <= auth.status && auth.status < 300) {
      return children;
    }
    if (auth.status === 0) {
      return <Loading />;
    }
    return <Error />;
  }

  setAuthStatus(status: number) {
    this.setState(
      update(this.state, {
        auth: {
          status: { $set: status },
        },
      })
    );
  }

  resetAuthStatus() {
    this.setAuthStatus(0);
  }

  render() {
    return (
      <BrowserRouter>
        <Navbar>
          <Container>
            <Navbar.Toggle aria-controls="basic-navbar-nav" />
            <Navbar.Collapse id="basic-navbar-nav">
              <Nav className="me-auto">
                <LinkContainer to="/web">
                  <Nav.Link>Home</Nav.Link>
                </LinkContainer>
              </Nav>
              <Nav>
                <Nav.Link href="https://github.com/iwataka/mybot">
                  <FaGithub />
                </Nav.Link>
              </Nav>
            </Navbar.Collapse>
          </Container>
        </Navbar>
        <Container>
          <Routes>
            <Route path="/web" element={this.requireAuth(<Home />)} />
            <Route path="/web/config" element={this.requireAuth(<Config />)} />
            <Route
              path="/web/setup"
              element={<Setup resetAuthStatus={this.resetAuthStatus} />}
            />
            <Route path="/web/login" element={<Login />} />
          </Routes>
        </Container>
      </BrowserRouter>
    );
  }
}

class Loading extends React.Component<{}, {}> {
  render() {
    return <div>Loading...</div>;
  }
}

class Home extends React.Component<HomeProps, any> {
  constructor(props: HomeProps) {
    super(props);
    this.state = {
      statuses: {
        twitter_direct_message: null,
        twitter_timeline: null,
        twitter_polling: null,
        slack_channel: null,
      },
      imageAnalysisResult: {
        url: "",
        src: "",
        analysis_result: "",
        analysis_date: "",
      },
      error: "",
    };
  }

  componentDidMount() {
    fetch("/api/worker/status", {
      credentials: "same-origin",
    })
      .then((res) => {
        if (res.ok) {
          res.json().then((statuses) => {
            this.setState({ statuses: statuses });
          });
        } else {
          res.text().then((t) => {
            this.setState({ error: t });
          });
        }
      })
      .catch((err) => {
        this.setState({ error: err });
      });
  }

  statusBadge(status: boolean) {
    if (status === null) {
      return <Badge bg="secondary">Unknown</Badge>;
    }
    if (status) {
      return <Badge bg="success">Active</Badge>;
    }
    return <Badge bg="danger">Inactive</Badge>;
  }

  render() {
    let analysisResult = null;
    if (this.state.imageAnalysisResult.url) {
      analysisResult = (
        <Row>
          <Col>
            <h3>Image</h3>
            <Figure>
              <Figure.Image src={this.state.imageAnalysisResult.url} />
            </Figure>
          </Col>
          <Col>
            <h3>Analysis result</h3>
            <Alert variant="secondary">
              <pre>{this.state.imageAnalysisResult.analysis_result}</pre>
            </Alert>
          </Col>
        </Row>
      );
    } else {
      analysisResult = <Alert variant="info">Nothing to show currently</Alert>;
    }

    return (
      <div>
        <h1>Mybot</h1>
        <p className="lead">
          automatically collect and transfer any kinds of information for you
        </p>
        <h2 className="mt-5">Process Status</h2>
        <p>
          Mybot has the following processes.
          <br />
          If you find <Badge bg="danger">Inactive</Badge> process, please check
          your configuration or notify to administrators.
        </p>
        <Table responsive>
          <thead>
            <tr>
              <th>Category</th>
              <th>Process</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td rowSpan={3}>Twitter</td>
              <td>Direct Message</td>
              <td>
                {this.statusBadge(this.state.statuses.twitter_direct_message)}
              </td>
            </tr>
            <tr>
              <td>Timeline</td>
              <td>{this.statusBadge(this.state.statuses.twitter_timeline)}</td>
            </tr>
            <tr>
              <td>Polling (Search and Favorite)</td>
              <td>{this.statusBadge(this.state.statuses.twitter_polling)}</td>
            </tr>
            <tr>
              <td>Slack</td>
              <td>Channel Events</td>
              <td>{this.statusBadge(this.state.statuses.slack_channel)}</td>
            </tr>
          </tbody>
        </Table>
        <h2 className="mt-5">Image Analysis Result</h2>
        <p>
          Mybot has a feature to analyze image by AI (currently only Google
          Vision API is supported).
          <br />
          You can check the latest analysis result here.
        </p>
        {analysisResult}
      </div>
    );
  }
}

type HomeProps = {};

class Config extends React.Component<ConfigProps, any> {
  constructor(props: ConfigProps) {
    super(props);
    this.state = {
      config: {},
    };
  }

  render() {
    let config = this.state.config;

    let timelines = [];
    if (config.twitter != null && config.twitter.timelines != null) {
      for (let [i, val] of this.state.config.twitter.timelines.entries()) {
        timelines.push(<TimelineConfig eventKey={i} config={val} />);
      }
    }
    return (
      <div>
        <h1>Config</h1>
        <p className="lead">Customize your own bot as you want</p>
        <h2 className="mt-5">
          <FaTwitter /> Timeline
        </h2>
        <Accordion>{timelines}</Accordion>
        <h2 className="mt-5">
          <FaTwitter /> Favorite
        </h2>
        <h2 className="mt-5">
          <FaTwitter /> Search
        </h2>
        <h2 className="mt-5">
          <FaSlack /> Message
        </h2>
        <h2 className="mt-5">General</h2>
      </div>
    );
  }
}

type ConfigProps = {};

class TimelineConfig extends React.Component<TimelineConfigProps, any> {
  constructor(props: TimelineConfigProps) {
    super(props);
    this.state = {
      config: props.config,
    };
  }

  render() {
    let config = this.state.config;

    return (
      <div>
        <Accordion.Item eventKey={this.props.eventKey}>
          <Accordion.Header>{config.name}</Accordion.Header>
          <Accordion.Body></Accordion.Body>
        </Accordion.Item>
      </div>
    );
  }
}

type TimelineConfigProps = {
  eventKey: string;
  config: any;
};

class Setup extends React.Component<SetupProps, any> {
  constructor(props: SetupProps) {
    super(props);
    this.state = {
      credential: {
        twitter: {
          consumer_key: "",
          consumer_secret: "",
        },
        slack: {
          consumer_key: "",
          consumer_secret: "",
        },
      },
      error: "",
    };
    this.submit = this.submit.bind(this);
    this.updateState = this.updateState.bind(this);
  }

  submit() {
    fetch("/api/credential", {
      credentials: "same-origin",
      body: JSON.stringify(this.state.credential),
      method: "POST",
    })
      .then((res) => {
        if (!res.ok) {
          res.text().then((t) => this.setState({ error: t }));
        } else {
          this.props.resetAuthStatus();
        }
      })
      .catch((err) => {
        this.setState({ error: err });
      });
  }

  updateState(typ: string, field: string, value: string) {
    this.setState(
      update(this.state, {
        credential: {
          [typ]: {
            [field]: { $set: value },
          },
        },
      })
    );
  }

  render() {
    let errorAlert = null;
    if (this.state.error) {
      errorAlert = <Alert variant="danger">{this.state.error}</Alert>;
    }
    return (
      <div>
        <h1>Setup</h1>
        <p className="lead">
          You need to register Application information first
        </p>
        {errorAlert}

        <Form>
          <h2 className="mt-5">Twitter</h2>
          <p>
            refer to the Twitter App page{" "}
            <a href="https://apps.twitter.com/">here</a>
          </p>
          <Form.Group className="mb-3">
            <Form.Label>Consumer Key</Form.Label>
            <Form.Control
              type="password"
              placeholder="Enter consumer key"
              value={this.state.consumer_key}
              onChange={(e) => {
                this.updateState("twitter", "consumer_key", e.target.value);
              }}
            />
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Consumer Secret</Form.Label>
            <Form.Control
              type="password"
              placeholder="Enter consumer secret"
              value={this.state.consumer_secret}
              onChange={(e) => {
                this.updateState("twitter", "consumer_secret", e.target.value);
              }}
            />
          </Form.Group>

          <h2 className="mt-5">Slack</h2>
          <p>
            refer to the Slack App page{" "}
            <a href="https://api.slack.com/apps">here</a>
          </p>
          <Form.Group className="mb-3">
            <Form.Label>Consumer Key</Form.Label>
            <Form.Control
              type="password"
              placeholder="Enter consumer key"
              value={this.state.consumer_key}
              onChange={(e) => {
                this.updateState("slack", "consumer_key", e.target.value);
              }}
            />
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Consumer Secret</Form.Label>
            <Form.Control
              type="password"
              placeholder="Enter consumer secret"
              value={this.state.consumer_secret}
              onChange={(e) => {
                this.updateState("slack", "consumer_secret", e.target.value);
              }}
            />
          </Form.Group>

          <Button variant="primary" onClick={this.submit}>
            Submit
          </Button>
        </Form>
      </div>
    );
  }
}

type SetupProps = {
  resetAuthStatus: VoidFunction;
};

class Login extends React.Component<{}, {}> {
  createCallbackURL(provider: string) {
    let location = window.location;
    let proto = location.protocol;
    let host = location.hostname;
    let port = location.port;
    let callbackUrl = `${proto}//${host}:${port}/api/auth/callback/${provider}`;
    return encodeURI(callbackUrl);
  }

  render() {
    return (
      <div>
        <h1>Login</h1>
        <p className="lead">Login with your social account</p>
        <Button
          className="me-3"
          href={`/api/auth/twitter?callback=${this.createCallbackURL(
            "twitter"
          )}`}
        >
          <FaTwitter /> Login with Twitter
        </Button>
        <Button
          className="me-3"
          href={`/api/auth/slack?callback=${this.createCallbackURL("slack")}`}
        >
          <FaSlack /> Login with Slack
        </Button>
      </div>
    );
  }
}

class Error extends React.Component<{}, {}> {
  render() {
    return (
      <div>
        <Alert variant="danger">
          <Alert.Heading>Sorry, something went wrong.</Alert.Heading>
          <p>Please contact to administrators.</p>
        </Alert>
      </div>
    );
  }
}

export default App;
