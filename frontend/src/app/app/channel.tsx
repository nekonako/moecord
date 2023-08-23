"use client";

import { Button, useTheme } from "@mui/material";
import Grid2 from "@mui/material/Unstable_Grid2/Grid2";
import { useAtom } from "jotai";
import { useRouter } from "next/navigation";
import { Suspense, useEffect, useState } from "react";
import { selectedServerIDAtom } from "./atom";

type ApiResponse = {
  code: number;
  message: string;
  data: Array<Channel>;
};

type Channel = {
  id: string;
  name: string;
};

const fetchChannel = async (id: string) => {
  const response = await fetch(`/api/channels/${id}`, {
    method: "GET",
  });
  const data: ApiResponse = await response.json();
  return data;
};

export default function Channel() {
  const [channels, setChannels] = useState<Array<Channel>>([]);
  const router = useRouter();
  const [selectedServerID] = useAtom(selectedServerIDAtom);
  const theme = useTheme();

  useEffect(() => {
    (async () => {
      const { data, code } = await fetchChannel(selectedServerID);
      if (code != 200) {
        router.push("/oauth");
      }
      setChannels(data);
    })();
  }, [selectedServerID]);

  return (
    <Grid2 direction={"row"} xs={2} bgcolor={theme.palette.grey[100]}>
      <Suspense fallback={<div>loading....</div>}>
        {channels.map((val) => (
          <div>
            <Button>{val.name}</Button>
          </div>
        ))}
      </Suspense>
    </Grid2>
  );
}
