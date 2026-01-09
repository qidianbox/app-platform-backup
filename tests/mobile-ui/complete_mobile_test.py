#!/usr/bin/env python3
"""
完整的移动端UI自动化测试脚本
测试所有模块界面的移动端适配情况
"""

import os
import time
import json
from datetime import datetime
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

# 配置
BASE_URL = "http://localhost:5173"
API_URL = "http://localhost:8080"
SCREENSHOT_DIR = "/home/ubuntu/app-platform/tests/mobile-ui/module_screenshots"

# 测试设备配置
DEVICES = {
    "iPhone_12": {"width": 390, "height": 844, "pixel_ratio": 3},
    "iPhone_SE": {"width": 375, "height": 667, "pixel_ratio": 2},
}

# 测试页面配置
TEST_PAGES = [
    {"name": "登录页面", "path": "/login", "need_login": False},
    {"name": "仪表盘", "path": "/dashboard", "need_login": True},
    {"name": "APP列表", "path": "/apps", "need_login": True},
    {"name": "模块管理", "path": "/modules", "need_login": True},
    # APP配置子页面
    {"name": "APP配置-概览", "path": "/apps/1/config/dashboard", "need_login": True},
    {"name": "APP配置-用户中心", "path": "/apps/1/config/user", "need_login": True},
    {"name": "APP配置-消息推送", "path": "/apps/1/config/message", "need_login": True},
    {"name": "APP配置-支付", "path": "/apps/1/config/payment", "need_login": True},
    {"name": "APP配置-数据统计", "path": "/apps/1/config/analytics", "need_login": True},
    {"name": "APP配置-安全", "path": "/apps/1/config/security", "need_login": True},
    {"name": "APP配置-版本", "path": "/apps/1/config/version", "need_login": True},
    {"name": "APP配置-模块", "path": "/apps/1/config/modules", "need_login": True},
]

class MobileUITester:
    def __init__(self):
        self.results = []
        self.driver = None
        
    def setup_driver(self, device_name, device_config):
        """设置Chrome驱动为移动端模式"""
        options = Options()
        options.add_argument('--headless')
        options.add_argument('--no-sandbox')
        options.add_argument('--disable-dev-shm-usage')
        options.add_argument('--disable-gpu')
        
        # 设置移动端模拟
        mobile_emulation = {
            "deviceMetrics": {
                "width": device_config["width"],
                "height": device_config["height"],
                "pixelRatio": device_config["pixel_ratio"]
            },
            "userAgent": "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15"
        }
        options.add_experimental_option("mobileEmulation", mobile_emulation)
        
        self.driver = webdriver.Chrome(options=options)
        self.driver.set_page_load_timeout(30)
        
    def login(self):
        """执行登录"""
        try:
            self.driver.get(f"{BASE_URL}/login")
            time.sleep(2)
            
            # 设置token直接登录
            self.driver.execute_script(f"""
                localStorage.setItem('token', 'test_token');
            """)
            return True
        except Exception as e:
            print(f"登录失败: {e}")
            return False
    
    def take_screenshot(self, device_name, page_name):
        """截图"""
        device_dir = os.path.join(SCREENSHOT_DIR, device_name)
        os.makedirs(device_dir, exist_ok=True)
        
        timestamp = datetime.now().strftime("%H%M%S")
        filename = f"{page_name}_{timestamp}.png"
        filepath = os.path.join(device_dir, filename)
        
        self.driver.save_screenshot(filepath)
        return filepath
    
    def check_ui_issues(self):
        """检查UI问题"""
        issues = []
        
        try:
            # 检查水平滚动
            scroll_width = self.driver.execute_script("return document.body.scrollWidth")
            client_width = self.driver.execute_script("return document.body.clientWidth")
            if scroll_width > client_width + 10:
                issues.append(f"水平溢出: scrollWidth={scroll_width}, clientWidth={client_width}")
            
            # 检查元素溢出
            overflow_elements = self.driver.execute_script("""
                var elements = document.querySelectorAll('*');
                var overflowing = [];
                for (var i = 0; i < elements.length; i++) {
                    var el = elements[i];
                    if (el.scrollWidth > el.clientWidth + 5) {
                        overflowing.push(el.tagName + '.' + el.className);
                    }
                }
                return overflowing.slice(0, 5);
            """)
            if overflow_elements:
                issues.append(f"溢出元素: {overflow_elements}")
                
        except Exception as e:
            issues.append(f"检查异常: {str(e)}")
        
        return issues
    
    def test_page(self, device_name, page_config):
        """测试单个页面"""
        page_name = page_config["name"]
        page_path = page_config["path"]
        
        result = {
            "device": device_name,
            "page": page_name,
            "path": page_path,
            "status": "PASS",
            "issues": [],
            "screenshot": ""
        }
        
        try:
            self.driver.get(f"{BASE_URL}{page_path}")
            time.sleep(2)
            
            # 检查UI问题
            issues = self.check_ui_issues()
            if issues:
                result["issues"] = issues
                result["status"] = "WARNING"
            
            # 截图
            result["screenshot"] = self.take_screenshot(device_name, page_name.replace("/", "_"))
            
        except Exception as e:
            result["status"] = "FAIL"
            result["issues"].append(str(e))
        
        return result
    
    def run_tests(self):
        """运行所有测试"""
        print("=" * 60)
        print("开始移动端UI完整测试")
        print("=" * 60)
        
        for device_name, device_config in DEVICES.items():
            print(f"\n测试设备: {device_name} ({device_config['width']}x{device_config['height']})")
            print("-" * 40)
            
            self.setup_driver(device_name, device_config)
            
            # 登录
            self.login()
            
            for page_config in TEST_PAGES:
                result = self.test_page(device_name, page_config)
                self.results.append(result)
                
                status_icon = "✅" if result["status"] == "PASS" else ("⚠️" if result["status"] == "WARNING" else "❌")
                print(f"  {status_icon} {result['page']}: {result['status']}")
                if result["issues"]:
                    for issue in result["issues"]:
                        print(f"      - {issue}")
            
            self.driver.quit()
        
        self.generate_report()
    
    def generate_report(self):
        """生成测试报告"""
        report_path = os.path.join(SCREENSHOT_DIR, "MODULE_TEST_REPORT.md")
        
        # 统计
        total = len(self.results)
        passed = sum(1 for r in self.results if r["status"] == "PASS")
        warnings = sum(1 for r in self.results if r["status"] == "WARNING")
        failed = sum(1 for r in self.results if r["status"] == "FAIL")
        
        with open(report_path, "w") as f:
            f.write("# 移动端模块界面测试报告\n\n")
            f.write(f"**测试时间**: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n\n")
            
            f.write("## 测试概览\n\n")
            f.write("| 统计项 | 数量 |\n")
            f.write("|:---|:---|\n")
            f.write(f"| ✅ 通过 | {passed} |\n")
            f.write(f"| ⚠️ 警告 | {warnings} |\n")
            f.write(f"| ❌ 失败 | {failed} |\n")
            f.write(f"| **总计** | {total} |\n\n")
            
            f.write("## 测试设备\n\n")
            for device_name, config in DEVICES.items():
                f.write(f"- **{device_name}**: {config['width']}x{config['height']}\n")
            f.write("\n")
            
            f.write("## 详细结果\n\n")
            
            for device_name in DEVICES.keys():
                f.write(f"### {device_name}\n\n")
                f.write("| 页面 | 状态 | 问题 |\n")
                f.write("|:---|:---|:---|\n")
                
                device_results = [r for r in self.results if r["device"] == device_name]
                for r in device_results:
                    status_icon = "✅" if r["status"] == "PASS" else ("⚠️" if r["status"] == "WARNING" else "❌")
                    issues_str = ", ".join(r["issues"][:2]) if r["issues"] else "-"
                    f.write(f"| {r['page']} | {status_icon} {r['status']} | {issues_str} |\n")
                f.write("\n")
            
            f.write("## 截图\n\n")
            for device_name in DEVICES.keys():
                f.write(f"### {device_name}\n\n")
                device_results = [r for r in self.results if r["device"] == device_name]
                for r in device_results:
                    if r["screenshot"]:
                        rel_path = os.path.relpath(r["screenshot"], SCREENSHOT_DIR)
                        f.write(f"#### {r['page']}\n")
                        f.write(f"![{r['page']}]({rel_path})\n\n")
        
        print(f"\n测试报告已生成: {report_path}")
        print(f"截图目录: {SCREENSHOT_DIR}")

if __name__ == "__main__":
    os.makedirs(SCREENSHOT_DIR, exist_ok=True)
    tester = MobileUITester()
    tester.run_tests()
