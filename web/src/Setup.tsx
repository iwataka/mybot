import produce from "immer";
import { Button, Form } from "react-bootstrap";
import { BaseComponent, BaseProps } from "./base";

class Setup extends BaseComponent<SetupProps, any> {
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
    };
    this.submit = this.submit.bind(this);
    this.updateState = this.updateState.bind(this);
  }

  componentDidMount() {
    this.getAPI(
      "/api/auth/credential",
      (res) => res.json().then((data) => this.setState({ credential: data })),
      (res) => this.props.handleErrorRespopnse(res),
      (err) => this.props.handleError(err)
    );
  }

  submit() {
    this.fetchAPI(
      "/api/auth/credential",
      "POST",
      JSON.stringify(this.state.credential),
      (_) => this.props.resetAuthStatus(),
      (res) => this.props.handleErrorRespopnse(res),
      (err) => this.props.handleError(err)
    );
  }

  updateState(typ: string, field: string, value: string) {
    this.setState(
      produce(this.state, (draft) => {
        draft.credential[typ][field] = value;
      })
    );
  }

  render() {
    return (
      <div>
        <h1>Setup</h1>
        <p className="lead">
          You need to register Application information first
        </p>
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
            Save
          </Button>
        </Form>
      </div>
    );
  }
}

type SetupProps = {
  resetAuthStatus: VoidFunction;
} & BaseProps;

export default Setup;
