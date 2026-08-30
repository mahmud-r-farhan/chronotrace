"use client";

import { useState, useEffect, useRef } from "react";
import Image from "next/image";

export default function Home() {
  const [selectedPeriod, setSelectedPeriod] = useState("today");
  const [copiedCmd, setCopiedCmd] = useState(false);
  const [userOS, setUserOS] = useState("windows");
  const [showDropdown, setShowDropdown] = useState(false);
  const dropdownRef = useRef(null);

  const installCommand = "git clone https://github.com/mahmud-r-farhan/chronotrace.git && cd chronotrace && make daemon";

  const handleCopy = () => {
    navigator.clipboard.writeText(installCommand);
    setCopiedCmd(true);
    setTimeout(() => setCopiedCmd(false), 2000);
  };

  useEffect(() => {
    if (typeof window !== "undefined") {
      const ua = window.navigator.userAgent.toLowerCase();
      if (ua.includes("mac")) {
        setUserOS("macos");
      } else if (ua.includes("linux")) {
        setUserOS("linux");
      } else {
        setUserOS("windows");
      }
    }
  }, []);

  useEffect(() => {
    function handleClickOutside(event) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setShowDropdown(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const releaseBase = "https://github.com/mahmud-r-farhan/chronotrace/releases/download/v0.1.0";
  const downloads = {
    windows: {
      name: "Windows (x64)",
      shortName: "Windows",
      icon: "🪟",
      filename: "ChronoTrace-Windows-x64-Setup.zip",
      url: `${releaseBase}/ChronoTrace-Windows-x64-Setup.zip`,
      badge: "1-Click Setup (.zip)",
      desc: "Autostart Registry Integration • Zero Terminal Popup",
      ext: ".zip",
    },
    macos: {
      name: "macOS (Universal)",
      shortName: "macOS",
      icon: "🍎",
      filename: "ChronoTrace-macOS-Universal.zip",
      url: `${releaseBase}/ChronoTrace-macOS-Universal.zip`,
      badge: "Apple Silicon & Intel",
      desc: "Native AppKit Hooks • launchd Background Service",
      ext: ".zip",
    },
    linux: {
      name: "Linux (x64)",
      shortName: "Linux",
      icon: "🐧",
      filename: "ChronoTrace-Linux-x64.tar.gz",
      url: `${releaseBase}/ChronoTrace-Linux-x64.tar.gz`,
      badge: "tar.gz + Systemd Installer",
      desc: "X11 & Wayland Support • Systemd User Unit",
      ext: ".tar.gz",
    },
  };

  const appUsageData = {
    today: [
      { name: "Visual Studio Code", duration: "4h 18m", pct: 85, color: "from-purple-500 to-indigo-500", sessions: 24 },
      { name: "Google Chrome", duration: "2h 45m", pct: 60, color: "from-blue-500 to-cyan-500", sessions: 42 },
      { name: "Windows Terminal", duration: "1h 12m", pct: 30, color: "from-emerald-500 to-teal-500", sessions: 18 },
      { name: "Slack", duration: "45m", pct: 18, color: "from-amber-500 to-rose-500", sessions: 9 },
      { name: "Figma", duration: "32m", pct: 12, color: "from-pink-500 to-rose-500", sessions: 5 },
    ],
    week: [
      { name: "Visual Studio Code", duration: "28h 40m", pct: 90, color: "from-purple-500 to-indigo-500", sessions: 142 },
      { name: "Google Chrome", duration: "16h 15m", pct: 55, color: "from-blue-500 to-cyan-500", sessions: 210 },
      { name: "Windows Terminal", duration: "8h 30m", pct: 28, color: "from-emerald-500 to-teal-500", sessions: 85 },
      { name: "Slack", duration: "5h 10m", pct: 16, color: "from-amber-500 to-rose-500", sessions: 54 },
      { name: "Figma", duration: "3h 45m", pct: 11, color: "from-pink-500 to-rose-500", sessions: 22 },
    ],
    month: [
      { name: "Visual Studio Code", duration: "112h 10m", pct: 88, color: "from-purple-500 to-indigo-500", sessions: 580 },
      { name: "Google Chrome", duration: "64h 20m", pct: 50, color: "from-blue-500 to-cyan-500", sessions: 890 },
      { name: "Windows Terminal", duration: "32h 45m", pct: 25, color: "from-emerald-500 to-teal-500", sessions: 320 },
      { name: "Slack", duration: "21h 00m", pct: 16, color: "from-amber-500 to-rose-500", sessions: 210 },
      { name: "Figma", duration: "14h 30m", pct: 10, color: "from-pink-500 to-rose-500", sessions: 95 },
    ],
  };

  const currentPlatform = downloads[userOS] || downloads.windows;

  return (
    <div className="min-h-screen bg-[#09090d] text-[#f5f5fc] selection:bg-purple-500/30 selection:text-purple-200">
      {/* Background Gradients */}
      <div className="fixed inset-0 bg-grid pointer-events-none opacity-40"></div>
      <div className="fixed top-0 left-1/2 -translate-x-1/2 w-[1000px] h-[550px] bg-radial-glow pointer-events-none"></div>

      {/* Navigation */}
      <nav className="sticky top-0 z-50 backdrop-blur-xl bg-[#09090d]/80 border-b border-white/10 px-6 py-4">
        <div className="max-w-7xl mx-auto flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="relative w-9 h-9 rounded-xl overflow-hidden shadow-[0_0_20px_rgba(139,92,246,0.4)] border border-purple-500/30">
              <Image src="/logo.jpg" alt="ChronoTrace Logo" width={36} height={36} className="w-full h-full object-cover" priority />
            </div>
            <div>
              <span className="font-bold text-lg tracking-tight bg-gradient-to-r from-white via-purple-200 to-cyan-300 bg-clip-text text-transparent">
                ChronoTrace
              </span>
              <span className="ml-2 text-[10px] uppercase font-semibold tracking-wider px-2 py-0.5 rounded-full bg-purple-500/10 text-purple-400 border border-purple-500/20">
                v0.1.0 Open Source
              </span>
            </div>
          </div>

          <div className="hidden md:flex items-center gap-8 text-sm font-medium text-[#9494b8]">
            <a href="#demo" className="hover:text-white transition-colors">Interactive Demo</a>
            <a href="#features" className="hover:text-white transition-colors">Features</a>
            <a href="#architecture" className="hover:text-white transition-colors">Architecture</a>
            <a href="#compare" className="hover:text-white transition-colors">Comparison</a>
            <a href="#api" className="hover:text-white transition-colors">REST API</a>
            <a href="#faq" className="hover:text-white transition-colors">FAQ</a>
          </div>

          <div className="flex items-center gap-3">
            <a
              href="https://github.com/mahmud-r-farhan/chronotrace"
              target="_blank"
              rel="noreferrer"
              className="flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-medium bg-white/5 hover:bg-white/10 border border-white/10 transition-all hover:scale-[1.02] active:scale-[0.98]"
            >
              <svg className="w-4 h-4 fill-current" viewBox="0 0 24 24">
                <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z" />
              </svg>
              <span>GitHub</span>
            </a>
            <a
              href={currentPlatform.url}
              className="hidden sm:inline-flex px-4 py-2 rounded-xl text-sm font-semibold bg-gradient-to-r from-purple-500 to-cyan-500 hover:from-purple-400 hover:to-cyan-400 text-white shadow-[0_0_20px_rgba(139,92,246,0.3)] transition-all hover:scale-[1.02] active:scale-[0.98]"
            >
              Download
            </a>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="relative pt-20 pb-16 px-6 max-w-7xl mx-auto text-center">
        {/* Status Pill */}
        <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-white/5 border border-white/10 mb-8 backdrop-blur-md">
          <span className="w-2 h-2 rounded-full bg-emerald-400 shadow-[0_0_8px_#34d399] animate-pulse"></span>
          <span className="text-xs font-medium text-[#9494b8]">Ultra-Lightweight Background Daemon &bull; &lt; 15MB RAM</span>
        </div>

        {/* Main Heading */}
        <h1 className="text-4xl sm:text-6xl md:text-7xl font-extrabold tracking-tight max-w-4xl mx-auto leading-[1.1] mb-6">
          Track Your Screen Time.{" "}
          <span className="text-gradient">Zero Telemetry.</span>
        </h1>

        <p className="text-base sm:text-lg md:text-xl text-[#9494b8] max-w-2xl mx-auto mb-10 leading-relaxed">
          ChronoTrace is a privacy-first, cross-platform app usage tracker. A headless background daemon writes silently to your local SQLite database with ~0% CPU and &lt; 15MB RAM.
        </p>

        {/* Dynamic OS Download CTA with Split Dropdown */}
        <div className="relative inline-flex flex-col items-center gap-3 mb-10" ref={dropdownRef}>
          <div className="inline-flex items-stretch rounded-2xl shadow-[0_0_35px_rgba(139,92,246,0.35)] bg-gradient-to-r from-purple-600 via-indigo-600 to-cyan-600 p-[1px]">
            {/* Primary Detected OS Action */}
            <a
              href={currentPlatform.url}
              className="flex items-center gap-3 px-6 sm:px-8 py-4 rounded-l-2xl font-semibold text-sm sm:text-base bg-gradient-to-r from-purple-500 via-indigo-500 to-cyan-500 hover:from-purple-400 hover:to-cyan-400 text-white transition-all hover:brightness-110 active:scale-[0.99]"
            >
              <span className="text-xl">{currentPlatform.icon}</span>
              <span className="text-left">
                <span className="block font-bold">Download for {currentPlatform.shortName}</span>
                <span className="block text-[11px] opacity-80 font-normal">{currentPlatform.badge}</span>
              </span>
            </a>

            {/* Platform Selector Chevron */}
            <button
              onClick={() => setShowDropdown(!showDropdown)}
              className="px-4 rounded-r-2xl bg-[#0e0e18] hover:bg-[#181828] text-white/80 hover:text-white border-l border-white/10 transition-colors flex items-center justify-center cursor-pointer"
              title="Select another operating system"
              aria-label="Select platform"
            >
              <svg className={`w-5 h-5 transition-transform duration-200 ${showDropdown ? "rotate-180" : ""}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
              </svg>
            </button>
          </div>

          {/* Dropdown Menu */}
          {showDropdown && (
            <div className="absolute top-full mt-3 w-80 rounded-2xl glass border border-white/15 bg-[#12121f]/95 backdrop-blur-2xl shadow-2xl p-2 z-50 text-left animate-in fade-in slide-in-from-top-2 duration-150">
              <div className="px-3 py-2 text-[11px] font-bold uppercase tracking-wider text-[#5e5e7a] border-b border-white/10">
                Choose Operating System
              </div>
              {Object.entries(downloads).map(([key, item]) => (
                <a
                  key={key}
                  href={item.url}
                  onClick={() => setShowDropdown(false)}
                  className={`flex items-center justify-between p-3 rounded-xl transition-all ${
                    key === userOS ? "bg-purple-500/15 border border-purple-500/30 text-white" : "hover:bg-white/5 text-[#9494b8] hover:text-white"
                  }`}
                >
                  <div className="flex items-center gap-3">
                    <span className="text-xl">{item.icon}</span>
                    <div>
                      <div className="font-semibold text-sm text-white flex items-center gap-2">
                        {item.name}
                        {key === userOS && <span className="text-[10px] px-1.5 py-0.2 rounded bg-purple-500/20 text-purple-300">Detected</span>}
                      </div>
                      <div className="text-xs text-[#9494b8]">{item.badge}</div>
                    </div>
                  </div>
                  <svg className="w-4 h-4 text-[#5e5e7a]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                  </svg>
                </a>
              ))}
              <div className="pt-2 mt-1 border-t border-white/10">
                <a
                  href="https://github.com/mahmud-r-farhan/chronotrace/releases"
                  target="_blank"
                  rel="noreferrer"
                  onClick={() => setShowDropdown(false)}
                  className="flex items-center justify-between p-2.5 rounded-xl hover:bg-white/5 text-xs text-purple-400 hover:text-purple-300 transition-colors"
                >
                  <span>View All Releases &amp; Standalone Daemons</span>
                  <span>↗</span>
                </a>
              </div>
            </div>
          )}

          <div className="text-xs text-[#5e5e7a]">
            Free &amp; Open Source forever &bull; MIT License &bull; Zero Telemetry
          </div>
        </div>

        {/* Copyable Quick Install */}
        <div className="inline-flex items-center gap-3 px-4 py-2.5 rounded-xl bg-[#151522] border border-white/10 font-mono text-xs text-[#9494b8] max-w-xl mx-auto">
          <span className="text-purple-400">$</span>
          <span className="truncate">{installCommand}</span>
          <button
            onClick={handleCopy}
            className="ml-2 px-2.5 py-1 rounded-md bg-white/5 hover:bg-white/10 text-white transition-colors"
          >
            {copiedCmd ? "✓ Copied" : "Copy"}
          </button>
        </div>

        {/* Live Metrics Row */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 max-w-4xl mx-auto mt-14 text-left">
          <div className="p-5 rounded-2xl glass">
            <span className="text-xs uppercase font-bold tracking-wider text-purple-400">RAM Consumption</span>
            <div className="text-2xl sm:text-3xl font-extrabold text-white mt-1">&lt; 15 MB</div>
            <span className="text-xs text-[#9494b8]">~2MB typical background idle</span>
          </div>

          <div className="p-5 rounded-2xl glass">
            <span className="text-xs uppercase font-bold tracking-wider text-cyan-400">CPU Overhead</span>
            <div className="text-2xl sm:text-3xl font-extrabold text-white mt-1">~0.0%</div>
            <span className="text-xs text-[#9494b8]">Jittered OS hook polling</span>
          </div>

          <div className="p-5 rounded-2xl glass">
            <span className="text-xs uppercase font-bold tracking-wider text-emerald-400">Privacy Policy</span>
            <div className="text-2xl sm:text-3xl font-extrabold text-white mt-1">100% Local</div>
            <span className="text-xs text-[#9494b8]">Zero external HTTP requests</span>
          </div>

          <div className="p-5 rounded-2xl glass">
            <span className="text-xs uppercase font-bold tracking-wider text-amber-400">Disk Storage</span>
            <div className="text-2xl sm:text-3xl font-extrabold text-white mt-1">Pure SQLite</div>
            <span className="text-xs text-[#9494b8]">WAL mode in user directory</span>
          </div>
        </div>
      </section>

      {/* Interactive Demo Section */}
      <section id="demo" className="py-16 px-6 max-w-7xl mx-auto border-t border-white/10">
        <div className="text-center mb-12">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-purple-500/10 text-purple-400 text-xs font-bold uppercase tracking-wider mb-3 border border-purple-500/20">
            Interactive Live Preview
          </div>
          <h2 className="text-3xl sm:text-5xl font-extrabold tracking-tight mb-4">
            See ChronoTrace In Action
          </h2>
          <p className="text-[#9494b8] max-w-xl mx-auto text-sm sm:text-base">
            This dashboard runs locally on your machine via Wails v2 without continuous memory footprint.
          </p>
        </div>

        {/* Dashboard Frame Container */}
        <div className="rounded-3xl glass border border-white/10 p-6 md:p-8 max-w-5xl mx-auto shadow-2xl">
          {/* Header Controls */}
          <div className="flex flex-col sm:flex-row items-center justify-between gap-4 mb-8 border-b border-white/10 pb-6">
            <div>
              <h3 className="text-xl font-bold text-white">Application Usage Summary</h3>
              <p className="text-xs text-[#9494b8] mt-1">Aggregated window active time &bull; 127.0.0.1:42069</p>
            </div>

            {/* Time Filter Tabs */}
            <div className="flex items-center gap-1 bg-[#0e0e17] p-1.5 rounded-xl border border-white/10">
              {["today", "week", "month"].map((period) => (
                <button
                  key={period}
                  onClick={() => setSelectedPeriod(period)}
                  className={`px-4 py-1.5 rounded-lg text-xs font-semibold capitalize transition-all ${
                    selectedPeriod === period
                      ? "bg-purple-500 text-white shadow-[0_0_12px_rgba(139,92,246,0.4)]"
                      : "text-[#9494b8] hover:text-white"
                  }`}
                >
                  {period}
                </button>
              ))}
            </div>
          </div>

          {/* Usage Chart Horizontal Bars */}
          <div className="space-y-4 mb-10">
            {appUsageData[selectedPeriod].map((app) => (
              <div key={app.name} className="p-4 rounded-xl bg-white/[0.02] border border-white/5 hover:border-white/15 transition-all">
                <div className="flex items-center justify-between text-sm mb-2">
                  <div className="flex items-center gap-3">
                    <span className="font-semibold text-white">{app.name}</span>
                    <span className="text-xs text-[#5e5e7a]">{app.sessions} sessions</span>
                  </div>
                  <span className="font-mono text-xs font-bold text-purple-300">{app.duration}</span>
                </div>
                <div className="w-full h-2.5 rounded-full bg-white/5 overflow-hidden">
                  <div
                    className={`h-full rounded-full bg-gradient-to-r ${app.color} transition-all duration-700`}
                    style={{ width: `${app.pct}%` }}
                  ></div>
                </div>
              </div>
            ))}
          </div>

          {/* Hourly Heatmap Preview */}
          <div>
            <h4 className="text-xs uppercase font-bold tracking-wider text-[#9494b8] mb-3">
              Hourly Activity Heatmap (24-Hour Timeline)
            </h4>
            <div className="grid grid-cols-12 gap-2 text-center">
              {Array.from({ length: 24 }).map((_, i) => {
                const active = i >= 8 && i <= 20;
                const opacity = active ? 0.3 + (i % 5) * 0.15 : 0.05;
                return (
                  <div key={i} className="flex flex-col items-center gap-1.5">
                    <div
                      className="w-full aspect-square rounded-lg border border-white/10 transition-all hover:scale-110"
                      style={{
                        backgroundColor: `rgba(139, 92, 246, ${opacity})`,
                        boxShadow: active ? "0 0 10px rgba(139, 92, 246, 0.2)" : "none",
                      }}
                      title={`${i}:00 — ${active ? (i % 4) + 1 + "h active" : "Idle"}`}
                    ></div>
                    <span className="text-[10px] font-mono text-[#5e5e7a]">{i}</span>
                  </div>
                );
              })}
            </div>
          </div>
        </div>
      </section>

      {/* Features Grid */}
      <section id="features" className="py-20 px-6 max-w-7xl mx-auto border-t border-white/10">
        <div className="text-center mb-16">
          <h2 className="text-3xl sm:text-5xl font-extrabold tracking-tight mb-4">
            Why ChronoTrace?
          </h2>
          <p className="text-[#9494b8] max-w-xl mx-auto text-sm sm:text-base">
            Built for developers, professionals, and students who care about their device performance and privacy.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="p-8 rounded-3xl glass border border-white/10 hover:border-purple-500/40 transition-all hover:-translate-y-1">
            <div className="text-3xl mb-4">⚡</div>
            <h3 className="text-xl font-bold text-white mb-2">Zero CGO &bull; Tiny Binary</h3>
            <p className="text-sm text-[#9494b8] leading-relaxed">
              Compiled into pure native Go binaries with no CGO dependencies. Consumes under 15MB RAM and 0% CPU at idle.
            </p>
          </div>

          <div className="p-8 rounded-3xl glass border border-white/10 hover:border-cyan-500/40 transition-all hover:-translate-y-1">
            <div className="text-3xl mb-4">🔒</div>
            <h3 className="text-xl font-bold text-white mb-2">100% Offline &amp; Private</h3>
            <p className="text-sm text-[#9494b8] leading-relaxed">
              No cloud accounts, no tracking scripts, and no telemetry. Your logs remain on your disk in a local SQLite file.
            </p>
          </div>

          <div className="p-8 rounded-3xl glass border border-white/10 hover:border-emerald-500/40 transition-all hover:-translate-y-1">
            <div className="text-3xl mb-4">🔌</div>
            <h3 className="text-xl font-bold text-white mb-2">Local REST API</h3>
            <p className="text-sm text-[#9494b8] leading-relaxed">
              Query your data anytime on <code>http://127.0.0.1:42069</code> with curl, Python, or automated scripts.
            </p>
          </div>
        </div>
      </section>

      {/* Architecture Section */}
      <section id="architecture" className="py-20 px-6 max-w-7xl mx-auto border-t border-white/10">
        <div className="text-center mb-16">
          <h2 className="text-3xl sm:text-5xl font-extrabold tracking-tight mb-4">
            Decoupled Architecture
          </h2>
          <p className="text-[#9494b8] max-w-xl mx-auto text-sm sm:text-base">
            The daemon runs quietly in the background. The desktop dashboard is only opened when you want to view reports.
          </p>
        </div>

        <div className="p-8 rounded-3xl glass border border-white/10 max-w-4xl mx-auto font-mono text-xs sm:text-sm text-[#9494b8] overflow-x-auto">
          <pre className="text-center leading-loose">
{`┌──────────────────────────────────────┐        ┌──────────────────────────────────────┐
│       Headless Background Daemon     │  ───►  │           Local SQLite DB            │
│         (Under 15MB RAM, 0% CPU)     │        │      (~/.local/share, %APPDATA%)     │
└──────────────────────────────────────┘        └──────────────────┬───────────────────┘
                    ▲                                              │
                    │ (Auto-Spawns / OS Start)                     ▼
┌──────────────────────────────────────┐        ┌──────────────────────────────────────┐
│        ChronoTrace GUI (Wails)       │  ◄───  │           Local REST API             │
│        (Optional Desktop Window)     │        │         (127.0.0.1:42069)            │
└──────────────────────────────────────┘        └──────────────────────────────────────┘`}
          </pre>
        </div>
      </section>

      {/* Comparison Table */}
      <section id="compare" className="py-20 px-6 max-w-5xl mx-auto border-t border-white/10">
        <div className="text-center mb-16">
          <h2 className="text-3xl sm:text-5xl font-extrabold tracking-tight mb-4">
            How ChronoTrace Compares
          </h2>
        </div>

        <div className="overflow-x-auto rounded-2xl glass border border-white/10">
          <table className="w-full text-left text-sm">
            <thead className="border-b border-white/10 bg-white/5">
              <tr>
                <th className="p-4 font-bold text-white">Feature</th>
                <th className="p-4 font-bold text-purple-400">ChronoTrace</th>
                <th className="p-4 font-medium text-[#9494b8]">ActivityWatch</th>
                <th className="p-4 font-medium text-[#9494b8]">RescueTime</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5 text-xs sm:text-sm">
              <tr>
                <td className="p-4 font-medium text-white">RAM Usage</td>
                <td className="p-4 font-bold text-emerald-400">&lt; 15 MB</td>
                <td className="p-4 text-[#9494b8]">100–250 MB</td>
                <td className="p-4 text-[#9494b8]">80–180 MB</td>
              </tr>
              <tr>
                <td className="p-4 font-medium text-white">Telemetry &amp; Privacy</td>
                <td className="p-4 font-bold text-emerald-400">0% Telemetry (100% Local)</td>
                <td className="p-4 text-[#9494b8]">Local first</td>
                <td className="p-4 text-red-400">Cloud dependent</td>
              </tr>
              <tr>
                <td className="p-4 font-medium text-white">Headless Tracking</td>
                <td className="p-4 font-bold text-emerald-400">Yes (Daemon only)</td>
                <td className="p-4 text-[#9494b8]">Requires multiple processes</td>
                <td className="p-4 text-[#9494b8]">No</td>
              </tr>
              <tr>
                <td className="p-4 font-medium text-white">Local REST API</td>
                <td className="p-4 font-bold text-emerald-400">Yes (127.0.0.1:42069)</td>
                <td className="p-4 text-[#9494b8]">Yes</td>
                <td className="p-4 text-red-400">No</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      {/* Frequently Asked Questions (FAQ) */}
      <section id="faq" className="py-20 px-6 max-w-5xl mx-auto border-t border-white/10">
        <div className="text-center mb-16">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-purple-500/10 text-purple-400 text-xs font-bold uppercase tracking-wider mb-4 border border-purple-500/20">
            Got Questions?
          </div>
          <h2 className="text-3xl sm:text-5xl font-extrabold tracking-tight mb-4">
            Frequently Asked Questions
          </h2>
        </div>

        <div className="space-y-4">
          <div className="p-6 rounded-2xl glass border border-white/10 hover:border-purple-500/30 transition-colors">
            <h3 className="text-base sm:text-lg font-bold text-white mb-2">How does ChronoTrace achieve &lt; 15MB RAM and ~0% CPU?</h3>
            <p className="text-sm text-[#9494b8] leading-relaxed">
              Unlike bloated Electron or Python alternatives, ChronoTrace&apos;s background daemon is compiled into a lightweight native Go binary with pure Go SQLite (zero CGO). It uses jittered 2–3 second OS foreground window hooks and in-memory batching that flushes to disk only every 45 seconds.
            </p>
          </div>

          <div className="p-6 rounded-2xl glass border border-white/10 hover:border-purple-500/30 transition-colors">
            <h3 className="text-base sm:text-lg font-bold text-white mb-2">Does ChronoTrace send any data to external servers or the cloud?</h3>
            <p className="text-sm text-[#9494b8] leading-relaxed">
              Never. ChronoTrace has zero telemetry, zero analytics scripts, and zero cloud accounts. All recorded data stays strictly on your machine in a local SQLite file stored in your operating system&apos;s user directory.
            </p>
          </div>

          <div className="p-6 rounded-2xl glass border border-white/10 hover:border-purple-500/30 transition-colors">
            <h3 className="text-base sm:text-lg font-bold text-white mb-2">How does the 1-Click install work?</h3>
            <p className="text-sm text-[#9494b8] leading-relaxed">
              Download the package for Windows, Linux, or macOS. Extract the zip file and run the installer script (`install-windows.bat`, `install-linux.sh`, or `install-macos.sh`). It copies the files, registers the background daemon in your OS autostart on boot, and opens the GUI.
            </p>
          </div>

          <div className="p-6 rounded-2xl glass border border-white/10 hover:border-purple-500/30 transition-colors">
            <h3 className="text-base sm:text-lg font-bold text-white mb-2">Can I query the data programmatically with scripts?</h3>
            <p className="text-sm text-[#9494b8] leading-relaxed">
              Yes! The daemon exposes a local JSON REST API on <code>http://127.0.0.1:42069</code>. You can query your hourly timelines, daily summaries, and per-app usage directly using curl or Python.
            </p>
          </div>
        </div>
      </section>

      {/* Downloads / Release Section */}
      <section id="downloads" className="py-20 px-6 max-w-7xl mx-auto border-t border-white/10 text-center">
        <h2 className="text-3xl sm:text-5xl font-extrabold tracking-tight mb-4">
          Download ChronoTrace v0.1.0
        </h2>
        <p className="text-[#9494b8] max-w-xl mx-auto mb-12 text-sm sm:text-base">
          Free and open-source forever. Download 1-click packages with full autostart support for your platform.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-4xl mx-auto text-left">
          {/* Windows */}
          <div className="p-6 rounded-2xl glass hover:border-purple-500/40 transition-all flex flex-col justify-between">
            <div>
              <div className="text-2xl mb-2">🪟</div>
              <h3 className="text-lg font-bold text-white">Windows (x64)</h3>
              <p className="text-xs text-[#9494b8] mt-1 mb-4">
                1-Click Setup &bull; HKCU Autostart Registry &bull; Zero popup
              </p>
            </div>
            <a
              href={downloads.windows.url}
              className="w-full py-3 rounded-xl bg-purple-600/30 hover:bg-purple-600/50 border border-purple-500/40 text-center text-xs font-semibold text-white transition-colors"
            >
              Download Windows .zip
            </a>
          </div>

          {/* macOS */}
          <div className="p-6 rounded-2xl glass hover:border-purple-500/40 transition-all flex flex-col justify-between">
            <div>
              <div className="text-2xl mb-2">🍎</div>
              <h3 className="text-lg font-bold text-white">macOS (Universal)</h3>
              <p className="text-xs text-[#9494b8] mt-1 mb-4">
                Apple Silicon &amp; Intel &bull; Native launchd LaunchAgent
              </p>
            </div>
            <a
              href={downloads.macos.url}
              className="w-full py-3 rounded-xl bg-purple-600/30 hover:bg-purple-600/50 border border-purple-500/40 text-center text-xs font-semibold text-white transition-colors"
            >
              Download macOS .zip
            </a>
          </div>

          {/* Linux */}
          <div className="p-6 rounded-2xl glass hover:border-purple-500/40 transition-all flex flex-col justify-between">
            <div>
              <div className="text-2xl mb-2">🐧</div>
              <h3 className="text-lg font-bold text-white">Linux (x64)</h3>
              <p className="text-xs text-[#9494b8] mt-1 mb-4">
                X11 &amp; Wayland &bull; Systemd user service installer
              </p>
            </div>
            <a
              href={downloads.linux.url}
              className="w-full py-3 rounded-xl bg-purple-600/30 hover:bg-purple-600/50 border border-purple-500/40 text-center text-xs font-semibold text-white transition-colors"
            >
              Download Linux .tar.gz
            </a>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-white/10 py-12 px-6 bg-[#07070a]">
        <div className="max-w-7xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-6">
          <div className="flex items-center gap-3">
            <div className="relative w-7 h-7 rounded-lg overflow-hidden border border-purple-500/30">
              <Image src="/logo.jpg" alt="ChronoTrace Logo" width={28} height={28} className="w-full h-full object-cover" />
            </div>
            <span className="font-bold text-sm bg-gradient-to-r from-white to-purple-300 bg-clip-text text-transparent">
              ChronoTrace
            </span>
            <span className="text-xs text-[#5e5e7a]">
              &copy; 2026 ChronoTrace Contributors. MIT License.
            </span>
          </div>

          <div className="flex items-center gap-6 text-xs text-[#9494b8]">
            <a href="https://github.com/mahmud-r-farhan/chronotrace" className="hover:text-white transition-colors">GitHub</a>
            <a href="https://github.com/mahmud-r-farhan/chronotrace/blob/main/LICENSE" className="hover:text-white transition-colors">MIT License</a>
            <a href="https://github.com/mahmud-r-farhan/chronotrace/releases" className="hover:text-white transition-colors">Releases</a>
            <a href="/llms.txt" className="hover:text-white transition-colors">llms.txt</a>
          </div>
        </div>
      </footer>
    </div>
  );
}
