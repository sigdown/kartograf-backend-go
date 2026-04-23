create table "user" (
    user_id bigserial primary key,
    username varchar(64) not null,
    display_name varchar(128),
    email varchar(255) not null,
    password_hash text not null,
    role varchar(16) not null check (role in ('USER', 'ADMIN')),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    constraint chk_user_username_not_blank
        check (btrim(username) <> ''),
    constraint chk_user_email_not_blank
        check (btrim(email) <> ''),
    constraint chk_user_display_name_not_blank
        check (display_name is null or btrim(display_name) <> '')
);

create table map (
    uuid uuid primary key,
    created_by bigint not null,
    current_archive_id uuid null,
    slug varchar(128) not null,
    title varchar(255) not null,
    description text,
    year smallint,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    constraint fk_map_created_by
        foreign key (created_by)
        references "user"(user_id)
        on delete restrict,

    constraint chk_map_slug_not_blank
        check (btrim(slug) <> ''),
    constraint chk_map_title_not_blank
        check (btrim(title) <> ''),
    constraint chk_map_year_range
        check (year is null or year between 1 and 2100)
);

create table map_archive (
    archive_id uuid primary key,
    map_id uuid not null,
    bucket varchar(128) not null,
    storage_key varchar(512) not null,
    uploaded_by bigint not null,
    size_bytes bigint,
    checksum varchar(128),
    status varchar(16) not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    constraint fk_map_archive_map_id
        foreign key (map_id)
        references map(uuid)
        on delete cascade,

    constraint fk_map_archive_uploaded_by
        foreign key (uploaded_by)
        references "user"(user_id)
        on delete restrict,

    constraint chk_map_archive_bucket_not_blank
        check (btrim(bucket) <> ''),
    constraint chk_map_archive_storage_key_not_blank
        check (btrim(storage_key) <> ''),
    constraint chk_map_archive_checksum_not_blank
        check (checksum is null or btrim(checksum) <> ''),
    constraint chk_map_archive_size_bytes_nonnegative
        check (size_bytes is null or size_bytes >= 0),
    constraint chk_map_archive_status
        check (status in ('UPLOADED', 'ACTIVE', 'REPLACED', 'DELETED')),

    constraint uq_map_archive_map_id_archive_id
        unique (map_id, archive_id)
);

alter table map
    add constraint fk_map_current_archive_belongs_to_map
    foreign key (uuid, current_archive_id)
    references map_archive(map_id, archive_id)
    on delete set null
    deferrable initially deferred;

create table point (
    point_id bigserial primary key,
    owner_id bigint not null,
    name varchar(255) not null,
    description text,
    lat double precision not null,
    lon double precision not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    constraint fk_point_owner_id
        foreign key (owner_id)
        references "user"(user_id)
        on delete cascade,

    constraint chk_point_name_not_blank
        check (btrim(name) <> ''),
    constraint chk_point_lat_range
        check (lat >= -90 and lat <= 90),
    constraint chk_point_lon_range
        check (lon >= -180 and lon <= 180)
);

create unique index uq_user_username_lower on "user"(lower(username));
create unique index uq_user_email_lower on "user"(lower(email));
create unique index uq_map_slug_lower on map(lower(slug));
create unique index uq_map_archive_storage_key on map_archive(storage_key);

create index idx_map_created_by on map(created_by);
create index idx_map_current_archive_id on map(current_archive_id);

create index idx_map_archive_map_id on map_archive(map_id);
create index idx_map_archive_uploaded_by on map_archive(uploaded_by);
create index idx_map_archive_status on map_archive(status);

create index idx_point_owner_id on point(owner_id);
create index idx_point_lat_lon on point(lat, lon);