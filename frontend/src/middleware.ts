import { NextResponse } from "next/server";
import { NextRequest } from "next/server";

export function middleware(request: NextRequest) {
  if (!request.cookies.get("access_token") || request.nextUrl.pathname == "/") {
    return NextResponse.redirect(new URL("/home", request.url));
  }

  if (request.nextUrl.pathname.startsWith("/api")) {
    const accessTokne = request.cookies.get("access_token");
    request.headers.set("Authorization", "Bearer " + accessTokne?.value);
  }

  const response = NextResponse.next({
    request: request,
  });

  if (response.status == 401) {
    return NextResponse.redirect(new URL("/oauth", request.url));
  }

  return response;
}

export const config = {
  matcher: ["/((?!home|oauth|_next/static|_next/image|favicon.ico).*)"],
};
