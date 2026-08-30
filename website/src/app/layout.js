import { Inter, JetBrains_Mono } from "next/font/google";
import "./globals.css";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-sans",
  display: "swap",
});

const jetbrainsMono = JetBrains_Mono({
  subsets: ["latin"],
  variable: "--font-mono",
  display: "swap",
});

export const metadata = {
  metadataBase: new URL("https://github.com/mahmud-r-farhan/chronotrace"),
  title: {
    default: "ChronoTrace — Ultra-Lightweight Privacy-First Screen Time Tracker",
    template: "%s | ChronoTrace",
  },
  description: "Cross-platform, open-source screen time and app activity tracker consuming under 15MB RAM with zero telemetry. 100% offline local SQLite storage for Windows, macOS, and Linux.",
  keywords: [
    "screen time tracker",
    "open source activity tracker",
    "privacy first screen time",
    "lightweight app tracker",
    "chronotrace",
    "wails desktop app",
    "golang screen time tracker",
    "local sqlite time tracker",
    "activitywatch alternative",
    "rescuetime alternative",
    "zero telemetry tracker",
    "windows screen time",
    "macos app usage tracker",
    "linux screen time monitor",
  ],
  authors: [{ name: "Farhan & ChronoTrace Contributors", url: "https://github.com/mahmud-r-farhan/chronotrace" }],
  creator: "Farhan",
  publisher: "ChronoTrace",
  formatDetection: {
    email: false,
    address: false,
    telephone: false,
  },
  alternates: {
    canonical: "/",
  },
  openGraph: {
    title: "ChronoTrace — Ultra-Lightweight Privacy-First Screen Time Tracker",
    description: "Tracks active window time locally into SQLite with < 15MB RAM and 0% CPU. Zero telemetry, 100% offline for Windows, macOS & Linux.",
    url: "https://github.com/mahmud-r-farhan/chronotrace",
    siteName: "ChronoTrace",
    images: [
      {
        url: "/logo.png",
        width: 1024,
        height: 1024,
        alt: "ChronoTrace Logo — Ultra-Lightweight Screen Time Tracker",
      },
    ],
    locale: "en_US",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "ChronoTrace — Ultra-Lightweight Screen Time Tracker",
    description: "Privacy-first, open-source, cross-platform app usage tracker with < 15MB RAM.",
    images: ["/logo.png"],
    creator: "@chronotrace",
  },
  icons: {
    icon: "/favicon.png",
    shortcut: "/favicon.png",
    apple: "/logo.png",
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      "max-video-preview": -1,
      "max-image-preview": "large",
      "max-snippet": -1,
    },
  },
};

const jsonLd = {
  "@context": "https://schema.org",
  "@graph": [
    {
      "@type": "SoftwareApplication",
      "name": "ChronoTrace",
      "applicationCategory": "UtilitiesApplication",
      "operatingSystem": "Windows 10, Windows 11, macOS, Linux",
      "softwareVersion": "0.1.0",
      "description": "Privacy-first, ultra-lightweight screen time and app usage tracker consuming under 15MB RAM with zero telemetry.",
      "url": "https://github.com/mahmud-r-farhan/chronotrace",
      "license": "https://opensource.org/licenses/MIT",
      "isAccessibleForFree": true,
      "offers": {
        "@type": "Offer",
        "price": "0",
        "priceCurrency": "USD",
      },
      "author": {
        "@type": "Organization",
        "name": "ChronoTrace Open Source Contributors",
        "url": "https://github.com/mahmud-r-farhan/chronotrace",
      },
    },
    {
      "@type": "FAQPage",
      "mainEntity": [
        {
          "@type": "Question",
          "name": "What is ChronoTrace?",
          "acceptedAnswer": {
            "@type": "Answer",
            "text": "ChronoTrace is an ultra-lightweight, privacy-first screen time and application usage tracker that stores all activity logs locally in SQLite without cloud telemetry.",
          },
        },
        {
          "@type": "Question",
          "name": "How much RAM and CPU does ChronoTrace consume?",
          "acceptedAnswer": {
            "@type": "Answer",
            "text": "The ChronoTrace background daemon consumes under 15MB of RAM (typically measured at ~2MB) and near 0% CPU by utilizing jittered 2–3 second window polling and batched SQLite disk writes.",
          },
        },
        {
          "@type": "Question",
          "name": "Does ChronoTrace upload my data or window titles to the cloud?",
          "acceptedAnswer": {
            "@type": "Answer",
            "text": "No. ChronoTrace has zero telemetry, zero analytics, and zero cloud dependencies. All data is saved exclusively on your local machine.",
          },
        },
        {
          "@type": "Question",
          "name": "Which operating systems are supported by ChronoTrace?",
          "acceptedAnswer": {
            "@type": "Answer",
            "text": "ChronoTrace supports Windows 10/11, macOS (Intel and Apple Silicon), and Linux (X11 and Wayland compositors).",
          },
        },
      ],
    },
  ],
};

export default function RootLayout({ children }) {
  return (
    <html lang="en" className={`dark ${inter.variable} ${jetbrainsMono.variable}`} suppressHydrationWarning>
      <head>
        <link rel="author" href="https://github.com/mahmud-r-farhan/chronotrace" />
        <link rel="help" href="/llms.txt" />
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
        />
      </head>
      <body className="antialiased min-h-screen bg-[#09090d] text-[#f5f5fc] font-sans" suppressHydrationWarning>
        {children}
      </body>
    </html>
  );
}
