"use client";

import axios from "axios";
import { useSearchParams } from "next/navigation";
import { useEffect } from "react";

export default function Oauth({ params }: { params: { provider: string } }) {
  const param = useSearchParams();

  const tokenExchange = async () => {
    try {
      console.log(params);
      const state = param.get("state");
      const code = param.get("code");
      const payload = {
        state: state,
        authorization_code: code,
        provider: params.provider,
      };
      const res = await axios.post(
        `http://localhost:4000/v1/login/oauth/callback/${params.provider}`,
        payload
      );
      console.log(res);
    } catch (err) {
      console.log(err);
    }
  };

  useEffect(() => {
    tokenExchange();
  }, []);

  return <>hello world</>;
}
