CREATE TABLE IF NOT EXISTS public.chats
(
    chat_id bigserial NOT NULL,
    chat_name character varying(64),
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    PRIMARY KEY (chat_id)
);