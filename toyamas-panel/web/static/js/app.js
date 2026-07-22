// Toyamas Panel - Dashboard, App Store & AI Assistant Frontend Logic

let ws = null;
let cpuChart = null;
let ramChart = null;
let allApps = [];

const maxChartDataPoints = 20;
const chartLabels = Array(maxChartDataPoints).fill('');
const cpuData = Array(maxChartDataPoints).fill(0);
const ramData = Array(maxChartDataPoints).fill(0);

// Initialize Chart.js Instances
function initCharts() {
    const chartConfig = (label, colorHex, dataArray) => ({
        type: 'line',
        data: {
            labels: chartLabels,
            datasets: [{
                label: label,
                data: dataArray,
                borderColor: colorHex,
                backgroundColor: colorHex + '15',
                borderWidth: 2,
                fill: true,
                tension: 0.3,
                pointRadius: 0
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            animation: { duration: 300 },
            scales: {
                x: { display: false },
                y: {
                    beginAtZero: true,
                    grid: { color: '#334155' },
                    ticks: { color: '#94a3b8', font: { size: 10 } }
                }
            },
            plugins: {
                legend: { display: false }
            }
        }
    });

    const ctxCpu = document.getElementById('cpuChart').getContext('2d');
    cpuChart = new Chart(ctxCpu, chartConfig('CPU Usage %', '#10b981', cpuData));

    const ctxRam = document.getElementById('ramChart').getContext('2d');
    ramChart = new Chart(ctxRam, chartConfig('RAM Usage MB', '#14b8a6', ramData));
}

// Connect WebSocket
function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws/metrics`;

    const statusBadge = document.getElementById('wsStatus');
    const statusText = document.getElementById('wsStatusText');

    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
        statusBadge.className = 'flex items-center gap-2 px-3 py-1 rounded-full text-xs font-medium bg-emerald-500/10 text-emerald-400 border border-emerald-500/20';
        statusText.textContent = 'Live Connected';
    };

    ws.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            if (data.metrics) updateMetrics(data.metrics);
            if (data.containers) updateContainers(data.containers);
        } catch (e) {
            console.error('Failed to parse WS payload:', e);
        }
    };

    ws.onclose = () => {
        statusBadge.className = 'flex items-center gap-2 px-3 py-1 rounded-full text-xs font-medium bg-red-500/10 text-red-400 border border-red-500/20';
        statusText.textContent = 'Disconnected (Reconnecting...)';
        setTimeout(connectWebSocket, 3000);
    };

    ws.onerror = (err) => {
        console.error('WS Error:', err);
        ws.close();
    };
}

function updateMetrics(m) {
    document.getElementById('hostName').textContent = m.hostname || '-';
    document.getElementById('hostOS').textContent = `${m.os || 'Linux'}`;
    document.getElementById('hostUptime').textContent = formatUptime(m.uptime_sec);
    if (m.cpu.load_avg) {
        document.getElementById('hostLoad').textContent = m.cpu.load_avg.map(n => n.toFixed(2)).join(', ');
    }

    const cpuPct = m.cpu.usage_percent || 0;
    document.getElementById('cpuPercent').textContent = `${cpuPct.toFixed(1)}%`;
    document.getElementById('cpuBar').style.width = `${Math.min(cpuPct, 100)}%`;
    document.getElementById('cpuCores').textContent = m.cpu.cores;

    const ramPct = m.ram.percent || 0;
    document.getElementById('ramPercent').textContent = `${ramPct.toFixed(1)}%`;
    document.getElementById('ramBar').style.width = `${Math.min(ramPct, 100)}%`;
    document.getElementById('ramUsed').textContent = `${m.ram.used_mb.toFixed(0)} MB`;
    document.getElementById('ramTotal').textContent = `/ ${m.ram.total_mb.toFixed(0)} MB`;

    const diskPct = m.disk.percent || 0;
    document.getElementById('diskPercent').textContent = `${diskPct.toFixed(1)}%`;
    document.getElementById('diskBar').style.width = `${Math.min(diskPct, 100)}%`;
    document.getElementById('diskUsed').textContent = `${m.disk.used_gb.toFixed(1)} GB`;
    document.getElementById('diskTotal').textContent = `/ ${m.disk.total_gb.toFixed(1)} GB`;

    document.getElementById('netRx').textContent = formatSpeed(m.network.rx_bytes_sec);
    document.getElementById('netTx').textContent = formatSpeed(m.network.tx_bytes_sec);
    document.getElementById('netTotalRx').textContent = `Total RX: ${m.network.total_rx_mb.toFixed(1)} MB`;
    document.getElementById('netTotalTx').textContent = `Total TX: ${m.network.total_tx_mb.toFixed(1)} MB`;

    cpuData.shift();
    cpuData.push(cpuPct);
    cpuChart.update();

    ramData.shift();
    ramData.push(m.ram.used_mb);
    ramChart.update();
}

function updateContainers(containers) {
    const tbody = document.getElementById('containerList');
    const containerCount = document.getElementById('containerCount');

    if (!containers || containers.length === 0) {
        tbody.innerHTML = `<tr><td colspan="4" class="p-4 text-center text-slate-500">No Docker containers detected.</td></tr>`;
        containerCount.textContent = '0 Containers';
        return;
    }

    containerCount.textContent = `${containers.length} Container(s)`;

    tbody.innerHTML = containers.map(c => {
        const name = (c.names && c.names.length > 0) ? c.names[0].replace('/', '') : c.id.substring(0, 12);
        const isRunning = c.state === 'running';
        const badgeClass = isRunning ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' : 'bg-red-500/10 text-red-400 border-red-500/20';

        return `
            <tr class="hover:bg-slate-800/40 transition-all">
                <td class="p-3 font-semibold text-slate-200">${name}</td>
                <td class="p-3 text-slate-400 text-xs truncate max-w-[150px]">${c.image}</td>
                <td class="p-3">
                    <span class="px-2 py-0.5 text-xs rounded-full border ${badgeClass}">${c.status}</span>
                </td>
                <td class="p-3 text-right space-x-1">
                    ${isRunning ? 
                        `<button onclick="dockerAction('${c.id}', 'restart')" class="px-2.5 py-1 text-xs bg-slate-800 hover:bg-slate-700 text-slate-300 rounded border border-slate-700">Restart</button>
                         <button onclick="dockerAction('${c.id}', 'stop')" class="px-2.5 py-1 text-xs bg-red-500/20 hover:bg-red-500/30 text-red-300 rounded border border-red-500/30">Stop</button>` :
                        `<button onclick="dockerAction('${c.id}', 'start')" class="px-2.5 py-1 text-xs bg-emerald-500/20 hover:bg-emerald-500/30 text-emerald-300 rounded border border-emerald-500/30">Start</button>`
                    }
                </td>
            </tr>
        `;
    }).join('');
}

async function dockerAction(id, action) {
    try {
        const res = await fetch('/api/docker/action', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id, action })
        });
        const data = await res.json();
        if (!res.ok) alert(`Error: ${data.error}`);
    } catch (e) {
        alert('Failed to send Docker action');
    }
}

async function fetchServices() {
    const list = document.getElementById('serviceList');
    try {
        const res = await fetch('/api/services');
        const data = await res.json();

        if (!data || data.length === 0) {
            list.innerHTML = `<div class="text-xs text-slate-500 text-center py-2">No key services found.</div>`;
            return;
        }

        list.innerHTML = data.map(s => {
            const isActive = s.active === 'active';
            const badgeClass = isActive ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' : 'bg-slate-800 text-slate-400 border-slate-700';

            return `
                <div class="flex items-center justify-between p-2.5 rounded-xl bg-slate-800/40 border border-slate-800 hover:border-slate-700 transition-all">
                    <div>
                        <div class="font-semibold text-sm text-slate-200">${s.name}</div>
                        <span class="inline-block px-2 py-0.5 text-[10px] uppercase font-bold rounded-full border ${badgeClass}">${s.active}</span>
                    </div>
                    <button onclick="restartService('${s.name}')" class="px-2.5 py-1 text-xs bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-lg border border-slate-700">
                        Restart
                    </button>
                </div>
            `;
        }).join('');
    } catch (e) {
        list.innerHTML = `<div class="text-xs text-red-400 text-center py-2">Failed to load services.</div>`;
    }
}

async function restartService(name) {
    if (!confirm(`Are you sure you want to restart service '${name}'?`)) return;
    try {
        const res = await fetch('/api/services/restart', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name })
        });
        const data = await res.json();
        if (!res.ok) alert(`Error: ${data.error}`);
        else fetchServices();
    } catch (e) {
        alert('Failed to send service restart request');
    }
}

// APP STORE ENGINE LOGIC
async function fetchApps() {
    const grid = document.getElementById('appGrid');
    try {
        const res = await fetch('/api/apps');
        allApps = await res.json();
        renderApps(allApps);
    } catch (e) {
        grid.innerHTML = `<div class="col-span-full text-center py-8 text-red-400">Failed to load App Store manifests.</div>`;
    }
}

function renderApps(apps) {
    const grid = document.getElementById('appGrid');
    if (!apps || apps.length === 0) {
        grid.innerHTML = `<div class="col-span-full text-center py-8 text-slate-500">No applications match your search.</div>`;
        return;
    }

    grid.innerHTML = apps.map(a => {
        const isInstalled = a.status === 'running' || a.status === 'installed' || a.status === 'stopped';
        const isRunning = a.status === 'running';

        return `
            <div class="bg-slate-900/80 border border-slate-800 hover:border-slate-700/80 rounded-2xl p-5 shadow-xl flex flex-col justify-between transition-all group">
                <div>
                    <div class="flex items-start justify-between mb-3">
                        <div class="flex items-center gap-3">
                            <div class="text-3xl p-2 rounded-xl bg-slate-800/80 border border-slate-700/60">${a.icon || '📦'}</div>
                            <div>
                                <h4 class="font-bold text-white text-base group-hover:text-emerald-400 transition-colors">${a.name}</h4>
                                <div class="flex items-center gap-2 mt-0.5">
                                    <span class="text-[11px] text-slate-400 px-2 py-0.5 bg-slate-800 rounded-md border border-slate-700">${a.category}</span>
                                    <span class="text-xs text-slate-500">v${a.version}</span>
                                </div>
                            </div>
                        </div>
                        ${isInstalled ? 
                            `<span class="px-2.5 py-0.5 text-[10px] uppercase font-bold rounded-full border ${isRunning ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' : 'bg-amber-500/10 text-amber-400 border-amber-500/20'}">${isRunning ? 'Running' : 'Installed'}</span>` : 
                            `<span class="px-2.5 py-0.5 text-[10px] uppercase font-bold rounded-full border bg-slate-800 text-slate-400 border-slate-700">Available</span>`
                        }
                    </div>
                    <p class="text-xs text-slate-400 line-clamp-2 leading-relaxed mb-4">${a.description}</p>
                </div>

                <div class="pt-3 border-t border-slate-800/80 flex items-center justify-between">
                    <span class="text-xs text-slate-500">Default Port: <strong class="text-slate-300">${a.default_port}</strong></span>
                    <div class="flex items-center gap-1.5">
                        ${isInstalled ? `
                            <button onclick="updateApp('${a.id}')" class="px-3 py-1.5 text-xs font-semibold bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-xl border border-slate-700 transition-all">Update</button>
                            <button onclick="uninstallApp('${a.id}')" class="px-3 py-1.5 text-xs font-semibold bg-red-500/20 hover:bg-red-500/30 text-red-300 rounded-xl border border-red-500/30 transition-all">Uninstall</button>
                        ` : `
                            <button onclick="openInstallModal('${a.id}')" class="px-4 py-1.5 text-xs font-bold bg-emerald-500 hover:bg-emerald-400 text-slate-950 rounded-xl shadow-lg shadow-emerald-500/20 transition-all">Install</button>
                        `}
                    </div>
                </div>
            </div>
        `;
    }).join('');
}

function openInstallModal(appID) {
    const app = allApps.find(a => a.id === appID);
    if (!app) return;

    document.getElementById('modalAppID').value = app.id;
    document.getElementById('modalAppName').innerHTML = `<span>${app.icon || '⚙️'}</span> Install ${app.name}`;

    const fieldsContainer = document.getElementById('modalEnvFields');
    fieldsContainer.innerHTML = (app.env_vars || []).map(env => `
        <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">${env.label}</label>
            <input type="${env.type || 'text'}" name="${env.name}" value="${env.default}" required
                class="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-xl text-sm text-white focus:outline-none focus:ring-2 focus:ring-emerald-500">
            ${env.description ? `<p class="text-[10px] text-slate-500 mt-1">${env.description}</p>` : ''}
        </div>
    `).join('');

    document.getElementById('appModal').classList.remove('hidden');
}

function closeInstallModal() {
    document.getElementById('appModal').classList.add('hidden');
}

async function handleInstallSubmit(e) {
    e.preventDefault();
    const appID = document.getElementById('modalAppID').value;
    const form = document.getElementById('modalInstallForm');
    const formData = new FormData(form);

    const env = {};
    formData.forEach((val, key) => {
        if (key !== 'app_id') env[key] = val;
    });

    closeInstallModal();

    try {
        const res = await fetch('/api/apps/install', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ app_id: appID, env })
        });
        const data = await res.json();
        if (!res.ok) alert(`Installation failed: ${data.error}`);
        else fetchApps();
    } catch (err) {
        alert('Installation request error');
    }
}

async function uninstallApp(appID) {
    if (!confirm(`Are you sure you want to uninstall application '${appID}'? This will stop and remove its container.`)) return;

    try {
        const res = await fetch('/api/apps/uninstall', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ app_id: appID })
        });
        const data = await res.json();
        if (!res.ok) alert(`Uninstall failed: ${data.error}`);
        else fetchApps();
    } catch (err) {
        alert('Uninstall request error');
    }
}

async function updateApp(appID) {
    try {
        const res = await fetch('/api/apps/update', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ app_id: appID })
        });
        const data = await res.json();
        if (!res.ok) alert(`Update failed: ${data.error}`);
        else {
            alert(`Application '${appID}' updated successfully!`);
            fetchApps();
        }
    } catch (err) {
        alert('Update request error');
    }
}

// AI ASSISTANT FRONTEND ENGINE
function sendQuickPrompt(text) {
    document.getElementById('aiInputPrompt').value = text;
    document.getElementById('aiPromptForm').dispatchEvent(new Event('submit'));
}

async function handleAISubmit(e) {
    e.preventDefault();
    const input = document.getElementById('aiInputPrompt');
    const prompt = input.value.trim();
    if (!prompt) return;

    const chatBox = document.getElementById('aiChatBox');

    // Append User Message
    chatBox.innerHTML += `
        <div class="flex items-start justify-end gap-3">
            <div class="bg-emerald-500/10 border border-emerald-500/20 p-4 rounded-2xl text-sm text-emerald-300 leading-relaxed max-w-xl">
                ${escapeHTML(prompt)}
            </div>
            <div class="w-8 h-8 rounded-xl bg-slate-800 border border-slate-700 text-slate-300 font-bold flex items-center justify-center text-xs">YOU</div>
        </div>
    `;
    input.value = '';
    chatBox.scrollTop = chatBox.scrollHeight;

    // Append Loading Indicator
    const loadingId = 'aiLoading-' + Date.now();
    chatBox.innerHTML += `
        <div id="${loadingId}" class="flex items-start gap-3">
            <div class="w-8 h-8 rounded-xl bg-emerald-500 text-slate-950 font-bold flex items-center justify-center text-sm">AI</div>
            <div class="bg-slate-800/80 border border-slate-700/60 p-4 rounded-2xl text-sm text-slate-400 animate-pulse">
                Analyzing request & system context...
            </div>
        </div>
    `;
    chatBox.scrollTop = chatBox.scrollHeight;

    try {
        const res = await fetch('/api/ai/chat', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ prompt })
        });
        const data = await res.json();
        document.getElementById(loadingId)?.remove();

        let aiHTML = `<div class="bg-slate-800/80 border border-slate-700/60 p-4 rounded-2xl text-sm leading-relaxed max-w-2xl text-slate-200">
            ${formatMarkdown(data.message)}
        `;

        if (data.requires_user_confirmation && data.action_plan) {
            const plan = data.action_plan;
            aiHTML += `
                <div class="mt-4 p-4 rounded-xl bg-amber-500/10 border border-amber-500/30 text-amber-200 space-y-3">
                    <div class="font-bold text-sm flex items-center gap-2">
                        <span>⚠️</span> Proposed Action Plan: ${plan.target}
                    </div>
                    <ul class="space-y-1.5 text-xs text-amber-300/90 pl-2">
                        ${plan.steps.map(s => `<li>• <strong>${s.title}</strong>: ${s.description}</li>`).join('')}
                    </ul>
                    <div class="pt-2 flex justify-end">
                        <button onclick="confirmAIAction('${plan.token}')" class="px-4 py-2 bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs rounded-xl shadow-lg transition-all flex items-center gap-1.5">
                            <span>✔</span> Confirm & Execute Action
                        </button>
                    </div>
                </div>
            `;
        }

        aiHTML += `</div>`;

        chatBox.innerHTML += `
            <div class="flex items-start gap-3">
                <div class="w-8 h-8 rounded-xl bg-emerald-500 text-slate-950 font-bold flex items-center justify-center text-sm">AI</div>
                ${aiHTML}
            </div>
        `;
        chatBox.scrollTop = chatBox.scrollHeight;

    } catch (e) {
        document.getElementById(loadingId)?.remove();
        chatBox.innerHTML += `
            <div class="flex items-start gap-3">
                <div class="w-8 h-8 rounded-xl bg-emerald-500 text-slate-950 font-bold flex items-center justify-center text-sm">AI</div>
                <div class="bg-red-500/10 border border-red-500/20 p-4 rounded-2xl text-sm text-red-400">
                    Failed to communicate with AI Assistant. Ensure server backend is running.
                </div>
            </div>
        `;
    }
}

async function confirmAIAction(token) {
    const chatBox = document.getElementById('aiChatBox');
    try {
        const res = await fetch('/api/ai/confirm', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ token })
        });
        const data = await res.json();

        chatBox.innerHTML += `
            <div class="flex items-start gap-3">
                <div class="w-8 h-8 rounded-xl bg-emerald-500 text-slate-950 font-bold flex items-center justify-center text-sm">AI</div>
                <div class="bg-emerald-500/10 border border-emerald-500/30 p-4 rounded-2xl text-sm text-emerald-300 leading-relaxed max-w-2xl">
                    ${formatMarkdown(data.message)}
                </div>
            </div>
        `;
        chatBox.scrollTop = chatBox.scrollHeight;
    } catch (e) {
        alert('Confirmation execution failed');
    }
}

// Helpers
function escapeHTML(str) {
    return str.replace(/[&<>'"]/g, 
        tag => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', "'": '&#39;', '"': '&quot;' }[tag] || tag)
    );
}

function formatMarkdown(text) {
    if (!text) return '';
    return text.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
               .replace(/\n/g, '<br>');
}

function formatUptime(seconds) {
    if (!seconds) return '-';
    const hrs = Math.floor(seconds / 3600);
    const mins = Math.floor((seconds % 3600) / 60);
    return `${hrs}h ${mins}m`;
}

function formatSpeed(bytesPerSec) {
    if (!bytesPerSec || bytesPerSec === 0) return '0 B/s';
    const k = 1024;
    const sizes = ['B/s', 'KB/s', 'MB/s', 'GB/s'];
    const i = Math.floor(Math.log(bytesPerSec) / Math.log(k));
    return (bytesPerSec / Math.pow(k, i)).toFixed(1) + ' ' + sizes[i];
}

// Initialize on DOM Ready
document.addEventListener('DOMContentLoaded', () => {
    initCharts();
    connectWebSocket();
    fetchServices();
    fetchApps();

    // Tab Switching
    const overviewBtn = document.getElementById('tabOverviewBtn');
    const appStoreBtn = document.getElementById('tabAppStoreBtn');
    const aiBtn = document.getElementById('tabAIBtn');
    const overviewSec = document.getElementById('sectionOverview');
    const appStoreSec = document.getElementById('sectionAppStore');
    const aiSec = document.getElementById('sectionAIAssistant');

    overviewBtn.addEventListener('click', () => {
        overviewSec.classList.remove('hidden');
        appStoreSec.classList.add('hidden');
        aiSec.classList.add('hidden');
        overviewBtn.className = 'px-4 py-1.5 rounded-lg text-xs font-semibold bg-emerald-500 text-slate-950 transition-all shadow';
        appStoreBtn.className = 'px-4 py-1.5 rounded-lg text-xs font-semibold text-slate-300 hover:text-white transition-all';
        aiBtn.className = 'px-4 py-1.5 rounded-lg text-xs font-semibold text-slate-300 hover:text-white transition-all';
    });

    appStoreBtn.addEventListener('click', () => {
        appStoreSec.classList.remove('hidden');
        overviewSec.classList.add('hidden');
        aiSec.classList.add('hidden');
        appStoreBtn.className = 'px-4 py-1.5 rounded-lg text-xs font-semibold bg-emerald-500 text-slate-950 transition-all shadow';
        overviewBtn.className = 'px-4 py-1.5 rounded-lg text-xs font-semibold text-slate-300 hover:text-white transition-all';
        aiBtn.className = 'px-4 py-1.5 rounded-lg text-xs font-semibold text-slate-300 hover:text-white transition-all';
    });

    aiBtn.addEventListener('click', () => {
        aiSec.classList.remove('hidden');
        overviewSec.classList.add('hidden');
        appStoreSec.classList.add('hidden');
        aiBtn.className = 'px-4 py-1.5 rounded-lg text-xs font-semibold bg-emerald-500 text-slate-950 transition-all shadow';
        overviewBtn.className = 'px-4 py-1.5 rounded-lg text-xs font-semibold text-slate-300 hover:text-white transition-all';
        appStoreBtn.className = 'px-4 py-1.5 rounded-lg text-xs font-semibold text-slate-300 hover:text-white transition-all';
    });

    // App Search Input
    document.getElementById('appSearchInput').addEventListener('input', (e) => {
        const query = e.target.value.toLowerCase();
        const filtered = allApps.filter(a => a.name.toLowerCase().includes(query) || a.category.toLowerCase().includes(query) || a.description.toLowerCase().includes(query));
        renderApps(filtered);
    });

    // AI Form Events
    document.getElementById('aiPromptForm').addEventListener('submit', handleAISubmit);

    // Modal Events
    document.getElementById('modalCloseBtn').addEventListener('click', closeInstallModal);
    document.getElementById('modalCancelBtn').addEventListener('click', closeInstallModal);
    document.getElementById('modalInstallForm').addEventListener('submit', handleInstallSubmit);

    document.getElementById('refreshSvcsBtn').addEventListener('click', fetchServices);

    document.getElementById('logoutBtn').addEventListener('click', async () => {
        await fetch('/api/logout', { method: 'POST' });
        window.location.href = '/login';
    });
});
