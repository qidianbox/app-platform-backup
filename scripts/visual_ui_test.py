#!/usr/bin/env python3
"""
自动视觉UI测试脚本
覆盖所有页面的截图和视觉回归测试
"""

import os
import json
import time
import hashlib
from datetime import datetime
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from PIL import Image, ImageChops, ImageDraw, ImageFont
import math

# 配置
BASE_URL = "http://localhost:5173"
RESULTS_DIR = "/home/ubuntu/app-platform/test_results/visual_test"
BASELINE_DIR = os.path.join(RESULTS_DIR, "baseline")
CURRENT_DIR = os.path.join(RESULTS_DIR, "current")
DIFF_DIR = os.path.join(RESULTS_DIR, "diff")

# 测试页面配置
TEST_PAGES = [
    {
        "name": "登录页",
        "path": "/login",
        "requires_auth": False,
        "wait_for": ".login-container, .login-form, form"
    },
    {
        "name": "APP列表页",
        "path": "/apps",
        "requires_auth": True,
        "wait_for": ".app-list, .apps-container, table"
    },
    {
        "name": "APP详情-概览",
        "path": "/apps/2/config",
        "requires_auth": True,
        "wait_for": ".page-content, .stats-cards"
    },
    {
        "name": "APP详情-基础配置",
        "path": "/apps/2/config",
        "requires_auth": True,
        "wait_for": ".config-form",
        "action": "click_basic_config"
    },
    {
        "name": "APP详情-工作台",
        "path": "/apps/2/config",
        "requires_auth": True,
        "wait_for": ".workspace-content",
        "action": "click_workspace"
    }
]

# 设备配置
DEVICES = [
    {"name": "Desktop_1920x1080", "width": 1920, "height": 1080},
    {"name": "Laptop_1366x768", "width": 1366, "height": 768},
    {"name": "Tablet_768x1024", "width": 768, "height": 1024},
    {"name": "Mobile_375x667", "width": 375, "height": 667}
]

class VisualUITester:
    def __init__(self):
        self.driver = None
        self.results = []
        self.timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        
        # 创建目录
        for dir_path in [RESULTS_DIR, BASELINE_DIR, CURRENT_DIR, DIFF_DIR]:
            os.makedirs(dir_path, exist_ok=True)
    
    def setup_driver(self, width, height):
        """设置浏览器驱动"""
        chrome_options = Options()
        chrome_options.add_argument("--headless")
        chrome_options.add_argument("--no-sandbox")
        chrome_options.add_argument("--disable-dev-shm-usage")
        chrome_options.add_argument("--disable-gpu")
        chrome_options.add_argument(f"--window-size={width},{height}")
        chrome_options.add_argument("--force-device-scale-factor=1")
        
        service = Service("/usr/bin/chromedriver")
        self.driver = webdriver.Chrome(service=service, options=chrome_options)
        self.driver.set_window_size(width, height)
    
    def login(self):
        """执行登录"""
        try:
            self.driver.get(f"{BASE_URL}/login")
            time.sleep(2)
            
            # 查找并填写登录表单
            username_input = self.driver.find_element(By.CSS_SELECTOR, "input[type='text'], input[placeholder*='用户'], input[name='username']")
            password_input = self.driver.find_element(By.CSS_SELECTOR, "input[type='password']")
            
            username_input.clear()
            username_input.send_keys("admin")
            password_input.clear()
            password_input.send_keys("admin123")
            
            # 点击登录按钮
            login_btn = self.driver.find_element(By.CSS_SELECTOR, "button[type='submit'], .login-btn, button.el-button--primary")
            login_btn.click()
            
            time.sleep(3)
            return True
        except Exception as e:
            print(f"登录失败: {e}")
            return False
    
    def wait_for_element(self, selector, timeout=10):
        """等待元素出现"""
        try:
            selectors = selector.split(", ")
            for sel in selectors:
                try:
                    WebDriverWait(self.driver, timeout).until(
                        EC.presence_of_element_located((By.CSS_SELECTOR, sel.strip()))
                    )
                    return True
                except:
                    continue
            return False
        except:
            return False
    
    def execute_action(self, action):
        """执行页面操作"""
        try:
            if action == "click_basic_config":
                # 点击基础配置菜单
                elem = self.driver.find_element(By.XPATH, "//*[contains(text(), '基础配置')]")
                elem.click()
                time.sleep(1)
            elif action == "click_workspace":
                # 点击工作台Tab
                elem = self.driver.find_element(By.XPATH, "//*[contains(text(), '工作台')]")
                elem.click()
                time.sleep(1)
        except Exception as e:
            print(f"执行操作失败: {action}, {e}")
    
    def take_screenshot(self, name, device_name):
        """截图"""
        filename = f"{name}_{device_name}_{self.timestamp}.png"
        filepath = os.path.join(CURRENT_DIR, filename)
        self.driver.save_screenshot(filepath)
        return filepath
    
    def calculate_image_diff(self, img1_path, img2_path):
        """计算两张图片的差异"""
        try:
            img1 = Image.open(img1_path).convert('RGB')
            img2 = Image.open(img2_path).convert('RGB')
            
            # 调整尺寸一致
            if img1.size != img2.size:
                img2 = img2.resize(img1.size, Image.Resampling.LANCZOS)
            
            # 计算差异
            diff = ImageChops.difference(img1, img2)
            
            # 计算差异百分比
            diff_pixels = 0
            total_pixels = img1.size[0] * img1.size[1]
            
            for pixel in diff.getdata():
                if pixel != (0, 0, 0):
                    diff_pixels += 1
            
            diff_percentage = (diff_pixels / total_pixels) * 100
            
            return diff_percentage, diff
        except Exception as e:
            print(f"计算图片差异失败: {e}")
            return -1, None
    
    def create_diff_image(self, img1_path, img2_path, diff_path):
        """创建差异对比图"""
        try:
            img1 = Image.open(img1_path).convert('RGB')
            img2 = Image.open(img2_path).convert('RGB')
            
            if img1.size != img2.size:
                img2 = img2.resize(img1.size, Image.Resampling.LANCZOS)
            
            # 创建差异图
            diff = ImageChops.difference(img1, img2)
            
            # 增强差异可见度
            diff = diff.point(lambda x: min(255, x * 10))
            
            # 创建并排对比图
            width = img1.size[0] * 3
            height = img1.size[1]
            comparison = Image.new('RGB', (width, height))
            
            comparison.paste(img1, (0, 0))
            comparison.paste(img2, (img1.size[0], 0))
            comparison.paste(diff, (img1.size[0] * 2, 0))
            
            # 添加标签
            draw = ImageDraw.Draw(comparison)
            try:
                font = ImageFont.truetype("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 20)
            except:
                font = ImageFont.load_default()
            
            draw.text((10, 10), "Baseline", fill="white", font=font)
            draw.text((img1.size[0] + 10, 10), "Current", fill="white", font=font)
            draw.text((img1.size[0] * 2 + 10, 10), "Diff", fill="white", font=font)
            
            comparison.save(diff_path)
            return True
        except Exception as e:
            print(f"创建差异图失败: {e}")
            return False
    
    def run_test(self, page, device):
        """运行单个测试"""
        result = {
            "page": page["name"],
            "device": device["name"],
            "status": "unknown",
            "diff_percentage": 0,
            "screenshot": "",
            "baseline": "",
            "diff_image": "",
            "error": ""
        }
        
        try:
            # 设置浏览器
            self.setup_driver(device["width"], device["height"])
            
            # 登录（如果需要）
            if page["requires_auth"]:
                if not self.login():
                    result["status"] = "failed"
                    result["error"] = "登录失败"
                    return result
            
            # 访问页面
            self.driver.get(f"{BASE_URL}{page['path']}")
            time.sleep(2)
            
            # 等待元素
            self.wait_for_element(page["wait_for"])
            
            # 执行操作（如果有）
            if "action" in page:
                self.execute_action(page["action"])
                time.sleep(1)
            
            # 截图
            screenshot_path = self.take_screenshot(
                page["name"].replace(" ", "_").replace("-", "_"),
                device["name"]
            )
            result["screenshot"] = screenshot_path
            
            # 查找基准图
            baseline_pattern = f"{page['name'].replace(' ', '_').replace('-', '_')}_{device['name']}_"
            baseline_files = [f for f in os.listdir(BASELINE_DIR) if f.startswith(baseline_pattern)]
            
            if baseline_files:
                # 有基准图，进行对比
                baseline_path = os.path.join(BASELINE_DIR, sorted(baseline_files)[-1])
                result["baseline"] = baseline_path
                
                diff_percentage, _ = self.calculate_image_diff(baseline_path, screenshot_path)
                result["diff_percentage"] = round(diff_percentage, 2)
                
                if diff_percentage < 0:
                    result["status"] = "error"
                    result["error"] = "图片对比失败"
                elif diff_percentage < 1:
                    result["status"] = "passed"
                elif diff_percentage < 5:
                    result["status"] = "warning"
                else:
                    result["status"] = "failed"
                
                # 创建差异图
                if diff_percentage > 0:
                    diff_filename = f"diff_{page['name'].replace(' ', '_')}_{device['name']}_{self.timestamp}.png"
                    diff_path = os.path.join(DIFF_DIR, diff_filename)
                    self.create_diff_image(baseline_path, screenshot_path, diff_path)
                    result["diff_image"] = diff_path
            else:
                # 没有基准图，保存为新基准
                import shutil
                baseline_filename = f"{page['name'].replace(' ', '_').replace('-', '_')}_{device['name']}_{self.timestamp}.png"
                baseline_path = os.path.join(BASELINE_DIR, baseline_filename)
                shutil.copy(screenshot_path, baseline_path)
                result["baseline"] = baseline_path
                result["status"] = "new_baseline"
            
        except Exception as e:
            result["status"] = "error"
            result["error"] = str(e)
        finally:
            if self.driver:
                self.driver.quit()
                self.driver = None
        
        return result
    
    def run_all_tests(self):
        """运行所有测试"""
        print("=" * 60)
        print("开始视觉UI测试")
        print("=" * 60)
        
        total_tests = len(TEST_PAGES) * len(DEVICES)
        completed = 0
        
        for page in TEST_PAGES:
            for device in DEVICES:
                completed += 1
                print(f"\n[{completed}/{total_tests}] 测试: {page['name']} @ {device['name']}")
                
                result = self.run_test(page, device)
                self.results.append(result)
                
                status_emoji = {
                    "passed": "✅",
                    "warning": "⚠️",
                    "failed": "❌",
                    "new_baseline": "🆕",
                    "error": "💥"
                }.get(result["status"], "❓")
                
                print(f"   状态: {status_emoji} {result['status']}")
                if result["diff_percentage"] > 0:
                    print(f"   差异: {result['diff_percentage']}%")
                if result["error"]:
                    print(f"   错误: {result['error']}")
        
        return self.results
    
    def generate_report(self):
        """生成测试报告"""
        report_path = os.path.join(RESULTS_DIR, f"visual_test_report_{self.timestamp}.md")
        
        # 统计
        total = len(self.results)
        passed = len([r for r in self.results if r["status"] == "passed"])
        warnings = len([r for r in self.results if r["status"] == "warning"])
        failed = len([r for r in self.results if r["status"] == "failed"])
        new_baselines = len([r for r in self.results if r["status"] == "new_baseline"])
        errors = len([r for r in self.results if r["status"] == "error"])
        
        pass_rate = (passed / total * 100) if total > 0 else 0
        
        report = f"""# 视觉UI测试报告

**生成时间**: {datetime.now().strftime("%Y-%m-%d %H:%M:%S")}

## 测试概览

| 指标 | 数值 |
|------|------|
| 总测试数 | {total} |
| 通过 | {passed} ✅ |
| 警告 | {warnings} ⚠️ |
| 失败 | {failed} ❌ |
| 新基准 | {new_baselines} 🆕 |
| 错误 | {errors} 💥 |
| **通过率** | **{pass_rate:.1f}%** |

## 测试页面

| 页面 | 路径 | 需要认证 |
|------|------|----------|
"""
        for page in TEST_PAGES:
            auth = "是" if page["requires_auth"] else "否"
            report += f"| {page['name']} | {page['path']} | {auth} |\n"
        
        report += f"""
## 测试设备

| 设备 | 分辨率 |
|------|--------|
"""
        for device in DEVICES:
            report += f"| {device['name']} | {device['width']}x{device['height']} |\n"
        
        report += """
## 详细结果

"""
        # 按页面分组
        pages_results = {}
        for result in self.results:
            page_name = result["page"]
            if page_name not in pages_results:
                pages_results[page_name] = []
            pages_results[page_name].append(result)
        
        for page_name, results in pages_results.items():
            report += f"### {page_name}\n\n"
            report += "| 设备 | 状态 | 差异 | 说明 |\n"
            report += "|------|------|------|------|\n"
            
            for r in results:
                status_emoji = {
                    "passed": "✅ 通过",
                    "warning": "⚠️ 警告",
                    "failed": "❌ 失败",
                    "new_baseline": "🆕 新基准",
                    "error": "💥 错误"
                }.get(r["status"], "❓ 未知")
                
                diff = f"{r['diff_percentage']}%" if r["diff_percentage"] > 0 else "-"
                note = r["error"] if r["error"] else "-"
                
                report += f"| {r['device']} | {status_emoji} | {diff} | {note} |\n"
            
            report += "\n"
        
        # 失败和警告详情
        issues = [r for r in self.results if r["status"] in ["failed", "warning"]]
        if issues:
            report += """## 需要关注的问题

"""
            for issue in issues:
                report += f"""### {issue['page']} @ {issue['device']}

- **状态**: {issue['status']}
- **差异百分比**: {issue['diff_percentage']}%
- **当前截图**: `{os.path.basename(issue['screenshot'])}`
- **基准截图**: `{os.path.basename(issue['baseline'])}`
"""
                if issue["diff_image"]:
                    report += f"- **差异图**: `{os.path.basename(issue['diff_image'])}`\n"
                report += "\n"
        
        report += f"""
## 文件位置

- **基准图目录**: `{BASELINE_DIR}`
- **当前截图目录**: `{CURRENT_DIR}`
- **差异图目录**: `{DIFF_DIR}`

## 使用说明

1. **首次运行**: 会自动创建基准截图
2. **后续运行**: 与基准截图对比，检测视觉变化
3. **更新基准**: 将current目录中的截图复制到baseline目录
4. **差异阈值**: 
   - < 1%: 通过
   - 1-5%: 警告
   - > 5%: 失败

---
*报告由自动视觉UI测试工具生成*
"""
        
        with open(report_path, 'w', encoding='utf-8') as f:
            f.write(report)
        
        print(f"\n报告已生成: {report_path}")
        return report_path
    
    def save_results_json(self):
        """保存JSON结果"""
        json_path = os.path.join(RESULTS_DIR, f"visual_test_results_{self.timestamp}.json")
        with open(json_path, 'w', encoding='utf-8') as f:
            json.dump(self.results, f, ensure_ascii=False, indent=2)
        return json_path


def main():
    tester = VisualUITester()
    
    # 运行所有测试
    results = tester.run_all_tests()
    
    # 生成报告
    report_path = tester.generate_report()
    json_path = tester.save_results_json()
    
    # 打印摘要
    print("\n" + "=" * 60)
    print("测试完成!")
    print("=" * 60)
    
    total = len(results)
    passed = len([r for r in results if r["status"] == "passed"])
    new_baselines = len([r for r in results if r["status"] == "new_baseline"])
    
    print(f"总测试: {total}")
    print(f"通过: {passed}")
    print(f"新基准: {new_baselines}")
    print(f"通过率: {(passed / total * 100) if total > 0 else 0:.1f}%")
    print(f"\n报告: {report_path}")
    print(f"JSON: {json_path}")


if __name__ == "__main__":
    main()
