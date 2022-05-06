import React from "react";
import { Accordion, Form, Table } from "react-bootstrap";
import { FaSlack, FaTwitter } from "react-icons/fa";

const filterSchema = {
  has_media: "boolean",
  favorite_threshold: "number",
  retweet_threshold: "number",
  lang: "string",
  patterns: "array<string>",
  url_patterns: "array<string>",
  vision: {
    label: "string",
    face: {
      anger_likelihood: "string",
      bluerred_likelihood: "string",
      headwear_likelihood: "string",
      joy_likelihood: "string",
    },
    text: "array<string>",
    landmark: "array<string>",
    logo: "array<string>",
  },
  language: {
    min_sentiment: "number",
    max_sentiment: "number",
  },
};

const actionSchema = {
  twitter: {
    tweet: "boolean",
    retweet: "boolean",
    favorite: "boolean",
    collections: "array<string>",
  },
  slack: {
    pin: "boolean",
    star: "boolean",
    reactions: "array<string>",
    channels: "array<string>",
  },
};

// null value means end sign of schema
const twitterTimelineSchema = {
  name: "string",
  screen_names: "array<string>",
  exclude_replies: "boolean",
  include_rts: "boolean",
  count: "number",
  filter: filterSchema,
  action: actionSchema,
};

const twitterFavoriteSchema = {
  name: "string",
  screen_names: "array<string>",
  count: "number",
  filter: filterSchema,
  action: actionSchema,
};

const twitterSearchSchema = {
  name: "string",
  queries: "array<string>",
  result_type: "string",
  count: "number",
  filter: filterSchema,
  action: actionSchema,
};

const slackMessageSchema = {
  name: "string",
  channels: "array<string>",
  filter: filterSchema,
  action: actionSchema,
};

const generalSchema = {
  duration: "string",
};

class Config extends React.Component<ConfigProps, any> {
  constructor(props: ConfigProps) {
    super(props);
    this.state = {
      config: {},
    };
  }

  componentDidMount() {
    fetch("/api/config", {
      credentials: "same-origin",
    }).then((res) => {
      if (res.ok) {
        res.json().then((data) => {
          this.setState({ config: data });
        });
      }
    });
  }

  render() {
    let config = this.state.config;
    return (
      <div>
        <h1>Config</h1>
        <p className="lead">Customize your own bot as you want</p>
        <h2 className="mt-5">
          <FaTwitter /> Twitter
        </h2>
        <h3 className="mt-5">Timeline</h3>
        {config && config.twitter ? (
          <ConfigTableList
            configList={config.twitter.timelines}
            schema={twitterTimelineSchema}
          />
        ) : null}
        <h3 className="mt-5">Favorite</h3>
        {config && config.twitter ? (
          <ConfigTableList
            configList={config.twitter.favorites}
            schema={twitterFavoriteSchema}
          />
        ) : null}
        <h3 className="mt-5">Search</h3>
        {config && config.twitter ? (
          <ConfigTableList
            configList={config.twitter.searches}
            schema={twitterSearchSchema}
          />
        ) : null}
        <h2 className="mt-5">
          <FaSlack /> Slack
        </h2>
        <h3 className="mt-5">Message</h3>
        {config && config.slack ? (
          <ConfigTableList
            configList={config.slack.messages}
            schema={slackMessageSchema}
          />
        ) : null}
        <h2 className="mt-5">General</h2>
        {config ? (
          <ConfigTable
            eventKey="general"
            config={config}
            schema={generalSchema}
          />
        ) : null}
      </div>
    );
  }
}

type ConfigProps = {};

class ConfigTableList extends React.Component<ConfigTableListProps, any> {
  render() {
    let configList: JSX.Element[] = [];
    if (this.props.configList !== null) {
      configList = Object.entries(this.props.configList).map(([i, val]) => {
        return (
          <Accordion.Item key={i} eventKey={i}>
            <Accordion.Header>{val.name}</Accordion.Header>
            <Accordion.Body>
              <ConfigTable
                key={i}
                eventKey={i}
                config={val}
                schema={this.props.schema}
              />
            </Accordion.Body>
          </Accordion.Item>
        );
      });
    }
    if (configList.length > 0) {
      return <Accordion>{configList}</Accordion>;
    }
    return null;
  }
}

type ConfigTableListProps = {
  configList: any[];
  schema: any;
};

class ConfigTable extends React.Component<ConfigTableProps, any> {
  render() {
    return this.renderConfigTable(this.props.config);
  }

  private renderConfigTable(config: any): JSX.Element {
    let numOfFieldCols = this.calcDepth(this.props.schema, 0);
    let tableRows = this.renderConfigRows(
      [],
      this.props.schema,
      this.calcRowSpans([], this.props.schema),
      {},
      config,
      numOfFieldCols
    );
    return (
      <Table>
        <thead>
          <tr>
            <th colSpan={numOfFieldCols}>config item</th>
            <th>value</th>
          </tr>
        </thead>
        <tbody>{tableRows}</tbody>
      </Table>
    );
  }

  private calcDepth(schema: any, depth: number): number {
    if (this.isSchemaEnd(schema)) {
      return depth;
    }
    return Math.max(
      ...Object.values(schema).map((val) => this.calcDepth(val, depth + 1))
    );
  }

  private calcRowSpans(
    schemaStack: string[],
    curSchema: any
  ): { [key: string]: number } {
    if (this.isSchemaEnd(curSchema)) {
      return { [schemaStack.join(".")]: 1 };
    }

    let result: { [key: string]: number } = {};
    let rowSpan = 0;
    for (const [key, value] of Object.entries(curSchema)) {
      let newStack = schemaStack.slice();
      newStack.push(key);
      let schemaToRowSpan = this.calcRowSpans(newStack, value);
      rowSpan += Math.max(...Object.values(schemaToRowSpan));
      for (const [k, v] of Object.entries(schemaToRowSpan)) {
        result[k] = v;
      }
    }
    if (schemaStack.length > 0) {
      result[schemaStack.join(".")] = rowSpan;
    }
    return result;
  }

  private renderConfigRows(
    schemaStack: string[],
    curSchema: any,
    schemaToRowSpan: { [key: string]: number },
    schemaIsRendered: { [key: string]: boolean },
    config: any,
    numOfFieldCols: number
  ): JSX.Element[] {
    if (this.isSchemaEnd(curSchema)) {
      let schema: string[] = [];
      let field_cols = schemaStack.map((key, index) => {
        schema.push(key);
        if (schemaIsRendered[schema.join(".")] || false) {
          return null;
        }
        let colSpan =
          index === schemaStack.length - 1
            ? numOfFieldCols - schemaStack.length + 1
            : 1;
        let rowSpan = schemaToRowSpan[schema.join(".")];
        schemaIsRendered[schema.join(".")] = true;
        return (
          <td key={key} colSpan={colSpan} rowSpan={rowSpan}>
            {key}
          </td>
        );
      });
      return [
        <tr key={schemaStack.join(".")}>
          {field_cols}
          <td>{this.renderValue(config, curSchema)}</td>
        </tr>,
      ];
    }

    return Object.entries(curSchema).flatMap(([key, value]) => {
      let new_schema_stack = schemaStack.slice();
      new_schema_stack.push(key);
      return this.renderConfigRows(
        new_schema_stack,
        value,
        schemaToRowSpan,
        schemaIsRendered,
        config[key],
        numOfFieldCols
      );
    });
  }

  private renderValue(value: any, typ: string): JSX.Element {
    if (typ === "string") {
      return (
        <Form.Control value={value ? value : ""} type="text" readOnly={true} />
      );
    }
    if (typ === "boolean") {
      return <Form.Check type="switch" checked={value} readOnly={true} />;
    }
    if (typ === "array<string>") {
      return (
        <Form.Control
          value={value ? value.join(",") : ""}
          type="text"
          readOnly={true}
        />
      );
    }
    if (typ === "number") {
      return (
        <Form.Control
          value={value ? value.toString() : ""}
          type="number"
          readOnly={true}
        />
      );
    }
    return <div>{value}</div>;
  }

  private isSchemaEnd(value: any): boolean {
    return value === null || typeof value === "string";
  }
}

type ConfigTableProps = {
  eventKey: string;
  config: any;
  schema: any;
};

export default Config;
