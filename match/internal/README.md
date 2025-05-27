# internal

## 資料流

API 層 (DTO)
    ↓
Service 層 → DTO 轉 Entity
    ↓
Biz 層 (Entity with business logic)
    ↓
Repo interface
    ↓
Data 層 → Entity 轉 DAO → 寫入 DB

## 資料夾結構

- biz
    - 業務邏輯層，負責 usecase 與 entity，沒有基礎建設依賴，通常分為下列層級：
        - Entity
            - 表示業務核心模型，不含邏輯或依賴
            - 代表業務的核心狀態 + 規則 + 行為
            - 是 “乾淨” 的，不依賴任何框架或 infra。
            - 可以含有行為（方法）、驗證邏輯。
        - DTO(Data Transfer Object)
            - 「轉給外部/上層通訊的資料格式。」
            - 無邏輯，純資料。
            - 出現在 api/、data/dto/、service/ 之間。
            - 常用來序列化/反序列化 HTTP 或 gRPC 請求。
        - DAO (Data Access Object)
            - 「對應資料表結構、與 DB 互動的物件。」
            - 對應 DB 欄位，可能有 ORM tag。
            - 出現在 data/model/ 或 data/mysql/。
            - 不應含業務邏輯，只處理資料持久化。
        - UseCase
            - 封裝業務邏輯，依賴 repo 介面
        - Repo Interface
            - 定義資料存取的抽象介面

- data
    - 資料存取層（如 DB、RPC、API 存取），通常分為兩層：
        - Repo Impl
            - 實作 biz 中定義的介面（如 UserRepo）。
            - 與 ORM（如 gorm）或 RPC client 做實際資料互動。
            - 通常每個資源一個檔案，如 user.go。
        - Data Provider
            - 初始化和統一管理外部資源（DB、Cache、其他 service client）。
            - 提供依賴注入來源（如 *gorm.DB、*redis.Client 等）。
            - 一般集中放在 data.go

- service
    - gRPC/HTTP server handler（調用 usecase），也可稱 controller。

- server
    - 初始化 gRPC、HTTP server 的設定
