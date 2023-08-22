"use client";

import Grid from "@mui/material/Unstable_Grid2"; // Grid version 2
import Image from "next/image";
import { Button, Container, Typography } from "@mui/material";
import { useRouter } from "next/navigation";

export default function Home() {
  const router = useRouter();
  return (
    <Container maxWidth="lg">
      <Grid container padding={2} minWidth={"100%"}>
        <Grid display="flex" justifyContent="center" alignItems="center">
          <Image src="/logo.png" alt="logo" width={35} height={35} />
          <Typography paddingLeft={2} fontWeight={600}>
            Moecord
          </Typography>
        </Grid>
        <Grid
          xs
          display="flex"
          justifyContent="center"
          alignItems="center"
        ></Grid>
        <Grid display="flex" justifyContent="center" alignItems="center">
          <Button variant="contained" onClick={() => router.push("/oauth")}>
            <Typography textTransform={"none"} color={"white"}>
              Login
            </Typography>
          </Button>
        </Grid>
      </Grid>
    </Container>
  );
}
