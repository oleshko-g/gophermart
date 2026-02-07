INSERT INTO
  users (
    id,
    login,
    hashed_password,
    created_at,
    updated_at,
    deleted_at
  )
VALUES
  (
    '019c22a8e55179a4af55ace1ea6cdf68',
    'user_Test_pocessAccrual',
    '$2a$10$uIdSMKOJ/dLQfsPDhE7/mOvGcqu0CgSAJ.US9VsmOeq.IJl5RoRMy',
    '2026-02-03 11:40:20.049632+03',
    '2026-02-03 11:40:20.049632+03',
    NULL
  );

-- user_Test_pocessAccrual
INSERT INTO
  orders (id, number, user_id, STATUS, created_at)
VALUES
  -- no accrual
  (
    '019c22a8e94f7d75bcb973b05f070535',
    '568082882086285',
    '019c22a8e55179a4af55ace1ea6cdf68',
    'NEW',
    '2026-02-03 11:40:21.071876+03'
  ),
  -- invalid accrual
  (
    '019c22a8ea2974ab87a39959d1301bf0',
    '73568464702',
    '019c22a8e55179a4af55ace1ea6cdf68',
    'NEW',
    '2026-02-03 11:40:21.289304+03'
  ),
  -- PROCESSED no accrual
  (
    '019c22a8e55774419eae11fbe458db69',
    '167862623743756',
    '019c22a8e55179a4af55ace1ea6cdf68',
    'NEW',
    '2026-02-03 11:40:20.055278+03'
  ),
  -- PROCESSED with accrual
  (
    '019c3807-87e8-7a77-b9f8-d392ab5797df',
    '2402188500348',
    '019c22a8e55179a4af55ace1ea6cdf68',
    'NEW',
    '2026-02-07 15:15:43.592683+03'
  );