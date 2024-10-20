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
    user_id uuid NOT NULL
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

--
-- Name: message; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.message (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    chat_id uuid NOT NULL,
    author_id uuid NOT NULL,
    message text,
    sent_at timestamp with time zone NOT NULL,
    is_redacted boolean DEFAULT false NOT NULL,
    sticker_path text
);


ALTER TABLE public.message OWNER TO postgres;

--
-- Name: message_payload; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.message_payload (
    payload_path text NOT NULL,
    id uuid NOT NULL,
    message_id uuid NOT NULL
);


ALTER TABLE public.message_payload OWNER TO postgres;

--
-- Name: user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."user" (
    username text NOT NULL,
    version integer NOT NULL,
    password text NOT NULL,
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

--
-- Name: profile; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."profile" (
	user_id uuid NOT NULL,
    name text NOT NULL,
    bio text,
    birthdate timestamp with time zone,
    avatar_path text,
    id uuid NOT NULL
);


ALTER TABLE public."profile" OWNER TO postgres;

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
-- Name: contact contact_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.contact
    ADD CONSTRAINT contact_pkey PRIMARY KEY (id);


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
-- Name: profile profile_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."profile"
    ADD CONSTRAINT profile_pkey PRIMARY KEY (id);


--
-- Name: user_role user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_role
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (id);


--
-- Name: message author_id_fk_messages_chat_id_pk_chat_users; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.message
    ADD CONSTRAINT author_id_fk_messages_chat_id_pk_chat_users FOREIGN KEY (author_id, chat_id) REFERENCES public.chat_user(user_id, chat_id);


--
-- Name: chat_user chat_id_fk_chat_users_chat_id_pk_chats; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chat_user
    ADD CONSTRAINT chat_id_fk_chat_users_chat_id_pk_chats FOREIGN KEY (chat_id) REFERENCES public.chat(id);


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


--
-- Name: profile profile_fk_user_pk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.profile
    ADD CONSTRAINT profile_fk_user_pk FOREIGN KEY (user_id) REFERENCES public.user(id);
   
--
-- PostgreSQL database dump complete
--


INSERT INTO public.chat_type (value) VALUES
('personal'),
('group'),
('channel');

INSERT INTO  public.user_role ( value) VALUES
('none'),
('owner'),
('admin');


--
-- Insert test data to user
--

INSERT INTO public."user" (id, username, version, password) VALUES
    (gen_random_uuid(), 'user11', 0, 'e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648'),
    (gen_random_uuid(), 'user22', 0, 'e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648'),
    (gen_random_uuid(), 'user33', 0, 'e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648'),
    (gen_random_uuid(), 'user44', 0, 'e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648');

--
-- Insert test data to profile
--

INSERT INTO public."profile" (id, user_id, name) VALUES
    (gen_random_uuid(), (SELECT id FROM public."user" WHERE username = 'user11'), 'Бал Матье'),
    (gen_random_uuid(), (SELECT id FROM public."user" WHERE username = 'user22'), 'Жабка Пепе'),
    (gen_random_uuid(), (SELECT id FROM public."user" WHERE username = 'user33'), 'Dr Peper'),
    (gen_random_uuid(), (SELECT id FROM public."user" WHERE username = 'user44'), 'Vincent Vega');