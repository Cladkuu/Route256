-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.status_mapping(
    id_from smallint,
    id_to smallint,
    constraint pk_status_mapping primary key (id_from,id_to)
    );
insert into public.status_mapping values
    (1, 2),
    (1, 5),
    (2, 3),
    (2, 5),
    (3, 4),
    (3, 5)
ON CONFLICT (id_from, id_to) DO UPDATE SET id_from = EXCLUDED.id_from, id_to= EXCLUDED.id_to;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.status_mapping
DROP CONSTRAINT IF EXISTS pk_status_mapping;
DROP TABLE IF EXISTS public.status_mapping;
-- +goose StatementEnd
