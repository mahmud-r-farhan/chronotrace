/**
 * ChronoTrace — Frontend Application
 * Pure JS, no heavy framework. All Wails bindings via window.go.*
 */

// ────────────────────────────────────────────────────────────────
// Wails bridge — wraps window.go calls with error handling.
// In development (served via browser), we use mock data.
// ────────────────────────────────────────────────────────────────
const isWails = () => typeof window.go !== 'undefined';

const go = {
  async GetStatus()         { return isWails() ? window.go.main.App.GetStatus()         : mockStatus(); },
  async GetUsageToday()     { return isWails() ? window.go.main.App.GetUsageToday()     : mockUsage(); },
  async GetUsageWeek()      { return isWails() ? window.go.main.App.GetUsageWeek()      : mockUsage(); },
  async GetUsageMonth()     { return isWails() ? window.go.main.App.GetUsageMonth()     : mockUsage(); },
  async GetTimeline(date)   { return isWails() ? window.go.main.App.GetTimeline(date)   : mockTimeline(); },
  async GetSummary(date)    { return isWails() ? window.go.main.App.GetSummary(date)    : mockSummary(); },
  async EnableAutostart()   { return isWails() ? window.go.main.App.EnableAutostart()   : null; },
  async DisableAutostart()  { return isWails() ? window.go.main.App.DisableAutostart()  : null; },
  async IsAutostartEnabled(){ return isWails() ? window.go.main.App.IsAutostartEnabled(): false; },
};

// ────────────────────────────────────────────────────────────────
// State
// ────────────────────────────────────────────────────────────────
const state = {
  view: 'dashboard',
  period: 'today',       // today | week | month
  usageCache: {},        // period → []AppUsage
  timelineCache: {},     // date  → []TimelineSlot
  refreshInterval: null,
};

// ────────────────────────────────────────────────────────────────
// Router — switch between views
// ────────────────────────────────────────────────────────────────
function navigateTo(view) {
  document.querySelectorAll('.view').forEach(el => el.classList.remove('active'));
  document.querySelectorAll('.nav-item').forEach(el => el.classList.remove('active'));

  const viewEl = document.getElementById(`view-${view}`);
  const navEl  = document.getElementById(`nav-${view}`);
  if (viewEl) viewEl.classList.add('active');
  if (navEl)  navEl.classList.add('active');

  state.view = view;

  if (view === 'dashboard') loadDashboard();
  if (view === 'apps')      loadApps();
  if (view === 'timeline')  loadTimeline();
  if (view === 'settings')  loadSettings();
}

// ────────────────────────────────────────────────────────────────
// Status poller
// ────────────────────────────────────────────────────────────────
async function pollStatus() {
  try {
    const s = await go.GetStatus();
    const badge = document.getElementById('status-badge');
    if (s.connected) {
      badge.className = 'status-badge online';
      badge.querySelector('.status-text').textContent = 'Daemon connected';
    } else {
      badge.className = 'status-badge offline';
      badge.querySelector('.status-text').textContent = 'Daemon offline';
    }
  } catch {
    const badge = document.getElementById('status-badge');
    badge.className = 'status-badge offline';
    badge.querySelector('.status-text').textContent = 'Daemon offline';
  }
}

// ────────────────────────────────────────────────────────────────
// Dashboard
// ────────────────────────────────────────────────────────────────
async function loadDashboard() {
  updateTodayDate();
  await Promise.all([
    loadUsageChart(),
    loadSummaryStats(),
    loadTimelineWidget(),
  ]);
}

function updateTodayDate() {
  const el = document.getElementById('today-date');
  if (!el) return;
  const now = new Date();
  el.textContent = now.toLocaleDateString('en-US', {
    weekday: 'long', year: 'numeric', month: 'long', day: 'numeric'
  });
}

async function loadSummaryStats() {
  try {
    const [summary, status] = await Promise.all([
      go.GetSummary(''),
      go.GetStatus(),
    ]);

    setEl('stat-total', summary.total_seconds != null ? formatDuration(summary.total_seconds) : '--');
    setEl('stat-apps',  summary.app_count != null ? summary.app_count : '--');
    if (summary.top_apps && summary.top_apps.length > 0) {
      setEl('stat-top', summary.top_apps[0].app_name);
    }
    if (status.connected) {
      setEl('stat-active', 'Tracking');
    } else {
      setEl('stat-active', 'Offline');
    }
  } catch (e) {
    console.warn('loadSummaryStats:', e);
  }
}

async function loadUsageChart() {
  const container = document.getElementById('usage-chart');
  if (!container) return;

  container.innerHTML = '<div class="empty-state"><div class="spinner"></div><p>Loading...</p></div>';

  try {
    let data;
    if (state.usageCache[state.period]) {
      data = state.usageCache[state.period];
    } else {
      data = await fetchUsage(state.period);
      state.usageCache[state.period] = data;
    }

    if (!data || data.length === 0) {
      container.innerHTML = '<div class="empty-state"><p>No data yet. Keep working! 🚀</p></div>';
      document.getElementById('app-count-badge').textContent = '0 apps';
      return;
    }

    document.getElementById('app-count-badge').textContent = `${data.length} app${data.length !== 1 ? 's' : ''}`;

    const max = data[0].total_seconds || 1;
    container.innerHTML = '';

    data.forEach((app, i) => {
      const pct = ((app.total_seconds / max) * 100).toFixed(1);
      const row = document.createElement('div');
      row.className = 'usage-bar-row';
      row.innerHTML = `
        <span class="usage-bar-label" title="${app.app_name}">${app.app_name}</span>
        <div class="usage-bar-track">
          <div class="usage-bar-fill bar-color-${i % 8}" style="width:0%" data-pct="${pct}"></div>
        </div>
        <span class="usage-bar-time">${app.formatted_time || formatDuration(app.total_seconds)}</span>
      `;
      container.appendChild(row);
    });

    // Animate bars after DOM insert
    requestAnimationFrame(() => {
      container.querySelectorAll('.usage-bar-fill').forEach(el => {
        el.style.width = el.dataset.pct + '%';
      });
    });

  } catch (e) {
    console.error('loadUsageChart:', e);
    container.innerHTML = '<div class="empty-state"><p>Failed to load data</p></div>';
  }
}

async function loadTimelineWidget() {
  const container = document.getElementById('timeline-chart');
  if (!container) return;
  container.innerHTML = '';

  try {
    const slots = await go.GetTimeline('');
    renderTimelineHeatmap(container, slots);
  } catch (e) {
    container.innerHTML = '<div class="empty-state"><p>No timeline data</p></div>';
  }
}

function renderTimelineHeatmap(container, slots) {
  const max = Math.max(...slots.map(s => s.total_seconds), 1);

  // Split into two rows: AM (0-11) and PM (12-23)
  ['AM (0–11)', 'PM (12–23)'].forEach((label, half) => {
    const rowLabel = document.createElement('div');
    rowLabel.style.cssText = 'font-size:10px;color:var(--text-muted);margin-bottom:4px;';
    rowLabel.textContent = label;
    container.appendChild(rowLabel);

    const grid = document.createElement('div');
    grid.className = 'timeline-grid';

    const labelRow = document.createElement('div');
    labelRow.className = 'timeline-labels';

    for (let i = 0; i < 12; i++) {
      const hour = half * 12 + i;
      const slot = slots[hour] || { hour, total_seconds: 0 };
      const intensity = slot.total_seconds / max;

      const cell = document.createElement('div');
      cell.className = 'timeline-cell';
      cell.title = `${hour}:00 — ${formatDuration(slot.total_seconds)}`;
      cell.style.background = intensity > 0
        ? `rgba(139,92,246,${Math.max(0.1, intensity * 0.9)})`
        : 'var(--bg-glass)';
      if (intensity > 0) cell.dataset.active = 'true';
      grid.appendChild(cell);

      const lbl = document.createElement('div');
      lbl.className = 'timeline-label';
      lbl.textContent = `${hour}`;
      labelRow.appendChild(lbl);
    }

    container.appendChild(grid);
    container.appendChild(labelRow);
  });

  // Legend
  const legend = document.createElement('div');
  legend.className = 'timeline-legend';
  legend.innerHTML = `
    <span>Less</span>
    ${[0.1, 0.3, 0.55, 0.75, 1].map(v =>
      `<div class="legend-cell" style="background:rgba(139,92,246,${v})"></div>`
    ).join('')}
    <span>More</span>
  `;
  container.appendChild(legend);
}

// ────────────────────────────────────────────────────────────────
// Apps View
// ────────────────────────────────────────────────────────────────
async function loadApps() {
  const container = document.getElementById('apps-list');
  if (!container) return;
  container.innerHTML = '<div class="empty-state"><div class="spinner"></div><p>Loading apps...</p></div>';

  try {
    const data = await fetchUsage(state.period);
    if (!data || data.length === 0) {
      container.innerHTML = '<div class="empty-state"><p>No data recorded yet</p></div>';
      return;
    }

    const max = data[0].total_seconds || 1;
    container.innerHTML = '';

    data.forEach((app, i) => {
      const pct = ((app.total_seconds / max) * 100).toFixed(1);
      const row = document.createElement('div');
      row.className = 'app-row';
      row.innerHTML = `
        <span class="app-rank">#${i + 1}</span>
        <div class="app-info">
          <div class="app-name">${app.app_name}</div>
          <div class="app-sessions">${app.session_count || 0} session${app.session_count !== 1 ? 's' : ''}</div>
        </div>
        <div class="app-bar-col">
          <div class="app-bar-track">
            <div class="app-bar-fill bar-color-${i % 8}" style="width:0%" data-pct="${pct}"></div>
          </div>
        </div>
        <span class="app-time">${app.formatted_time || formatDuration(app.total_seconds)}</span>
      `;
      container.appendChild(row);
    });

    requestAnimationFrame(() => {
      container.querySelectorAll('.app-bar-fill').forEach(el => {
        el.style.width = el.dataset.pct + '%';
      });
    });
  } catch (e) {
    console.error('loadApps:', e);
    container.innerHTML = '<div class="empty-state"><p>Failed to load data</p></div>';
  }
}

// ────────────────────────────────────────────────────────────────
// Timeline View
// ────────────────────────────────────────────────────────────────
async function loadTimeline() {
  const picker = document.getElementById('timeline-date-picker');
  const container = document.getElementById('timeline-detail');
  if (!container) return;

  const date = picker?.value || new Date().toISOString().slice(0, 10);

  container.innerHTML = '<div class="empty-state"><div class="spinner"></div><p>Loading...</p></div>';

  try {
    const slots = await go.GetTimeline(date);
    renderTimelineDetail(container, slots);
  } catch (e) {
    container.innerHTML = '<div class="empty-state"><p>Failed to load timeline</p></div>';
  }
}

function renderTimelineDetail(container, slots) {
  const max = Math.max(...slots.map(s => s.total_seconds), 1);
  container.innerHTML = '';

  slots.forEach(slot => {
    const pct = ((slot.total_seconds / max) * 100).toFixed(1);
    const row = document.createElement('div');
    row.className = 'timeline-hour-row';
    row.innerHTML = `
      <span class="timeline-hour-label">${String(slot.hour).padStart(2, '0')}:00</span>
      <div class="timeline-hour-bar">
        <div class="timeline-hour-fill" style="width:0%" data-pct="${pct}"></div>
      </div>
      <span class="timeline-hour-time">${slot.total_seconds > 0 ? formatDuration(slot.total_seconds) : ''}</span>
    `;
    container.appendChild(row);
  });

  requestAnimationFrame(() => {
    container.querySelectorAll('.timeline-hour-fill').forEach(el => {
      el.style.width = el.dataset.pct + '%';
    });
  });
}

// ────────────────────────────────────────────────────────────────
// Settings View
// ────────────────────────────────────────────────────────────────
async function loadSettings() {
  try {
    const [status, autostartEnabled] = await Promise.all([
      go.GetStatus(),
      go.IsAutostartEnabled(),
    ]);

    if (status.version) setEl('setting-version', `v${status.version}`);
    if (status.addr)    setEl('setting-addr', status.addr);
    if (status.uptime_seconds != null) {
      setEl('setting-uptime', formatDuration(status.uptime_seconds));
    }

    const toggle = document.getElementById('autostart-toggle');
    if (toggle) toggle.checked = !!autostartEnabled;

  } catch (e) {
    console.warn('loadSettings:', e);
  }
}

// ────────────────────────────────────────────────────────────────
// Helpers
// ────────────────────────────────────────────────────────────────
async function fetchUsage(period) {
  switch (period) {
    case 'week':  return go.GetUsageWeek();
    case 'month': return go.GetUsageMonth();
    default:      return go.GetUsageToday();
  }
}

function setEl(id, text) {
  const el = document.getElementById(id);
  if (el) el.textContent = text;
}

function formatDuration(seconds) {
  if (!seconds || seconds <= 0) return '0s';
  if (seconds < 60) return `${seconds}s`;
  const m = Math.floor(seconds / 60);
  if (m < 60) return `${m}m`;
  const h = Math.floor(m / 60);
  const rem = m % 60;
  return rem === 0 ? `${h}h` : `${h}h ${rem}m`;
}

// ────────────────────────────────────────────────────────────────
// Mock data for browser development / when daemon is offline
// ────────────────────────────────────────────────────────────────
function mockStatus() {
  return { connected: true, version: '0.1.0', uptime_seconds: 3600, addr: '127.0.0.1:42069' };
}

function mockUsage() {
  const apps = [
    { app_name: 'VS Code',        total_seconds: 7200, session_count: 14, formatted_time: '2h' },
    { app_name: 'Google Chrome',  total_seconds: 5400, session_count: 32, formatted_time: '1h 30m' },
    { app_name: 'Slack',          total_seconds: 2700, session_count: 8,  formatted_time: '45m' },
    { app_name: 'Terminal',       total_seconds: 1800, session_count: 21, formatted_time: '30m' },
    { app_name: 'Figma',          total_seconds: 1200, session_count: 4,  formatted_time: '20m' },
    { app_name: 'Spotify',        total_seconds:  900, session_count: 2,  formatted_time: '15m' },
    { app_name: 'Notion',         total_seconds:  600, session_count: 6,  formatted_time: '10m' },
    { app_name: 'Discord',        total_seconds:  360, session_count: 3,  formatted_time: '6m' },
  ];
  return apps;
}

function mockTimeline() {
  return Array.from({ length: 24 }, (_, i) => ({
    hour: i,
    total_seconds: i >= 9 && i <= 20 ? Math.floor(Math.random() * 3000) : 0,
  }));
}

function mockSummary() {
  return {
    date: new Date().toISOString().slice(0, 10),
    total_seconds: 19260,
    app_count: 8,
    top_apps: mockUsage().slice(0, 5),
  };
}

// ────────────────────────────────────────────────────────────────
// Event Listeners
// ────────────────────────────────────────────────────────────────
function initNavigation() {
  document.querySelectorAll('.nav-item').forEach(item => {
    item.addEventListener('click', e => {
      e.preventDefault();
      navigateTo(item.dataset.view);
    });
  });
}

function initPeriodTabs() {
  document.querySelectorAll('.period-tab').forEach(tab => {
    tab.addEventListener('click', () => {
      // Update active tab in same parent group
      tab.closest('.period-tabs').querySelectorAll('.period-tab').forEach(t => t.classList.remove('active'));
      tab.classList.add('active');

      state.period = tab.dataset.period;
      state.usageCache = {}; // Bust cache on period change

      if (state.view === 'dashboard') loadUsageChart();
      if (state.view === 'apps')      loadApps();
    });
  });
}

function initAutostart() {
  const toggle = document.getElementById('autostart-toggle');
  if (!toggle) return;
  toggle.addEventListener('change', async () => {
    try {
      if (toggle.checked) {
        await go.EnableAutostart();
      } else {
        await go.DisableAutostart();
      }
    } catch (e) {
      console.error('autostart toggle:', e);
      toggle.checked = !toggle.checked; // revert on error
    }
  });
}

function initDatePicker() {
  const picker = document.getElementById('timeline-date-picker');
  if (!picker) return;
  picker.value = new Date().toISOString().slice(0, 10);
  picker.addEventListener('change', () => {
    if (state.view === 'timeline') loadTimeline();
  });
}

// ────────────────────────────────────────────────────────────────
// Boot
// ────────────────────────────────────────────────────────────────
async function init() {
  initNavigation();
  initPeriodTabs();
  initAutostart();
  initDatePicker();

  // Initial load
  await pollStatus();
  navigateTo('dashboard');

  // Auto-refresh status every 15s, data every 60s
  setInterval(pollStatus, 15_000);
  setInterval(() => {
    state.usageCache = {};
    if (state.view === 'dashboard') loadDashboard();
    if (state.view === 'apps')      loadApps();
  }, 60_000);
}

// Wait for Wails runtime to be ready
if (typeof window.runtime !== 'undefined') {
  window.runtime.EventsOn('wails:ready', init);
} else {
  // Browser dev mode
  document.addEventListener('DOMContentLoaded', init);
}
