#!/usr/bin/env python3
"""
移动端UI自动化测试脚本 (Selenium版本)
使用Selenium模拟多种移动设备进行测试
"""

import os
import time
import requests
from datetime import datetime
from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

# 测试配置
BASE_URL = "http://localhost:5173"
API_URL = "http://localhost:8080"
SCREENSHOTS_DIR = "/home/ubuntu/app-platform/tests/mobile-ui/screenshots"
REPORT_FILE = "/home/ubuntu/app-platform/tests/mobile-ui/test_report.md"

# 测试账号
TEST_USERNAME = "admin"
TEST_PASSWORD = "admin123"

# 要测试的移动设备配置
MOBILE_DEVICES = {
    "iPhone_12": {
        "width": 390,
        "height": 844,
        "pixel_ratio": 3,
        "user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1"
    },
    "iPhone_SE": {
        "width": 375,
        "height": 667,
        "pixel_ratio": 2,
        "user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1"
    },
    "Pixel_5": {
        "width": 393,
        "height": 851,
        "pixel_ratio": 2.75,
        "user_agent": "Mozilla/5.0 (Linux; Android 11; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36"
    },
    "Galaxy_S9": {
        "width": 360,
        "height": 740,
        "pixel_ratio": 4,
        "user_agent": "Mozilla/5.0 (Linux; Android 8.0.0; SM-G960F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36"
    },
    "iPad_Mini": {
        "width": 768,
        "height": 1024,
        "pixel_ratio": 2,
        "user_agent": "Mozilla/5.0 (iPad; CPU OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1"
    }
}

# 测试结果
test_results = []

def setup_directories():
    """创建截图目录"""
    os.makedirs(SCREENSHOTS_DIR, exist_ok=True)
    for device in MOBILE_DEVICES.keys():
        device_dir = os.path.join(SCREENSHOTS_DIR, device)
        os.makedirs(device_dir, exist_ok=True)

def create_mobile_driver(device_name, device_config):
    """创建模拟移动设备的Chrome驱动"""
    options = Options()
    options.add_argument("--headless")
    options.add_argument("--no-sandbox")
    options.add_argument("--disable-dev-shm-usage")
    options.add_argument("--disable-gpu")
    options.add_argument(f"--window-size={device_config['width']},{device_config['height']}")
    options.add_argument(f"--user-agent={device_config['user_agent']}")
    
    # 启用移动设备模拟
    mobile_emulation = {
        "deviceMetrics": {
            "width": device_config["width"],
            "height": device_config["height"],
            "pixelRatio": device_config["pixel_ratio"]
        },
        "userAgent": device_config["user_agent"]
    }
    options.add_experimental_option("mobileEmulation", mobile_emulation)
    
    driver = webdriver.Chrome(options=options)
    return driver

def take_screenshot(driver, device_name, page_name):
    """截图并保存"""
    device_dir = os.path.join(SCREENSHOTS_DIR, device_name)
    filename = f"{page_name}_{datetime.now().strftime('%H%M%S')}.png"
    filepath = os.path.join(device_dir, filename)
    driver.save_screenshot(filepath)
    return filepath

def get_auth_token():
    """获取认证Token"""
    try:
        response = requests.post(
            f"{API_URL}/api/v1/admin/login",
            json={"username": TEST_USERNAME, "password": TEST_PASSWORD}
        )
        data = response.json()
        if data.get("code") == 0:
            return data["data"]["token"]
    except Exception as e:
        print(f"获取Token失败: {e}")
    return None

def check_horizontal_overflow(driver):
    """检查页面是否有水平溢出"""
    viewport_width = driver.execute_script("return window.innerWidth")
    body_width = driver.execute_script("return document.body.scrollWidth")
    if body_width > viewport_width + 10:
        return f"页面有水平溢出: {body_width}px > {viewport_width}px"
    return None

def test_login_page(driver, device_name):
    """测试登录页面"""
    result = {"device": device_name, "page": "登录页面", "status": "PASS", "issues": []}
    
    try:
        driver.get(BASE_URL)
        time.sleep(2)
        
        # 截图
        screenshot = take_screenshot(driver, device_name, "01_login")
        result["screenshot"] = screenshot
        
        # 检查登录表单元素
        try:
            driver.find_element(By.CSS_SELECTOR, "input[placeholder*='用户名'], input[placeholder*='账号']")
        except:
            result["issues"].append("用户名输入框不可见")
            
        try:
            driver.find_element(By.CSS_SELECTOR, "input[placeholder*='密码']")
        except:
            result["issues"].append("密码输入框不可见")
        
        # 检查水平溢出
        overflow = check_horizontal_overflow(driver)
        if overflow:
            result["issues"].append(overflow)
            
    except Exception as e:
        result["status"] = "FAIL"
        result["issues"].append(str(e))
    
    if result["issues"]:
        result["status"] = "WARN" if result["status"] == "PASS" else result["status"]
    
    return result

def test_dashboard_page(driver, device_name, token):
    """测试仪表盘页面"""
    result = {"device": device_name, "page": "仪表盘", "status": "PASS", "issues": []}
    
    try:
        # 设置token
        driver.get(BASE_URL)
        driver.execute_script(f"localStorage.setItem('token', '{token}')")
        driver.get(f"{BASE_URL}/#/dashboard")
        time.sleep(3)
        
        # 截图
        screenshot = take_screenshot(driver, device_name, "02_dashboard")
        result["screenshot"] = screenshot
        
        # 检查统计卡片
        stat_cards = driver.find_elements(By.CSS_SELECTOR, ".stat-card")
        if len(stat_cards) < 4:
            result["issues"].append(f"统计卡片数量不足: {len(stat_cards)}/4")
        
        # 检查水平溢出
        overflow = check_horizontal_overflow(driver)
        if overflow:
            result["issues"].append(overflow)
            
    except Exception as e:
        result["status"] = "FAIL"
        result["issues"].append(str(e))
    
    if result["issues"]:
        result["status"] = "WARN" if result["status"] == "PASS" else result["status"]
    
    return result

def test_app_list_page(driver, device_name, token):
    """测试APP列表页面"""
    result = {"device": device_name, "page": "APP列表", "status": "PASS", "issues": []}
    
    try:
        driver.execute_script(f"localStorage.setItem('token', '{token}')")
        driver.get(f"{BASE_URL}/#/app/list")
        time.sleep(3)
        
        # 截图
        screenshot = take_screenshot(driver, device_name, "03_app_list")
        result["screenshot"] = screenshot
        
        # 检查APP卡片
        app_cards = driver.find_elements(By.CSS_SELECTOR, ".app-card")
        if len(app_cards) < 1:
            result["issues"].append("没有找到APP卡片")
        
        # 检查水平溢出
        overflow = check_horizontal_overflow(driver)
        if overflow:
            result["issues"].append(overflow)
            
    except Exception as e:
        result["status"] = "FAIL"
        result["issues"].append(str(e))
    
    if result["issues"]:
        result["status"] = "WARN" if result["status"] == "PASS" else result["status"]
    
    return result

def test_create_app_dialog(driver, device_name, token):
    """测试创建APP弹窗"""
    result = {"device": device_name, "page": "创建APP弹窗", "status": "PASS", "issues": []}
    
    try:
        driver.execute_script(f"localStorage.setItem('token', '{token}')")
        driver.get(f"{BASE_URL}/#/app/list")
        time.sleep(3)
        
        # 点击创建按钮
        create_cards = driver.find_elements(By.CSS_SELECTOR, ".create-card")
        if create_cards:
            create_cards[0].click()
            time.sleep(2)
            
            # 截图
            screenshot = take_screenshot(driver, device_name, "04_create_dialog")
            result["screenshot"] = screenshot
            
            # 检查弹窗
            dialogs = driver.find_elements(By.CSS_SELECTOR, ".el-dialog")
            if not dialogs:
                result["issues"].append("创建弹窗未显示")
            else:
                # 检查弹窗宽度
                dialog_width = driver.execute_script("return document.querySelector('.el-dialog').offsetWidth")
                viewport_width = driver.execute_script("return window.innerWidth")
                if dialog_width > viewport_width:
                    result["issues"].append(f"弹窗宽度超出屏幕: {dialog_width}px > {viewport_width}px")
        else:
            result["issues"].append("创建卡片不可见")
            
    except Exception as e:
        result["status"] = "FAIL"
        result["issues"].append(str(e))
    
    if result["issues"]:
        result["status"] = "WARN" if result["status"] == "PASS" else result["status"]
    
    return result

def test_app_config_page(driver, device_name, token):
    """测试APP配置页面"""
    result = {"device": device_name, "page": "APP配置", "status": "PASS", "issues": []}
    
    try:
        driver.execute_script(f"localStorage.setItem('token', '{token}')")
        driver.get(f"{BASE_URL}/#/app/6/config")
        time.sleep(3)
        
        # 截图
        screenshot = take_screenshot(driver, device_name, "05_app_config")
        result["screenshot"] = screenshot
        
        # 检查水平溢出
        overflow = check_horizontal_overflow(driver)
        if overflow:
            result["issues"].append(overflow)
            
    except Exception as e:
        result["status"] = "FAIL"
        result["issues"].append(str(e))
    
    if result["issues"]:
        result["status"] = "WARN" if result["status"] == "PASS" else result["status"]
    
    return result

def generate_report(results):
    """生成测试报告"""
    report = f"""# 移动端UI自动化测试报告

**测试时间**: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}

**测试URL**: {BASE_URL}

## 测试概览

| 设备 | 页面 | 状态 | 问题 |
|:-----|:-----|:-----|:-----|
"""
    
    pass_count = 0
    warn_count = 0
    fail_count = 0
    
    for result in results:
        status_emoji = "✅" if result["status"] == "PASS" else ("⚠️" if result["status"] == "WARN" else "❌")
        issues = ", ".join(result["issues"]) if result["issues"] else "-"
        report += f"| {result['device']} | {result['page']} | {status_emoji} {result['status']} | {issues} |\n"
        
        if result["status"] == "PASS":
            pass_count += 1
        elif result["status"] == "WARN":
            warn_count += 1
        else:
            fail_count += 1
    
    report += f"""
## 统计

- ✅ 通过: {pass_count}
- ⚠️ 警告: {warn_count}
- ❌ 失败: {fail_count}
- **总计**: {len(results)}

## 截图

"""
    
    # 按设备分组添加截图
    devices_screenshots = {}
    for result in results:
        device = result["device"]
        if device not in devices_screenshots:
            devices_screenshots[device] = []
        if "screenshot" in result:
            devices_screenshots[device].append({
                "page": result["page"],
                "path": result["screenshot"]
            })
    
    for device, screenshots in devices_screenshots.items():
        report += f"### {device}\n\n"
        for ss in screenshots:
            rel_path = os.path.basename(ss["path"])
            report += f"**{ss['page']}**\n\n"
            report += f"![{ss['page']}](screenshots/{device}/{rel_path})\n\n"
    
    # 保存报告
    with open(REPORT_FILE, "w", encoding="utf-8") as f:
        f.write(report)
    
    return report

def run_tests():
    """运行所有测试"""
    print("=" * 50)
    print("移动端UI自动化测试 (Selenium)")
    print("=" * 50)
    
    # 设置目录
    setup_directories()
    
    # 获取Token
    print("\n获取认证Token...")
    token = get_auth_token()
    if not token:
        print("❌ 无法获取Token，测试终止")
        return
    print("✅ Token获取成功")
    
    results = []
    
    for device_name, device_config in MOBILE_DEVICES.items():
        print(f"\n测试设备: {device_name} ({device_config['width']}x{device_config['height']})")
        print("-" * 40)
        
        driver = None
        try:
            driver = create_mobile_driver(device_name, device_config)
            
            # 测试登录页面
            print(f"  测试登录页面...")
            result = test_login_page(driver, device_name)
            results.append(result)
            print(f"    {result['status']}: {', '.join(result['issues']) if result['issues'] else '无问题'}")
            
            # 测试仪表盘
            print(f"  测试仪表盘...")
            result = test_dashboard_page(driver, device_name, token)
            results.append(result)
            print(f"    {result['status']}: {', '.join(result['issues']) if result['issues'] else '无问题'}")
            
            # 测试APP列表
            print(f"  测试APP列表...")
            result = test_app_list_page(driver, device_name, token)
            results.append(result)
            print(f"    {result['status']}: {', '.join(result['issues']) if result['issues'] else '无问题'}")
            
            # 测试创建APP弹窗
            print(f"  测试创建APP弹窗...")
            result = test_create_app_dialog(driver, device_name, token)
            results.append(result)
            print(f"    {result['status']}: {', '.join(result['issues']) if result['issues'] else '无问题'}")
            
            # 测试APP配置页面
            print(f"  测试APP配置页面...")
            result = test_app_config_page(driver, device_name, token)
            results.append(result)
            print(f"    {result['status']}: {', '.join(result['issues']) if result['issues'] else '无问题'}")
            
        except Exception as e:
            print(f"  ❌ 测试异常: {e}")
        finally:
            if driver:
                driver.quit()
    
    # 生成报告
    print("\n" + "=" * 50)
    print("生成测试报告...")
    report = generate_report(results)
    print(f"报告已保存: {REPORT_FILE}")
    print("=" * 50)
    
    # 打印统计
    pass_count = sum(1 for r in results if r["status"] == "PASS")
    warn_count = sum(1 for r in results if r["status"] == "WARN")
    fail_count = sum(1 for r in results if r["status"] == "FAIL")
    
    print(f"\n测试结果统计:")
    print(f"  ✅ 通过: {pass_count}")
    print(f"  ⚠️ 警告: {warn_count}")
    print(f"  ❌ 失败: {fail_count}")
    print(f"  总计: {len(results)}")

if __name__ == "__main__":
    run_tests()
