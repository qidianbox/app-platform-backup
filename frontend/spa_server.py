#!/usr/bin/env python3
"""
Simple HTTP server for Single Page Applications (SPA)
Falls back to index.html for any route not found
"""

import http.server
import os
import socketserver

PORT = 5174
DIRECTORY = "dist"

class SPAHandler(http.server.SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=DIRECTORY, **kwargs)
    
    def do_GET(self):
        # Get the file path
        path = self.translate_path(self.path)
        
        # If the path exists and is a file, serve it
        if os.path.exists(path) and os.path.isfile(path):
            return super().do_GET()
        
        # If it's a directory with index.html, serve that
        if os.path.isdir(path):
            index_path = os.path.join(path, 'index.html')
            if os.path.exists(index_path):
                return super().do_GET()
        
        # For SPA routes, serve the root index.html
        # But not for actual file requests (with extensions)
        if '.' not in os.path.basename(self.path) or self.path.endswith('/'):
            self.path = '/index.html'
            return super().do_GET()
        
        # Otherwise, return 404
        return super().do_GET()

if __name__ == "__main__":
    with socketserver.TCPServer(("", PORT), SPAHandler) as httpd:
        print(f"SPA Server running on http://localhost:{PORT}/")
        httpd.serve_forever()
