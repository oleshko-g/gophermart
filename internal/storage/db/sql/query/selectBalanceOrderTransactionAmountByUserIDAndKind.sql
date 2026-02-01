-- name: SelectBalanceOrderTransactionAmountByUserIDAndKind :many
WITH
  cte_transactions AS (
    SELECT
      transactions.created_at,
      transactions.amount,
      transactions.order_id
    FROM
      transactions
    WHERE
      transactions.user_id = $1
      AND kind = $2
  )
SELECT
  orders.number,
  cte_transactions.created_at,
  cte_transactions.amount
FROM
  orders
  JOIN cte_transactions ON orders.id = cte_transactions.order_id;
