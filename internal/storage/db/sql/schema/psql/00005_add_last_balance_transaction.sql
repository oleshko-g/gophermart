-- +goose Up
DROP VIEW user_balances;

CREATE VIEW user_balances (
user_id,
current,
withdrawn_sum,
last_transaction_id
) AS
WITH
cte_users AS (
  SELECT
    id AS user_id
  FROM
  users
  ),
accruals AS (
  SELECT
    user_id,
    sum(amount) AS amount_sum
  FROM
  transactions
  WHERE
  kind = 'ACCRUAL'
  GROUP BY
  user_id
  ),
withdrawals AS (
  SELECT
    user_id,
    sum(amount) AS amount_sum
  FROM
  transactions
  WHERE
  kind = 'WITHDRAWAL'
  GROUP BY
  user_id
  ),
cte_ranked_transactions AS (
  SELECT
    rank() OVER (
    PARTITION BY
    user_id
    ORDER BY
    transactions.id DESC
    ),
    transactions.id,
    transactions.user_id
  FROM
  transactions
  ),
cte_last_transactions AS (
  SELECT
    cte_ranked_transactions.user_id,
    cte_ranked_transactions.id
  FROM
  cte_ranked_transactions
  WHERE
  cte_ranked_transactions.rank = 1
  ),
cte_balances AS (
  SELECT
    cte_users.user_id,
    COALESCE(accruals.amount_sum, 0) AS accrued_sum,
    COALESCE(withdrawals.amount_sum, 0) AS withdrawn_sum,
    COALESCE(accruals.amount_sum, 0) - COALESCE(withdrawals.amount_sum, 0) AS current,
    cte_last_transactions.id AS last_transaction_id
  FROM
  cte_users
  LEFT JOIN accruals USING (user_id)
  LEFT JOIN withdrawals USING (user_id)
  LEFT JOIN cte_last_transactions USING (user_id)
  )
SELECT
  user_id,
  accrued_sum - withdrawn_sum AS current,
  withdrawn_sum,
  last_transaction_id
FROM
cte_balances;


-- +goose Down
DROP VIEW user_balances;


CREATE VIEW user_balances (user_id, current, withdrawn_sum) AS
WITH
cte_users AS (
  SELECT
    id AS user_id
  FROM
  users
  ),
accruals AS (
  SELECT
    user_id,
    sum(amount) AS amount_sum
  FROM
  transactions
  WHERE
  kind = 'ACCRUAL'
  GROUP BY
  user_id
  ),
withdrawals AS (
  SELECT
    user_id,
    sum(amount) AS amount_sum
  FROM
  transactions
  WHERE
  kind = 'WITHDRAWAL'
  GROUP BY
  user_id
  ),
balances AS (
  SELECT
    cte_users.user_id,
    COALESCE(accruals.amount_sum, 0) AS accrued_sum,
    COALESCE(withdrawals.amount_sum, 0) AS withdrawn_sum,
    COALESCE(accruals.amount_sum, 0) - COALESCE(withdrawals.amount_sum, 0) AS current
  FROM
  cte_users
  LEFT JOIN accruals USING (user_id)
  LEFT JOIN withdrawals USING (user_id)
  )
SELECT
  user_id,
  accrued_sum - withdrawn_sum AS current,
  withdrawn_sum
FROM
BALANCES;
