"use client";

import { Box, Button } from "@mui/material";
import Grid2 from "@mui/material/Unstable_Grid2/Grid2";
import { atom, useAtom } from "jotai";
import { useRouter } from "next/navigation";
import { Suspense, useEffect, useState } from "react";
import { selectedServerIDAtom } from "./atom";
import theme from "../theme";

type ApiResponse = {
  code: number;
  message: string;
  data: Array<Server>;
};

type Server = {
  id: string;
  name: string;
  direct_message: boolean;
};

const fetchServer = async () => {
  const response = await fetch("/api/servers", {
    method: "GET",
  });
  const data: ApiResponse = await response.json();
  return data;
};

export default function Server() {
  const [servers, setServers] = useState<Array<Server>>([]);
  const router = useRouter();
  const [, setSelectedServerID] = useAtom(selectedServerIDAtom);

  useEffect(() => {
    (async () => {
      const { data, code } = await fetchServer();
      if (code != 200) {
        router.push("/oauth");
      }
      setSelectedServerID(data[0].id);
      setServers(data);
    })();
  }, []);

  return (
    <Grid2 xs={0.5}>
      <Suspense fallback={<div>loading....</div>}>
        {servers.map((val) => (
          <Button variant="outlined">
            <Box
              bgcolor={theme.palette.primary.contrastText}
              borderRadius={30}
              height={40}
              onClick={() => setSelectedServerID(val.id)}
            ></Box>
          </Button>
        ))}
      </Suspense>
    </Grid2>
  );
}
