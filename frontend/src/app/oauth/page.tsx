"use client";

import { GitHub, Google } from "@mui/icons-material";
import {
  Box,
  Button,
  ButtonGroup,
  Container,
  Typography,
  useTheme,
} from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2/Grid2";
import Image from "next/image";

export default function Oauth() {
  const theme = useTheme();
  return (
    <Container maxWidth="sm">
      <Box
        bgcolor={theme.palette.common.white}
        alignItems={"center"}
        alignSelf={"center"}
        justifyContent="center"
        padding={5}
        boxShadow={1}
        marginTop={"50%"}
      >
        <Grid
          display="flex"
          justifyContent="center"
          alignItems="center"
          marginBottom={3}
        >
          <Image src="/logo.png" alt="logo" width={40} height={40} />
          <Typography paddingLeft={2} fontWeight={600} fontSize={25}>
            Moecord
          </Typography>
        </Grid>
        <Grid
          display="flex"
          justifyContent="center"
          alignItems="center"
          marginBottom={1}
          marginTop={10}
        >
          Login with
        </Grid>
        <Box textAlign={"center"}>
          <ButtonGroup
            size="small"
            aria-label="small button group"
            orientation="horizontal"
          >
            <Button
              color="inherit"
              href="http://localhost:4000/v1/login/oauth/authorization/google"
            >
              <Google fontSize={"small"} color="warning" />
              <span style={{ marginLeft: 4 }}>Google</span>
            </Button>
            <Button
              color="inherit"
              href="http://localhost:4000/v1/login/oauth/authorization/github"
            >
              <GitHub fontSize={"small"} />
              <span style={{ marginLeft: 4 }}>Github</span>
            </Button>
            <Button
              color="inherit"
              href="http://localhost:4000/v1/login/oauth/authorization/discord"
            >
              <Image src="/discord.svg" alt="discord" width={20} height={20} />
              <span style={{ marginLeft: 4 }}>Discord</span>
            </Button>
          </ButtonGroup>
        </Box>
      </Box>
    </Container>
  );
}
