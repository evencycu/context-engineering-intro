# 根目錄腳本

這個目錄包含從根目錄搬遷過來的腳本檔案。

## 檔案說明

- `dev-deploy.sh` - 完整開發部署腳本
- `quick-redeploy.sh` - 快速重構部署腳本
- `smart-deploy.sh` - 智能部署腳本

## 使用方法

### 開發部署
```bash
./scripts/root/dev-deploy.sh
```

### 快速重構
```bash
./scripts/root/quick-redeploy.sh
```

### 智能部署
```bash
# 本地進程模式
./scripts/root/smart-deploy.sh local-process

# 本地 Docker 模式
./scripts/root/smart-deploy.sh local-docker

# 持久化 Docker 模式
./scripts/root/smart-deploy.sh persistent-docker
```

## 注意事項

這些腳本原本位於根目錄，現在已搬遷到此目錄以保持專案結構的整潔。
