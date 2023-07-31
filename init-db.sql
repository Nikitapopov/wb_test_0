create table if not exists public.order (
    order_uid          varchar(512) primary key not null,
	track_number       varchar(512)             not null,
	entry              varchar(512)             not null,
    delivery           jsonb                    not null,
    payment            jsonb                    not null,
    items              jsonb                    not null,
	locale             varchar(512)             not null,
	internal_signature varchar(512)             not null,
	customer_id        varchar(512)             not null,
	delivery_service   varchar(512)             not null,
	shardkey           varchar(512)             not null,
	sm_id              integer                  not null,
	date_created       varchar(512)             not null,
	oof_shard          varchar(512)             not null
);