create table users (
	id bigserial,
	username varchar not null,
	password varchar not null,
	is_admin boolean not null default false,
	created_at timestamp not null default now(),
	primary key (id)
);
create unique index users_username_idx on users (username);

create table auth_tokens (
	id varchar,
	user_id bigint not null,
	created_at timestamp not null default now(),
	primary key (id),
	foreign key (user_id) references users (id)
);

create table shops (
	id bigserial,
	name varchar not null,
	description varchar not null,
	photos varchar[] not null,
	created_at timestamp not null default now(),
	primary key (id)
);

create table reviews (
	id bigserial,
	shop_id bigint not null,
	user_id bigint not null,
	rating smallint not null,
	comment varchar not null,
	photos varchar[] not null,
	created_at timestamp not null default now(),
	primary key (id),
	foreign key (shop_id) references shops (id),
	foreign key (user_id) references users (id)
);
