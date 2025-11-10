**依赖安装**：在项目目录下运行`go mod tidy`命令会自动下载并安装所有必需的依赖包
**运行方式**：使用`go run .`
**生成go文件**: `abigen --abi bin\contracts\Counter.abi --pkg main --type Counter --out Counter.go`