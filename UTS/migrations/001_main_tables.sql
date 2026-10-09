--
-- PostgreSQL database dump
--

\restrict blcQzR6rzCSGwr10LiNZBeKAWSfFhapJI7rUUf5yrxnSpuZstxuwC97RnvPQaRl

-- Dumped from database version 18.2
-- Dumped by pg_dump version 18.2

-- Started on 2026-10-09 21:28:19

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 4 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: pg_database_owner
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO pg_database_owner;

--
-- TOC entry 5084 (class 0 OID 0)
-- Dependencies: 4
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: pg_database_owner
--

COMMENT ON SCHEMA public IS 'standard public schema';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 224 (class 1259 OID 17253)
-- Name: courses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.courses (
    id integer CONSTRAINT course_id_not_null NOT NULL,
    kode_mk character varying(10) CONSTRAINT course_kode_mk_not_null NOT NULL,
    nama_mk character varying(100) CONSTRAINT course_nama_mk_not_null NOT NULL,
    sks integer CONSTRAINT course_sks_not_null NOT NULL,
    semester integer CONSTRAINT course_semester_not_null NOT NULL,
    kuota integer CONSTRAINT course_kuota_not_null NOT NULL
);


ALTER TABLE public.courses OWNER TO postgres;

--
-- TOC entry 223 (class 1259 OID 17252)
-- Name: course_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.course_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.course_id_seq OWNER TO postgres;

--
-- TOC entry 5085 (class 0 OID 0)
-- Dependencies: 223
-- Name: course_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.course_id_seq OWNED BY public.courses.id;


--
-- TOC entry 226 (class 1259 OID 17267)
-- Name: enrollments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.enrollments (
    id integer NOT NULL,
    student_id integer NOT NULL,
    course_id integer NOT NULL,
    tahun_akademik character varying(20),
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.enrollments OWNER TO postgres;

--
-- TOC entry 225 (class 1259 OID 17266)
-- Name: enrollments_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.enrollments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.enrollments_id_seq OWNER TO postgres;

--
-- TOC entry 5086 (class 0 OID 0)
-- Dependencies: 225
-- Name: enrollments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.enrollments_id_seq OWNED BY public.enrollments.id;


--
-- TOC entry 227 (class 1259 OID 17289)
-- Name: permissions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.permissions (
    name character varying(50) NOT NULL,
    description character varying(150) NOT NULL
);


ALTER TABLE public.permissions OWNER TO postgres;

--
-- TOC entry 230 (class 1259 OID 17311)
-- Name: refresh_tokens; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.refresh_tokens (
    id bigint NOT NULL,
    user_id integer NOT NULL,
    token_hash text NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    revoked_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.refresh_tokens OWNER TO postgres;

--
-- TOC entry 229 (class 1259 OID 17310)
-- Name: refresh_tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.refresh_tokens_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.refresh_tokens_id_seq OWNER TO postgres;

--
-- TOC entry 5087 (class 0 OID 0)
-- Dependencies: 229
-- Name: refresh_tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.refresh_tokens_id_seq OWNED BY public.refresh_tokens.id;


--
-- TOC entry 228 (class 1259 OID 17297)
-- Name: role_permissions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.role_permissions (
    role_name character varying(20) NOT NULL,
    permissions_name character varying(50) NOT NULL
);


ALTER TABLE public.role_permissions OWNER TO postgres;

--
-- TOC entry 222 (class 1259 OID 17234)
-- Name: students; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.students (
    id integer NOT NULL,
    user_id integer NOT NULL,
    nim character varying(50) NOT NULL,
    nama character varying(100) NOT NULL,
    prodi character varying(30) NOT NULL,
    angkatan character varying(4) NOT NULL,
    ipk_terakhir double precision,
    deleted_at timestamp with time zone
);


ALTER TABLE public.students OWNER TO postgres;

--
-- TOC entry 221 (class 1259 OID 17233)
-- Name: students_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.students_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.students_id_seq OWNER TO postgres;

--
-- TOC entry 5088 (class 0 OID 0)
-- Dependencies: 221
-- Name: students_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.students_id_seq OWNED BY public.students.id;


--
-- TOC entry 220 (class 1259 OID 17220)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    role character varying(10) NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 17219)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- TOC entry 5089 (class 0 OID 0)
-- Dependencies: 219
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 4886 (class 2604 OID 17256)
-- Name: courses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.courses ALTER COLUMN id SET DEFAULT nextval('public.course_id_seq'::regclass);


--
-- TOC entry 4887 (class 2604 OID 17270)
-- Name: enrollments id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.enrollments ALTER COLUMN id SET DEFAULT nextval('public.enrollments_id_seq'::regclass);


--
-- TOC entry 4889 (class 2604 OID 17314)
-- Name: refresh_tokens id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.refresh_tokens ALTER COLUMN id SET DEFAULT nextval('public.refresh_tokens_id_seq'::regclass);


--
-- TOC entry 4885 (class 2604 OID 17237)
-- Name: students id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.students ALTER COLUMN id SET DEFAULT nextval('public.students_id_seq'::regclass);


--
-- TOC entry 4884 (class 2604 OID 17223)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 5072 (class 0 OID 17253)
-- Dependencies: 224
-- Data for Name: courses; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.courses (id, kode_mk, nama_mk, sks, semester, kuota) VALUES (1, 'TI-001', 'Algoritma dan Pemrograman Teori', 3, 1, 60);
INSERT INTO public.courses (id, kode_mk, nama_mk, sks, semester, kuota) VALUES (2, 'TI-002', 'Algoritma dan Pemrograman Praktikum', 2, 1, 40);
INSERT INTO public.courses (id, kode_mk, nama_mk, sks, semester, kuota) VALUES (3, 'TI-003', 'Struktur Data dan Algoritma', 3, 2, 60);
INSERT INTO public.courses (id, kode_mk, nama_mk, sks, semester, kuota) VALUES (4, 'TI-004', 'Pemrograman Basis Data Teori', 3, 3, 60);
INSERT INTO public.courses (id, kode_mk, nama_mk, sks, semester, kuota) VALUES (5, 'TI-005', 'Pemrograman Basis Data Praktikum', 2, 3, 40);
INSERT INTO public.courses (id, kode_mk, nama_mk, sks, semester, kuota) VALUES (6, 'TI-006', 'Pemrograman Web Lanjut', 3, 4, 50);
INSERT INTO public.courses (id, kode_mk, nama_mk, sks, semester, kuota) VALUES (7, 'TI-007', 'Pemrograman Backend Lanjut Teori', 3, 4, 60);
INSERT INTO public.courses (id, kode_mk, nama_mk, sks, semester, kuota) VALUES (8, 'TI-008', 'Pemrograman Backend Lanjut Praktikum', 2, 4, 40);
INSERT INTO public.courses (id, kode_mk, nama_mk, sks, semester, kuota) VALUES (9, 'TI-009', 'Rekayasa Perangkat Lunak', 3, 5, 60);
INSERT INTO public.courses (id, kode_mk, nama_mk, sks, semester, kuota) VALUES (11, 'TI-011', 'Kecerdasan Buatan Teori', 3, 6, 60);
INSERT INTO public.courses (id, kode_mk, nama_mk, sks, semester, kuota) VALUES (12, 'TI-012', 'Kecerdasan Buatan Praktikum', 2, 6, 40);
INSERT INTO public.courses (id, kode_mk, nama_mk, sks, semester, kuota) VALUES (13, 'TI-013', 'Keamanan Jaringan Komputer', 3, 7, 50);
INSERT INTO public.courses (id, kode_mk, nama_mk, sks, semester, kuota) VALUES (14, 'TI-014', 'Manajemen Proyek TI', 2, 7, 60);
INSERT INTO public.courses (id, kode_mk, nama_mk, sks, semester, kuota) VALUES (15, 'TI-015', 'Etika Profesi dan Bisnis', 2, 8, 60);
INSERT INTO public.courses (id, kode_mk, nama_mk, sks, semester, kuota) VALUES (10, 'TI-010', 'Desain UI/UX', 2, 5, 1);


--
-- TOC entry 5074 (class 0 OID 17267)
-- Dependencies: 226
-- Data for Name: enrollments; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.enrollments (id, student_id, course_id, tahun_akademik, created_at) VALUES (14, 14, 7, '2026/2027-Ganjil', '2026-10-09 21:00:40.834277+07');


--
-- TOC entry 5075 (class 0 OID 17289)
-- Dependencies: 227
-- Data for Name: permissions; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.permissions (name, description) VALUES ('student:list', 'Melihat daftar seluruh mahasiswa dengan pagination dan filter');
INSERT INTO public.permissions (name, description) VALUES ('student:create', 'Menambah mahasiswa baru sekaligus akun usernya');
INSERT INTO public.permissions (name, description) VALUES ('student:read:any', 'Melihat detail data mahasiswa mana pun');
INSERT INTO public.permissions (name, description) VALUES ('student:read:own', 'Melihat detail data mahasiswa milik sendiri');
INSERT INTO public.permissions (name, description) VALUES ('student:update:any', 'Memperbarui data mahasiswa mana pun');
INSERT INTO public.permissions (name, description) VALUES ('student:delete:any', 'Menghapus (soft delete) data mahasiswa mana pun');
INSERT INTO public.permissions (name, description) VALUES ('course:list', 'Melihat daftar mata kuliah beserta sisa kuota');
INSERT INTO public.permissions (name, description) VALUES ('enrollment:create', 'Mengambil mata kuliah (menambah KRS)');
INSERT INTO public.permissions (name, description) VALUES ('enrollment:delete:own', 'Membatalkan mata kuliah dari KRS milik sendiri');


--
-- TOC entry 5078 (class 0 OID 17311)
-- Dependencies: 230
-- Data for Name: refresh_tokens; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (1, 8, '035b8621c9cc1bdf5b4c72fbc3fb188cb6304848a36fa56d3a42f9cdc56c6323', '2026-10-15 18:10:27.180277+07', NULL, '2026-10-08 18:10:27.182127+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (2, 8, '7cd8ed93219591511dfed8376ccbd0a089836e4794ad70c1d398a0df82e388ea', '2026-10-15 18:12:33.911498+07', NULL, '2026-10-08 18:12:33.912419+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (3, 8, '6d08dadfabbff63a8fb7441fa7244a01ad3ff4ddfae3ab8e06d177ac062a5803', '2026-10-15 18:42:27.570574+07', NULL, '2026-10-08 18:42:27.571541+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (5, 8, '148c8fe64be90554919a640db01d4347c7d849b12bf3dd10f6f56b8d9c65a395', '2026-10-15 18:59:45.601227+07', '2026-10-08 19:20:18.909398+07', '2026-10-08 18:59:45.602408+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (6, 8, '79589a748a9799bab164165607c48f0096b4b503edf2b852d820c14846a4e40a', '2026-10-15 19:20:18.909252+07', NULL, '2026-10-08 19:20:18.909943+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (8, 8, '12722536e593789085badc0bad65aad717625b636b34b9a5cb919e040bad0c67', '2026-10-15 20:25:24.122679+07', NULL, '2026-10-08 20:25:24.123382+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (9, 8, '57380601acef97255798c9f456dfd5dbd746e8c3c0fcb68bf54918dace22dd10', '2026-10-15 21:55:49.112605+07', NULL, '2026-10-08 21:55:49.113522+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (11, 8, 'bb8d4e06c3aec114aa596fc0f54942838d73afb41daa708d77ee0077033361d0', '2026-10-16 06:37:49.627428+07', NULL, '2026-10-09 06:37:49.628544+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (13, 8, '5ea058a105962eb48f95abff10d1ec29a9511757372fbdb41af84fdf68c5e053', '2026-10-16 07:27:39.281407+07', '2026-10-09 07:43:11.059098+07', '2026-10-09 07:27:39.283086+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (14, 8, 'eebb3a54de183abdcec5f5c47cf8335486b1d47b9f2c9bf2f2b4f70a51d472d2', '2026-10-16 07:43:11.065734+07', '2026-10-09 08:07:09.632036+07', '2026-10-09 07:43:11.066119+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (15, 8, '5b6070a5f6502b715163ed7de19fbca06007fbd14f7bd6176b47fb146e875d6b', '2026-10-16 08:07:09.638567+07', '2026-10-09 08:30:01.185108+07', '2026-10-09 08:07:09.639002+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (16, 8, '82a8ec4639d21a65f5b23522ec1756fd5d4d948131c8812efc6a1f20ce996a8f', '2026-10-16 08:30:01.189948+07', '2026-10-09 08:51:39.432578+07', '2026-10-09 08:30:01.190517+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (17, 8, 'bf46c4cf0a38655c343796e05750a718152a2790ca5dffbf9702ac740a5bb076', '2026-10-16 08:51:39.432503+07', '2026-10-09 13:52:44.528374+07', '2026-10-09 08:51:39.433032+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (18, 8, 'db4f4f0273bec29e94f0d7fd04b7a9f059683cf6d217b5f5b67cbbc1ce945dc8', '2026-10-16 13:52:44.529474+07', NULL, '2026-10-09 13:52:44.53005+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (20, 17, '245fadec0c9dbb736dcd6ab58d38cf964250fc23f404053605eed4be8a26db6d', '2026-10-16 14:48:58.712744+07', NULL, '2026-10-09 14:48:58.714964+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (19, 20, 'd89627f6510a71652d2e5fa95b50a79f85ddb9c829c447e784da88b559b7d017', '2026-10-16 14:04:50.971313+07', '2026-10-09 15:09:28.195217+07', '2026-10-09 14:04:50.972696+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (21, 20, '19b9ed1e48c8ef0449eca4e42db77637fdc553bc82d56be82b13eb4649248932', '2026-10-16 15:09:28.196304+07', '2026-10-09 15:21:09.476339+07', '2026-10-09 15:09:28.196625+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (23, 8, '9b5a649f0b9402e6d9d6aa76860ce51a83cea08c479c355f5fa1ff5948fcda77', '2026-10-16 20:17:54.47024+07', NULL, '2026-10-09 20:17:54.471479+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (22, 20, '0c4aa6b6988d7bf4ea0534519d54074132f94a7f6928921b46e3cff5ee98293f', '2026-10-16 15:21:09.477082+07', '2026-10-09 20:22:41.907154+07', '2026-10-09 15:21:09.478462+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (24, 20, '44b8f404c18d186a579a87eedc7422fb751e45fba1db4a7140392387280fabe4', '2026-10-16 20:22:41.909781+07', NULL, '2026-10-09 20:22:41.910189+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (25, 8, 'cdc30292f83115d83fbbeecd3adfa939afed9fa3719d693431611794ec0f8aa5', '2026-10-16 20:23:43.206585+07', NULL, '2026-10-09 20:23:43.207537+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (26, 20, '523f9eb87429cf12acfa44c2da0386d998389266305ba97056f4cb846ce6ed74', '2026-10-16 20:38:44.367348+07', NULL, '2026-10-09 20:38:44.367773+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (27, 14, '07b7ca8f319375176053a23ef34d4147090a3618050aeb5be5f3f60c634cf762', '2026-10-16 20:45:55.877363+07', NULL, '2026-10-09 20:45:55.878263+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (28, 8, 'feb56607785511cb74f0cb237852c8a710c9564eb3b1f4ded6edab1a947ff2c4', '2026-10-16 20:47:46.331785+07', NULL, '2026-10-09 20:47:46.332162+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (29, 20, '40cc7a1faa329ea29ad170c7286cc915611e504440c895afb343e61a176e8bb9', '2026-10-16 20:59:00.668176+07', NULL, '2026-10-09 20:59:00.668601+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (30, 17, '72246b4f2c2d67e43b14e10ebf5f48ee47c22a395dba658b7328656e67138451', '2026-10-16 20:59:38.637494+07', NULL, '2026-10-09 20:59:38.638171+07');
INSERT INTO public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) VALUES (31, 8, 'f705820fbd33528111584dea8e942f1b4ae5ed74a91a851f829251ab5b55f642', '2026-10-16 21:13:54.605501+07', NULL, '2026-10-09 21:13:54.609281+07');


--
-- TOC entry 5076 (class 0 OID 17297)
-- Dependencies: 228
-- Data for Name: role_permissions; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.role_permissions (role_name, permissions_name) VALUES ('admin', 'student:list');
INSERT INTO public.role_permissions (role_name, permissions_name) VALUES ('admin', 'student:create');
INSERT INTO public.role_permissions (role_name, permissions_name) VALUES ('admin', 'student:read:any');
INSERT INTO public.role_permissions (role_name, permissions_name) VALUES ('admin', 'student:update:any');
INSERT INTO public.role_permissions (role_name, permissions_name) VALUES ('admin', 'student:delete:any');
INSERT INTO public.role_permissions (role_name, permissions_name) VALUES ('admin', 'course:list');
INSERT INTO public.role_permissions (role_name, permissions_name) VALUES ('mahasiswa', 'student:read:own');
INSERT INTO public.role_permissions (role_name, permissions_name) VALUES ('mahasiswa', 'course:list');
INSERT INTO public.role_permissions (role_name, permissions_name) VALUES ('mahasiswa', 'enrollment:create');
INSERT INTO public.role_permissions (role_name, permissions_name) VALUES ('mahasiswa', 'enrollment:delete:own');


--
-- TOC entry 5070 (class 0 OID 17234)
-- Dependencies: 222
-- Data for Name: students; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (2, 8, '00000', 'Admin', 'Unknown', '0000', NULL, NULL);
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (8, 14, '202310010002', 'Citra Lestari', 'Teknik Informatika', '2023', NULL, NULL);
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (9, 15, '202311010001', 'Dimas Anggara', 'Ilmu Komputer', '2023', NULL, NULL);
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (10, 16, '202311010002', 'Eka Ramadhani', 'Ilmu Komputer', '2023', NULL, NULL);
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (11, 17, '202312010001', 'Fajar Nugroho', 'Teknik Elektro', '2023', NULL, NULL);
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (14, 20, '202312010002', 'Gita Savitri', 'Teknik Elektro', '2023', NULL, NULL);
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (16, 22, '202511010001', 'Qori Maharani', 'Ilmu Komputer', '2025', NULL, NULL);
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (17, 24, '202211010001', 'Miswa Dama', 'Ilmu Komputer', '2022', NULL, '2026-10-09 20:36:53.469603+07');
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (7, 13, '202310010001', 'Bima Saputra', 'Teknik Informatika', '2023', 2, NULL);
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (18, 25, '202411010001', 'Kiki Amalia', 'Ilmu Komputer', '2024', NULL, NULL);
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (19, 26, '202610010001', 'Surya Paloh', 'Teknik Informatika', '2026', NULL, NULL);
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (20, 28, '202410010002', 'Indah Permata', 'Teknik Informatika', '2024', NULL, NULL);
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (21, 29, '202410010003', 'Joko Widodo', 'Teknik Informatika', '2024', NULL, NULL);
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (22, 30, '202411010002', 'Lukman Hakim', 'Ilmu Komputer', '2024', NULL, NULL);
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (23, 31, '202510010001', 'Oscar Darmawan', 'Teknik Informatika', '2025', NULL, NULL);
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (24, 32, '202510010002', 'Putri Tanjung', 'Teknik Informatika', '2025', NULL, NULL);
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (25, 33, '202412010001', 'Mira Lesmana', 'Teknik Elektro', '2024', NULL, NULL);
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (26, 34, '202412010002', 'Nana Mirdad', 'Teknik Elektro', '2024', NULL, NULL);
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (27, 35, '202512010001', 'Rizky Febian', 'Teknik Elektro', '2025', NULL, NULL);
INSERT INTO public.students (id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at) VALUES (28, 36, '202611010001', 'Tari Eka', 'Ilmu Komputer', '2026', NULL, NULL);


--
-- TOC entry 5068 (class 0 OID 17220)
-- Dependencies: 220
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.users (id, email, password, role) VALUES (8, 'admin@mail.com', '$2a$10$v2n.92px/LhnNAcsoVkBguxMYm6C0qw.hap.HqOgk3uCTReDXcc9K', 'admin');
INSERT INTO public.users (id, email, password, role) VALUES (13, 'bima.s@student.kampus.ac.id', '$2a$10$62oMRnjBHnrCR5bOpOhZjubJ9aTF1Mln4w3Yuxs0MXPC2GfZVn3B6', 'mahasiswa');
INSERT INTO public.users (id, email, password, role) VALUES (14, 'citra.l@student.kampus.ac.id', '$2a$10$zQpTij3JJp2.W7kBfK25pOH/qI/M/ehwGVkTUcxrpvQnycYEL.gxa', 'mahasiswa');
INSERT INTO public.users (id, email, password, role) VALUES (15, 'dimas.a@student.kampus.ac.id', '$2a$10$kYNk0FD82.6KfrG01.RBOOxSXIRN0S8b4YE./GPt2uoRdRX9hubGq', 'mahasiswa');
INSERT INTO public.users (id, email, password, role) VALUES (16, 'eka.r@student.kampus.ac.id', '$2a$10$gZzqRQZNmI8iiOLGyeDJMeXHH.vRmolM7fTyaOe3rx4Y9uMyW6FDy', 'mahasiswa');
INSERT INTO public.users (id, email, password, role) VALUES (17, 'fajar.n@student.kampus.ac.id', '$2a$10$VcqxUuNx2IfZBahJAxDB6Od2XXuqLQjOQBVpQBz2wJT1E/vip3RPa', 'mahasiswa');
INSERT INTO public.users (id, email, password, role) VALUES (20, 'gita.s@student.kampus.ac.id', '$2a$10$DouCLrD.CK5CkAXjGRvud.po1EccAwcV8Co1o/41jvvtLD1AIneBC', 'mahasiswa');
INSERT INTO public.users (id, email, password, role) VALUES (22, 'qori.m@student.kampus.ac.id', '$2a$10$MH7sJh7xNf/Kk1q14Zw4.OOWJj6EPyn7a/OjZ3IopG/8tZCOa5/mm', 'mahasiswa');
INSERT INTO public.users (id, email, password, role) VALUES (24, 'miswa.d@student.kampus.ac.id', '$2a$10$88pwAqSHHX8CL.UZrL7br.IkUJElqB6pe0IK9CnOQoR3AVx8nqPGm', 'mahasiswa');
INSERT INTO public.users (id, email, password, role) VALUES (25, 'kiki.a@student.kampus.ac.id', '$2a$10$XzQQlyUdN.09WPFV62nuxuUOhQoYInsDK4vU8VOGXWYGLnP7fhfIW', 'mahasiswa');
INSERT INTO public.users (id, email, password, role) VALUES (26, 'surya.p@student.kampus.ac.id', '$2a$10$2GAo5NxhKWph3yMunobM2ekV1Ey.rzm2nvWXFzAjSWagVCwFF1qdG', 'mahasiswa');
INSERT INTO public.users (id, email, password, role) VALUES (28, 'indah.p@student.kampus.ac.id', '$2a$10$L817pxbeHDKuelYqQK2aveJCCW3IcRDPo0NTDVU3gkI/NSeo0sTDS', 'mahasiswa');
INSERT INTO public.users (id, email, password, role) VALUES (29, 'joko.w@student.kampus.ac.id', '$2a$10$WkrlEUZuFA21fbicOX/Ly.x8dxkMWM9zp5BTmOOb/Ejxl8NIba2ja', 'mahasiswa');
INSERT INTO public.users (id, email, password, role) VALUES (30, 'lukman.h@student.kampus.ac.id', '$2a$10$Wd2etJUulieC/s1ubAo4nOf7dia0bD7Bdo160E.MWwqJld0hWEcme', 'mahasiswa');
INSERT INTO public.users (id, email, password, role) VALUES (31, 'oscar.d@student.kampus.ac.id', '$2a$10$Gem7q/gorI8mAfjJWlwItOH05bGAmwKmnrhfnLihPF2Jk0CTj3W/q', 'mahasiswa');
INSERT INTO public.users (id, email, password, role) VALUES (32, 'putri.t@student.kampus.ac.id', '$2a$10$yU0DbcEMPLCGVNqr3Z0qF.tYMGr8j7kqJw8UDw3wJanmKS6AgFVZu', 'mahasiswa');
INSERT INTO public.users (id, email, password, role) VALUES (33, 'mira.l@student.kampus.ac.id', '$2a$10$gvUyy5FIL4L3m8bwEb93/uYoov9r/FuY6avJOrSXzE2hGdLzDgUqe', 'mahasiswa');
INSERT INTO public.users (id, email, password, role) VALUES (34, 'nana.m@student.kampus.ac.id', '$2a$10$zbmdcPovDQHY.NI3mtaTb.GkGUn8kX.atg/rM2pES0tqLzoXEki82', 'mahasiswa');
INSERT INTO public.users (id, email, password, role) VALUES (35, 'rizky.f@student.kampus.ac.id', '$2a$10$FvCjBonuHWlif2pgKeHu0..vlgD4/Ai6EpQTrrOL0HHRhPsGGGQom', 'mahasiswa');
INSERT INTO public.users (id, email, password, role) VALUES (36, 'tari.e@student.kampus.ac.id', '$2a$10$agBvNXleM4Yn0IMOoDLNzODsiPFi8UffncFzQgFV0GxneBCRM6/uO', 'mahasiswa');


--
-- TOC entry 5090 (class 0 OID 0)
-- Dependencies: 223
-- Name: course_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.course_id_seq', 15, true);


--
-- TOC entry 5091 (class 0 OID 0)
-- Dependencies: 225
-- Name: enrollments_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.enrollments_id_seq', 14, true);


--
-- TOC entry 5092 (class 0 OID 0)
-- Dependencies: 229
-- Name: refresh_tokens_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.refresh_tokens_id_seq', 31, true);


--
-- TOC entry 5093 (class 0 OID 0)
-- Dependencies: 221
-- Name: students_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.students_id_seq', 28, true);


--
-- TOC entry 5094 (class 0 OID 0)
-- Dependencies: 219
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 36, true);


--
-- TOC entry 4900 (class 2606 OID 17264)
-- Name: courses course_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.courses
    ADD CONSTRAINT course_pkey PRIMARY KEY (id);


--
-- TOC entry 4902 (class 2606 OID 17277)
-- Name: enrollments enrollments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.enrollments
    ADD CONSTRAINT enrollments_pkey PRIMARY KEY (id);


--
-- TOC entry 4906 (class 2606 OID 17295)
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (name);


--
-- TOC entry 4911 (class 2606 OID 17324)
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);


--
-- TOC entry 4913 (class 2606 OID 17326)
-- Name: refresh_tokens refresh_tokens_token_hash_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_token_hash_key UNIQUE (token_hash);


--
-- TOC entry 4908 (class 2606 OID 17303)
-- Name: role_permissions role_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (role_name, permissions_name);


--
-- TOC entry 4897 (class 2606 OID 17245)
-- Name: students students_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.students
    ADD CONSTRAINT students_pkey PRIMARY KEY (id);


--
-- TOC entry 4894 (class 2606 OID 17231)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 4898 (class 1259 OID 17265)
-- Name: course_kode_mk; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX course_kode_mk ON public.courses USING btree (kode_mk);


--
-- TOC entry 4903 (class 1259 OID 25411)
-- Name: enrollments_students_courses_year; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX enrollments_students_courses_year ON public.enrollments USING btree (student_id, course_id, tahun_akademik);


--
-- TOC entry 4904 (class 1259 OID 17296)
-- Name: permissions_name_lower; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX permissions_name_lower ON public.permissions USING btree (lower((name)::text));


--
-- TOC entry 4914 (class 1259 OID 17332)
-- Name: refresh_tokens_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX refresh_tokens_user_id_idx ON public.refresh_tokens USING btree (user_id);


--
-- TOC entry 4909 (class 1259 OID 17309)
-- Name: role_permissions_role_lower; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX role_permissions_role_lower ON public.role_permissions USING btree (lower((role_name)::text));


--
-- TOC entry 4895 (class 1259 OID 17246)
-- Name: students_nim; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX students_nim ON public.students USING btree (nim);


--
-- TOC entry 4891 (class 1259 OID 25410)
-- Name: users_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX users_email ON public.users USING btree (email);


--
-- TOC entry 4892 (class 1259 OID 17232)
-- Name: users_email_lower_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX users_email_lower_idx ON public.users USING btree (lower((email)::text));


--
-- TOC entry 4916 (class 2606 OID 17283)
-- Name: enrollments enrollments_fkey_courses; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.enrollments
    ADD CONSTRAINT enrollments_fkey_courses FOREIGN KEY (course_id) REFERENCES public.courses(id);


--
-- TOC entry 4917 (class 2606 OID 17278)
-- Name: enrollments enrollments_fkey_students; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.enrollments
    ADD CONSTRAINT enrollments_fkey_students FOREIGN KEY (student_id) REFERENCES public.students(id);


--
-- TOC entry 4919 (class 2606 OID 17327)
-- Name: refresh_tokens refresh_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 4918 (class 2606 OID 17304)
-- Name: role_permissions role_permissions_permissions_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_permissions_name_fkey FOREIGN KEY (permissions_name) REFERENCES public.permissions(name) ON DELETE CASCADE;


--
-- TOC entry 4915 (class 2606 OID 17247)
-- Name: students students_fkey_users; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.students
    ADD CONSTRAINT students_fkey_users FOREIGN KEY (user_id) REFERENCES public.users(id);


-- Completed on 2026-10-09 21:28:20

--
-- PostgreSQL database dump complete
--

\unrestrict blcQzR6rzCSGwr10LiNZBeKAWSfFhapJI7rUUf5yrxnSpuZstxuwC97RnvPQaRl

