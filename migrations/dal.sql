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

    SELECT id INTO d_order_id
    FROM public.orders 
    WHERE user_id = p_user_id AND is_paid = FALSE;

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

CREATE OR REPLACE PROCEDURE pay_for_order(p_user_id INT, p_order_id INT)
AS $$
DECLARE
    d_total_price DECIMAL(10, 2);
    d_user_balance INT;
    d_is_paid BOOLEAN;
BEGIN
    SELECT total_price, is_paid INTO d_total_price, d_is_paid
    FROM public.orders
    WHERE id = p_order_id;

    IF d_total_price IS NULL OR is_paid = TRUE THEN
        ROLLBACK;
    END IF;

    SELECT balance INTO d_user_balance
    FROM public.users
    WHERE id = p_user_id;

    IF d_user_balance IS NULL OR d_user_balance < d_total_price THEN
        ROLLBACK;
    END IF;

    UPDATE TABLE public.orders
    SET is_paid = TRUE
    WHERE id = p_order_id;

    UPDATE TABLE public.users
    SET balance = balance - d_total_price
    WHERE id = p_user_id;
END;
$$ LANGUAGE PLPGSQL;

CREATE OR REPLACE FUNCTION log_paid_order()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_paid = TRUE AND OLD.is_paid = FALSE THEN
        INSERT INTO public.buy_logs (buyer_id, album_id)
        SELECT NEW.user_id, oi.album_id
        FROM public.order_items oi
        WHERE oi.order_id = NEW.id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER orders_paid_trigger
AFTER UPDATE ON public.orders
FOR EACH ROW
EXECUTE FUNCTION log_paid_order();