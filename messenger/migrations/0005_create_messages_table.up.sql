CREATE TABLE IF NOT EXISTS public.messages
(
    message_id character varying(36) NOT NULL,
    user_id bigint NOT NULL,
    chat_id bigint NOT NULL,
    message_body text,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
    PRIMARY KEY (message_id),
    
    CONSTRAINT fk_chat FOREIGN KEY (chat_id) REFERENCES public.chats (chat_id) ON DELETE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.users (user_id) ON DELETE CASCADE
);