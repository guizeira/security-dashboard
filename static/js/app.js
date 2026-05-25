let activeEventSource = null;

function normalizeTarget(value) {
    if (!value) return "";
    return value
        .trim()
        .replace(/^https?:\/\//i, "")
        .replace(/^www\./i, "")
        .replace(/\/.*$/, "")
        .replace(/:\d+$/, "")
        .toLowerCase();
}

function isValidTarget(target) {
    const domainRegex = /^(?!-)[a-z0-9-]{1,63}(?<!-)(\.[a-z]{2,})+$/i;
    const ipRegex = /^(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}$/;
    return domainRegex.test(target) || ipRegex.test(target);
}

function createSkeletonCard(title, idPrefix, iconName) {
    return `
        <div class="card-enterprise" id="${idPrefix}">
            <div class="card-header">
                <div class="card-title-group">
                    <div class="card-icon">
                        <i data-lucide="${iconName}"></i>
                    </div>
                    <div>
                        <h3>${title}</h3>
                        <span class="card-subtitle">Live Intelligence</span>
                    </div>
                </div>
            </div>
            <div id="${idPrefix}-body" class="card-body">
                <div class="skeleton-line"></div>
                <div class="skeleton-line short"></div>
            </div>
        </div>
    `;
}

// ======================================================
// SSE PROGRESS & SIMULATED PROGRESS
// ======================================================

function initRealProgress(target, elementId) {
    const el = document.getElementById(elementId);
    el.innerHTML = `
        <div class="progress-wrapper">
            <div class="progress-container">
                <div id="progress-bar-${elementId}" class="progress-bar"></div>
            </div>
            <div class="progress-meta">
                <span class="progress-text" id="progress-text-${elementId}">Initializing scan...</span>
                <span class="progress-percent" id="progress-percent-${elementId}">0%</span>
            </div>
        </div>
    `;

    const bar = document.getElementById(`progress-bar-${elementId}`);
    const text = document.getElementById(`progress-text-${elementId}`);
    const percentText = document.getElementById(`progress-percent-${elementId}`);

    let currentPercent = 0;
    let fakeProgressInterval = null;

    const updateUI = (val) => {
        currentPercent = val;
        if (bar) bar.style.width = `${currentPercent}%`;
        if (percentText) percentText.innerText = `${currentPercent.toFixed(1)}%`;
    };

    fakeProgressInterval = setInterval(() => {
        if (currentPercent < 92) {
            const increment = Math.random() * 1.5;
            const nextVal = Math.min(currentPercent + increment, 92);
            updateUI(nextVal);
            if (text) text.innerText = "Deep scanning infrastructure...";
        } else {
            clearInterval(fakeProgressInterval);
        }
    }, 600);

    const eventSource = new EventSource(`/api/scan-stream?target=${encodeURIComponent(target)}`);

    eventSource.onerror = () => {
        clearInterval(fakeProgressInterval);
        eventSource.close();
    };

    eventSource.addEventListener("progress", (event) => {
        const realPercent = parseFloat(event.data);
        if (isNaN(realPercent)) return;
        if (realPercent > currentPercent) {
            updateUI(Math.min(realPercent, 99));
        }
    });

    eventSource.addEventListener("done", () => {
        clearInterval(fakeProgressInterval);
        updateUI(100);
        if (text) text.innerText = "Scan completed";
        eventSource.close();
    });

    eventSource.cleanup = () => {
        clearInterval(fakeProgressInterval);
        eventSource.close();
    };

    return eventSource;
}

// ======================================================
// LOOKUP
// ======================================================

async function lookup() {
    const input = document.getElementById("target");
    let target = normalizeTarget(input.value.trim());
    const resultsDiv = document.getElementById("results");

    if (!target) return;

    if (activeEventSource) {
        if (activeEventSource.cleanup) activeEventSource.cleanup();
        else activeEventSource.close();
    }

    input.value = target;

    if (!isValidTarget(target)) {
        resultsDiv.innerHTML = `
            <div class="card-enterprise">
                <div class="card-body">
                    <span class="data-value status-danger">Invalid domain or IP</span>
                </div>
            </div>
        `;
        return;
    }

    resultsDiv.innerHTML = `
        ${createSkeletonCard("WHOIS Registry", "whois-card", "file-text")}
        ${createSkeletonCard("Location", "loc-card", "map-pin")}
        ${createSkeletonCard("Exposed Ports", "ports-card", "cpu")}
        ${createSkeletonCard("DNS Records", "dns-card", "globe")}
        ${createSkeletonCard("MX Records", "mx-card", "mail")}
        ${createSkeletonCard("SSL Certificate", "ssl-card", "shield-check")}
    `;

    lucide.createIcons();

    activeEventSource = initRealProgress(target, "ports-card-body");
    const portsContainer = document.getElementById("ports-card-body");
    let portsRendered = false;

    activeEventSource.addEventListener("port", (event) => {
        const port = JSON.parse(event.data);
        if (!portsRendered) portsRendered = true;

        const portId = `port-${port.port}`;
        if (!document.getElementById(portId)) {
            const portHtml = `
                <div class="port-item" id="${portId}">
                    <div class="port-info">
                        <span class="port-number">PORT ${port.port}</span>
                        <span class="port-desc">${port.service || "Unknown service"}</span>
                        ${port.version ? `<span class="port-version">${port.version}</span>` : ""}
                    </div>
                    <span class="badge-open">${port.status.toUpperCase()}</span>
                </div>
            `;
            portsContainer.insertAdjacentHTML('beforeend', portHtml);
        }
    });

    activeEventSource.addEventListener("result", (event) => {
        const data = JSON.parse(event.data);
        if (!data.length && !portsRendered) {
            portsContainer.insertAdjacentHTML('beforeend', `<span class="data-value">No exposed ports detected</span>`);
        }
    });

    // ======================================================
    // WHOIS - AJUSTADO PARA SEPARAR ORG E AS
    // ======================================================
    fetch("/api/whois", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ target })
    })
    .then(res => res.json())
    .then(data => {
        const whoisBody = document.getElementById("whois-card-body");
        if (whoisBody) {
            // Extrai apenas o AS (ex: AS28604) do campo 'as' que vem da API
            let asNumber = "N/A";
            if (data.as) {
                const match = data.as.match(/AS\d+/i);
                asNumber = match ? match[0] : data.as;
            }

            whoisBody.innerHTML = `
                <div class="data-row">
                    <span class="data-label">Organization</span>
                    <span class="data-value">${data.owner || "Unknown"}</span>
                </div>
                <div class="data-row">
                    <span class="data-label">AS Number</span>
                    <span class="data-value">${asNumber}</span>
                </div>
                <div class="data-row">
                    <span class="data-label">Country</span>
                    <span class="data-value">${data.country || "Unknown"}</span>
                </div>
            `;
        }
        const locBody = document.getElementById("loc-card-body");
        if (locBody) {
            locBody.innerHTML = `
                <div class="data-row">
                    <span class="data-label">City</span>
                    <span class="data-value">${data.city || "Unknown"}</span>
                </div>
                <div class="data-row">
                    <span class="data-label">Region</span>
                    <span class="data-value">${data.region || "Unknown"}</span>
                </div>
            `;
        }
    });

    // DNS e SSL continuam iguais
    fetch("/api/dns", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ target })
    })
    .then(res => res.json())
    .then(data => {
        const dnsBody = document.getElementById("dns-card-body");
        if (dnsBody) {
            dnsBody.innerHTML = (data.ips || []).map(ip => `
                <div class="data-row">
                    <span class="data-label">A Record</span>
                    <span class="data-value">${ip}</span>
                </div>
            `).join("") || '<span class="data-value">None</span>';
        }
        const mxBody = document.getElementById("mx-card-body");
        if (mxBody) {
            mxBody.innerHTML = (data.mx || []).map(mx => `
                <div class="data-row">
                    <span class="data-label">MX</span>
                    <span class="data-value">${mx}</span>
                </div>
            `).join("") || '<span class="data-value">None</span>';
        }
    });

    fetch("/api/ssl", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ target })
    })
    .then(res => res.json())
    .then(data => {
        const container = document.getElementById("ssl-card-body");
        if (!container) return;
        if (!data.active) {
            container.innerHTML = `<span class="data-value status-danger">SSL not detected</span>`;
            return;
        }
        container.innerHTML = `
            <div class="data-row">
                <span class="data-label">Issuer</span>
                <span class="data-value">${data.issuer}</span>
            </div>
            <div class="data-row">
                <span class="data-label">Expiration</span>
                <span class="data-value">${data.days_left} days</span>
            </div>
        `;
    });
}

document.getElementById("target").addEventListener("keydown", (e) => {
    if (e.key === "Enter") {
        lookup();
    }
});
