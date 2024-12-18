DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

CREATE TABLE public.chat (
    chat_name text NOT NULL,
    chat_type_id integer NOT NULL,
    avatar_path text,
    chat_link_name text,
    id uuid NOT NULL
);


ALTER TABLE public.chat OWNER TO postgres;

--
-- Name: chat_type; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chat_type (
    id integer NOT NULL,
    value text NOT NULL
);


ALTER TABLE public.chat_type OWNER TO postgres;

--
-- Name: chat_types_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.chat_type ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.chat_types_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: chat_user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chat_user (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_role_id integer NOT NULL,
    chat_id uuid NOT NULL,
    user_id uuid NOT NULL,
    send_notifications boolean DEFAULT true NOT NULL 
);


ALTER TABLE public.chat_user OWNER TO postgres;

--
-- Name: contact; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.contact (
    user_id uuid NOT NULL,
    contact_id uuid NOT NULL,
    id uuid DEFAULT gen_random_uuid() NOT NULL
);


ALTER TABLE public.contact OWNER TO postgres;

CREATE TABLE public.sticker (
	id uuid NOT NULL,
	sticker_path text NOT NULL,
	CONSTRAINT sticker_pk PRIMARY KEY (id),
	CONSTRAINT sticker_unique UNIQUE (sticker_path)
);

CREATE TABLE public.sticker_pack (
	id uuid NOT NULL,
	"name" text NULL,
	photo text NOT NULL,
	CONSTRAINT sticker_pack_pk PRIMARY KEY (id)
);

CREATE TABLE public.sticker_sticker_pack (
	id uuid NOT NULL,
	sticker uuid NOT NULL,
	pack uuid NOT NULL,
	CONSTRAINT sticker_sticker_pack_pk PRIMARY KEY (id),
	CONSTRAINT sticker_sticker_pack_sticker_fk FOREIGN KEY (sticker) REFERENCES public.sticker(id),
	CONSTRAINT sticker_sticker_pack_sticker_pack_fk FOREIGN KEY (pack) REFERENCES public.sticker_pack(id)
);

--
-- Name: message; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.message_type (
    id integer NOT NULL,
    value text NOT NULL
);

ALTER TABLE ONLY public.message_type
    ADD CONSTRAINT message_type_pkey PRIMARY KEY (id);

CREATE TABLE public.message (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    chat_id uuid NOT NULL,
    author_id uuid NOT NULL,
    branch_id uuid,
    message text,
    sent_at timestamp with time zone NOT NULL,
    is_redacted boolean DEFAULT false NOT NULL,
    sticker_path text,
    message_type_id integer DEFAULT 1 NOT NULL,
    CONSTRAINT message_sticker_path_fk FOREIGN KEY (sticker_path) REFERENCES public.sticker(sticker_path)
);


ALTER TABLE public.message OWNER TO postgres;

ALTER TABLE public.message_type ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.message_type_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

ALTER TABLE ONLY public.message
    ADD CONSTRAINT message_type_id_fk FOREIGN KEY (message_type_id) REFERENCES public.message_type(id) NOT VALID;


CREATE TABLE public.payload_type (
	id integer GENERATED ALWAYS AS IDENTITY NOT NULL,
	value text NOT NULL,
	CONSTRAINT payload_type_pk PRIMARY KEY (id),
	CONSTRAINT payload_type_unique UNIQUE (value)
);

ALTER TABLE public.payload_type OWNER TO postgres;

--
-- Name: message_payload; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.message_payload (
    payload_path text NOT NULL,
    id uuid NOT NULL,
    message_id uuid NOT NULL,
    payload_type integer DEFAULT 1 NOT NULL,
    filename text NOT NULL,
    size int NOT NULL
);

ALTER TABLE public.message_payload OWNER TO postgres;

ALTER TABLE ONLY public.message_payload
    ADD CONSTRAINT payload_type_id_fk FOREIGN KEY (payload_type) REFERENCES public.payload_type(id) NOT VALID;


--
-- Name: user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."user" (
    username text NOT NULL,
    version integer NOT NULL,
    password text NOT NULL,
    name text NOT NULL,
    bio text,
    birthdate timestamp with time zone,
    avatar_path text,
    id uuid NOT NULL
);


ALTER TABLE public."user" OWNER TO postgres;

--
-- Name: user_role; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_role (
    id integer NOT NULL,
    value text NOT NULL
);


ALTER TABLE public.user_role OWNER TO postgres;

--
-- Name: user_roles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.user_role ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.user_roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


ALTER TABLE ONLY public.chat_user
    ADD CONSTRAINT chat_id_and_user_id_uniq UNIQUE (chat_id, user_id);


ALTER TABLE ONLY public.chat
    ADD CONSTRAINT chat_link_name_uniq UNIQUE (chat_link_name);


--
-- Name: chat chat_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chat
    ADD CONSTRAINT chat_pkey PRIMARY KEY (id);


--
-- Name: chat_type chat_types_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chat_type
    ADD CONSTRAINT chat_types_pkey PRIMARY KEY (id);


--
-- Name: chat_user chat_user_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chat_user
    ADD CONSTRAINT chat_user_pkey PRIMARY KEY (id);

--
-- Name: message branch_id_fk_messages_chat_id_pk_chat; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.message
    ADD CONSTRAINT branch_id_fk_messages_chat_id_pk_chat FOREIGN KEY (branch_id) REFERENCES public.chat(id)
    ON DELETE CASCADE;  

   
--
-- Name: user uniq_branch_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.message
    ADD CONSTRAINT uniq_branch_id UNIQUE (branch_id); 


--
-- Name: contact contact_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.contact
    ADD CONSTRAINT contact_pkey PRIMARY KEY (id);
   
   
--
-- Name: contact unique_user_contact_pair; Type: CONSTRAINT; Schema: public; Owner: postgres
--  
 
ALTER TABLE ONLY public.contact
	ADD CONSTRAINT unique_user_contact_pair UNIQUE (user_id, contact_id);


--
-- Name: contact user_contact_not_equal; Type: CONSTRAINT; Schema: public; Owner: postgres
--  
 
ALTER TABLE ONLY public.contact
	ADD CONSTRAINT user_contact_not_equal CHECK (user_id <> contact_id);


--
-- Name: message message_or_sticker_has_to_be_null; Type: CHECK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE public.message
    ADD CONSTRAINT message_or_sticker_has_to_be_null CHECK ((((message IS NULL) AND (sticker_path IS NOT NULL)) OR ((message IS NOT NULL) AND (sticker_path IS NULL)))) NOT VALID;


--
-- Name: message_payload message_payload_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.message_payload
    ADD CONSTRAINT message_payload_pkey PRIMARY KEY (id);


--
-- Name: message message_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.message
    ADD CONSTRAINT message_pkey PRIMARY KEY (id);


--
-- Name: user uniq_username; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT uniq_username UNIQUE (username);


--
-- Name: user user_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);


--
-- Name: user_role user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_role
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (id);


--
-- Name: message author_id_fk_messages_chat_id_pk_chat_users; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.message
    ADD CONSTRAINT author_id_fk_messages_chat_id_pk_chat_users FOREIGN KEY (author_id, chat_id) REFERENCES public.chat_user(user_id, chat_id)
    ON DELETE CASCADE;


--
-- Name: chat_user chat_id_fk_chat_users_chat_id_pk_chats; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chat_user
    ADD CONSTRAINT chat_id_fk_chat_users_chat_id_pk_chats FOREIGN KEY (chat_id) REFERENCES public.chat(id)
    ON DELETE CASCADE;


--
-- Name: chat chats_fk_chats_type_pk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chat
    ADD CONSTRAINT chats_fk_chats_type_pk FOREIGN KEY (chat_type_id) REFERENCES public.chat_type(id);


--
-- Name: contact contact_id_fk_contacts_user_id_pk_users; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.contact
    ADD CONSTRAINT contact_id_fk_contacts_user_id_pk_users FOREIGN KEY (contact_id) REFERENCES public."user"(id);


--
-- Name: message_payload message_id_fk_message_payload_id_pk_messages; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.message_payload
    ADD CONSTRAINT message_id_fk_message_payload_id_pk_messages FOREIGN KEY (message_id) REFERENCES public.message(id);


--
-- Name: chat_user user_id_fk_chat_users_user_id_pk_users; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chat_user
    ADD CONSTRAINT user_id_fk_chat_users_user_id_pk_users FOREIGN KEY (user_id) REFERENCES public."user"(id);


--
-- Name: contact user_id_fk_contacts_user_id_pk_users; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.contact
    ADD CONSTRAINT user_id_fk_contacts_user_id_pk_users FOREIGN KEY (user_id) REFERENCES public."user"(id);


--
-- Name: chat_user user_role_id_fk_chat_users_chat_id_pk_user_roles; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chat_user
    ADD CONSTRAINT user_role_id_fk_chat_users_chat_id_pk_user_roles FOREIGN KEY (user_role_id) REFERENCES public.user_role(id);

   
CREATE TABLE public.chat_branch
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    chat_id uuid NOT NULL,
    branch_id uuid NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT chat_id_fk_for_branch FOREIGN KEY (chat_id)
        REFERENCES public.chat (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT branch_id_fk_for_branch FOREIGN KEY (branch_id)
        REFERENCES public.chat (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

ALTER TABLE IF EXISTS public.chat_branch
    OWNER to postgres;

COMMENT ON TABLE public.chat_branch
    IS 'Таблица в которой хранятся чаты и их ветки';

--
-- PostgreSQL database dump complete
--

INSERT INTO public.chat_type (value) VALUES
('personal'),
('group'),
('channel'),
('branch');

INSERT INTO  public.user_role ( value) VALUES
('none'),
('owner'),
('admin');

--
-- Insert test data to user
--
INSERT INTO message_type (value) VALUES 
    ('default'),
    ('informational'),
    ('with_payload'),
    ('sticker');

INSERT INTO  payload_type (value) VALUES
('file'),
('photo');

INSERT INTO public."user" (id, username, version, password, name, bio, birthdate, avatar_path) VALUES
    ('39a9aea0-d461-437d-b4eb-bf030a0efc80', 'user11', 0, 'e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648', 'Бал Матье', 'Люблю путешествия 🌍', '1990-05-15T00:00:00Z', '/uploads/avatar/642c5a57-ebc7-49d0-ac2d-f2f1f474bee7.png'),
    ('fa4e08e4-1024-49cb-a799-4aa2a4f3a9df', 'user22', 0, 'e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648', 'Жабка Пепе', 'Кулинар и знаток природы 🍽️🦎', '1992-08-28T00:00:00Z', '/uploads/avatar/d60053d3-e3a9-4a30-b9a3-cdfdc3431fde.png'),
    (gen_random_uuid(), 'user33', 0, 'e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648', 'Dr Peper', 'Люблю газированные напитки 🥤', '1988-12-01T00:00:00Z', NULL),
    (gen_random_uuid(), 'user44', 0, 'e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648', 'Vincent Vega', 'Фанат кино 🎬', '1985-07-14T00:00:00Z', '/uploads/avatar/8027453b-fb36-452d-92dc-c356075fabef.png');


--
-- Insert test data to contacts
--

INSERT INTO contact (id, user_id, contact_id) VALUES 
    ('a0a0aaa0-d461-437d-b4eb-bf030a0efc80', (SELECT id FROM public."user" WHERE username = 'user11'), (SELECT id FROM public."user" WHERE username = 'user22')),
    ('b0a0aaa0-d461-437d-b4eb-bf030a0efc80', (SELECT id FROM public."user" WHERE username = 'user11'), (SELECT id FROM public."user" WHERE username = 'user33')),
    ('c0a0aaa0-d461-437d-b4eb-bf030a0efc80', (SELECT id FROM public."user" WHERE username = 'user11'), (SELECT id FROM public."user" WHERE username = 'user44')),
    ('d0a0aaa0-d461-437d-b4eb-bf030a0efc80', (SELECT id FROM public."user" WHERE username = 'user22'), (SELECT id FROM public."user" WHERE username = 'user11')),
    ('e0a0aaa0-d461-437d-b4eb-bf030a0efc80', (SELECT id FROM public."user" WHERE username = 'user22'), (SELECT id FROM public."user" WHERE username = 'user33')),
    ('f0a0aaa0-d461-437d-b4eb-bf030a0efc80', (SELECT id FROM public."user" WHERE username = 'user33'), (SELECT id FROM public."user" WHERE username = 'user22'));


INSERT INTO chat (chat_name, chat_type_id, id) VALUES
    ('oleg', 1, 'a9a9aea0-d461-437d-b4eb-bf030a0efc80'),
    ('kizaru', 1, 'b9a9aea0-d461-437d-b4eb-bf030a0efc80'),
    ('marsel', 2, 'c9a9aea0-d461-437d-b4eb-bf030a0efc80'),
    ('funny channel', 3, 'd9a9aea0-d461-437d-b4eb-bf030a0efc80'),
    ('not funny channel', 3, 'e9a9aea0-d461-437d-b4eb-bf030a0efc80'),
    ('my little channel', 3, 'f9a9aea0-d461-437d-b4eb-bf030a0efc80');

INSERT INTO chat_user (id, user_role_id, chat_id, user_id) VALUES
    ('a0a0aaa0-d461-437d-b4eb-bf030a0efc80', 2, (SELECT id FROM public.chat WHERE chat_name = 'oleg'), (SELECT id FROM public.user where username ='user11')),
    ('b0a0aaa0-d461-437d-b4eb-bf030a0efc80', 2,(SELECT id FROM public.chat WHERE chat_name = 'oleg'),  (SELECT id FROM public.user where username ='user22')),
    ('c0a0aaa0-d461-437d-b4eb-bf030a0efc80', 2,(SELECT id FROM public.chat WHERE chat_name = 'kizaru'), (SELECT id FROM public.user where username ='user11')),
    ('d0a0aaa0-d461-437d-b4eb-bf030a0efc80', 2,(SELECT id FROM public.chat WHERE chat_name = 'kizaru'), (SELECT id FROM public.user where username ='user44')),
    ('e0a0aaa0-d461-437d-b4eb-bf030a0efc80', 2,(SELECT id FROM public.chat WHERE chat_name = 'marsel'), (SELECT id FROM public.user where username ='user11')),
    ('f0a0aaa0-d461-437d-b4eb-bf030a0efc80', 1,(SELECT id FROM public.chat WHERE chat_name = 'marsel'), (SELECT id FROM public.user where username ='user22')),
    ('f1a0aaa0-d461-437d-b4eb-bf030a0efc80', 1,(SELECT id FROM public.chat WHERE chat_name = 'marsel'), (SELECT id FROM public.user where username ='user33')),
    ('f2a0aaa0-d461-437d-b4eb-bf030a0efc80', 3,(SELECT id FROM public.chat WHERE chat_name = 'marsel'), (SELECT id FROM public.user where username ='user44')),
    ('f4a0aaa0-d461-437d-b4eb-bf030a0efc80', 3,(SELECT id FROM public.chat WHERE chat_name = 'funny channel'), (SELECT id FROM public.user where username ='user11')),
    ('f5a0aaa0-d461-437d-b4eb-bf030a0efc80', 2,(SELECT id FROM public.chat WHERE chat_name = 'funny channel'), (SELECT id FROM public.user where username ='user22')),
    ('f6a0aaa0-d461-437d-b4eb-bf030a0efc80', 1,(SELECT id FROM public.chat WHERE chat_name = 'funny channel'), (SELECT id FROM public.user where username ='user33')),
    ('f7a0aaa0-d461-437d-b4eb-bf030a0efc80', 1,(SELECT id FROM public.chat WHERE chat_name = 'not funny channel'), (SELECT id FROM public.user where username ='user22')),
    ('f8a0aaa0-d461-437d-b4eb-bf030a0efc80', 3,(SELECT id FROM public.chat WHERE chat_name = 'not funny channel'), (SELECT id FROM public.user where username ='user44')),
    ('f9a0aaa0-d461-437d-b4eb-bf030a0efc80', 1,(SELECT id FROM public.chat WHERE chat_name = 'not funny channel'), (SELECT id FROM public.user where username ='user33')),
    ('a1a0aaa0-d461-437d-b4eb-bf030a0efc80', 3,(SELECT id FROM public.chat WHERE chat_name = 'my little channel'), (SELECT id FROM public.user where username ='user44'));

-- /files/675f2ea013dbaf51a93aa2d3
-- /files/675f466313dbaf51a93aa2e4
-- /files/675f391413dbaf51a93aa2db
INSERT INTO sticker (id, sticker_path) VALUES
    ('a0a0aaa0-d461-437d-b4eb-bf030a0efc80', '/files/675f2ea013dbaf51a93aa2d3'),
    ('b0a0aaa0-d461-437d-b4eb-bf030a0efc80', '/files/675f466313dbaf51a93aa2e4'),
    ('c0a0aaa0-d461-437d-b4eb-bf030a0efc80', '/files/675f391413dbaf51a93aa2db'),
    ('d0a0aaa0-d461-437d-b4eb-bf030a0efc80', '/files/6762d25b5803e3d181d0ecc4'),
    ('e0a0aaa0-d461-437d-b4eb-bf030a0efc80', '/files/6762d4535803e3d181d0ecc6'),
    ('f0a0aaa0-d461-437d-b4eb-bf030a0efc80', '/files/6762d4545803e3d181d0ecc7'),
    ('a1a0aaa0-d461-437d-b4eb-bf030a0efc80', '/files/6762d5135803e3d181d0ecc8'),
    ('f1a0aaa0-d461-437d-b4eb-bf030a0efc80', '/files/6762d5505803e3d181d0ecc9'),
    ('f2a0aaa0-d461-437d-b4eb-bf030a0efc80', '/files/6762d7f95803e3d181d0ecca'),
    ('f3a0aaa0-d461-437d-b4eb-bf030a0efc80', '/files/6762d8aa5803e3d181d0eccb'),
    ('f4a0aaa0-d461-437d-b4eb-bf030a0efc80', '/files/6762d8d85803e3d181d0eccc'),
    ('f5a0aaa0-d461-437d-b4eb-bf030a0efc80', '/files/6762d8f45803e3d181d0eccd'),
    ('f6a0aaa0-d461-437d-b4eb-bf030a0efc80', '/files/6762d90e5803e3d181d0ecce'),
    ('f7a0aaa0-d461-437d-b4eb-bf030a0efc80', '/files/6762d9215803e3d181d0eccf');



INSERT INTO sticker_pack (id, photo) VALUES
    ('a0a0aaa0-d461-437d-b4eb-bf030a0efc80', '/files/675f2ea013dbaf51a93aa2d3'),
    ('a1a0aaa0-d461-437d-b4eb-bf030a0efc80', '/files/6762d7f95803e3d181d0ecca');

INSERT INTO sticker_sticker_pack (id, sticker, pack) VALUES
    ('a0a0aaa0-d461-437d-b4eb-bf030a0efc80','a0a0aaa0-d461-437d-b4eb-bf030a0efc80', 'a0a0aaa0-d461-437d-b4eb-bf030a0efc80'),
    ('b0a0aaa0-d461-437d-b4eb-bf030a0efc80','b0a0aaa0-d461-437d-b4eb-bf030a0efc80', 'a0a0aaa0-d461-437d-b4eb-bf030a0efc80'),
    ('c0a0aaa0-d461-437d-b4eb-bf030a0efc80','c0a0aaa0-d461-437d-b4eb-bf030a0efc80', 'a0a0aaa0-d461-437d-b4eb-bf030a0efc80'),
    ('d0a0aaa0-d461-437d-b4eb-bf030a0efc80','d0a0aaa0-d461-437d-b4eb-bf030a0efc80', 'a0a0aaa0-d461-437d-b4eb-bf030a0efc80'),
    ('e0a0aaa0-d461-437d-b4eb-bf030a0efc80','e0a0aaa0-d461-437d-b4eb-bf030a0efc80', 'a0a0aaa0-d461-437d-b4eb-bf030a0efc80'),
    ('f0a0aaa0-d461-437d-b4eb-bf030a0efc80','f0a0aaa0-d461-437d-b4eb-bf030a0efc80', 'a0a0aaa0-d461-437d-b4eb-bf030a0efc80'),
    ('a1a0aaa0-d461-437d-b4eb-bf030a0efc80','a1a0aaa0-d461-437d-b4eb-bf030a0efc80', 'a0a0aaa0-d461-437d-b4eb-bf030a0efc80'),
    ('f1a0aaa0-d461-437d-b4eb-bf030a0efc80','f1a0aaa0-d461-437d-b4eb-bf030a0efc80', 'a1a0aaa0-d461-437d-b4eb-bf030a0efc80'),
    ('f2a0aaa0-d461-437d-b4eb-bf030a0efc80','f2a0aaa0-d461-437d-b4eb-bf030a0efc80', 'a1a0aaa0-d461-437d-b4eb-bf030a0efc80'),
    ('f3a0aaa0-d461-437d-b4eb-bf030a0efc80','f3a0aaa0-d461-437d-b4eb-bf030a0efc80', 'a1a0aaa0-d461-437d-b4eb-bf030a0efc80'),
    ('f4a0aaa0-d461-437d-b4eb-bf030a0efc80','f4a0aaa0-d461-437d-b4eb-bf030a0efc80', 'a1a0aaa0-d461-437d-b4eb-bf030a0efc80'),
    ('f5a0aaa0-d461-437d-b4eb-bf030a0efc80','f5a0aaa0-d461-437d-b4eb-bf030a0efc80', 'a1a0aaa0-d461-437d-b4eb-bf030a0efc80'),
    ('f6a0aaa0-d461-437d-b4eb-bf030a0efc80','f6a0aaa0-d461-437d-b4eb-bf030a0efc80', 'a1a0aaa0-d461-437d-b4eb-bf030a0efc80'),
    ('f7a0aaa0-d461-437d-b4eb-bf030a0efc80','f7a0aaa0-d461-437d-b4eb-bf030a0efc80', 'a1a0aaa0-d461-437d-b4eb-bf030a0efc80');
