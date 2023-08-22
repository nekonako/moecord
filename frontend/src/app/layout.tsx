"use client";

import React from "react";
import ThemeRegistry from "./theme_registry";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html>
      {/* <ThemeRegistry options={{ key: "mui" }}> */}
      <body>{children}</body>
      {/* </ThemeRegistry> */}
    </html>
  );
}
