CREATE TABLE IF NOT EXISTS "example_orders" (
    "id" BIGSERIAL PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at" TIMESTAMPTZ,
    "user_id" BIGINT NOT NULL,
    "order_number" VARCHAR(100) UNIQUE NOT NULL,
    "total_amount" DECIMAL(12, 2) NOT NULL,
    "status" VARCHAR(50) NOT NULL DEFAULT 'pending',
    "shipping_address" JSONB,
    "payment_details" JSONB,

    CONSTRAINT fk_user
        FOREIGN KEY(user_id) 
        REFERENCES example_users(id)
        ON DELETE RESTRICT, -- ป้องกันการลบ User ถ้ายังมี Order ค้างอยู่

    CONSTRAINT check_order_status CHECK (status IN ('pending', 'processing', 'shipped', 'completed', 'cancelled', 'refunded'))
);

CREATE INDEX IF NOT EXISTS "idx_example_orders_user_id" ON "example_orders" ("user_id");
CREATE INDEX IF NOT EXISTS "idx_example_orders_status" ON "example_orders" ("status");