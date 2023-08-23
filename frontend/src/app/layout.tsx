import React from "react";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html>
      {/* <ThemeRegistry options={{ key: "mui" }}> */}
      <body style={{ margin: 0 }}>{children}</body>
      {/* </ThemeRegistry> */}
    </html>
  );
}
