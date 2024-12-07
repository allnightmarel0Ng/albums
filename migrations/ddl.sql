DROP TABLE IF EXISTS public.notifications CASCADE;
DROP TABLE IF EXISTS public.buy_logs CASCADE;
DROP TABLE IF EXISTS public.order_items CASCADE;
DROP TABLE IF EXISTS public.orders CASCADE;
DROP TABLE IF EXISTS public.purchased_albums CASCADE;
DROP TABLE IF EXISTS public.tracks CASCADE;
DROP TABLE IF EXISTS public.albums CASCADE;
DROP TABLE IF EXISTS public.artists CASCADE;
DROP TABLE IF EXISTS public.credentials CASCADE;
DROP TABLE IF EXISTS public.users CASCADE;

CREATE TABLE public.users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) NOT NULL UNIQUE,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    nickname VARCHAR(30) NOT NULL,
    balance DECIMAL(10, 2) NOT NULL DEFAULT 0,
    image_url VARCHAR(255) NOT NULL DEFAULT '-'
);

CREATE TABLE public.credentials (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES public.users(id) ON DELETE CASCADE,
    password_hash VARCHAR(70) NOT NULL
);

CREATE TABLE public.artists (
    id SERIAL PRIMARY KEY,
    name VARCHAR(512) NOT NULL,
    genre VARCHAR(64) NOT NULL,
    image_url VARCHAR(128) NOT NULL
);

CREATE TABLE public.albums (
    id SERIAL PRIMARY KEY,
    name VARCHAR(512) NOT NULL,
    artist_id INT REFERENCES public.artists(id) ON DELETE SET NULL,
    image_url VARCHAR(128) NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);

CREATE TABLE public.tracks (
    id SERIAL PRIMARY KEY,
    album_id INT REFERENCES public.albums(id) ON DELETE SET NULL,
    name VARCHAR(512) NOT NULL,
    number INT NOT NULL
);

CREATE TABLE public.purchased_albums (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES public.users(id) ON DELETE CASCADE,
    album_id INT REFERENCES public.albums(id) ON DELETE CASCADE 
);

CREATE TABLE public.orders (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES public.users(id) ON DELETE SET NULL,
    date TIMESTAMP NOT NULL DEFAULT NOW(),
    total_price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    is_paid BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE public.order_items (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES public.orders(id) ON DELETE CASCADE,
    album_id INT REFERENCES public.albums(id) ON DELETE SET NULL
);

CREATE TABLE public.buy_logs (
    id SERIAL PRIMARY KEY,
    buyer_id INT REFERENCES public.users(id) ON DELETE SET NULL,
    album_id INT REFERENCES public.albums(id) ON DELETE SET NULL,
    logging_time TIMESTAMP NOT NULL DEFAULT NOW()
);
