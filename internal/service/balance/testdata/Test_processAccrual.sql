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
    'pWNAiTM8LVTdB',
    '$2a$10$uIdSMKOJ/dLQfsPDhE7/mOvGcqu0CgSAJ.US9VsmOeq.IJl5RoRMy',
    '2026-02-03 11:40:20.049632+03',
    '2026-02-03 11:40:20.049632+03',
    NULL
  ),
  (
    '019c22a8e9ab76fd9f93b5d5b536709d',
    'Hrdi3',
    '$2a$10$ffSwQQCvgyiCYp1mbCcKdOnBlzpJSw4Ek4n1vSEhTa1oDRfjlxZtm',
    '2026-02-03 11:40:21.163458+03',
    '2026-02-03 11:40:21.163458+03',
    NULL
  ),
  (
    '019c22a8ea257d36968db83f4cc750ae',
    'gs4wBnQA0H',
    '$2a$10$lKwKcm7JQ3XkTf.Np/UxMuTbEo0fwqsDXXJFDx9j1KeFm02ytbEU6',
    '2026-02-03 11:40:21.285866+03',
    '2026-02-03 11:40:21.285866+03',
    NULL
  ),
  (
    '019c22a8ea687ad9be7c83069f5a0870',
    'oseklbtFgg0T',
    '$2a$10$teDimQlKmvvuzqa5Z30Zd.taXZm1gzkbGhK9zHSRUS9UHF/dmECFG',
    '2026-02-03 11:40:21.352711+03',
    '2026-02-03 11:40:21.352711+03',
    NULL
  );

INSERT INTO
  orders (id, number, user_id, STATUS, created_at)
VALUES
  (
    '019c22a8e94f7d75bcb973b05f070535',
    '568082882086285',
    '019c22a8e55179a4af55ace1ea6cdf68',
    'PROCESSED',
    '2026-02-03 11:40:21.071876+03'
  ),
  (
    '019c22a8ea2974ab87a39959d1301bf0',
    '73568464702',
    '019c22a8ea257d36968db83f4cc750ae',
    'NEW',
    '2026-02-03 11:40:21.289304+03'
  ),
  (
    '019c22a8e55774419eae11fbe458db69',
    '167862623743756',
    '019c22a8e55179a4af55ace1ea6cdf68',
    'PROCESSED',
    '2026-02-03 11:40:20.055278+03'
  );

INSERT INTO
  transactions (id, user_id, order_id, kind, amount, created_at)
VALUES
  (
    '019c22a8e8987cb2b8e5b4732c25bbf1',
    '019c22a8e55179a4af55ace1ea6cdf68',
    '019c22a8e55774419eae11fbe458db69',
    'ACCRUAL',
    72998,
    '2026-02-03 11:40:20.888833+03'
  ),
  (
    '019c22a8e95076399f3d1ea92858ff11',
    '019c22a8e55179a4af55ace1ea6cdf68',
    '019c22a8e94f7d75bcb973b05f070535',
    'WITHDRAWAL',
    10989,
    '2026-02-03 11:40:21.072408+03'
  );