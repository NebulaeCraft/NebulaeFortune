# NebulaeFortune

星云娘(Nebulae Robot)的今日人品&今日运势实现

# Usage
```bash
go mod tidy  // 安装依赖
go build     // 编译
./main       // 运行
```

在`localhost:1145`下提供服务

- `/jrrp` 今日人品
  -  `/jrrp?rp=114` 指定人品
- `/jrys` 今日运势
  - `/jrys?ys=冰魂&emo=1` 指定运势、 表情

# Config
表情所用图片放置在`/emotion`下，格式需为`.png`。

为了避免被夹，最后输出格式为`.gif`。同时带来`gif`的限制，文件中只包含256色。请提前将文件调色板限制到256色中。否则使用内置[Floyd–Steinberg dithering](https://en.wikipedia.org/wiki/Floyd%E2%80%93Steinberg_dithering)算法进行颜色舍入`Plan9`空间中。

# License
[MIT](LICENSE)

# Thanks
感谢晶蓝姐姐的图！

**Welcome for PRs and Issues!**

