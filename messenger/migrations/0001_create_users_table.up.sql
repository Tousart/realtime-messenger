CREATE TABLE IF NOT EXISTS public.users (
    user_id bigserial NOT NULL,
    user_name character varying(30) NOT NULL,
    password character varying(72) NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    PRIMARY KEY (user_id),
    CONSTRAINT users_user_name_unique UNIQUE (user_name) INCLUDE (user_name)
);