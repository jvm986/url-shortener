import React from "react";
import ReactDOM from "react-dom";
import CssBaseline from "@material-ui/core/CssBaseline";
import { createMuiTheme, ThemeProvider } from "@material-ui/core/styles";
import App from "./App";

const theme = createMuiTheme({
  palette: {
    type: window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light",
  },
});

ReactDOM.render(
  <ThemeProvider theme={theme}>
    {/* CssBaseline kickstart an elegant, consistent, and simple baseline to build upon. */}
    <CssBaseline />
    <App />
  </ThemeProvider>,
  document.querySelector("#root")
);
