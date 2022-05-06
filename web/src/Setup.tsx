import produce from "immer";
import React from "react";
import { Alert, Button, Form } from "react-bootstrap";

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

  componentDidMount() {
    fetch("/api/auth/credential", {
      credentials: "same-origin",
    })
      .then((res) => {
        if (res.ok) {
          res.json().then((data) => {
            this.setState({ credential: data });
          });
        } else {
          res.text().then((t) => this.setState({ error: t }));
        }
      })
      .catch((err) => {
        this.setState({ error: err });
      });
  }

  submit() {
    fetch("/api/auth/credential", {
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
      produce(this.state, (draft) => {
        draft.credential[typ][field] = value;
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
              type="text"
              placeholder="Enter consumer key"
              value={this.state.credential.twitter.consumer_key}
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
              value={this.state.credential.twitter.consumer_secret}
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
              type="text"
              placeholder="Enter consumer key"
              value={this.state.credential.slack.consumer_key}
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
              value={this.state.credential.slack.consumer_secret}
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

export default Setup;
