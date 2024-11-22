DROP TABLE IF EXISTS public.order_items;
DROP TABLE IF EXISTS public.orders;
DROP TABLE IF EXISTS public.tracks;
DROP TABLE IF EXISTS public.albums;
DROP TABLE IF EXISTS public.customers;
DROP TABLE IF EXISTS public.artists;
DROP TABLE IF EXISTS public.users;
DROP TYPE IF EXISTS role;

CREATE TYPE role AS ENUM ('artist', 'customer');
CREATE TABLE public.users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role role NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE public.customers (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES public.users(id) ON DELETE CASCADE NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    balance MONEY NOT NULL,
);

CREATE TABLE public.artists (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES public.users(id) ON DELETE CASCADE NOT NULL,
    name VARCHAR(100) NOT NULL,
    profile_picture_url VARCHAR(100),
    bio TEXT
);

CREATE TABLE public.albums (
    id SERIAL PRIMARY KEY,
    artist_id INT REFERENCES public.artists(id) ON DELETE SET NULL,
    name VARCHAR(50) NOT NULL,
    release_date DATE NOT NULL,
    cover_art_url VARCHAR(255),
    price MONEY NOT NULL,
    genre VARCHAR(50)
);

CREATE TABLE public.tracks (
    id SERIAL PRIMARY KEY,
    album_id INT REFERENCES public.artists(id) ON DELETE SET NULL,
    name VARCHAR(50) NOT NULL,
    number INT NOT NULL,
    duration INT NOT NULL,
    audio_file_url VARCHAR(255) NOT NULL
);

CREATE TABLE public.orders (
    id SERIAL PRIMARY KEY,
    customer_id INT REFERENCES public.customers(id) ON DELETE SET NULL,
    date TIMESTAMP NOT NULL DEFAULT NOW(),
    total_price MONEY NOT NULL
);

CREATE TABLE public.order_items (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES public.orders(id) ON DELETE CASCADE,
    album_id INT REFERENCES public.albums(id) ON DELETE SET NULL,
    quantity INT NOT NULL,
    price MONEY NOT NULL
);