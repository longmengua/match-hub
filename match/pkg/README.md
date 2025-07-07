# pkg

- Go 社群整體推薦使用 pkg/ 命名，用來放 跨領域通用邏輯，包含：
    - 公共 function（例如 ParseTimeRFC3339、ToFloatSafe）
    - 公共 struct 封裝（如 log, config 擴展）
    - 中介錯誤包裝（error wrapping with stack trace）
    - 驗證工具（如 email/phone 格式、UUID 等）

pkg/
├── logger/         # 封裝 Logger 並支援 zap/slog
├── timeu/          # 時間處理（轉字串、加減、區間）
├── errcode/        # 錯誤碼包裝 + gRPC + HTTP 對應
├── randx/          # 隨機工具（字串、數字、UUID）
├── configx/        # 設定讀取 helper