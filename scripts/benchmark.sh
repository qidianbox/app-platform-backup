#!/bin/bash

# APP中台系统压力测试脚本
# 使用hey工具进行HTTP压力测试

BASE_URL="http://localhost:8080"
RESULTS_DIR="/home/ubuntu/app-platform/benchmark_results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# 创建结果目录
mkdir -p $RESULTS_DIR

# 获取Token
echo "=== 获取认证Token ==="
TOKEN=$(curl -s -X POST "$BASE_URL/api/v1/admin/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "获取Token失败，退出"
  exit 1
fi
echo "Token获取成功"

# 压测参数
REQUESTS=500        # 总请求数
CONCURRENCY=50      # 并发数
DURATION="30s"      # 持续时间

echo ""
echo "=== 压测配置 ==="
echo "总请求数: $REQUESTS"
echo "并发数: $CONCURRENCY"
echo "持续时间: $DURATION"
echo ""

# 结果文件
RESULT_FILE="$RESULTS_DIR/benchmark_$TIMESTAMP.txt"

echo "APP中台系统压力测试报告" > $RESULT_FILE
echo "测试时间: $(date)" >> $RESULT_FILE
echo "========================================" >> $RESULT_FILE
echo "" >> $RESULT_FILE

# 测试函数
run_test() {
  local name=$1
  local method=$2
  local url=$3
  local data=$4
  
  echo ">>> 测试: $name"
  echo "" >> $RESULT_FILE
  echo "### $name ###" >> $RESULT_FILE
  echo "URL: $url" >> $RESULT_FILE
  echo "Method: $method" >> $RESULT_FILE
  echo "" >> $RESULT_FILE
  
  if [ "$method" == "GET" ]; then
    hey -n $REQUESTS -c $CONCURRENCY \
      -H "Authorization: Bearer $TOKEN" \
      "$url" >> $RESULT_FILE 2>&1
  else
    hey -n $REQUESTS -c $CONCURRENCY \
      -m $method \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "$data" \
      "$url" >> $RESULT_FILE 2>&1
  fi
  
  echo "" >> $RESULT_FILE
  echo "----------------------------------------" >> $RESULT_FILE
}

# 1. 登录接口压测（无需Token）
echo ""
echo "=== 1. 登录接口压测 ==="
echo "### 登录接口 ###" >> $RESULT_FILE
hey -n $REQUESTS -c $CONCURRENCY \
  -m POST \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  "$BASE_URL/api/v1/admin/login" >> $RESULT_FILE 2>&1
echo "" >> $RESULT_FILE
echo "----------------------------------------" >> $RESULT_FILE

# 2. APP列表接口
echo ""
echo "=== 2. APP列表接口压测 ==="
run_test "APP列表" "GET" "$BASE_URL/api/v1/apps"

# 3. 用户列表接口
echo ""
echo "=== 3. 用户列表接口压测 ==="
run_test "用户列表" "GET" "$BASE_URL/api/v1/users?page=1&size=20"

# 4. 消息列表接口
echo ""
echo "=== 4. 消息列表接口压测 ==="
run_test "消息列表" "GET" "$BASE_URL/api/v1/messages?page=1&size=20"

# 5. 日志查询接口
echo ""
echo "=== 5. 日志查询接口压测 ==="
run_test "日志查询" "GET" "$BASE_URL/api/v1/logs?page=1&size=20"

# 6. 监控指标接口
echo ""
echo "=== 6. 监控指标接口压测 ==="
run_test "监控指标" "GET" "$BASE_URL/api/v1/monitor/metrics?app_id=1"

# 7. 监控统计接口
echo ""
echo "=== 7. 监控统计接口压测 ==="
run_test "监控统计" "GET" "$BASE_URL/api/v1/monitor/stats?app_id=1"

# 8. 事件列表接口
echo ""
echo "=== 8. 事件列表接口压测 ==="
run_test "事件列表" "GET" "$BASE_URL/api/v1/events?app_id=1&page=1&size=20"

# 9. 版本列表接口
echo ""
echo "=== 9. 版本列表接口压测 ==="
run_test "版本列表" "GET" "$BASE_URL/api/v1/versions?app_id=1"

# 10. 推送列表接口
echo ""
echo "=== 10. 推送列表接口压测 ==="
run_test "推送列表" "GET" "$BASE_URL/api/v1/push?app_id=1&page=1&size=20"

# 11. 文件列表接口
echo ""
echo "=== 11. 文件列表接口压测 ==="
run_test "文件列表" "GET" "$BASE_URL/api/v1/files?app_id=1&page=1&size=20"

# 12. 告警列表接口
echo ""
echo "=== 12. 告警列表接口压测 ==="
run_test "告警列表" "GET" "$BASE_URL/api/v1/monitor/alerts?app_id=1"

# 13. 审计日志接口
echo ""
echo "=== 13. 审计日志接口压测 ==="
run_test "审计日志" "GET" "$BASE_URL/api/v1/audit?page=1&size=20"

echo ""
echo "=== 压测完成 ==="
echo "结果已保存到: $RESULT_FILE"

# 生成摘要
echo ""
echo "=== 生成测试摘要 ==="
SUMMARY_FILE="$RESULTS_DIR/summary_$TIMESTAMP.txt"

echo "APP中台系统压测摘要" > $SUMMARY_FILE
echo "测试时间: $(date)" >> $SUMMARY_FILE
echo "========================================" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE

# 提取关键指标
grep -A 20 "###" $RESULT_FILE | grep -E "(Requests/sec|Average|Fastest|Slowest|Status code)" >> $SUMMARY_FILE

echo ""
echo "摘要已保存到: $SUMMARY_FILE"
echo ""
cat $SUMMARY_FILE
