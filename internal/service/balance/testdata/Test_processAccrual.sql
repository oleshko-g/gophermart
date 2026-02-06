-- +goose Up

--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: gennadyoleshko
--
COPY public.orders (id, number, user_id, status, created_at) FROM stdin;
019c22a8-e557-7441-9eae-11fbe458db69	167862623743756	019c22a8-e551-79a4-af55-ace1ea6cdf68	PROCESSED	2026-02-03 11:40:20.055278+03
019c22a8-e94f-7d75-bcb9-73b05f070535	568082882086285	019c22a8-e551-79a4-af55-ace1ea6cdf68	PROCESSED	2026-02-03 11:40:21.071876+03
019c22a8-ea29-74ab-87a3-9959d1301bf0	73568464702	019c22a8-ea25-7d36-968d-b83f4cc750ae	NEW	2026-02-03 11:40:21.289304+03
\.


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: gennadyoleshko
--

COPY public.transactions (id, user_id, order_id, kind, amount, created_at) FROM stdin;
019c22a8-e898-7cb2-b8e5-b4732c25bbf1	019c22a8-e551-79a4-af55-ace1ea6cdf68	019c22a8-e557-7441-9eae-11fbe458db69	ACCRUAL	72998	2026-02-03 11:40:20.888833+03
019c22a8-e950-7639-9f3d-1ea92858ff11	019c22a8-e551-79a4-af55-ace1ea6cdf68	019c22a8-e94f-7d75-bcb9-73b05f070535	WITHDRAWAL	10989	2026-02-03 11:40:21.072408+03
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: gennadyoleshko
--

COPY public.users (id, login, hashed_password, created_at, updated_at, deleted_at) FROM stdin;
019c22a8-e551-79a4-af55-ace1ea6cdf68	pWNAiTM8LVTdB	$2a$10$uIdSMKOJ/dLQfsPDhE7/mOvGcqu0CgSAJ.US9VsmOeq.IJl5RoRMy	2026-02-03 11:40:20.049632+03	2026-02-03 11:40:20.049632+03	\N
019c22a8-e9ab-76fd-9f93-b5d5b536709d	Hrdi3	$2a$10$ffSwQQCvgyiCYp1mbCcKdOnBlzpJSw4Ek4n1vSEhTa1oDRfjlxZtm	2026-02-03 11:40:21.163458+03	2026-02-03 11:40:21.163458+03	\N
019c22a8-ea25-7d36-968d-b83f4cc750ae	gs4wBnQA0H	$2a$10$lKwKcm7JQ3XkTf.Np/UxMuTbEo0fwqsDXXJFDx9j1KeFm02ytbEU6	2026-02-03 11:40:21.285866+03	2026-02-03 11:40:21.285866+03	\N
019c22a8-ea68-7ad9-be7c-83069f5a0870	oseklbtFgg0T	$2a$10$teDimQlKmvvuzqa5Z30Zd.taXZm1gzkbGhK9zHSRUS9UHF/dmECFG	2026-02-03 11:40:21.352711+03	2026-02-03 11:40:21.352711+03	\N
\.
