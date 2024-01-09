create table users (
id               uuid primary key not null,
first_name       varchar(25),
last_name        varchar(25),
email            varchar(25),
phone            varchar(13)
);
create table orders (
id               uuid primary key not null,
amount           varchar(10),
user_id          uuid references users(id),
created_at       timestamp default current_timestamp
);
create table products (
id               uuid primary key not null,
price            integer,
product_name     varchar(15)
);
create table order_products (
id               uuid primary key not null,
order_id         uuid references orders(id),
product_id       uuid references products(id),
quantity         integer,
price            integer
);