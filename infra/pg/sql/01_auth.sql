create extension if not exists pgcrypto;

create table users
(
    id bigserial not null primary key,
    username text unique,
    admin bool default false,
    pass_hash text not null
);

create table authorizations
(
    id bigserial not null primary key,
    auth text not null unique,
    description text not null default ''
);

create table user_authz
(
    user_id bigint not null,
    auth_id bigint not null,
    granted timestamp not null default current_timestamp,
    revoked timestamp null,
    foreign key (auth_id) references authorizations(id) on delete cascade,
    foreign key (user_id) references users(id) on delete cascade,
    primary key (user_id, auth_id)
);

create view auth_grants as
select
    ua.user_id,
    u.username,
    ua.auth_id,
    a.auth,
    ua.granted,
    ua.revoked
from
    user_authz ua
    join users u on ua.user_id = u.id
    join authorizations a on ua.auth_id = a.id
where
    (ua.revoked is null or ua.revoked > current_timestamp)
;

create or replace function gen_passwd(p_plaint_text text) returns text as $$
begin
    return crypt(p_plaint_text, gen_salt('bf', 15));
end;
$$ language plpgsql;

create or replace function check_passwd(p_username text, p_pass text) returns boolean as $$
declare r_hash text;
begin
    select pass_hash from users where username = p_username into r_hash;
    if r_hash is null then
        -- Adding this to prevent timing-based data leaks.
        -- Without it, an invalid username results in a very quick return, which is a dead giveaway that the username is invalid.
        select crypt(p_pass, gen_salt('bf', 15)) into r_hash;
        return false;
    end if;
    return crypt(p_pass, r_hash) = r_hash;
end;
$$ language plpgsql;

create table user_audit
(
    username text not null,
    action text not null,
    event_time timestamp not null default current_timestamp
);

create table session
(
    id bigserial not null primary key,
    user_id bigint not null,
    session_key text not null default encode(gen_random_bytes(32), 'hex'),
    created_at timestamp not null default current_timestamp,
    revoked_at timestamp not null default current_timestamp + interval '30 minutes',
    max_ttl timestamp not null default current_timestamp + interval '2 hours',
    foreign key (user_id) references users(id) on delete cascade
);

create or replace procedure update_session_ttl(p_session_key text) as $$
declare r_max_ttl timestamp;
    declare new_revoked timestamp;
begin
    select max_ttl
    from session
    where session_key = p_session_key
      and revoked_at > current_timestamp
      and max_ttl > current_timestamp
    into r_max_ttl;
    if r_max_ttl is null then
        delete from session where session_key = p_session_key;
        raise exception 'No live sessions exist with this session key';
    end if;

    new_revoked := current_timestamp + interval '30 minutes';
    if new_revoked > r_max_ttl then
        update session
        set revoked_at = max_ttl
        where session_key = p_session_key;
    else
        update session
        set revoked_at = new_revoked
        where session_key = p_session_key;
    end if;
end;
$$ language plpgsql;

create or replace function create_session(p_username text) returns text as $$
    declare r_session_key text;
    declare r_user_id bigint;
    begin
        -- Get users with this username
        select id from users where username = p_username into r_user_id;
        if r_user_id is null then
            raise exception 'No user with that username exists';
        end if;

        select
            session_key
        from session s
        where s.revoked_at > current_timestamp
            and s.max_ttl > current_timestamp
            and s.user_id = r_user_id
        into r_session_key;
        if r_session_key is not null then
            call update_session_ttl(r_session_key);
            return r_session_key;
        end if;

        insert into session (user_id) values (r_user_id) returning session_key into r_session_key;
        -- Delete old sessions while we're here.
        delete from session where revoked_at < current_timestamp;
        return r_session_key;
    end;
$$ language plpgsql;
