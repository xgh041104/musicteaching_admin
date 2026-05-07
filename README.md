# 学科课程平台

### 介绍
主要用于跨学课课程教学使用

### 软件架构
- 前端基于umi@3框架构建

### 前端开发配置

- 安装nodejs v16.16.0
- 安装cnpm:`npm install -g cnpm`
- 下载代码`git clone git@gitee.com:YiBuJiaoYuKeJi/subject-course-platform.git`
- `cd subject-course-platform/web_react`
- 使用vscode打开时注意打开`web_react`前端项目目录!!!否则在vscode下运行下面命令会有问题！
- `cnpm install`
- `npm run start`
#### 前端项目说明
1. 在`src/pages`目录下新建页面
    1. 页面文件夹中带`_`符号可建立组件不被渲染成路由
    2. 页面文件夹中新建model.js文件可使用dva管理页面数据流
2. 在`src/components`目录下可建立通用组件;
3. 在`src/layouts`目录下可建立布局组件;
4. 在`src/models`目录下可使用dva管理全局数据流;

### 后端开发配置
[后端配置使用说明](/server_go/README.md)





