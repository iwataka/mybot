import React from "react";
import { Alert } from "react-bootstrap";

export class BaseComponent<P, S> extends React.Component<P, S> {
  renderErrorAlert(error: any) {
    if (error) {
      return <Alert variant="danger">{error}</Alert>;
    }
    return null;
  }

  fetchAPI(
    endpoint: string,
    method: string,
    body: string | null,
    ifOK: (res: Response) => void,
    ifFail: (res: Response) => void,
    ifError: (err: Error) => void
  ): void {
    fetch(endpoint, {
      credentials: "same-origin",
      method: method,
      body: body,
    })
      .then((res) => {
        if (res.ok) {
          ifOK(res);
        } else {
          ifFail(res);
        }
      })
      .catch((err) => {
        ifError(err);
      });
  }

  getAPI(
    endpoint: string,
    ifOK: (res: Response) => void,
    ifFail: (res: Response) => void,
    ifError: (err: Error) => void
  ): void {
    this.fetchAPI(endpoint, "GET", null, ifOK, ifFail, ifError);
  }

  fetchJsonAPI(
    endpoint: string,
    method: string,
    body: string | null,
    setResponseData: (data: any) => void,
    setError: (err: string | Error) => void
  ): void {
    this.fetchAPI(
      endpoint,
      method,
      body,
      (res) => {
        res.json().then((data) => {
          setResponseData(data);
        });
      },
      (res) => {
        res.text().then((text) => {
          setError(text);
        });
      },
      (err) => {
        setError(err);
      }
    );
  }

  getJsonAPI(
    endpoint: string,
    setResponseData: (data: any) => void,
    setError: (err: string | Error) => void
  ): void {
    this.fetchJsonAPI(endpoint, "GET", null, setResponseData, setError);
  }
}

export type BaseProps = {
  setError: (err: string | Error) => void;
};
