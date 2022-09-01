CREATE TABLE IF NOT EXISTS public.currency(
    iso_code varchar(3) primary key not null
    );
insert into public.currency values
                                ('RUB'),
                                ('USD'),
                                ('BYN'),
                                ('EUR'),
                                ('CNY')
    ON CONFLICT (iso_code) DO UPDATE SET iso_code = EXCLUDED.iso_code;