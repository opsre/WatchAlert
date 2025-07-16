# 🤝 贡献指南（Contributing Guide）

欢迎为 **WatchAlert** 项目做贡献！我们非常感谢你对这个开源项目的热情和帮助。在开始之前，请花一点时间阅读以下指南，以便更好地理解如何提交代码、报告问题以及与社区互动。

## 📌 如何贡献

有多种方式可以为本项目做出贡献：

- 提交 bug 报告
- 提出新功能建议
- 编写文档
- 改进测试用例
- 提交代码修复 bug 或实现新功能
- 审查 Pull Requests

无论哪种方式，我们都欢迎并感激你的贡献！

---

## 🐛 报告 Bug

如果你发现了 bug，请通过 GitHub Issues 提交报告。请尽量提供以下信息：

- 描述清晰的问题现象
- 复现步骤（尽可能详细）
- 环境信息（操作系统、Go 版本、依赖版本等）
- 错误日志或截图（如有）

---

## 💡 提出功能请求（Feature Request）

如果你有一个改进 WatchAlert 的想法，请先查看是否有类似的功能请求已存在。如果没有，请创建一个新的 Issue，并说明：

- 功能用途
- 使用场景
- 建议的实现方式（可选）

---

## 🧪 开发流程

### 1. Fork 仓库

点击右上角的 "Fork" 按钮将仓库 fork 到你自己的账号下。

### 2. 克隆仓库

```bash
git clone https://github.com/your-username/WatchAlert.git 
cd WatchAlert
```

### 3. 创建分支

请基于`master`分支创建你的开发分支：

```
git checkout -b feature/your-feature-name
```

建议使用有意义的分支名，如：

- `feature/add-email-notification`
- `bugfix/fix-login-issue`

### 4. 编码规范
- 使用 Go Modules 管理依赖
- 遵循 Go 官方编码规范（gofmt）
- 添加必要的注释和文档更新
- 确保所有测试通过

### 5. 提交 Commit

使用清晰且描述性的 commit message，推荐使用[Conventional Commits](https://www.conventionalcommits.org/?spm=a2ty_o01.29997173.0.0.36bfc921cGwRPu)格式：

```
🚀 Feat: add email notification support
🚧 Fix: resolve panic when config is missing
📄 Docs: update README with installation steps
```

### 6. Push 到远程分支

```
git push origin feature/your-feature-name
```

### 7. 创建 Pull Request

前往 GitHub 页面创建`Pull Request`，并填写 PR 描述：

- 解决了什么问题或实现了什么功能？
- 是否包含配置变更或数据库迁移？
- 是否需要更新文档？

PR 将由维护者审核，可能要求你进行修改后合并。

##  📚 文档贡献

如果你希望改进文档，请确保：

- 修改内容清晰准确
- 保持与现有风格一致
- 更新相关的示例或截图（如适用）

## 🧑‍💻 成为维护者

长期积极参与项目、提交高质量 PR 并协助 issue 回复的贡献者，可能会被邀请成为项目维护者。

## ❤️ 感谢贡献者

我们衷心感谢每一位贡献者！你让 WatchAlert 变得更好 🎉
