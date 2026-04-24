import "antd/dist/reset.css";
import "./globals.css";
import { AppShell } from "@/components/AppShell";

export const metadata = {
  title: "RelayOps",
  description: "Event-driven workflow platform dashboard",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <AppShell>{children}</AppShell>
      </body>
    </html>
  );
}