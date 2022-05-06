import "bootstrap/dist/css/bootstrap.css";
import produce from "immer";
import React from "react";
import Alert from "react-bootstrap/Alert";
import Container from "react-bootstrap/Container";
import Nav from "react-bootstrap/Nav";
import Navbar from "react-bootstrap/Navbar";
import { FaGithub } from "react-icons/fa";
import { LinkContainer } from "react-router-bootstrap";
import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import "./App.css";
import { BaseComponent } from "./base";
import Config from "./Config";
import Home from "./Home";
import Login from "./Login";
import Setup from "./Setup";

const httpStatusNotAuthenticated = 498;
const httpStatusNotSetup = 499;

class App extends React.Component<{}, {}> {
  render() {
    return (
      <BrowserRouter>
        <AppWithoutRouter />
      </BrowserRouter>
    );
  }
}

class AppWithoutRouter extends BaseComponent<{}, any> {
  constructor(props: {}) {
    super(props);
    this.state = {
      auth: {
        status: 0,
      },
      error: "",
    };
    this.requireAuth = this.requireAuth.bind(this);
    this.setAuthStatus = this.setAuthStatus.bind(this);
    this.resetAuthStatus = this.resetAuthStatus.bind(this);
    this.handleError = this.handleError.bind(this);
    this.handleErrorResponse = this.handleErrorResponse.bind(this);
  }

  requireAuth(children: JSX.Element) {
    let auth = this.state.auth;

    if (auth.status === 0) {
      this.getAPI(
        "/api/auth/status",
        (res) => this.setAuthStatus(res.status),
        (res) => this.setAuthStatus(res.status),
        (err) => this.handleError(err)
      );
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
    return <ErrorView />;
  }

  setAuthStatus(status: number) {
    this.setState(
      produce(this.state, (draft) => {
        draft.auth.status = status;
      })
    );
  }

  handleError(err: Error) {
    this.setState({ error: err });
  }

  handleErrorResponse(res: Response) {
    if (
      res.status === httpStatusNotAuthenticated ||
      res.status === httpStatusNotSetup
    ) {
      this.setAuthStatus(res.status);
    }
    res.text().then((text) => {
      this.setState({ error: text });
    });
  }

  resetAuthStatus() {
    this.setAuthStatus(0);
  }

  render() {
    return (
      <div>
        <Navbar>
          <Container>
            <Navbar.Toggle aria-controls="basic-navbar-nav" />
            <Navbar.Collapse id="basic-navbar-nav">
              <Nav className="me-auto">
                <LinkContainer to="/web">
                  <Nav.Link>Home</Nav.Link>
                </LinkContainer>
                <LinkContainer to="/web/config">
                  <Nav.Link>Config</Nav.Link>
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
          {this.renderErrorAlert(this.state.error)}
          <Routes>
            <Route
              path="/web"
              element={this.requireAuth(
                <Home
                  handleError={this.handleError}
                  handleErrorRespopnse={this.handleErrorResponse}
                />
              )}
            />
            <Route
              path="/web/config"
              element={this.requireAuth(
                <Config
                  handleError={this.handleError}
                  handleErrorRespopnse={this.handleErrorResponse}
                />
              )}
            />
            <Route
              path="/web/setup"
              element={
                <Setup
                  resetAuthStatus={this.resetAuthStatus}
                  handleError={this.handleError}
                  handleErrorRespopnse={this.handleErrorResponse}
                />
              }
            />
            <Route
              path="/web/login"
              element={
                <Login
                  handleError={this.handleError}
                  handleErrorRespopnse={this.handleErrorResponse}
                />
              }
            />
          </Routes>
        </Container>
      </div>
    );
  }
}

class Loading extends React.Component<{}, {}> {
  render() {
    return <div>Loading...</div>;
  }
}

// NOTE: class name "Error" is conflicted with buil-in Error.
class ErrorView extends React.Component<{}, {}> {
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

export { App as default, AppWithoutRouter };
