import { NextRequest } from "next/server";
import { redirect } from "next/navigation";
import { cookies } from "next/headers";

type ApiResponse = {
  data: {
    access_token: string;
    refresh_token: string;
  };
  error: string;
  code: number;
};

type NextParams = {
  params: {
    provider: string;
  };
};

export async function GET(request: NextRequest, nextParams: NextParams) {
  const { searchParams } = new URL(request.url);
  const provider = nextParams.params.provider;
  const url = `http://localhost:4000/v1/login/oauth/callback/${provider}`;
  const payload = {
    authorization_code: searchParams.get("code"),
    state: searchParams.get("state"),
    provider: provider,
  };

  const response = await fetch(url, {
    body: JSON.stringify(payload),
    method: "POST",
  });

  if (response.status != 200) {
    redirect("/oauth");
  }

  const data: ApiResponse = await response.json();

  cookies().set("access_token", data.data.access_token, {
    httpOnly: true,
    secure: true,
  });

  cookies().set("refresh_token", data.data.refresh_token, {
    httpOnly: true,
    secure: true,
  });

  redirect("/app");
}
