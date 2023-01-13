CREATE TABLE sellers (
    id serial primary key ,
    username text not null unique ,
    firstname text,
    lastname text,
    pass text not null ,
    seller_id int not null unique ,
    seller_key text not null unique
);
CREATE TABLE category(
    id serial primary key ,
    title_ru text unique not null ,
    title_eng text,
    description text,
    user_id int references sellers (id)
);
CREATE TABLE subcategory (
    id serial primary key ,
    title text not null unique ,
    category_id int references category (id)
);
CREATE TABLE products (
    id serial primary key ,
    content_key text not null unique ,
    subcategory_id int references subcategory (id)
);

CREATE TABLE transactions(
    id serial primary key ,
    category_name text ,
    subcategory_name text  ,
    client_email text,
    content_key text,
    amount int,
    profit int,
    amount_usd int,
    count int,
    unique_inv int,
    user_id int references sellers (id),
    unique_code text unique,
    date_check text,
    date_delivery text,
    date_confirmed text,
    state text
);