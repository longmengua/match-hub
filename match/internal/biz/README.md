# biz

- usecase
    - 放使用業務邏輯的地方

- interface
    - repo 的 interface，放業務邏輯。

- entity
    - 「代表業務的核心狀態 + 規則 + 行為。」
    - 只存在於 biz/entity.go 中。
    - 是 “乾淨” 的，不依賴任何框架或 infra。
    - 可以含有行為（方法）、驗證邏輯。

- dto
    - 「轉給外部/上層通訊的資料格式。」
    - 無邏輯，純資料。
    - 出現在 api/、data/dto/、service/ 之間。
    - 常用來序列化/反序列化 HTTP 或 gRPC 請求。

- dao
    - 「對應資料表結構、與 DB 互動的物件。」
    - 對應 DB 欄位，可能有 ORM tag。
    - 出現在 data/model/ 或 data/mysql/。
    - 不應含業務邏輯，只處理資料持久化。

- 資料流

API 層 (DTO)
    ↓
Service 層 → DTO 轉 Entity
    ↓
Biz 層 (Entity with business logic)
    ↓
Repo interface
    ↓
Data 層 → Entity 轉 DAO → 寫入 DB