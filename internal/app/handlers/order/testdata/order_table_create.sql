CREATE TABLE IF NOT EXISTS public.orders(
                                            id serial PRIMARY KEY,
                                            status_id smallint not null,
                                            price smallint not null,
                                            currency varchar(3) not null,
    order_code varchar(10) not null
    );