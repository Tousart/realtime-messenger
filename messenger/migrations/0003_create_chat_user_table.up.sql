CREATE TABLE IF NOT EXISTS public.chat_user
(
    chat_id  bigint NOT NULL,
    user_id  bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    
    PRIMARY KEY (chat_id, user_id),
    
    CONSTRAINT fk_chat FOREIGN KEY (chat_id) REFERENCES public.chats (chat_id) ON DELETE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.users (user_id) ON DELETE CASCADE
);