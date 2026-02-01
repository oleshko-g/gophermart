-- +goose Up
ALTER TABLE transactions
ADD COLUMN created_at TIMESTAMPTZ;


WITH
  cte_processed_orders AS (
    SELECT
      orders.id,
      orders.created_at
    FROM
      orders
    WHERE
      orders.status = 'PROCESSED'
  )
UPDATE transactions
SET
  created_at = cte_processed_orders.created_at
FROM
  cte_processed_orders
WHERE
  cte_processed_orders.id = transactions.order_id;


ALTER TABLE transactions
ALTER COLUMN created_at
SET NOT NULL;


-- +goose Down
ALTER TABLE transactions
DROP COLUMN created_at;
