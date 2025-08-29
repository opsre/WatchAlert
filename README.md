
<p align="center">
  <a href="https://github.com/w8t-io/WatchAlert"> 
    <img src="WatchAlert.png" alt="cloud native monitoring" width="200" height="auto" /></a>
</p>

<p align="center">
  <b>WatchAlert 开源一站式多数据源监控告警引擎</b>
</p>

<p align="center">
<a href="https://github.com/w8t-io/WatchAlert/graphs/contributors">
  <img alt="GitHub contributors" src="https://img.shields.io/github/contributors-anon/w8t-io/WatchAlert"/></a>
<img alt="GitHub Repo stars" src="https://img.shields.io/github/stars/w8t-io/WatchAlert">
<img alt="GitHub forks" src="https://img.shields.io/github/forks/w8t-io/WatchAlert">
<img alt="GitHub Repo issues" src="https://img.shields.io/github/issues/w8t-io/WatchAlert">
<img alt="GitHub Repo issues closed" src="https://img.shields.io/github/issues-closed/w8t-io/WatchAlert">
<img alt="License" src="https://img.shields.io/badge/license-Apache--2.0-blue"/>

## 💎 项目介绍
WatchAlert 是一款为云原生环境量身打造的轻量级监控告警引擎，专注于**可观测稳定性**主题，提供全面的监控与告警支持。

AI + WatchAlert 实现智能化告警分析 高效处理故障告警;

**架构图**
![img.png](assets/architecture.png)

**能力**
- AI 智能分析
  - 针对`Metrics` `Logs` `Traces` 告警内容做内容分析，高效定位告警根因，并提供排查思路和解决方案;
- Metrics 监控
    - 集成：Prometheus、VictoriaMetrics
- Logs 监控
    - 集成：Loki、ElasticSearch、阿里云日志服务 (AliCloud SLS)
- Traces 监控
    - 集成：Jaeger
- Events 监控
    - 集成：Kubernetes
- Network 监控
    - 集成：HTTP、ICMP、TCP、SSL
- 告警通知
    - 飞书、钉钉、企业微信、邮件、自定义Hook


**为什么选择 WatchAlert？**
- 针对云原生环境优化，轻量易用。
- 灵活的告警规则配置，支持多种数据源。
- 专注于稳定性场景，助力快速发现与解决问题。
- 提供稳定可靠的监控告警解决方案，助力用户提升运维效率，降低维护成本。

## 📄 项目文档

- 功能介绍：https://cairry.github.io/docs/
- 快速部署：https://cairry.github.io/docs/install/docker.html


## 🎉 项目预览
- 演示环境：http://8.147.234.89/login
  （admin/123）

|    ![登陆页](assets/login.png)    | ![首页](assets/home.png)                 |
|:------------------------------:|--------------------------------|
|   ![rules.png](assets/rules.png)   | ![faultcenter.png](assets/faultcenter.png) |
| ![img_2.png](assets/img_2.png) | ![img_3.png](assets/img_3.png) |
| ![img_5.png](assets/img_5.png) | ![img_4.png](assets/img_4.png) |
| ![img_6.png](assets/img_6.png) | ![img_7.png](assets/img_7.png) |

## 🎸 支持
- 如果你觉得 WatchAlert 还不错，可以通过 Star 来表示你的喜欢
- 在公司或个人项目中使用 WatchAlert，并帮忙推广给伙伴使用

## 🧑‍💻 交流渠道
- [点击我](https://cairry.github.io/docs/#%E4%BA%A4%E6%B5%81%E7%BE%A4-%E8%81%94%E7%B3%BB%E6%88%91)