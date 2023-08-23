"use client";

import { useTheme } from "@mui/material/styles";
import Grid from "@mui/material/Unstable_Grid2";
import Server from "./server";
import { Button } from "@mui/material";
import Channel from "./channel";

export default function Layout({ children }: { children: React.ReactNode }) {
  const theme = useTheme();
  return (
    <Grid
      container
      bgcolor={theme.palette.background.default}
      minHeight={"100vh"}
      minWidth={"100%"}
    >
      <Server />
      <Channel />
      <Grid xs bgcolor={theme.palette.grey[400]}>
        <Button>{children}</Button>
      </Grid>
    </Grid>
  );
}
