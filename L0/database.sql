CREATE TABLE delivery(
 uid uuid default gen_random_uuid() primary key,
 name text not null,
 phone varchar(30) not null,
 zip text not null,
 city text not null,
 email varchar(254) not null,
 address text not null,
 region text not null
)

CREATE TABLE payment(
 transaction uuid default gen_random_uuid() primary key,
 request_id text,
 currency varchar(10) not null,
 provider varchar(50) not null,
 amount double precision not null,
 payment_dt integer not null,
 bank varchar(50) not null,
 delivery_cost double precision not null,
 goods_total integer not null,
 custom_fee double precision not null
)

CREATE TABLE items(
 chrt_id integer not null,
 track_number text not null,
 price double precision not null,
 rid uuid default gen_random_uuid() primary key,
 name text not null,
 sale integer not null,
 size varchar(10) not null,
 total_price double precision not null,
 nm_id integer not null,
 brand varchar(150) not null,
 status integer not null,
 order_uid uuid not null references orders(order_uid) ON DELETE CASCADE
)

CREATE TABLE orders(
 order_uid uuid default gen_random_uuid() primary key,
 track_number text not null,
 entry text not null,
 delivery_uid uuid references delivery(uid) ON DELETE CASCADE,
 payment_uid uuid references payment(transaction) ON DELETE CASCADE,
 locale varchar(10) not null,
 internal_signature text,
 customer_id text not null,
 delivery_service text not null,
 shardkey text not null,
 sm_id integer not null,
 date_created timestamptz not null,
 oof_shard text not null
)

-- Для теста

INSERT INTO delivery (uid, name, phone, zip, city, email, address, region)
VALUES (
  'b563feb7-b2b8-4b6b-563f-eb7b2b84b612'::uuid,
  'Test Testov',
  '+9720000000',
  '2639809',
  'Kiryat Mozkin',
  'test@gmail.com',
  'Ploshad Mira 15',
  'Kraiot'
);

INSERT INTO payment (
  transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
)
VALUES (
  'b563feb7-b2b8-4b6b-563f-eb7b2b84b612'::uuid,
  '',
  'USD',
  'wbpay',
  1817,
  1637907727,
  'alpha',
  1500,
  317,
  0
);

INSERT INTO items (
  chrt_id,
  track_number,
  price,
  rid,
  name,
  sale,
  size,
  total_price,
  nm_id,
  brand,
  status,
  order_uid
)
VALUES (
  9934930,
  'WBILMTESTTRACK',
  453,
  'ab421908-7a76-4ae0-ab42-19087a764ae0'::uuid,
  'Mascaras',
  30,
  '0',
  317,
  2389212,
  'Vivienne Sabo',
  202,
  'b563feb7-b2b8-4b6b-563f-eb7b2b84b612'::uuid
);

INSERT INTO orders (
  order_uid, track_number, entry, delivery_uid, payment_uid,
  locale, internal_signature, customer_id, delivery_service, shardkey,
  sm_id, date_created, oof_shard
)
VALUES (
  'b563feb7-b2b8-4b6b-563f-eb7b2b84b612'::uuid,
  'WBILMTESTTRACK',
  'WBIL',
  'b563feb7-b2b8-4b6b-563f-eb7b2b84b612'::uuid, -- delivery_uid
  'b563feb7-b2b8-4b6b-563f-eb7b2b84b612'::uuid, -- payment_uid
  'en',
  '',
  'test',
  'meest',
  '9',
  99,
  '2021-11-26T06:22:19Z'::timestamptz,
  '1'
);

