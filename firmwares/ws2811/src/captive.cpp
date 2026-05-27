#include "captive.h"
#include "state.h"

static const char CAPTIVE_HTML[] PROGMEM = R"HTML(<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Smart Home setup</title>
<style>
  *,*::before,*::after{box-sizing:border-box}
  html,body{margin:0;padding:0;background:#0f0f10;color:#e5e5e5;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;min-height:100%}
  body{display:flex;justify-content:center;padding:24px 16px;-webkit-font-smoothing:antialiased}
  .card{width:100%;max-width:380px;background:#1a1a1a;border:1px solid #2a2a2a;border-radius:16px;padding:24px;box-shadow:0 24px 64px rgba(0,0,0,.5)}
  .brand{display:flex;align-items:center;gap:10px;margin-bottom:6px}
  .brand-icon{width:36px;height:36px;border-radius:10px;background:#6366f1;display:flex;align-items:center;justify-content:center}
  .brand-icon svg{width:20px;height:20px;stroke:#fff;fill:none;stroke-width:2;stroke-linecap:round;stroke-linejoin:round}
  .title{font-size:18px;font-weight:600;color:#fff;line-height:1.2}
  .subtitle{font-size:13px;color:#737373;margin-top:2px}
  .meta{margin-top:16px;padding:10px 12px;background:#111;border:1px solid #222;border-radius:10px;font-family:ui-monospace,SFMono-Regular,Menlo,monospace;font-size:12px;color:#a3a3a3;word-break:break-all}
  .meta-label{color:#525252;font-size:11px;text-transform:uppercase;letter-spacing:.5px;margin-right:6px}
  form{margin-top:20px;display:flex;flex-direction:column;gap:14px}
  .field{display:flex;flex-direction:column;gap:6px}
  label{font-size:12px;color:#a3a3a3;font-weight:500}
  .input-wrap{position:relative;display:flex;align-items:center}
  input,select{width:100%;background:#111;border:1px solid #2a2a2a;color:#e5e5e5;border-radius:10px;padding:11px 12px;font-size:14px;font-family:inherit;outline:none;transition:border-color .15s,background .15s;-webkit-appearance:none;appearance:none}
  input:focus,select:focus{border-color:#6366f1;background:#141416}
  input::placeholder{color:#404040}
  .scan-btn{position:absolute;right:6px;height:30px;padding:0 10px;font-size:11px;color:#a3a3a3;background:#1f1f22;border:1px solid #2a2a2a;border-radius:7px;cursor:pointer;font-family:inherit;text-transform:uppercase;letter-spacing:.5px;transition:background .15s,color .15s}
  .scan-btn:hover{background:#26262a;color:#e5e5e5}
  .scan-btn:disabled{opacity:.5;cursor:default}
  .hint{font-size:11px;color:#525252;margin-top:2px}
  button.primary{width:100%;background:#6366f1;color:#fff;border:1px solid #6366f1;border-radius:10px;padding:12px;font-size:14px;font-weight:600;cursor:pointer;font-family:inherit;margin-top:6px;transition:background .15s}
  button.primary:hover{background:#4f46e5}
  button.primary:active{transform:translateY(1px)}
  .footer{margin-top:18px;font-size:11px;color:#404040;text-align:center;line-height:1.5}
  .spinner{display:inline-block;width:12px;height:12px;border:2px solid #2a2a2a;border-top-color:#6366f1;border-radius:50%;animation:spin .7s linear infinite;vertical-align:-2px;margin-right:6px}
  @keyframes spin{to{transform:rotate(360deg)}}
</style>
</head>
<body>
<div class="card">
  <div class="brand">
    <div class="brand-icon">
      <svg viewBox="0 0 24 24"><path d="M3 12l9-9 9 9"/><path d="M5 10v10h14V10"/></svg>
    </div>
    <div>
      <div class="title">Smart Home</div>
      <div class="subtitle">Device setup</div>
    </div>
  </div>

  <div class="meta">
    <span class="meta-label">Device</span><span id="uid">__UID__</span>
  </div>

  <form method="POST" action="/config" id="form">
    <div class="field">
      <label for="name">Device name</label>
      <input id="name" name="name" placeholder="e.g. Living room strip" maxlength="60" autocomplete="off" required>
      <div class="hint">How it will show up in the app</div>
    </div>

    <div class="field">
      <label for="ssid">Wi-Fi network</label>
      <div class="input-wrap">
        <input id="ssid" name="ssid" list="ssids" placeholder="Network name" autocomplete="off" autocapitalize="off" autocorrect="off" required>
        <button type="button" class="scan-btn" id="scanBtn">Scan</button>
      </div>
      <datalist id="ssids"></datalist>
    </div>

    <div class="field">
      <label for="pass">Wi-Fi password</label>
      <input id="pass" name="pass" type="password" placeholder="Leave blank for open networks" autocomplete="off">
    </div>

    <div class="field">
      <label for="code">Pairing code</label>
      <input id="code" name="code" placeholder="From the app" autocapitalize="characters" autocomplete="off" required>
      <div class="hint">Open Smart Home &rarr; Pair Device to get the code</div>
    </div>

    <button type="submit" class="primary" id="submit">Connect &amp; pair</button>
  </form>

  <div class="footer">The device will reboot after saving.</div>
</div>

<script>
(function(){
  var btn = document.getElementById('scanBtn');
  var dl  = document.getElementById('ssids');
  var ssid = document.getElementById('ssid');

  function loadScan(){
    btn.disabled = true;
    btn.textContent = '...';
    fetch('/scan').then(function(r){return r.json();}).then(function(list){
      dl.innerHTML = '';
      (list||[]).forEach(function(n){
        var o = document.createElement('option');
        o.value = n;
        dl.appendChild(o);
      });
      btn.textContent = (list && list.length) ? (list.length + ' found') : 'Scan';
      setTimeout(function(){ btn.disabled = false; btn.textContent = 'Scan'; }, 1500);
      if (list && list.length && !ssid.value) ssid.focus();
    }).catch(function(){
      btn.disabled = false;
      btn.textContent = 'Scan';
    });
  }

  btn.addEventListener('click', loadScan);
  // Auto-scan once on load (best effort)
  setTimeout(loadScan, 200);

  document.getElementById('form').addEventListener('submit', function(){
    var b = document.getElementById('submit');
    b.disabled = true;
    b.innerHTML = '<span class="spinner"></span>Saving...';
  });
})();
</script>
</body>
</html>)HTML";

static const char SAVED_HTML[] PROGMEM = R"HTML(<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Saved</title>
<style>
  *,*::before,*::after{box-sizing:border-box}
  html,body{margin:0;padding:0;background:#0f0f10;color:#e5e5e5;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;min-height:100%}
  body{display:flex;align-items:center;justify-content:center;padding:24px 16px;min-height:100vh}
  .card{width:100%;max-width:380px;background:#1a1a1a;border:1px solid #2a2a2a;border-radius:16px;padding:32px 24px;text-align:center}
  .icon{width:56px;height:56px;border-radius:14px;background:#052e16;display:inline-flex;align-items:center;justify-content:center;margin-bottom:16px}
  .icon svg{width:30px;height:30px;stroke:#4ade80;fill:none;stroke-width:2.5;stroke-linecap:round;stroke-linejoin:round}
  .title{font-size:18px;font-weight:600;color:#fff;margin:0 0 6px}
  .sub{font-size:14px;color:#737373;margin:0;line-height:1.5}
  .spinner{display:inline-block;width:14px;height:14px;border:2px solid #2a2a2a;border-top-color:#6366f1;border-radius:50%;animation:spin .7s linear infinite;vertical-align:-3px;margin-right:6px}
  @keyframes spin{to{transform:rotate(360deg)}}
  .footer{margin-top:20px;font-size:12px;color:#525252}
</style>
</head>
<body>
<div class="card">
  <div class="icon">
    <svg viewBox="0 0 24 24"><path d="M5 13l4 4L19 7"/></svg>
  </div>
  <h1 class="title">Settings saved</h1>
  <p class="sub">The device is rebooting and will join your network in a moment.</p>
  <div class="footer"><span class="spinner"></span>You can close this page.</div>
</div>
</body>
</html>)HTML";

static String renderCaptive() {
  String out = FPSTR(CAPTIVE_HTML);
  out.replace("__UID__", deviceUID);
  return out;
}

static String escapeJSON(const String &s) {
  String r;
  r.reserve(s.length() + 4);
  for (size_t i = 0; i < s.length(); ++i) {
    char c = s[i];
    switch (c) {
      case '"':  r += "\\\""; break;
      case '\\': r += "\\\\"; break;
      case '\n': r += "\\n";  break;
      case '\r': r += "\\r";  break;
      case '\t': r += "\\t";  break;
      default:
        if ((uint8_t)c < 0x20) {
          char buf[8];
          snprintf(buf, sizeof(buf), "\\u%04x", c);
          r += buf;
        } else {
          r += c;
        }
    }
  }
  return r;
}

static void handleScan() {
  LOG("[AP] GET /scan");
  int n = WiFi.scanNetworks(false, true);
  String body = "[";
  if (n > 0) {
    // Deduplicate SSIDs (multiple BSSIDs per SSID are common).
    String seen = "|";
    bool first = true;
    for (int i = 0; i < n; ++i) {
      String s = WiFi.SSID(i);
      if (s.length() == 0) continue;
      String key = "|" + s + "|";
      if (seen.indexOf(key) >= 0) continue;
      seen += s + "|";
      if (!first) body += ",";
      body += "\"" + escapeJSON(s) + "\"";
      first = false;
    }
  }
  body += "]";
  WiFi.scanDelete();
  server.sendHeader("Cache-Control", "no-store");
  server.send(200, "application/json", body);
}

void startCaptive(){
  LOG("[AP] Starting captive portal");

  deviceMode = MODE_PAIRING;

  WiFi.disconnect(true,true);
  WiFi.mode(WIFI_AP_STA);
  WiFi.softAPConfig(AP_IP_ADDR,AP_IP_ADDR,IPAddress(255,255,255,0));

  int start = 0;
  if(deviceUID.length() > 6) start = (int)deviceUID.length() - 6;
  String ap = "SmartHome-" + deviceUID.substring(start);

  WiFi.softAP(ap.c_str(), NULL);

  LOG("[AP] AP started: %s ip=%s", ap.c_str(), WiFi.softAPIP().toString().c_str());

  dnsServer.start(DNS_PORT,"*",AP_IP_ADDR);
  LOG("[AP] DNS captive started");

  server.on("/",HTTP_GET,[](){
    LOG("[AP] GET /");
    server.send(200,"text/html; charset=utf-8",renderCaptive());
  });

  server.on("/hotspot-detect.html",HTTP_GET,[](){
    LOG("[AP] iOS captive check");
    server.send(200,"text/html; charset=utf-8",renderCaptive());
  });

  server.on("/generate_204",HTTP_GET,[](){
    LOG("[AP] Android captive check");
    server.sendHeader("Location","/",true);
    server.send(302);
  });

  server.on("/scan", HTTP_GET, handleScan);

  server.on("/config",HTTP_POST,[](){
    String name = server.arg("name");
    String ssid = server.arg("ssid");
    String pass = server.arg("pass");
    String code = server.arg("code");

    LOG("[AP] Received config name='%s' ssid='%s' code='%s'", name.c_str(), ssid.c_str(), code.c_str());

    prefs.putString("name", name);
    prefs.putString("ssid", ssid);
    prefs.putString("pass", pass);
    prefs.putString("pair_code", code);
    prefs.putBool("paired", false);

    server.send(200,"text/html; charset=utf-8", FPSTR(SAVED_HTML));
    delay(500);
    ESP.restart();
  });

  server.onNotFound([](){
    server.sendHeader("Location","/",true);
    server.send(302);
  });

  server.begin();
  LOG("[AP] Web server started");
}
