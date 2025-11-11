create table if not exists orders (
    id uuid primary key,
    item varchar(500) not null,
    quantity integer not null check (quantity > 0)
);
