import requests
import streamlit as st

from main import get_authorization_header, GATEWAY_PORT

ORDERS_URL = "http://localhost:{GATEWAY_PORT}/orders"

def fetch_orders():
    headers = get_authorization_header()
    if not headers:
        st.error("You must be logged in to view your orders.")
        return None
    
    try:
        response = requests.get(ORDERS_URL.format(GATEWAY_PORT=GATEWAY_PORT), headers=headers)
        response.raise_for_status()
        data = response.json()
        return data.get("orders", [])
    except requests.exceptions.RequestException as e:
        st.error(f"Error fetching orders: {e}")
        return None

def buy_order():
    url = f"http://localhost:{GATEWAY_PORT}/buy"
    headers = get_authorization_header()

    try:
        response = requests.post(url, headers=headers)

        if response.status_code == 200:
            st.success(f"Order buying process has been started! Wait for info in notifications!")
            if "order" in st.session_state:
                st.session_state["order"].clear()
            st.session_state["bought"] = True
        else:
            data = response.json()
            st.error(f"Failed to buy order: {data.get('error', 'Unknown error')}")
    except requests.exceptions.RequestException as e:
        st.error(f"Error while buying order: {e}")

def display_order(order):
    if isinstance(order, dict):
        st.subheader(f"Order ID: {order['id']}")
        st.write(f"**Total Price:** ${order['totalPrice']:.2f}")
        st.write(f"**Is Paid:** {'Yes' if order['isPaid'] else 'No'}")
        st.write(f"**Date:** {order['date']}")
        
        if order.get("albums"):
            st.write("**Albums in this order:**")
            album_columns = st.columns(min(len(order['albums']), 5))
            for i, album in enumerate(order['albums']):
                with album_columns[i % 5]:
                    st.image(album["imageURL"], caption=album["name"], use_container_width=True)
        
        if order['isPaid'] == False:
            st.button('Buy', on_click=buy_order, disabled=st.session_state.get("bought", False))

def display_orders():
    orders_data = fetch_orders()
    if orders_data:
        st.header("Your Orders")
        if isinstance(orders_data, list):
            for order in orders_data:
                display_order(order)
        else:
            st.error("Invalid response format: Expected a list of orders.")

display_orders()
