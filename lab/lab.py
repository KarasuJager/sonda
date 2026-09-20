from http.server import BaseHTTPRequestHandler, HTTPServer
from urllib.parse import urlparse, parse_qs


class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        parsed = urlparse(self.path)
        params = parse_qs(parsed.query)

        query = params.get("q", [""])[0]

        if query == "admin":
            body = b"user=admin;role=administrator\n"
        elif query == "SONDA_TEST":
            body = b"no results\n"
        else:
            body = f"search={query}\n".encode()

        self.send_response(200)
        self.send_header("Content-Type", "text/plain")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()

        self.wfile.write(body)


server = HTTPServer(("127.0.0.1", 9090), Handler)

print("SONDA laboratory listening on http://127.0.0.1:9090")
server.serve_forever()

