import React from "react";
import { Alert, Badge, Col, Figure, Row, Table } from "react-bootstrap";

class Home extends React.Component<HomeProps, any> {
  constructor(props: HomeProps) {
    super(props);
    this.state = {
      workerStatus: {
        twitter_direct_message: null,
        twitter_timeline: null,
        twitter_polling: null,
        slack_channel: null,
      },
      imageAnalysisStatus: {
        google: null,
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
    this.fetchAndSet("/api/worker/status", "workerStatus");
    this.fetchAndSet("/api/analysis/image/status", "imageAnalysisStatus");
    this.fetchAndSet("/api/analysis/image/result", "imageAnalysisResult");
  }

  fetchAndSet(path: string, key: string) {
    fetch(path, {
      credentials: "same-origin",
    })
      .then((res) => {
        if (res.ok) {
          res.json().then((data) => {
            this.setState({ [key]: data });
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
        <h2 className="mt-5">Feature Status</h2>
        <p>
          Mybot mainly has the following features.
          <br />
          If you find <Badge bg="danger">Inactive</Badge> feature, please check
          your configuration or notify to administrators.
        </p>
        <Table responsive>
          <thead>
            <tr>
              <th>Category</th>
              <th>Feature</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td rowSpan={3}>Twitter</td>
              <td>Direct Message (retired by Twitter)</td>
              <td>
                {this.statusBadge(
                  this.state.workerStatus.twitter_direct_message
                )}
              </td>
            </tr>
            <tr>
              <td>Timeline</td>
              <td>
                {this.statusBadge(this.state.workerStatus.twitter_timeline)}
              </td>
            </tr>
            <tr>
              <td>Polling (Search and Favorite)</td>
              <td>
                {this.statusBadge(this.state.workerStatus.twitter_polling)}
              </td>
            </tr>
            <tr>
              <td>Slack</td>
              <td>Channel Events</td>
              <td>{this.statusBadge(this.state.workerStatus.slack_channel)}</td>
            </tr>
            <tr>
              <td>Google</td>
              <td>Vision API</td>
              <td>{this.statusBadge(this.state.imageAnalysisStatus.google)}</td>
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

export default Home;
