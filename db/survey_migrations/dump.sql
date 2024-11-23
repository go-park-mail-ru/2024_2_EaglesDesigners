DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

CREATE TABLE public.answer (
    id uuid NOT NULL,
    question_id uuid NOT NULL,
    text_answer text,
    numeric_answer integer,
    user_id uuid NOT NULL
);


ALTER TABLE public.answer OWNER TO postgres;

--
-- Name: question; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.question (
    id uuid NOT NULL,
    servey_id uuid NOT NULL,
    type_id integer NOT NULL,
    question_text text NOT NULL
);


ALTER TABLE public.question OWNER TO postgres;

--
-- Name: question_type; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.question_type (
    id integer NOT NULL,
    value text NOT NULL
);


ALTER TABLE public.question_type OWNER TO postgres;

--
-- Name: question_type_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.question_type ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.question_type_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: servey; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.servey (
    id uuid NOT NULL,
    topic text NOT NULL
);


ALTER TABLE public.servey OWNER TO postgres;

--
-- Data for Name: answer; Type: TABLE DATA; Schema: public; Owner: postgres
--


--
-- Data for Name: question; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.question (id, servey_id, type_id, question_text) VALUES
('1ca47e50-1640-4a44-9ac5-e34f2e12f6ab',	'c6c92ef5-65ff-4835-9e24-705abf8e00a8',	2,	'Нравятся ли вам горки в аквапарке?'),
('e37421bd-6621-429e-ab61-ec0656e426e3',	'c6c92ef5-65ff-4835-9e24-705abf8e00a8',	2,	'Нравятся ли вам бар в аквапарке?'),
('0c94656f-f4b7-4466-97fd-f2b981da5e52',	'02277257-bc5e-4264-b602-3891169a4ccb',	2,	'Понятен ли вам изучаемый материал?'),
('4cb4f0b2-e602-4272-ba41-bcdaec99bce2',	'02277257-bc5e-4264-b602-3891169a4ccb',	2,	'Нравится ли вам обучение?'),
('db04f9a5-ae80-43b7-8c75-5a4db98344c5',	'0bee6d65-15c2-473a-860e-acc370eeec40',	2,	'Есть ли у вас машина?'),
('d7c06b4c-13bd-4575-ac2c-83514bce13b1',	'0bee6d65-15c2-473a-860e-acc370eeec40',	2,	'В какое время нет места на парковке?'),
('dfdf1ff1-48d0-45a2-a72a-84f3c73350ec',	'c6c92ef5-65ff-4835-9e24-705abf8e00a8',	1,	'Насколько крут аквапарк?'),
('e5eb1921-fc4f-44c2-87ad-c7ef5e79ad30',	'02277257-bc5e-4264-b602-3891169a4ccb',	1,	'Насколько понятен материал?');



--
-- Data for Name: question_type; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO  public.question_type (value) VALUES
('numeric'),
('text');



--
-- Data for Name: servey; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO  public.servey (id, topic) VALUES
('c6c92ef5-65ff-4835-9e24-705abf8e00a8',	'В аквапарке реально охрененно?'),
('02277257-bc5e-4264-b602-3891169a4ccb',	'Довольны ли вы учебным процессом'),
('0bee6d65-15c2-473a-860e-acc370eeec40',	'Парковка у дгту виноделие');



--
-- Name: question_type_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

--
-- Name: question Question_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.question
    ADD CONSTRAINT "Question_pkey" PRIMARY KEY (id);


--
-- Name: servey Servey_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.servey
    ADD CONSTRAINT "Servey_pkey" PRIMARY KEY (id);


--
-- Name: question_type question_type_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.question_type
    ADD CONSTRAINT question_type_pkey PRIMARY KEY (id);


--
-- Name: answer question_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.answer
    ADD CONSTRAINT question_id_fk FOREIGN KEY (question_id) REFERENCES public.question(id);


--
-- Name: question question_servey_key_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.question
    ADD CONSTRAINT question_servey_key_fk FOREIGN KEY (servey_id) REFERENCES public.servey(id);


--
-- Name: question question_type_key_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.question
    ADD CONSTRAINT question_type_key_fk FOREIGN KEY (type_id) REFERENCES public.question_type(id);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

