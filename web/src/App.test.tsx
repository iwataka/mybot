import { render, screen } from "@testing-library/react";
import { AppWithoutRouter } from "./App";
import { Router } from "react-router-dom";
import { createMemoryHistory } from "history";

test("renders setup page", () => {
  const history = createMemoryHistory();
  history.push("/web/setup");
  render(
    <Router location={history.location} navigator={history}>
      <AppWithoutRouter />
    </Router>
  );
  const linkElement = screen.getByText(/Home/i);
  expect(linkElement).toBeInTheDocument();
});
