import "./globals.css";

export const metadata = {
  metadataBase: new URL("https://github.com/mahmud-r-farhan/chronotrace"),
  title: "ChronoTrace — Ultra-Lightweight Privacy-First Screen Time Tracker",
  description: "Cross-platform, open-source screen time and app activity tracker consuming under 15MB RAM with zero telemetry. Available for Windows, macOS, and Linux.",
  keywords: ["screen time tracker", "open source activity tracker", "privacy first screen time", "lightweight app tracker", "chronotrace", "wails", "golang screen time"],
  authors: [{ name: "Farhan & ChronoTrace Contributors" }],
  openGraph: {
    title: "ChronoTrace — Ultra-Lightweight Privacy-First Screen Time Tracker",
    description: "Tracks active window time locally into SQLite with < 15MB RAM and 0% CPU. Zero telemetry, 100% offline.",
    url: "https://github.com/mahmud-r-farhan/chronotrace",
    siteName: "ChronoTrace",
    images: [
      {
        url: "/logo.png",
        width: 1024,
        height: 1024,
        alt: "ChronoTrace Logo",
      },
    ],
    locale: "en_US",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "ChronoTrace — Ultra-Lightweight Screen Time Tracker",
    description: "Privacy-first, open-source, cross-platform app usage tracker.",
    images: ["/logo.png"],
  },
  icons: {
    icon: "/favicon.png",
    apple: "/logo.png",
  },
};

export default function RootLayout({ children }) {
  return (
    <html lang="en" className="dark">
      <body className="antialiased min-h-screen bg-[#09090d] text-[#f5f5fc]">
        {children}
      </body>
    </html>
  );
}
