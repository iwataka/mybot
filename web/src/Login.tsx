import React from "react";
import { Button } from "react-bootstrap";
import { FaTwitter, FaSlack } from "react-icons/fa";

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

export default Login;
