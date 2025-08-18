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
 transaction uuid default gen_random_uuid() primary key references orders(order_uid) ON DELETE CASCADE,
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
 track_number text not null references orders(track_number) ON DELETE CASCADE,
 price double precision not null,
 rid uuid default gen_random_uuid() primary key,
 name text not null,
 sale integer not null,
 size varchar(10) not null,
 total_price double precision not null,
 nm_id integer not null,
 brand varchar(150) not null,
 status integer not null
)

CREATE TABLE orders(
 order_uid uuid default gen_random_uuid() primary key,
 track_number text unique not null,
 entry text not null,
 delivery_uid uuid references delivery(uid) ON DELETE CASCADE,
 locale varchar(10) not null,
 internal_signature text,
 customer_id text not null,
 delivery_service text not null,
 shardkey text not null,
 sm_id integer not null,
 date_created timestamptz not null,
 oof_shard text not null
)

