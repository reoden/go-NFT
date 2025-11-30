env "gorm" {
  # 数据库连接
  url = "postgres://postgres:123456@localhost:5432/nft?sslmode=disable"

  # 指定迁移目录
  migrations = ["db/migrations/atlas"]
}
