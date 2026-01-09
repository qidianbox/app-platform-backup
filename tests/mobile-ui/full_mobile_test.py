#!/usr/bin/env python3
"""
完整移动端UI自动化测试
遍历所有功能页面，检查UI错位和溢出问题
"""

import os
import time
import json
import requests
from datetime import datetime
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException, NoSuchElementException

# 配置
BASE_URL = "http://localhost:5173"
API_URL = "http://localhost:8080"
SCREENSHOT_DIR = "/home/ubuntu/app-platform/tests/mobile-ui/full_screenshots"

# 移动设备配置
DEVICE = {
    "name": "iPhone_12",
    "width": 390,
    "height": 844,
    "user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1"
}

# 测试结果
test_results = []

def setup_driver():
    """设置Chrome驱动"""
    options = Options()
    options.add_argument('--headless')
    options.add_argument('--no-sandbox')
    options.add_argument('--disable-dev-shm-usage')
    options.add_argument('--disable-gpu')
    options.add_argument(f'--window-size={DEVICE["width"]},{DEVICE["height"]}')
    options.add_argument(f'--user-agent={DEVICE["user_agent"]}')
    
    # 模拟移动设备
    mobile_emulation = {
        "deviceMetrics": {"width": DEVICE["width"], "height": DEVICE["height"], "pixelRatio": 3.0},
        "userAgent": DEVICE["user_agent"]
    }
    options.add_experimental_option("mobileEmulation", mobile_emulation)
    
    driver = webdriver.Chrome(options=options)
    driver.set_window_size(DEVICE["width"], DEVICE["height"])
    return driver

def get_auth_token():
    """获取认证Token"""
    try:
        response = requests.post(f"{API_URL}/api/v1/admin/login", json={
            "username": "admin",
            "password": "admin123"
        })
        if response.status_code == 200:
            data = response.json()
            return data.get("data", {}).get("token")
    except Exception as e:
        print(f"获取Token失败: {e}")
    return None

def save_screenshot(driver, name, category="general"):
    """保存截图"""
    category_dir = os.path.join(SCREENSHOT_DIR, category)
    os.makedirs(category_dir, exist_ok=True)
    
    timestamp = datetime.now().strftime("%H%M%S")
    filename = f"{name}_{timestamp}.png"
    filepath = os.path.join(category_dir, filename)
    
    driver.save_screenshot(filepath)
    return filepath

def check_ui_issues(driver, page_name):
    """检查UI问题"""
    issues = []
    
    # 检查水平溢出
    try:
        body_width = driver.execute_script("return document.body.scrollWidth")
        viewport_width = driver.execute_script("return window.innerWidth")
        if body_width > viewport_width + 10:  # 允许10px误差
            issues.append(f"水平溢出: body宽度({body_width}px) > 视口宽度({viewport_width}px)")
    except:
        pass
    
    # 检查元素溢出
    try:
        overflow_elements = driver.execute_script("""
            var issues = [];
            var elements = document.querySelectorAll('*');
            var viewportWidth = window.innerWidth;
            for (var i = 0; i < elements.length; i++) {
                var rect = elements[i].getBoundingClientRect();
                if (rect.right > viewportWidth + 5) {
                    var tagName = elements[i].tagName.toLowerCase();
                    var className = elements[i].className;
                    if (className && typeof className === 'string') {
                        issues.push(tagName + '.' + className.split(' ')[0] + ' 溢出 ' + Math.round(rect.right - viewportWidth) + 'px');
                    }
                }
            }
            return issues.slice(0, 5);  // 只返回前5个
        """)
        if overflow_elements:
            issues.extend(overflow_elements)
    except:
        pass
    
    # 检查文字截断
    try:
        truncated = driver.execute_script("""
            var issues = [];
            var elements = document.querySelectorAll('*');
            for (var i = 0; i < elements.length; i++) {
                var style = window.getComputedStyle(elements[i]);
                if (style.overflow === 'hidden' && style.textOverflow === 'ellipsis') {
                    if (elements[i].scrollWidth > elements[i].clientWidth) {
                        var text = elements[i].innerText;
                        if (text && text.length > 0 && text.length < 50) {
                            issues.push('文字被截断: "' + text.substring(0, 20) + '..."');
                        }
                    }
                }
            }
            return issues.slice(0, 3);
        """)
        if truncated:
            issues.extend(truncated)
    except:
        pass
    
    return issues

def test_page(driver, url, page_name, category, wait_selector=None, actions=None):
    """测试单个页面"""
    print(f"  测试: {page_name}...")
    
    result = {
        "page": page_name,
        "category": category,
        "url": url,
        "status": "PASS",
        "issues": [],
        "screenshot": None
    }
    
    try:
        driver.get(url)
        time.sleep(2)  # 等待页面加载
        
        # 等待特定元素
        if wait_selector:
            try:
                WebDriverWait(driver, 5).until(
                    EC.presence_of_element_located((By.CSS_SELECTOR, wait_selector))
                )
            except TimeoutException:
                pass
        
        # 执行额外操作
        if actions:
            for action in actions:
                try:
                    if action["type"] == "click":
                        element = driver.find_element(By.CSS_SELECTOR, action["selector"])
                        element.click()
                        time.sleep(1)
                    elif action["type"] == "scroll":
                        driver.execute_script(f"window.scrollTo(0, {action['y']})")
                        time.sleep(0.5)
                except:
                    pass
        
        # 检查UI问题
        issues = check_ui_issues(driver, page_name)
        if issues:
            result["status"] = "WARN"
            result["issues"] = issues
        
        # 保存截图
        result["screenshot"] = save_screenshot(driver, page_name.replace(" ", "_").lower(), category)
        
    except Exception as e:
        result["status"] = "FAIL"
        result["issues"] = [str(e)]
    
    test_results.append(result)
    
    status_icon = "✅" if result["status"] == "PASS" else ("⚠️" if result["status"] == "WARN" else "❌")
    print(f"    {status_icon} {result['status']}: {', '.join(result['issues']) if result['issues'] else '无问题'}")
    
    return result

def inject_token(driver, token):
    """注入认证Token"""
    driver.execute_script(f"localStorage.setItem('token', '{token}')")

def main():
    print("=" * 60)
    print("完整移动端UI自动化测试")
    print(f"设备: {DEVICE['name']} ({DEVICE['width']}x{DEVICE['height']})")
    print("=" * 60)
    
    # 清理截图目录
    os.makedirs(SCREENSHOT_DIR, exist_ok=True)
    
    # 获取Token
    print("\n获取认证Token...")
    token = get_auth_token()
    if not token:
        print("❌ 无法获取Token，测试终止")
        return
    print("✅ Token获取成功")
    
    # 设置驱动
    driver = setup_driver()
    
    try:
        # ========== 1. 登录页面 ==========
        print("\n【1. 登录页面】")
        test_page(driver, f"{BASE_URL}/#/login", "登录页面", "01_login")
        
        # 注入Token并刷新
        driver.get(BASE_URL)
        inject_token(driver, token)
        driver.refresh()
        time.sleep(2)
        
        # ========== 2. 仪表盘 ==========
        print("\n【2. 仪表盘】")
        test_page(driver, f"{BASE_URL}/#/dashboard", "仪表盘", "02_dashboard", ".stat-card")
        
        # 滚动测试
        driver.execute_script("window.scrollTo(0, 500)")
        time.sleep(1)
        save_screenshot(driver, "dashboard_scrolled", "02_dashboard")
        
        # ========== 3. APP列表 ==========
        print("\n【3. APP列表】")
        test_page(driver, f"{BASE_URL}/#/apps", "APP列表", "03_app_list")
        
        # ========== 4. 创建APP弹窗 ==========
        print("\n【4. 创建APP弹窗】")
        driver.get(f"{BASE_URL}/#/apps")
        time.sleep(2)
        try:
            # 尝试点击创建按钮
            create_btn = driver.find_element(By.XPATH, "//button[contains(., '创建')]")
            create_btn.click()
            time.sleep(1)
            save_screenshot(driver, "create_app_dialog", "04_create_app")
            
            # 检查弹窗UI
            issues = check_ui_issues(driver, "创建APP弹窗")
            test_results.append({
                "page": "创建APP弹窗",
                "category": "04_create_app",
                "status": "WARN" if issues else "PASS",
                "issues": issues,
                "screenshot": f"{SCREENSHOT_DIR}/04_create_app/create_app_dialog.png"
            })
            print(f"    {'⚠️ WARN' if issues else '✅ PASS'}: {', '.join(issues) if issues else '无问题'}")
            
            # 关闭弹窗
            driver.execute_script("document.querySelector('.el-dialog__close')?.click()")
            time.sleep(0.5)
        except Exception as e:
            print(f"    ⚠️ 无法测试创建弹窗: {e}")
        
        # ========== 5. APP详情/配置 ==========
        print("\n【5. APP详情页面】")
        
        # 获取第一个APP的ID
        try:
            response = requests.get(f"{API_URL}/api/v1/apps", headers={"Authorization": f"Bearer {token}"})
            apps = response.json().get("data", {}).get("list", [])
            if apps:
                app_id = apps[0]["id"]
                
                # 测试APP配置主页
                test_page(driver, f"{BASE_URL}/#/apps/{app_id}/config", "APP配置-概览", "05_app_config")
                
                # 滚动查看更多内容
                driver.execute_script("window.scrollTo(0, 300)")
                time.sleep(0.5)
                save_screenshot(driver, "app_config_scrolled", "05_app_config")
                
                # ========== 6. 用户中心配置 ==========
                print("\n【6. 用户中心配置】")
                test_page(driver, f"{BASE_URL}/#/apps/{app_id}/config/user", "用户中心配置", "06_user_config")
                
                # 滚动查看
                driver.execute_script("window.scrollTo(0, 500)")
                time.sleep(0.5)
                save_screenshot(driver, "user_config_scrolled", "06_user_config")
                
                # ========== 7. 消息推送配置 ==========
                print("\n【7. 消息推送配置】")
                test_page(driver, f"{BASE_URL}/#/apps/{app_id}/config/message", "消息推送配置", "07_message_config")
                
                # ========== 8. 支付配置 ==========
                print("\n【8. 支付配置】")
                test_page(driver, f"{BASE_URL}/#/apps/{app_id}/config/payment", "支付配置", "08_payment_config")
                
                # ========== 9. 数据统计配置 ==========
                print("\n【9. 数据统计配置】")
                test_page(driver, f"{BASE_URL}/#/apps/{app_id}/config/analytics", "数据统计配置", "09_analytics_config")
                
                # ========== 10. 安全配置 ==========
                print("\n【10. 安全配置】")
                test_page(driver, f"{BASE_URL}/#/apps/{app_id}/config/security", "安全配置", "10_security_config")
                
                # ========== 11. 版本管理 ==========
                print("\n【11. 版本管理】")
                test_page(driver, f"{BASE_URL}/#/apps/{app_id}/config/version", "版本管理", "11_version_config")
                
                # ========== 12. 模块管理 ==========
                print("\n【12. 模块管理】")
                test_page(driver, f"{BASE_URL}/#/apps/{app_id}/modules", "模块管理", "12_modules")
                
        except Exception as e:
            print(f"    ❌ 获取APP失败: {e}")
        
    finally:
        driver.quit()
    
    # 生成报告
    generate_report()

def generate_report():
    """生成测试报告"""
    print("\n" + "=" * 60)
    print("生成测试报告...")
    
    report_path = os.path.join(SCREENSHOT_DIR, "FULL_TEST_REPORT.md")
    
    # 统计
    total = len(test_results)
    passed = sum(1 for r in test_results if r["status"] == "PASS")
    warned = sum(1 for r in test_results if r["status"] == "WARN")
    failed = sum(1 for r in test_results if r["status"] == "FAIL")
    
    with open(report_path, "w", encoding="utf-8") as f:
        f.write("# 完整移动端UI测试报告\n\n")
        f.write(f"**测试时间**: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n\n")
        f.write(f"**测试设备**: {DEVICE['name']} ({DEVICE['width']}x{DEVICE['height']})\n\n")
        f.write(f"**测试URL**: {BASE_URL}\n\n")
        
        f.write("## 测试概览\n\n")
        f.write("| 统计项 | 数量 |\n")
        f.write("|:---|:---|\n")
        f.write(f"| ✅ 通过 | {passed} |\n")
        f.write(f"| ⚠️ 警告 | {warned} |\n")
        f.write(f"| ❌ 失败 | {failed} |\n")
        f.write(f"| **总计** | {total} |\n\n")
        
        f.write("## 详细结果\n\n")
        f.write("| 页面 | 状态 | 问题 |\n")
        f.write("|:---|:---|:---|\n")
        
        for result in test_results:
            status_icon = "✅" if result["status"] == "PASS" else ("⚠️" if result["status"] == "WARN" else "❌")
            issues_str = "<br>".join(result["issues"]) if result["issues"] else "-"
            f.write(f"| {result['page']} | {status_icon} {result['status']} | {issues_str} |\n")
        
        f.write("\n## 截图\n\n")
        
        current_category = None
        for result in test_results:
            if result.get("screenshot"):
                if result["category"] != current_category:
                    current_category = result["category"]
                    f.write(f"\n### {result['page']}\n\n")
                
                # 使用相对路径
                rel_path = os.path.relpath(result["screenshot"], SCREENSHOT_DIR)
                f.write(f"![{result['page']}]({rel_path})\n\n")
    
    print(f"报告已保存: {report_path}")
    print("=" * 60)
    print(f"测试结果统计:")
    print(f"  ✅ 通过: {passed}")
    print(f"  ⚠️ 警告: {warned}")
    print(f"  ❌ 失败: {failed}")
    print(f"  总计: {total}")

if __name__ == "__main__":
    main()
