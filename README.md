# loot-chest

宝箱：钥匙消耗、周期内次数上限、到期刷新与开启幂等，只用 Go 标准库。

```
go run ./cmd/loot-chest --sample=keys
go run ./cmd/loot-chest --sample=cap
go run ./cmd/loot-chest --sample=refresh
go run ./cmd/loot-chest --sample=idem
go run ./cmd/loot-chest --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **钥匙**：每成功开一次消耗一把钥匙；没有钥匙时返回 false 且不改变任何状态。
- **次数上限**：每个宝箱在本周期内最多开 `maxOpens` 次，超出返回 false。
- **周期刷新**：距本周期起点满 `refreshMs` 时，已开次数清零、起点顺延到当前时刻，
  之后可以重新开满 `maxOpens` 次。
- **幂等**：同一个 `reqID` 只开一次，重复调用返回 false 且不扣钥匙。
- **不变量**：已开次数不超过 `maxOpens`；钥匙数不为负；
  同一串开启序列重复执行结果相同。
- **代价**：`Opens` 与 `Open` 不许遍历全部宝箱，`Scanned` 不随宝箱数线性增长。
- **规模**：10 万个宝箱、每秒 10 万次开启，单次摊还 O(1)。

## 输出契约（不改格式）

```
ok=.. keys=..
third=.. opens=..
before=.. after=..
second=.. keys=..
scanned=..
```

## 目录

```
chest.go              钥匙、次数、刷新与幂等
chest_test.go         单元测试
cmd/loot-chest/       命令行入口
```
