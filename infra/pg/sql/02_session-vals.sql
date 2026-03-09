create table session_vals
(
    ses_id bigint not null,
    kv_key text not null,
    kv_val text not null default '',
    foreign key (ses_id) references session (id) on delete cascade,
    constraint uq_ses_key unique (ses_id, kv_key)
);

create or replace procedure seskv_set(p_ses_key text, p_key text, p_val text) as $$
    declare v_ses_id bigint;
    declare v_key text;
    begin
        -- Check if session and key exists.
        select s.id, kv.kv_key
        from active_session s
            left join session_vals kv on s.id = kv.ses_id
        where s.session_key = p_ses_key
        into v_ses_id, v_key;
        if v_ses_id is null then
            raise exception 'No session with that ID exists';
        end if;
        -- Upsert
        if v_key is null then
            insert into session_vals (ses_id, kv_key, kv_val) values
            (
                v_ses_id, p_key, p_val
            );
        else
            update session_vals
            set kv_val = p_val
            where ses_id = v_ses_id
                and kv_key = v_key;
        end if;
    end;
$$ language plpgsql;

create or replace function seskv_get(p_ses_key text, p_key text) returns text as $$
    declare v_ses_id bigint;
    declare r_val text;
    begin
        select kv.kv_val
        from session_vals kv
            join session s on kv.ses_id = s.id
        where s.session_key = p_ses_key
            and kv_key = p_key
        into r_val;
        return r_val;
    end;
$$ language plpgsql;
