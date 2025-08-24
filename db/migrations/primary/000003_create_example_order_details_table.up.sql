CREATE TABLE IF NOT EXISTS "example_order_details" (
    "id" BIGSERIAL PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at" TIMESTAMPTZ,
    "order_id" BIGINT NOT NULL,
    "product_sku" VARCHAR(100) NOT NULL,
    "product_name" VARCHAR(255) NOT NULL,
    "quantity" INT NOT NULL CHECK (quantity > 0),
    "unit_price" DECIMAL(10, 2) NOT NULL,
    "total_price" DECIMAL(12, 2) NOT NULL,

    CONSTRAINT fk_order
        FOREIGN KEY(order_id)
        REFERENCES example_orders(id)
        ON DELETE CASCADE -- ถ้า Order ถูกลบ ให้ลบ Detail นี้ไปด้วยเลย
);

CREATE INDEX IF NOT EXISTS "idx_example_order_details_order_id" ON "example_order_details" ("order_id");
CREATE INDEX IF NOT EXISTS "idx_example_order_details_product_sku" ON "example_order_details" ("product_sku");