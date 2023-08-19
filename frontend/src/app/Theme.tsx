import { createTheme } from "@mui/material";

export default createTheme({
  palette: {
    mode: "dark",
    background: {
      default: "#070707",
      paper: "#0B0B0D",
    },
    primary: {
      main: "#ce4a3c",
      dark: "#140a0cc",
    },
    text: {
      primary: "#ffffff",
      secondary: "#b5b5b5",
    },
    grey: {
      "900": "#120a09",
    },
  },
});
