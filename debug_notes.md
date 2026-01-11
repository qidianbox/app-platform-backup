# 问题分析

## 问题描述
用户在移动端访问日志查询页面时，显示 "app_id is required" 和 "app_id 不能为空" 的错误。

## 分析
从截图可以看到：
1. 用户在移动端访问
2. 页面显示的是工作台的侧边栏菜单（数据概览、用户管理、消息推送等）
3. 当前选中的是"日志查询"
4. 错误提示 "app_id is required" 和 "app_id 不能为空"

## 可能的原因
1. 移动端菜单切换时，appId 没有正确传递
2. Workspace 组件的 props.appId 为空或 undefined
3. 路由参数没有正确获取

## 需要检查的代码
1. config/index.vue 中 appId 的计算属性
2. Workspace 组件中 props.appId 的使用
3. 移动端菜单切换逻辑
