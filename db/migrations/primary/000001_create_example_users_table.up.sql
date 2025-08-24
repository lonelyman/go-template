CREATE TABLE IF NOT EXISTS "example_users" (
    "id" BIGSERIAL PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at" TIMESTAMPTZ,
    "name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "status" VARCHAR(50) NOT NULL DEFAULT 'active',
    "role" VARCHAR(50) NOT NULL DEFAULT 'user',
    "last_login_at" TIMESTAMPTZ,

    CONSTRAINT check_user_status CHECK (status IN ('active', 'inactive', 'banned')),
    CONSTRAINT check_user_role CHECK (role IN ('user', 'admin'))
);

-- สร้าง Partial Unique Index
CREATE UNIQUE INDEX IF NOT EXISTS "unique_active_email"
ON "example_users" ("email")
WHERE "deleted_at" IS NULL;

-- สร้าง Index อื่นๆ ที่น่าจะได้ใช้บ่อย
CREATE INDEX IF NOT EXISTS "idx_example_users_status" ON "example_users" ("status");
CREATE INDEX IF NOT EXISTS "idx_example_users_role" ON "example_users" ("role");