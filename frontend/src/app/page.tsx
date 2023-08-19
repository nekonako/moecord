"use client";

import { Button } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import { useTheme } from "@mui/material/styles";
import { useEffect } from "react";

export default function Home() {
  useEffect(() => {
    const ws = new WebSocket("ws://localhost:4001");
    console.log(ws);
    ws.onopen = (e) => {
      console.log(e);
      ws.send("se");
    };
  }, []);

  const theme = useTheme();
  return (
    <Grid
      container
      bgcolor={theme.palette.background.default}
      minHeight={"100vh"}
      minWidth={"100%"}
    >
      <Grid xs={0.5}>
        <Button>xs</Button>
      </Grid>
      <Grid xs bgcolor={theme.palette.grey[900]}>
        <Button>1</Button>
      </Grid>
      <Grid
        xs={9}
        color={"ButtonText"}
        bgcolor={theme.palette.background.paper}
      >
        <Button>Hello world</Button>
      </Grid>
    </Grid>
  );
}
