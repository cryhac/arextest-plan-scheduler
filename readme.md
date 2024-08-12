# AREXTest Plan Scheduler 用于固化用例的回放测试。

AREXTest Plan Scheduler 是一个用 Go 编写的工具，用于定时的使用固化用例的回放测试。

## 功能

1. 连接到 MongoDB 数据库并查询固定的 API 调用。
2. 基于操作名称将这些 API 调用分组。
3. 查询与这些操作名称相对应的服务操作。
4. 创建并发送计划请求到指定的调度终端。

## 环境变量

该工具依赖以下环境变量：

- `APP_ID`：应用程序的标识符。
- `TARGET_HOST`：目标回放地址。
- `MONGO_URL`：指定的 MongoDB URL。如果未提供，将使用默认值 `mongodb://arex:iLoveArex@arex-helm-name-beta-arex-mongodb.arex.svc.cluster.local:27017/arex_storage_db`。
- `SCHEDULE_ENDPOINT`：指定的计划创建 API 终端点。如果未提供，将使用默认值 `http://arex-helm-name-beta-arex-schedule.arex.svc.cluster.local:8080/api/createPlan`。

## 快速开始

请按照以下步骤部署和运行该工具：

### 1. 克隆项目

```sh
git clone https://github.com/arextest/arex-plan-scheduler.git
cd arex-plan-scheduler
```

### 2. 制作镜像

```sh
docker build -t  registry.cn-hangzhou.aliyuncs.com/arexadmin01/arextest-plan-scheduler:0.6.5 .
docker push registry.cn-hangzhou.aliyuncs.com/arexadmin01/arextest-plan-scheduler:0.6.5
```
请按照自己的情况修改仓库镜像地址

### 3. 提交任务
```sh
➜  arextest-plan-scheduler kubectl apply -f your-app-name-job.yml  
cronjob.batch/your-app-name-job created
```
### 4. 查看任务和日志
查看任务:
![image-20240812172526864](https://test-1251091139.cos.ap-shanghai.myqcloud.com/picgoimage-20240812172526864.png)

查看日志:
![image-20240812173421601](https://test-1251091139.cos.ap-shanghai.myqcloud.com/picgoimage-20240812173421601.png)

### 5. 任务报告
![image-20240812180608047](https://test-1251091139.cos.ap-shanghai.myqcloud.com/picgoimage-20240812180608047.png)