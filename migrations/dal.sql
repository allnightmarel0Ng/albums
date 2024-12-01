CREATE OR REPLACE PROCEDURE add_album_to_user_order(p_user_id INT, p_album_id INT)
AS $$
DECLARE
    d_order_id INT;
    d_order_item_id INT;
    d_album_price DECIMAL(10, 2);
BEGIN
    SELECT price INTO d_album_price 
    FROM public.albums 
    WHERE id = p_album_id;

    IF d_album_price IS NULL THEN
        ROLLBACK;
    END IF;

    SELECT o.id INTO d_order_id
    FROM public.orders AS o
    WHERE o.user_id = p_user_id AND o.is_paid = FALSE;

    IF d_order_id IS NULL THEN
        INSERT INTO public.orders (user_id) 
        VALUES (p_user_id)
        RETURNING id INTO d_order_id;
    END IF;

    SELECT id INTO d_order_item_id 
    FROM public.order_items 
    WHERE order_id = d_order_id AND album_id = p_album_id;
    
    IF d_order_item_id IS NOT NULL THEN
        ROLLBACK;
    END IF;

    INSERT INTO public.order_items (order_id, album_id)
    VALUES (d_order_id, p_album_id);

    UPDATE public.orders 
    SET total_price = total_price + d_album_price 
    WHERE id = d_order_id;
END;
$$ LANGUAGE PLPGSQL;

CREATE OR REPLACE PROCEDURE delete_album_from_user_order(p_user_id INT, p_album_id INT)
AS $$
DECLARE
    d_order_id INT;
    d_order_item_id INT;
    d_album_price DECIMAL(10, 2);
BEGIN
    SELECT price INTO d_album_price 
    FROM public.albums 
    WHERE id = p_album_id;

    IF d_album_price IS NULL THEN
        ROLLBACK;
    END IF;

    SELECT id INTO d_order_id 
    FROM public.orders 
    WHERE user_id = p_user_id AND is_paid = FALSE;

    IF d_order_id IS NULL THEN
        ROLLBACK;
    END IF;

    SELECT id INTO d_order_item_id
    FROM public.order_items
    WHERE order_id = d_order_id;

    IF d_order_item_id IS NULL THEN
        ROLLBACK;
    END IF;

    DELETE 
    FROM public.order_items
    WHERE id = d_order_item_id;

    UPDATE public.orders
    SET total_price = total_price - d_album_price
    WHERE id = d_order_id;
END;
$$ LANGUAGE PLPGSQL;
