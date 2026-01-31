-- name: SelectOrdersByUserID :many
WITH
  cte_transactions AS (
    SELECT
      order_id,
      amount AS accrual
    FROM
      transactions
    WHERE
      transactions.user_id = $1
      AND transactions.kind = 'ACCRUAL'
  ),
  cte_orders AS (
    SELECT
      id AS order_id,
      number,
      status,
      created_at
    FROM
      orders
    WHERE
      orders.user_id = $1
  ),
  user_orders AS (
    SELECT
      cte_orders.number,
      cte_orders.status,
      cte_orders.created_at,
      cte_transactions.accrual
    FROM
      cte_orders
      LEFT JOIN cte_transactions ON cte_orders.order_id = cte_transactions.order_id
    ORDER BY
      created_at ASC
  )
SELECT
  *
FROM
  user_orders;
