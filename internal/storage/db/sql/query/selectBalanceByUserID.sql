-- name: SelectBalanceByUserID :one
SELECT
  current,
  withdrawn_sum,
  last_transaction_id
FROM
  user_balances
WHERE
  user_id = $1;
