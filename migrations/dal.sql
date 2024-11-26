CREATE OR REPLACE FUNCTION add_album_to_user_order(user_id INT, album_id INT)
RETURNS VOID AS $$
DECLARE
    order_id INT;
    order_item_id INT;
    album_price DECIMAL(10, 2);
BEGIN
    BEGIN
        SELECT price INTO album_price 
        FROM public.albums 
        WHERE id = album_id;

        IF album_price IS NULL THEN
            RAISE EXCEPTION 'Album with ID % does not exist.', album_id;
        END IF;

        SELECT id INTO order_id 
        FROM public.orders 
        WHERE user_id = user_id AND is_paid = FALSE;

        IF order_id IS NULL THEN
            INSERT INTO public.orders (user_id) 
            VALUES (user_id)
            RETURNING id INTO order_id;
        END IF;

        SELECT id INTO order_item_id 
        FROM public.order_items 
        WHERE order_id = order_id AND album_id = album_id;
        
        IF order_item_id IS NOT NULL THEN
            RAISE EXCEPTION 'Album with ID % already in order.', album_id;
        END IF;

        INSERT INTO public.order_items (order_id, album_id)
        VALUES (order_id, album_id);

        UPDATE public.orders 
        SET total_price = total_price + album_price 
        WHERE id = order_id;

        COMMIT;
    EXCEPTION
        WHEN OTHERS THEN
            ROLLBACK;
            RAISE;
    END;
END;
$$ LANGUAGE PLPGSQL;

CREATE OR REPLACE FUNCTION delete_album_from_user_order(user_id INT, album_id INT)
RETURNS VOID AS $$
DECLARE
    order_id INT;
    order_item_id INT;
    album_price DECIMAL(10, 2);
BEGIN
    BEGIN
        SELECT price INTO album_price 
        FROM public.albums 
        WHERE id = album_id;

        IF album_price IS NULL THEN
            RAISE EXCEPTION 'Album with ID % does not exist.', album_id;
        END IF;

        SELECT id INTO order_id 
        FROM public.orders 
        WHERE user_id = user_id AND is_paid = FALSE;

        IF order_id IS NULL THEN
            RAISE EXCEPTION 'No order found to delete from.';
        END IF;

        SELECT id INTO order_item_id
        FROM public.order_items
        WHERE order_id = order_id;

        IF order_item_id IS NULL THEN
            RAISE EXCEPTION 'Order item with % does not exist', album_id;
        END IF;

        DELETE 
        FROM public.order_items
        WHERE id = order_item_id;

        UPDATE public.orders
        SET total_price = total_price - album_price
        WHERE id = order_id;

        COMMIT;
    EXCEPTION
        WHEN OTHERS THEN
            ROLLBACK;
            RAISE;
    END;
END;
$$ LANGUAGE PLPGSQL;
