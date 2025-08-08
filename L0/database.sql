CREATE TABLE order(
 order_uid uuid default gen_random_uuid() primary key,
 track_number text not null
 entry text not null
 delivery_uid uuid references delivery(uid),
 payment_uid uuid references payment(transaction),
 items_rid uuid references items(rid),
 locale varchar(10) not null,
 internal_signature text,
 customer_id text not null,
 delivery_service text not null,
 shardkey text not null,
 sm_id integer not null,
 date_created timestamptz not null,
 oof_shard text not null,
)

CREATE TABLE delivery(
 uid uuid default gen_random_uuid() primary key,
 name text not null,
 phone varchar(30) not null,
 zip text not null,
 city text not null,
 email varchar(254) not null,
 address text not null,
 region text not null,
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
 custom_fee double precision not null,
)

CREATE TABLE items(
 items_uid uuid default gen_random_uuid() primary key,
 track_number text not null,
 rid uuid not null,
 status integer not null,
 product_id integer references product(nm_id),
)

CREATE TABLE product(
 nm_id integer primary key,
 chrt_id integer not null,
 price double precision not null,
 name text not null,
 sale integer not null,
 size varchar(10) not null,
 total_price double precision not null,
 brand varchar(150) not null,
)




















CREATE TABLE Item(
 ID SERIAL Primary key,
 Name varchar(50),
 Info varchar(300));

CREATE TABLE Person(
 ID serial Primary key,
 Type varchar(15),
 Name varchar(50),
 Adress varchar(50),
 Phone varchar(12));

CREATE TABLE Coming(
 ID serial Primary key,
 Data date,
 IdPerson integer references Person(ID) ON DELETE CASCADE,
 TotalPrice double precision);

CREATE TABLE ComingInfo(
 ID serial Primary key,
 IdComing integer references Coming(ID) ON DELETE CASCADE,
 IdItem integer references Item(ID) ON DELETE CASCADE,
 Quantity integer,
 Price double precision);

CREATE TABLE Expense(
 ID serial Primary key,
 Data date,
 IdPerson integer references Person(ID) ON DELETE CASCADE,
 TotalPrice double precision);

CREATE TABLE ExpenseInfo(
 ID serial Primary key,
 IdExpense integer references Expense(ID) ON DELETE CASCADE,
 IdItem integer references Item(ID) ON DELETE CASCADE,
 Quantity integer,
 Price double precision);
 
CREATE TABLE Payment(
 ID serial Primary key,
 IdExpense integer references Expense(ID) ON DELETE CASCADE,
 Data date,
 Payment double precision);

