CREATE TABLE IF NOT EXISTS "examples" (
   -- 4 Fields นี้มาจาก gorm.Model --
   "id" SERIAL PRIMARY KEY,
   "created_at" TIMESTAMPTZ,
   "updated_at" TIMESTAMPTZ,
   "deleted_at" TIMESTAMPTZ,

   -- 3 Fields นี้มาจากที่เราเพิ่มเข้าไปเอง --
   "name" VARCHAR(255) NOT NULL,
   "email" VARCHAR(255) UNIQUE NOT NULL,
   "status" VARCHAR(50) DEFAULT 'active'
);