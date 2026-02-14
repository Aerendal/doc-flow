#!/usr/bin/env python3
import json, sys, pathlib, html

def main(inp, outp):
    data = json.load(open(inp))
    rows = []
    for t in data.get("tasks", []):
        status = t.get("status", "")
        color = "#d1f7c4" if status == "READY" else "#fdd"
        rows.append(f"<tr class='data' style='background:{color}'><td>{html.escape(t.get('id',''))}</td><td>{html.escape(t.get('path',''))}</td><td>{html.escape(status)}</td><td>{html.escape(t.get('cache_status',''))}</td><td>{html.escape(str(t.get('reason','')))}</td><td>{html.escape(', '.join(t.get('violations') or []))}</td></tr>")
    rows_html = "\n".join(rows)
    body = f"""<html><head><meta charset='utf-8'><title>Queue Report</title>
    <style>
      table {{ border-collapse: collapse; }}
      th, td {{ border: 1px solid #ccc; padding: 4px 8px; }}
      th {{ cursor: pointer; }}
    </style>
    <script>
    function sortTable(n) {{
      const table = document.getElementById("qtable");
      let switching = true, dir = "asc", switchcount = 0;
      while (switching) {{
        switching = false;
        const rows = table.rows;
        for (let i = 1; i < rows.length - 1; i++) {{
          let shouldSwitch = false;
          const x = rows[i].getElementsByTagName("TD")[n];
          const y = rows[i+1].getElementsByTagName("TD")[n];
          if ((dir == "asc" && x.innerText.toLowerCase() > y.innerText.toLowerCase()) ||
              (dir == "desc" && x.innerText.toLowerCase() < y.innerText.toLowerCase())) {{
            shouldSwitch = true;
            break;
          }}
          if (shouldSwitch) {{
            rows[i].parentNode.insertBefore(rows[i+1], rows[i]);
            switching = true;
            switchcount ++;
          }}
        }}
        if (switchcount == 0 && dir == "asc") {{ dir = "desc"; switching = true; }}
      }}
    }}
    function filterStatus() {{
      const val = document.getElementById("statusFilter").value.toLowerCase();
      const rows = document.querySelectorAll("#qtable tr.data");
      rows.forEach(r => {{
        const s = r.children[2].innerText.toLowerCase();
        r.style.display = (!val || s == val) ? "" : "none";
      }});
    }}
    function filterCache() {{
      const val = document.getElementById("cacheFilter").value.toLowerCase();
      const rows = document.querySelectorAll("#qtable tr.data");
      rows.forEach(r => {{
        const s = r.children[3].innerText.toLowerCase();
        r.style.display = (!val || s == val) ? "" : "none";
      }});
    }}
    </script>
    </head><body>
    <h1>Queue Report</h1>
    <p>Ready: {data.get('ready',0)} | Blocked: {data.get('blocked',0)}</p>
    <label>Status: <select id="statusFilter" onchange="filterStatus()">
      <option value="">All</option><option value="ready">READY</option><option value="blocked">BLOCKED</option>
    </select></label>
    <label>Cache: <select id="cacheFilter" onchange="filterCache()">
      <option value="">All</option><option value="hit">hit</option><option value="miss">miss</option><option value="nocache">nocache</option>
    </select></label>
    <table id="qtable">
    <tr><th onclick="sortTable(0)">ID</th><th onclick="sortTable(1)">Path</th><th onclick="sortTable(2)">Status</th><th onclick="sortTable(3)">Cache</th><th onclick="sortTable(4)">Reason</th><th onclick="sortTable(5)">Violations</th></tr>
    {rows_html}
    </table>
    </body></html>"""
    pathlib.Path(outp).write_text(body, encoding='utf-8')

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("usage: queue_report_html.py input.json output.html")
        sys.exit(1)
    main(sys.argv[1], sys.argv[2])
