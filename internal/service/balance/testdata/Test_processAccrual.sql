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
  -- empty accrual
  (
    '019c22a8e94f7d75bcb973b05f070535',
    '568082882086285',
    '019c22a8e55179a4af55ace1ea6cdf68',
    'NEW',
    '2026-02-03 11:40:21.071876+03'
  ),
  -- accrual status INVALID accrual
  (
    '019c22a8ea2974ab87a39959d1301bf0',
    '73568464702',
    '019c22a8e55179a4af55ace1ea6cdf68',
    'NEW',
    '2026-02-03 11:40:21.289304+03'
  ),
  -- accrual status REGISTERED
  (
    '019c38078bdf7ed4a187d33b193aa2d9',
    '5441240171117',
    '019c22a8e55179a4af55ace1ea6cdf68',
    'NEW',
    '2026-02-07 15:15:44.607966+03'
  ),
  -- accrual status PROCESSING
  (
    '019c22a8e55774419eae11fbe458db69',
    '167862623743756',
    '019c22a8-e551-79a4-af55-ace1ea6cdf68',
    'NEW',
    '2026-02-03 11:40:20.055278+03'
  ),
  -- accrual status PROCESSED no accrual
  (
    '019c27b58f537dcb8aa9fd9ca0d61f93',
    '860548181160',
    '019c22a8e55179a4af55ace1ea6cdf68',
    'NEW',
    '2026-02-03 11:40:20.055278+03'
  ),
  -- accrual status PROCESSED with accrual
  (
    '019c380787e87a77b9f8d392ab5797df',
    '2402188500348',
    '019c22a8e55179a4af55ace1ea6cdf68',
    'NEW',
    '2026-02-07 15:15:43.592683+03'
  );