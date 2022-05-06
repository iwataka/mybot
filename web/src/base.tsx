import React from "react";
import { Alert } from "react-bootstrap";

export class BaseComponent<P, S> extends React.Component<P, S> {
  renderErrorAlert(error: any, onClose: VoidFunction) {
    if (error) {
      return (
        <Alert variant="danger" onClose={onClose} dismissible>
          <Alert.Heading>Error</Alert.Heading>
          <p>{error}</p>
        </Alert>
      );
    }
    return null;
  }

  fetchAPI(
    endpoint: string,
    method: string,
    body: string | null,
    handleResponse: (res: Response) => void,
    handleErrorResponse: (res: Response) => void,
    handleError: (err: Error) => void
  ): void {
    fetch(endpoint, {
      credentials: "same-origin",
      method: method,
      body: body,
    })
      .then((res) => {
        if (res.ok) {
          handleResponse(res);
        } else {
          handleErrorResponse(res);
        }
      })
      .catch((err) => {
        handleError(err);
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
}

export type BaseProps = {
  handleErrorRespopnse: (res: Response) => void;
  handleError: (err: Error) => void;
};
