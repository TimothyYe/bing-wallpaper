## 必应壁纸API

一个从Bing.com获取每日壁纸的RESTful API

## 使用方法

### API 
* API地址: [https://bing.biturl.top](https://bing.biturl.top/)
* 方法: `HTTP GET`

### 参数

* `format` 响应格式，可以是`json`或`image`。__如果响应格式设置为`image`，您将被直接重定向到壁纸图像__。
* `image_format` 壁纸图像的格式，可用值为`jpg`或`webp`。默认值为`jpg`。
* `index` 壁纸的索引，从0开始。默认情况下，`0`表示获取今天的图像，`1`表示获取昨天的图像，依此类推。或者您可以指定为`random`，以在0到7之间随机选择一个索引。
* `mkt` 区域参数，默认值为`zh-CN`，您还可以使用`en-US`、`ja-JP`、`en-AU`、`en-GB`、`de-DE`、`en-NZ`、`en-CA`、`en-IN`、`fr-FR`、`fr-CA`、`it-IT`、`es-ES`、`pt-BR`、`en-ROW`。或者，您可以将其设置为`random`，以随机选择区域。
* `resolution` 壁纸图像的分辨率。`1920`是默认值，您还可以使用`1366`和`3840`或`UHD`（4K分辨率）。

可用分辨率选项如下：
```
UHD
1920x1200
1920x1080
1366x768
1280x768
1024x768
800x600
800x480
768x1280 (Portrait mode)
720x1280 (Portrait mode)
640x480
480x800 (Portrait mode)
400x240
320x240
240x320 (Portrait mode)
```

### 示例

* 请求

```text
https://bing.biturl.top/?resolution=UHD&format=json&index=0&mkt=zh-CN
```

* 响应

```json
{
"start_date": "20240803",
"end_date": "20240804",
"url": "https://www.bing.com/th?id=OHR.ImpalaOxpecker_ZH-CN9652434873_UHD.jpg",
"copyright": "黑斑羚和红嘴牛椋鸟，南非 (© Matrishva Vyas/Getty Images)",
"copyright_link": "https://www.bing.com/search?q=%E5%8F%8B%E8%B0%8A%E6%97%A5&form=hpcapt&mkt=zh-cn"
}
```

### CSS背景图像

您还可以使用此API设置CSS背景图像：

```text
background-image: url(https://bing.biturl.top/?resolution=1920&format=image&index=0&mkt=zh-CN);
height: 100%;
background-position: center;
background-repeat: no-repeat;
background-size: cover;
```

__演示__  

[https://biturl.top](https://biturl.top)

![https://github.com/TimothyYe/biturl/blob/master/screenshots/1.jpg?raw=true](https://github.com/TimothyYe/biturl/blob/master/screenshots/1.jpg?raw=true)

## 使用Docker运行

获取最新版本的Docker镜像：

```bash
docker pull timothyye/bing:latest
```

使用镜像名称和标签（YYYYMMDD或latest）启动容器，例如：

```bash
docker run -d --name=bing-wallpaper --restart=always -p 9000:9000 timothyye/bing:latest
```

## 开发相关

### 构建
```bash
git clone https://github.com/TimothyYe/bing-wallpaper.git
make build
```

### 运行

```bash
bw/bw run
