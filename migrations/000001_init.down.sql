drop index if exists idx_point_lat_lon;
drop index if exists idx_point_owner_id;

drop index if exists idx_map_archive_status;
drop index if exists idx_map_archive_uploaded_by;
drop index if exists idx_map_archive_map_id;

drop index if exists idx_map_current_archive_id;
drop index if exists idx_map_created_by;

drop index if exists uq_map_archive_storage_key;
drop index if exists uq_map_slug_lower;
drop index if exists uq_user_email_lower;
drop index if exists uq_user_username_lower;

drop table if exists point;

alter table map drop constraint if exists fk_map_current_archive_belongs_to_map;

drop table if exists map_archive;
drop table if exists map;
drop table if exists "user";