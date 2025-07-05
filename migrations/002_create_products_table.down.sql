-- Drop trigger first
DROP TRIGGER IF EXISTS trigger_update_products_updated_at ON products;

-- Drop function
DROP FUNCTION IF EXISTS update_products_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_products_sku;
DROP INDEX IF EXISTS idx_products_category;
DROP INDEX IF EXISTS idx_products_price;
DROP INDEX IF EXISTS idx_products_stock;
DROP INDEX IF EXISTS idx_products_created_at;
DROP INDEX IF EXISTS idx_products_deleted_at;

-- Drop table
DROP TABLE IF EXISTS products; 