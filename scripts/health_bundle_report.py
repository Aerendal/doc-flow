#!/usr/bin/env python3
import argparse
import json
import pathlib

HTML_TMPL = """
<!doctype html>
<html><head><meta charset='utf-8'><title>docflow health bundle report</title></head>
<body>
<h1>docflow health bundle</h1>
<p><b>root:</b> {root}</p>
<ul>
  <li>health badge: <a href="{health_html}">{health_html}</a></li>
  <li>health report: <a href="{health_report}">{health_report}</a></li>
  <li>queue json: <a href="{queue_json}">{queue_json}</a></li>
  <li>queue html: <a href="{queue_html}">{queue_html}</a></li>
  <li>compliance: <a href="{compliance_json}">{compliance_json}</a></li>
</ul>
</body></html>
"""


def latest(glob):
    files = sorted(pathlib.Path().glob(glob), key=lambda p: p.stat().st_mtime, reverse=True)
    return files[0] if files else None


def main():
    p = argparse.ArgumentParser()
    p.add_argument('--root', required=True, help='root name (dir under dist/health)')
    p.add_argument('--out-html', required=True)
    p.add_argument('--out-json', required=True)
    args = p.parse_args()

    base = pathlib.Path('dist/health') / args.root
    data = {
        'root': args.root,
        'health_html': str((base / 'health_latest.html').resolve()),
        'health_report': str((base / 'health_report_latest.html').resolve()) if (base / 'health_report_latest.html').exists() else 'n/a',
        'queue_json': str((base / 'queue_latest.json').resolve()),
        'queue_html': str((base / 'queue_latest.html').resolve()),
        'compliance_json': str(latest(str(base / 'compliance_*.json')).resolve()) if list(base.glob('compliance_*.json')) else 'n/a',
    }

    pathlib.Path(args.out_html).write_text(HTML_TMPL.format(**data), encoding='utf-8')
    pathlib.Path(args.out_json).write_text(json.dumps(data, indent=2), encoding='utf-8')

if __name__ == '__main__':
    main()
