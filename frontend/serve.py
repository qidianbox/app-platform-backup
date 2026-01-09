#!/usr/bin/env python3
"""
SPA静态文件服务器 - 支持Vue Router history模式
所有非静态文件请求都返回index.html
"""

import http.server
import socketserver
import os

PORT = 5173
DIRECTORY = "dist"

class SPAHandler(http.server.SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=DIRECTORY, **kwargs)
    
    def do_GET(self):
        # 获取请求的文件路径
        path = self.translate_path(self.path)
        
        # 如果是静态资源文件（存在的文件），直接返回
        if os.path.isfile(path):
            return super().do_GET()
        
        # 如果是目录且包含index.html，返回index.html
        if os.path.isdir(path):
            index_path = os.path.join(path, 'index.html')
            if os.path.isfile(index_path):
                return super().do_GET()
        
        # 其他所有请求都返回根目录的index.html（SPA路由）
        self.path = '/index.html'
        return super().do_GET()

if __name__ == "__main__":
    os.chdir(os.path.dirname(os.path.abspath(__file__)))
    
    with socketserver.TCPServer(("", PORT), SPAHandler) as httpd:
        print(f"SPA Server running at http://localhost:{PORT}")
        httpd.serve_forever()
