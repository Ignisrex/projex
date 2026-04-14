import "./globals.css";

export const metadata = {
  title: "Project Tracker",
  description: "Basic project and task management scaffold"
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
