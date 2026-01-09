#!/usr/bin/env python3
"""
静态文件服务器，支持SPA路由和所有HTTP方法
"""
import http.server
import socketserver
import os
import sys

PORT = 5173
DIRECTORY = "dist"

class SPAHandler(http.server.SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=DIRECTORY, **kwargs)
    
    def do_GET(self):
        # 检查文件是否存在
        path = self.translate_path(self.path)
        if not os.path.exists(path) or os.path.isdir(path):
            # 如果是API请求，返回404
            if self.path.startswith('/api/'):
                self.send_error(404, "API endpoint not found - please use backend server")
                return
            # 对于SPA路由，返回index.html
            if not self.path.startswith('/assets/'):
                self.path = '/index.html'
        return super().do_GET()
    
    def end_headers(self):
        # 添加缓存控制头，防止浏览器缓存
        self.send_header('Cache-Control', 'no-cache, no-store, must-revalidate')
        self.send_header('Pragma', 'no-cache')
        self.send_header('Expires', '0')
        # 添加CORS头
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', '*')
        super().end_headers()

if __name__ == "__main__":
    os.chdir(os.path.dirname(os.path.abspath(__file__)))
    
    with socketserver.TCPServer(("0.0.0.0", PORT), SPAHandler) as httpd:
        print(f"Serving {DIRECTORY} at http://0.0.0.0:{PORT}")
        print("Press Ctrl+C to stop")
        try:
            httpd.serve_forever()
        except KeyboardInterrupt:
            print("\nServer stopped")
