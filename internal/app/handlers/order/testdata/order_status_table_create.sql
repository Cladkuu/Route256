CREATE TABLE IF NOT EXISTS public.status(
                                            id smallserial PRIMARY KEY,
                                            status varchar(55) not null unique
    );
insert into public.status values
                              (1, 'NEW'),
                              (2, 'IN_PACKAGING'),
                              (3, 'IN_DELIVERY'),
                              (4, 'RECEIVED'),
                              (5, 'CANCELLED')
    ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status;