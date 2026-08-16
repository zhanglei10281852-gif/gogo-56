# Bug Reproduction

## 包的性质

当前 test_model_fix 保存的是被测模型修复后的结果源码，不是初始含 Bug 源码。要复现原始缺陷，必须检出下面固定的 parent SHA；不要在当前修复结果源码上期待重新出现修复前失败。生成系统使用的可信验证补丁和完整验证日志仅在本地留存，不提交到结果分支。

## 问题现象

我们在一个西向磁偏角的洞里归算数据，靠近正北的腿全出问题：罗盘读 0.2 度、仪器改正加磁偏角合计 -0.6 度，改正后应该是 359.6 度（正北略偏西），可归算出来的方位角写成 0.00 度，东向分量正好是 0，支洞画得死贴正北。试了几个值，-0.1、-0.4、-0.75、-1 这种改正后只比 0 小一点的方位角都被报成 0，而 -45、-370 这种大负值仍然正确回绕成 315、350。请修复方位角规范化，让任何有限负值都回绕进 [0,360)，同时保持 0、360、720 仍折叠为 0、NaN 仍然传播、正数不受影响，并保证全量测试通过。

## 含 Bug 版本

- 仓库：zhanglei10281852-gif/gogo-56
- 仓库地址：https://github.com/zhanglei10281852-gif/gogo-56.git
- parent SHA：53da4bcf0750e47f25a926a53b62107dd3b09ac4

## 复现步骤

```bash
git clone -- https://github.com/zhanglei10281852-gif/gogo-56.git bug-repro
cd bug-repro
git checkout --detach 53da4bcf0750e47f25a926a53b62107dd3b09ac4
go test ./internal/units -run "^TestNormalizeAzimuthWrapsBearingsJustWestOfNorth$" -count=1 -v
```

## 双架构完整错误信息

### linux/amd64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/units -run "^TestNormalizeAzimuthWrapsBearingsJustWestOfNorth$" -count=1 -v
=== RUN   TestNormalizeAzimuthWrapsBearingsJustWestOfNorth
    azimuth_wrap_regression_test.go:26: NormalizeAzimuth(-0.1) = 0, want 359.9
--- FAIL: TestNormalizeAzimuthWrapsBearingsJustWestOfNorth (0.00s)
FAIL
FAIL	CaveLoop/internal/units	0.002s
FAIL

```

stderr：

```text
warning: internal/units/azimuth_wrap_regression_test.go has type 100755, expected 100644
warning: internal/units/azimuth_wrap_regression_test.go has type 100755, expected 100644

```

### linux/arm64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/units -run "^TestNormalizeAzimuthWrapsBearingsJustWestOfNorth$" -count=1 -v
=== RUN   TestNormalizeAzimuthWrapsBearingsJustWestOfNorth
    azimuth_wrap_regression_test.go:26: NormalizeAzimuth(-0.1) = 0, want 359.9
--- FAIL: TestNormalizeAzimuthWrapsBearingsJustWestOfNorth (0.01s)
FAIL
FAIL	CaveLoop/internal/units	0.127s
FAIL

```

stderr：

```text
warning: internal/units/azimuth_wrap_regression_test.go has type 100755, expected 100644
warning: internal/units/azimuth_wrap_regression_test.go has type 100755, expected 100644

```

## 通过条件

NormalizeAzimuth(-0.1)=359.9、(-0.4)=359.6、(-0.75)=359.25、(-1)=359、(-2.5)=357.5、(-45)=315、(-370)=350；负值结果始终落在 [0,360) 且与输入是同一个方位（AzimuthSeparation 为 0）、OppositeAzimuth 对原值与回绕值一致；0/360/720 仍折叠为 0，NaN 仍传播；定向测试、全量 go test ./... -count=1 与 go build ./... && go vet ./... 全部通过；校准与远端复跑均在 golang:1.22 linux/amd64 单架构完成。
