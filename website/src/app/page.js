import Image from "next/image";

const releaseBase =
  "https://github.com/mahmud-r-farhan/chronotrace/releases/download/v0.1.0";

const downloads = [
  {
    name: "Windows",
    url: `${releaseBase}/ChronoTrace-Windows-x64-Setup.zip`,
  },
  {
    name: "macOS",
    url: `${releaseBase}/ChronoTrace-macOS-Universal.zip`,
  },
  {
    name: "Linux",
    url: `${releaseBase}/ChronoTrace-Linux-x64.tar.gz`,
  },
];

export default function Home() {
  return (
    <div className="min-h-screen bg-white text-gray-900">
      <nav className="border-b border-gray-200">
        <div className="max-w-4xl mx-auto flex items-center justify-between px-6 py-3">
          <div className="flex items-center gap-2">
            <Image src="/logo.jpg" alt="" width={28} height={28} className="rounded" priority />
            <span className="font-semibold text-sm tracking-tight">ChronoTrace</span>
          </div>
          <div className="flex items-center gap-5 text-sm text-gray-500">
            <a href="#features" className="hover:text-gray-900">Features</a>
            <a href="#compare" className="hover:text-gray-900">Compare</a>
            <a href="#faq" className="hover:text-gray-900">FAQ</a>
            <a
              href="https://github.com/mahmud-r-farhan/chronotrace"
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-gray-900"
            >
              GitHub
            </a>
          </div>
        </div>
      </nav>

      <section className="max-w-3xl mx-auto text-center px-6 pt-24 pb-12">
        <h1 className="text-4xl md:text-5xl font-bold tracking-tight mb-4">
          Screen time tracker.
          <br />
          <span className="text-gray-400 font-normal">Zero telemetry.</span>
        </h1>
        <p className="text-gray-500 text-base md:text-lg max-w-lg mx-auto mb-8 leading-relaxed">
          Logs app usage to a local SQLite database. Under 15MB RAM. No cloud, no
          network, no accounts.
        </p>
        <div className="flex flex-col sm:flex-row items-center justify-center gap-2">
          <a
            href={downloads[0].url}
            className="w-full sm:w-auto px-5 py-2.5 rounded-md bg-gray-900 text-white text-sm font-medium hover:bg-gray-800"
          >
            Download for Windows
          </a>
          <a
            href={downloads[1].url}
            className="w-full sm:w-auto px-5 py-2.5 rounded-md bg-gray-100 text-gray-700 text-sm font-medium hover:bg-gray-200"
          >
            macOS
          </a>
          <a
            href={downloads[2].url}
            className="w-full sm:w-auto px-5 py-2.5 rounded-md bg-gray-100 text-gray-700 text-sm font-medium hover:bg-gray-200"
          >
            Linux
          </a>
        </div>
      </section>

      <section className="max-w-2xl mx-auto px-6 pb-12">
        <div className="px-4 py-2.5 rounded-md bg-gray-50 border border-gray-200 font-mono text-xs text-gray-500 overflow-x-auto">
          $ git clone https://github.com/mahmud-r-farhan/chronotrace.git &amp;&amp; cd
          chronotrace &amp;&amp; make daemon
        </div>
      </section>

      <section className="max-w-4xl mx-auto px-6 pb-16">
        <div className="grid grid-cols-3 gap-6 text-center">
          <div>
            <p className="text-2xl font-bold">&lt;15 MB</p>
            <p className="text-xs text-gray-400 mt-1">RAM</p>
          </div>
          <div>
            <p className="text-2xl font-bold">~0%</p>
            <p className="text-xs text-gray-400 mt-1">CPU</p>
          </div>
          <div>
            <p className="text-2xl font-bold">0%</p>
            <p className="text-xs text-gray-400 mt-1">Telemetry</p>
          </div>
        </div>
      </section>

      <section id="features" className="border-t border-gray-200">
        <div className="max-w-4xl mx-auto px-6 py-16">
          <h2 className="text-2xl font-bold mb-8">Features</h2>
          <div className="grid md:grid-cols-3 gap-6">
            <div>
              <h3 className="font-semibold text-sm mb-1">Lightweight daemon</h3>
              <p className="text-sm text-gray-500">
                Native Go binary. No CGO. Runs as a background process and
                auto-starts on login.
              </p>
            </div>
            <div>
              <h3 className="font-semibold text-sm mb-1">Fully offline</h3>
              <p className="text-sm text-gray-500">
                Zero network calls. All data stored locally in SQLite. No accounts
                or cloud sync.
              </p>
            </div>
            <div>
              <h3 className="font-semibold text-sm mb-1">REST API</h3>
              <p className="text-sm text-gray-500">
                Query your usage data at 127.0.0.1:42069. Works with curl, Python,
                or any HTTP client.
              </p>
            </div>
          </div>
        </div>
      </section>

      <section id="compare" className="border-t border-gray-200 bg-gray-50">
        <div className="max-w-4xl mx-auto px-6 py-16">
          <h2 className="text-2xl font-bold mb-8">Comparison</h2>
          <table className="w-full text-sm">
            <thead className="border-b border-gray-200">
              <tr>
                <th className="text-left py-2 font-semibold"></th>
                <th className="text-left py-2 font-semibold">ChronoTrace</th>
                <th className="text-left py-2 text-gray-400">ActivityWatch</th>
                <th className="text-left py-2 text-gray-400">RescueTime</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              <tr>
                <td className="py-2 text-gray-500">RAM</td>
                <td className="py-2 font-medium">&lt;15 MB</td>
                <td className="py-2 text-gray-500">100-250 MB</td>
                <td className="py-2 text-gray-500">80-180 MB</td>
              </tr>
              <tr>
                <td className="py-2 text-gray-500">Telemetry</td>
                <td className="py-2 font-medium">None</td>
                <td className="py-2 text-gray-500">Local first</td>
                <td className="py-2 text-gray-500">Cloud</td>
              </tr>
              <tr>
                <td className="py-2 text-gray-500">Headless</td>
                <td className="py-2 font-medium">Yes</td>
                <td className="py-2 text-gray-500">Multiple processes</td>
                <td className="py-2 text-gray-500">No</td>
              </tr>
              <tr>
                <td className="py-2 text-gray-500">REST API</td>
                <td className="py-2 font-medium">Yes</td>
                <td className="py-2 text-gray-500">Yes</td>
                <td className="py-2 text-gray-500">No</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section id="faq" className="border-t border-gray-200">
        <div className="max-w-4xl mx-auto px-6 py-16">
          <h2 className="text-2xl font-bold mb-8">FAQ</h2>
          <dl className="space-y-6">
            <div>
              <dt className="font-semibold text-sm mb-1">
                How does it stay under 15MB RAM?
              </dt>
              <dd className="text-sm text-gray-500 leading-relaxed">
                The daemon is a native Go binary with pure Go SQLite. It polls the
                foreground window every 2-3 seconds and flushes to disk every 45
                seconds.
              </dd>
            </div>
            <div>
              <dt className="font-semibold text-sm mb-1">
                Does it send data to any server?
              </dt>
              <dd className="text-sm text-gray-500 leading-relaxed">
                No. All data stays in a local SQLite file. There are zero network
                calls, zero analytics, and zero telemetry.
              </dd>
            </div>
            <div>
              <dt className="font-semibold text-sm mb-1">How do I install it?</dt>
              <dd className="text-sm text-gray-500 leading-relaxed">
                Download the package for your OS, extract it, and run the installer
                script. It registers the daemon to start on boot.
              </dd>
            </div>
            <div>
              <dt className="font-semibold text-sm mb-1">Can I use it from scripts?</dt>
              <dd className="text-sm text-gray-500 leading-relaxed">
                Yes. The REST API at 127.0.0.1:42069 returns JSON with hourly
                timelines, daily summaries, and per-app usage.
              </dd>
            </div>
          </dl>
        </div>
      </section>

      <footer className="border-t border-gray-200">
        <div className="max-w-4xl mx-auto px-6 py-6 flex items-center justify-between text-xs text-gray-400">
          <span>
            ChronoTrace &copy; 2026{" "}
            <a
              href="https://bengalbytes.com/"
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-gray-900 underline underline-offset-2"
            >
              Bengal Bytes
            </a>
            . MIT License.
          </span>
          <a
            href="https://github.com/mahmud-r-farhan/chronotrace"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-gray-900"
          >
            GitHub
          </a>
        </div>
      </footer>
    </div>
  );
}
