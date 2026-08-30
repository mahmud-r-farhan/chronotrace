"use client";

import { useState } from "react";
import Image from "next/image";

export default function Home() {
  const [selectedPeriod, setSelectedPeriod] = useState("today");
  const [copiedCmd, setCopiedCmd] = useState(false);
  const [activeTab, setActiveTab] = useState("overview");

  const installCommand = "git clone https://github.com/mahmud-r-farhan/chronotrace.git && cd chronotrace && make daemon";

  const handleCopy = () => {
    navigator.clipboard.writeText(installCommand);
    setCopiedCmd(true);
    setTimeout(() => setCopiedCmd(false), 2000);
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

  return (
    <div className="min-h-screen bg-[#09090d] text-[#f5f5fc] selection:bg-purple-500/30 selection:text-purple-200">
      {/* Background Gradients */}
      <div className="fixed inset-0 bg-grid pointer-events-none opacity-40"></div>
      <div className="fixed top-0 left-1/2 -translate-x-1/2 w-[1000px] h-[600px] bg-radial-glow pointer-events-none"></div>

      {/* Navigation */}
      <nav className="sticky top-0 z-50 backdrop-blur-xl bg-[#09090d]/80 border-b border-white/10 px-6 py-4">
        <div className="max-w-7xl mx-auto flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="relative w-10 h-10 rounded-xl overflow-hidden shadow-[0_0_20px_rgba(139,92,246,0.4)] border border-purple-500/30">
              <Image src="/logo.png" alt="ChronoTrace Logo" width={40} height={40} className="w-full h-full object-cover" priority />
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
            <a href="#features" className="hover:text-white transition-colors">Features</a>
            <a href="#demo" className="hover:text-white transition-colors">Interactive Demo</a>
            <a href="#architecture" className="hover:text-white transition-colors">Architecture</a>
            <a href="#compare" className="hover:text-white transition-colors">Comparison</a>
            <a href="#api" className="hover:text-white transition-colors">REST API</a>
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
              href="#downloads"
              className="px-4 py-2 rounded-xl text-sm font-semibold bg-gradient-to-r from-purple-500 to-cyan-500 hover:from-purple-400 hover:to-cyan-400 text-white shadow-[0_0_20px_rgba(139,92,246,0.3)] transition-all hover:scale-[1.02] active:scale-[0.98]"
            >
              Download Free
            </a>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="relative pt-24 pb-20 px-6 max-w-7xl mx-auto text-center">
        {/* Status Pill */}
        <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-white/5 border border-white/10 mb-8 backdrop-blur-md">
          <span className="w-2 h-2 rounded-full bg-emerald-400 shadow-[0_0_8px_#34d399] animate-pulse"></span>
          <span className="text-xs font-medium text-[#9494b8]">Ultra-Lightweight Background Daemon &bull; &lt; 15MB RAM</span>
        </div>

        {/* Main Heading */}
        <h1 className="text-5xl md:text-7xl font-extrabold tracking-tight max-w-4xl mx-auto leading-[1.1] mb-6">
          Track Your Screen Time.{" "}
          <span className="text-gradient">Zero Telemetry.</span>
        </h1>

        <p className="text-lg md:text-xl text-[#9494b8] max-w-2xl mx-auto mb-10 leading-relaxed">
          ChronoTrace is a privacy-first, cross-platform app usage tracker. A headless background daemon writes silently to your local SQLite database while consuming virtually zero CPU and RAM.
        </p>

        {/* CTA Buttons */}
        <div className="flex flex-wrap items-center justify-center gap-4 mb-12">
          <a
            href="https://github.com/mahmud-r-farhan/chronotrace/releases"
            className="flex items-center gap-3 px-8 py-4 rounded-2xl font-semibold text-base bg-gradient-to-r from-purple-500 via-indigo-500 to-cyan-500 hover:from-purple-400 hover:to-cyan-400 text-white shadow-[0_0_30px_rgba(139,92,246,0.4)] transition-all hover:scale-105 active:scale-95"
          >
            <svg className="w-5 h-5 fill-current" viewBox="0 0 24 24">
              <path d="M19.35 10.04C18.67 6.59 15.64 4 12 4 9.11 4 6.6 5.64 5.35 8.04 2.34 8.36 0 10.91 0 14c0 3.31 2.69 6 6 6h13c2.76 0 5-2.24 5-5 0-2.64-2.05-4.78-4.65-4.96zM17 13l-5 5-5-5h3V9h4v4h3z" />
            </svg>
            <span>Download for Windows, Mac &amp; Linux</span>
          </a>

          <a
            href="#demo"
            className="px-8 py-4 rounded-2xl font-semibold text-base bg-white/5 hover:bg-white/10 border border-white/10 transition-all hover:scale-105 active:scale-95 text-white"
          >
            Explore Interactive Demo &darr;
          </a>
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
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 max-w-4xl mx-auto mt-16 text-left">
          <div className="p-5 rounded-2xl glass">
            <span className="text-xs uppercase font-bold tracking-wider text-purple-400">RAM Consumption</span>
            <div className="text-3xl font-extrabold text-white mt-1">~2.0 MB</div>
            <span className="text-xs text-[#5e5e7a]">Measured working set</span>
          </div>
          <div className="p-5 rounded-2xl glass">
            <span className="text-xs uppercase font-bold tracking-wider text-cyan-400">CPU Usage</span>
            <div className="text-3xl font-extrabold text-white mt-1">~0.0%</div>
            <span className="text-xs text-[#5e5e7a]">Jittered 2-3s polling</span>
          </div>
          <div className="p-5 rounded-2xl glass">
            <span className="text-xs uppercase font-bold tracking-wider text-emerald-400">Telemetry &amp; Ads</span>
            <div className="text-3xl font-extrabold text-white mt-1">ZERO</div>
            <span className="text-xs text-[#5e5e7a]">100% Offline SQLite</span>
          </div>
          <div className="p-5 rounded-2xl glass">
            <span className="text-xs uppercase font-bold tracking-wider text-amber-400">Architecture</span>
            <div className="text-3xl font-extrabold text-white mt-1">Decoupled</div>
            <span className="text-xs text-[#5e5e7a]">Daemon + Optional UI</span>
          </div>
        </div>
      </section>

      {/* Interactive Mockup Demo Section */}
      <section id="demo" className="py-20 px-6 max-w-7xl mx-auto">
        <div className="text-center mb-12">
          <h2 className="text-3xl md:text-4xl font-extrabold tracking-tight mb-4">
            Interactive Dashboard Preview
          </h2>
          <p className="text-[#9494b8] max-w-xl mx-auto">
            Experience the responsive Wails desktop UI running locally on your device.
          </p>
        </div>

        {/* Window Container */}
        <div className="rounded-3xl glass border border-white/10 shadow-[0_0_50px_rgba(0,0,0,0.8)] overflow-hidden max-w-5xl mx-auto">
          {/* Window Title Bar */}
          <div className="bg-[#12121a] px-5 py-3.5 border-b border-white/10 flex items-center justify-between">
            <div className="flex items-center gap-2">
              <span className="w-3 h-3 rounded-full bg-rose-500/80"></span>
              <span className="w-3 h-3 rounded-full bg-amber-500/80"></span>
              <span className="w-3 h-3 rounded-full bg-emerald-500/80"></span>
              <span className="ml-3 text-xs font-mono text-[#5e5e7a]">ChronoTrace — Local Desktop Client</span>
            </div>
            <div className="flex items-center gap-2 text-xs font-mono text-emerald-400 bg-emerald-500/10 px-2.5 py-1 rounded-full border border-emerald-500/20">
              <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-ping"></span>
              <span>Daemon Connected (127.0.0.1:42069)</span>
            </div>
          </div>

          {/* Window Body */}
          <div className="p-6 md:p-8 bg-[#0d0d12]">
            {/* Dashboard Header */}
            <div className="flex flex-wrap items-center justify-between gap-4 mb-8 pb-6 border-b border-white/10">
              <div>
                <h3 className="text-2xl font-bold text-white">Daily Screen Time Summary</h3>
                <p className="text-xs text-[#5e5e7a] mt-1">Auto-aggregated from local SQLite database</p>
              </div>

              {/* Time Period Tabs */}
              <div className="flex items-center gap-1 p-1 rounded-xl bg-white/5 border border-white/10">
                {["today", "week", "month"].map((p) => (
                  <button
                    key={p}
                    onClick={() => setSelectedPeriod(p)}
                    className={`px-4 py-1.5 rounded-lg text-xs font-semibold capitalize transition-all ${
                      selectedPeriod === p
                        ? "bg-gradient-to-r from-purple-500 to-cyan-500 text-white shadow-[0_0_12px_rgba(139,92,246,0.3)]"
                        : "text-[#9494b8] hover:text-white"
                    }`}
                  >
                    {p}
                  </button>
                ))}
              </div>
            </div>

            {/* Stats Row */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
              <div className="p-5 rounded-2xl bg-white/5 border border-white/5">
                <span className="text-xs text-[#9494b8] uppercase font-bold tracking-wider">Total Active Time</span>
                <div className="text-2xl font-bold text-white mt-1">
                  {selectedPeriod === "today" ? "9h 22m" : selectedPeriod === "week" ? "62h 10m" : "244h 30m"}
                </div>
              </div>
              <div className="p-5 rounded-2xl bg-white/5 border border-white/5">
                <span className="text-xs text-[#9494b8] uppercase font-bold tracking-wider">Top Focused Application</span>
                <div className="text-2xl font-bold text-purple-400 mt-1">Visual Studio Code</div>
              </div>
              <div className="p-5 rounded-2xl bg-white/5 border border-white/5">
                <span className="text-xs text-[#9494b8] uppercase font-bold tracking-wider">Tracked Applications</span>
                <div className="text-2xl font-bold text-cyan-400 mt-1">
                  {selectedPeriod === "today" ? "14 Apps" : selectedPeriod === "week" ? "38 Apps" : "64 Apps"}
                </div>
              </div>
            </div>

            {/* App Usage Bars */}
            <div className="space-y-4">
              <div className="flex items-center justify-between text-xs font-bold text-[#5e5e7a] uppercase tracking-wider px-2">
                <span>Application Name</span>
                <span>Time Spent</span>
              </div>

              {appUsageData[selectedPeriod].map((app, idx) => (
                <div key={idx} className="p-4 rounded-xl bg-white/[0.02] border border-white/5 hover:border-white/10 transition-colors">
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-3">
                      <span className="text-xs font-mono font-bold text-[#5e5e7a]">#{idx + 1}</span>
                      <span className="font-semibold text-sm text-white">{app.name}</span>
                      <span className="text-[10px] text-[#5e5e7a]">({app.sessions} sessions)</span>
                    </div>
                    <span className="text-sm font-mono font-bold text-purple-300">{app.duration}</span>
                  </div>
                  <div className="h-2 rounded-full bg-white/5 overflow-hidden">
                    <div
                      className={`h-full rounded-full bg-gradient-to-r ${app.color} transition-all duration-700`}
                      style={{ width: `${app.pct}%` }}
                    ></div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* Feature Highlights Grid */}
      <section id="features" className="py-20 px-6 max-w-7xl mx-auto border-t border-white/10">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-5xl font-extrabold tracking-tight mb-4">
            Engineered for Precision &amp; Silence
          </h2>
          <p className="text-[#9494b8] max-w-2xl mx-auto">
            Unlike bloated telemetry tools that consume hundreds of megabytes of RAM and upload your private window titles to third-party clouds, ChronoTrace stays 100% on your device.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div className="p-8 rounded-3xl glass hover:border-purple-500/30 transition-all group">
            <div className="w-12 h-12 rounded-2xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-400 mb-6 group-hover:scale-110 transition-transform">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h3 className="text-xl font-bold mb-2 text-white">Ultra-Low Resource Daemon</h3>
            <p className="text-sm text-[#9494b8] leading-relaxed">
              Consuming &lt; 15MB RAM and ~0% CPU, the background service polls your active window using native Win32/AppKit/X11 hooks every 2–3s and flushes batch writes to SQLite.
            </p>
          </div>

          <div className="p-8 rounded-3xl glass hover:border-cyan-500/30 transition-all group">
            <div className="w-12 h-12 rounded-2xl bg-cyan-500/10 border border-cyan-500/20 flex items-center justify-center text-cyan-400 mb-6 group-hover:scale-110 transition-transform">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
            </div>
            <h3 className="text-xl font-bold mb-2 text-white">100% Private &amp; Offline</h3>
            <p className="text-sm text-[#9494b8] leading-relaxed">
              Zero cloud accounts. Zero tracking pixels. Zero telemetry. All window titles and timestamps reside strictly in an encrypted-capable local SQLite file in your OS user directory.
            </p>
          </div>

          <div className="p-8 rounded-3xl glass hover:border-emerald-500/30 transition-all group">
            <div className="w-12 h-12 rounded-2xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-400 mb-6 group-hover:scale-110 transition-transform">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
            </div>
            <h3 className="text-xl font-bold mb-2 text-white">Decoupled Architecture</h3>
            <p className="text-sm text-[#9494b8] leading-relaxed">
              The daemon auto-starts on login and runs silently in the background. The rich Wails GUI only starts when you open it — closing the window stops the UI while tracking never skips a beat.
            </p>
          </div>
        </div>
      </section>

      {/* Comparison Table */}
      <section id="compare" className="py-20 px-6 max-w-7xl mx-auto border-t border-white/10">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-5xl font-extrabold tracking-tight mb-4">
            How ChronoTrace Compares
          </h2>
          <p className="text-[#9494b8] max-w-xl mx-auto">
            Built with pure Go and zero CGO dependencies for instant cross-platform efficiency.
          </p>
        </div>

        <div className="rounded-3xl glass border border-white/10 overflow-x-auto">
          <table className="w-full text-left border-collapse text-sm">
            <thead>
              <tr className="border-b border-white/10 bg-white/[0.02]">
                <th className="p-5 font-bold text-[#f5f5fc]">Feature</th>
                <th className="p-5 font-bold text-purple-400">ChronoTrace</th>
                <th className="p-5 font-bold text-[#9494b8]">ActivityWatch</th>
                <th className="p-5 font-bold text-[#9494b8]">RescueTime</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5">
              <tr>
                <td className="p-5 font-medium text-white">Background RAM Footprint</td>
                <td className="p-5 font-bold text-emerald-400">&lt; 15 MB (Measured ~2MB)</td>
                <td className="p-5 text-[#9494b8]">100 MB – 250 MB (Python)</td>
                <td className="p-5 text-[#9494b8]">80 MB – 180 MB</td>
              </tr>
              <tr>
                <td className="p-5 font-medium text-white">Privacy &amp; Data Ownership</td>
                <td className="p-5 font-bold text-emerald-400">100% Local SQLite (Zero Telemetry)</td>
                <td className="p-5 text-emerald-400">100% Local</td>
                <td className="p-5 text-rose-400">Cloud Required (Uploads Data)</td>
              </tr>
              <tr>
                <td className="p-5 font-medium text-white">Decoupled Headless Daemon</td>
                <td className="p-5 font-bold text-emerald-400">Yes (Optional UI)</td>
                <td className="p-5 text-[#9494b8]">Partial (Multi-process)</td>
                <td className="p-5 text-[#9494b8]">Tray process</td>
              </tr>
              <tr>
                <td className="p-5 font-medium text-white">Cross-Platform Builds</td>
                <td className="p-5 font-bold text-emerald-400">Pure Go (Zero CGO)</td>
                <td className="p-5 text-[#9494b8]">Rust / Python</td>
                <td className="p-5 text-[#9494b8]">Proprietary</td>
              </tr>
              <tr>
                <td className="p-5 font-medium text-white">Local REST API Access</td>
                <td className="p-5 font-bold text-emerald-400">Yes (127.0.0.1:42069)</td>
                <td className="p-5 text-[#9494b8]">Yes (Port 5600)</td>
                <td className="p-5 text-rose-400">No</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      {/* REST API & Developer Integration */}
      <section id="api" className="py-20 px-6 max-w-7xl mx-auto border-t border-white/10">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          <div>
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-cyan-500/10 text-cyan-400 text-xs font-bold uppercase tracking-wider mb-4 border border-cyan-500/20">
              Developer First
            </div>
            <h2 className="text-3xl md:text-4xl font-extrabold tracking-tight mb-4">
              Local REST API for Custom Dashboards &amp; Automations
            </h2>
            <p className="text-[#9494b8] text-base mb-6 leading-relaxed">
              Every ChronoTrace daemon serves a fast JSON API on localhost. Build custom scripts, query your activity from your terminal, integrate with home automations, or export your logs.
            </p>
            <ul className="space-y-3 text-sm text-[#9494b8]">
              <li className="flex items-center gap-2">
                <span className="text-cyan-400 font-bold">&bull;</span>
                <code>GET /api/v1/status</code> &mdash; Daemon health &amp; uptime
              </li>
              <li className="flex items-center gap-2">
                <span className="text-cyan-400 font-bold">&bull;</span>
                <code>GET /api/v1/usage/today</code> &mdash; Today&apos;s aggregated app time
              </li>
              <li className="flex items-center gap-2">
                <span className="text-cyan-400 font-bold">&bull;</span>
                <code>GET /api/v1/usage/timeline</code> &mdash; Hourly activity heatmap
              </li>
            </ul>
          </div>

          {/* Code Box */}
          <div className="rounded-2xl glass p-6 font-mono text-xs text-purple-200 border border-white/10 shadow-2xl overflow-x-auto">
            <div className="text-[#5e5e7a] mb-2"># Fetch today&apos;s screen time in JSON</div>
            <div className="text-emerald-400 mb-4">$ curl http://127.0.0.1:42069/api/v1/usage/today</div>
            <pre className="text-[#9494b8]">
{`[
  {
    "app_name": "Visual Studio Code",
    "total_seconds": 15480,
    "session_count": 24,
    "formatted_time": "4h 18m"
  },
  {
    "app_name": "Google Chrome",
    "total_seconds": 9900,
    "session_count": 42,
    "formatted_time": "2h 45m"
  }
]`}
            </pre>
          </div>
        </div>
      </section>

      {/* Downloads / Release Section */}
      <section id="downloads" className="py-20 px-6 max-w-7xl mx-auto border-t border-white/10 text-center">
        <h2 className="text-3xl md:text-5xl font-extrabold tracking-tight mb-4">
          Get ChronoTrace Today
        </h2>
        <p className="text-[#9494b8] max-w-xl mx-auto mb-12">
          Free and open-source forever. Download standalone single-binary releases for your platform.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-4xl mx-auto text-left">
          {/* Windows */}
          <div className="p-6 rounded-2xl glass hover:border-purple-500/40 transition-all flex flex-col justify-between">
            <div>
              <div className="text-2xl mb-2">🪟</div>
              <h3 className="text-lg font-bold text-white">Windows (x64)</h3>
              <p className="text-xs text-[#9494b8] mt-1 mb-4">
                Headless background tracker &bull; HKCU Autostart registry integration
              </p>
            </div>
            <a
              href="https://github.com/mahmud-r-farhan/chronotrace/releases"
              className="w-full py-2.5 rounded-xl bg-white/5 hover:bg-white/10 border border-white/10 text-center text-xs font-semibold text-white transition-colors"
            >
              Download .exe
            </a>
          </div>

          {/* macOS */}
          <div className="p-6 rounded-2xl glass hover:border-purple-500/40 transition-all flex flex-col justify-between">
            <div>
              <div className="text-2xl mb-2">🍎</div>
              <h3 className="text-lg font-bold text-white">macOS (Apple Silicon &amp; Intel)</h3>
              <p className="text-xs text-[#9494b8] mt-1 mb-4">
                Native AppKit hooks &bull; launchd LaunchAgent autostart
              </p>
            </div>
            <a
              href="https://github.com/mahmud-r-farhan/chronotrace/releases"
              className="w-full py-2.5 rounded-xl bg-white/5 hover:bg-white/10 border border-white/10 text-center text-xs font-semibold text-white transition-colors"
            >
              Download .zip
            </a>
          </div>

          {/* Linux */}
          <div className="p-6 rounded-2xl glass hover:border-purple-500/40 transition-all flex flex-col justify-between">
            <div>
              <div className="text-2xl mb-2">🐧</div>
              <h3 className="text-lg font-bold text-white">Linux (x64 &amp; arm64)</h3>
              <p className="text-xs text-[#9494b8] mt-1 mb-4">
                X11 &amp; Wayland support &bull; systemd user service integration
              </p>
            </div>
            <a
              href="https://github.com/mahmud-r-farhan/chronotrace/releases"
              className="w-full py-2.5 rounded-xl bg-white/5 hover:bg-white/10 border border-white/10 text-center text-xs font-semibold text-white transition-colors"
            >
              Download Binary
            </a>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-white/10 py-12 px-6 bg-[#07070a]">
        <div className="max-w-7xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-6">
          <div className="flex items-center gap-3">
            <div className="relative w-7 h-7 rounded-lg overflow-hidden border border-purple-500/30">
              <Image src="/logo.png" alt="ChronoTrace Logo" width={28} height={28} className="w-full h-full object-cover" />
            </div>
            <span className="font-bold text-sm bg-gradient-to-r from-white to-purple-300 bg-clip-text text-transparent">
              ChronoTrace
            </span>
            <span className="text-xs text-[#5e5e7a]">
              &copy; 2026 ChronoTrace Contributors. MIT License.
            </span>
          </div>

          <div className="flex items-center gap-6 text-xs text-[#9494b8]">
            <a href="https://github.com/mahmud-r-farhan/chronotrace" className="hover:text-white transition-colors">GitHub Repository</a>
            <a href="https://github.com/mahmud-r-farhan/chronotrace/blob/main/LICENSE" className="hover:text-white transition-colors">MIT License</a>
            <a href="https://github.com/mahmud-r-farhan/chronotrace/releases" className="hover:text-white transition-colors">Releases</a>
          </div>
        </div>
      </footer>
    </div>
  );
}
