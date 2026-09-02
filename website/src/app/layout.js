import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-sans",
  display: "swap",
});

export const metadata = {
  metadataBase: new URL("https://chronotrace.org"),
  title: "ChronoTrace — Privacy-First Screen Time Tracker",
  description:
    "Open-source, cross-platform screen time tracker. Under 15MB RAM, zero telemetry. Local SQLite for Windows, macOS, and Linux.",
  keywords: [
    "screen time tracker",
    "open source",
    "privacy first",
    "lightweight",
    "chronotrace",
    "activitywatch alternative",
    "rescuetime alternative",
  ],
  openGraph: {
    title: "ChronoTrace — Privacy-First Screen Time Tracker",
    description:
      "Tracks app usage locally with < 15MB RAM and zero telemetry. Free and open source.",
    url: "https://github.com/mahmud-r-farhan/chronotrace",
    siteName: "ChronoTrace",
    images: [{ url: "/logo.png", width: 1024, height: 1024, alt: "ChronoTrace" }],
    type: "website",
  },
  icons: {
    icon: "/favicon.png",
    shortcut: "/favicon.png",
    apple: "/logo.png",
  },
};

export default function RootLayout({ children }) {
  return (
    <html lang="en" className={inter.variable} suppressHydrationWarning>
      <body className="antialiased min-h-screen bg-white text-gray-900 font-sans">
        {children}
      </body>
    </html>
  );
}
