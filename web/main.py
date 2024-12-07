import time
import requests
import streamlit as st
from itertools import islice
import base64
import json
from dotenv import load_dotenv
import os
import threading
import websocket

from globals import notification_list

load_dotenv()

NOTIFICATIONS_URL = "ws://localhost:{NOTIFICATIONS_PORT}/ws"
LOGIN_URL = "http://localhost:{GATEWAY_PORT}/login"
REGISTRATION_URL = "http://localhost:{GATEWAY_PORT}/registration"
LOGOUT_URL = "http://localhost:{GATEWAY_PORT}/logout"
SEARCH_URL = "http://localhost:{GATEWAY_PORT}/search"
ADD_TO_ORDER_URL = "http://localhost:{GATEWAY_PORT}/add/{album_id}"
REMOVE_FROM_ORDER_URL = "http://localhost:{GATEWAY_PORT}/remove/{album_id}"
ORDERS_URL = "http://localhost:{GATEWAY_PORT}/orders/"
DELETE_ALBUM = "http://localhost:{GATEWAY_PORT}/admin-panel/delete/{id}"

ADMIN_PASS = os.getenv("ADMIN_PASS", "")
GATEWAY_PORT = os.getenv("GATEWAY_PORT", "")
NOTIFICATIONS_PORT = os.getenv("NOTIFICATIONS_PORT", "")

def connect_to_notifications(jwt):
    ws_url = NOTIFICATIONS_URL.format(NOTIFICATIONS_PORT=NOTIFICATIONS_PORT)

    def on_message(ws, message):
        data = json.loads(message)
        notification_list.append(data)
        for notification in notification_list:
            print(notification)

    def on_error(ws, error):
        notification_list.append({"success": False, "message": f"WebSocket Error: {error}"})

    def on_close(ws, close_status_code, close_msg):
        notification_list.append({"success": False, "message": "WebSocket connection closed."})

    def on_open(ws):
        ws.send(json.dumps({"jwt": jwt}))

    ws = websocket.WebSocketApp(
        ws_url,
        on_message=on_message,
        on_error=on_error,
        on_close=on_close,
    )
    ws.on_open = on_open

    threading.Thread(target=ws.run_forever, daemon=True).start()

def fetch_orders():
    headers = get_authorization_header()
    if not headers:
        return None

    try:
        response = requests.get(ORDERS_URL.format(GATEWAY_PORT=GATEWAY_PORT), headers=headers)
        response.raise_for_status()
        data = response.json()
        all_orders = data.get("orders", [])
        unpaid_orders = [order for order in all_orders if not order.get("isPaid", False)]
        return unpaid_orders
    except requests.exceptions.RequestException as e:
        st.error(f"Error fetching orders: {e}")
        return None

def add_to_order(album_id):
    if "order" not in st.session_state:
        st.session_state["order"] = []
    
    if album_id not in st.session_state["order"]:
        st.session_state["order"].append(album_id)
        if "bought" in st.session_state:
            st.session_state["bought"] = False
        url = ADD_TO_ORDER_URL.format(GATEWAY_PORT=GATEWAY_PORT, album_id=album_id)
        headers = get_authorization_header()
        try:
            response = requests.post(url, headers=headers)
            response.raise_for_status()
            st.success(f"Album {album_id} added to order!")
            return True
        except requests.exceptions.RequestException as e:
            st.error(f"Error adding album to order: {e}")
            return False
    else:
        st.error(f"Album {album_id} is already in the order.")

def remove_from_order(album_id):
    if "order" in st.session_state and album_id in st.session_state["order"]:
        st.session_state["order"].remove(album_id)
        url = REMOVE_FROM_ORDER_URL.format(GATEWAY_PORT=GATEWAY_PORT, album_id=album_id)
        headers = get_authorization_header()
        try:
            response = requests.post(url, headers=headers)
            response.raise_for_status()
            st.success(f"Album {album_id} removed from order!")
            return True
        except requests.exceptions.RequestException as e:
            st.error(f"Error removing album from order: {e}")
            return False
    else:
        st.error(f"Album {album_id} is not in the order.")

def toggle_order_button(album_id, postfix: str = 'mainpage'):
    is_in_order = album_id in st.session_state.get("order", [])
    
    if is_in_order:
        st.button("Remove from order", key=f"remove_{album_id}_{postfix}", on_click=remove_from_order, args=(album_id,))
    else:
        st.button("Add to order", key=f"add_{album_id}_{postfix}", on_click=add_to_order, args=(album_id,))

def fetch_data():
    url = f"http://localhost:{GATEWAY_PORT}/"
    payload = {"albumsCount": 10, "artistsCount": 5}
    headers = get_authorization_header()
    try:
        response = requests.post(url, json=payload, headers=headers)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        st.error(f"Error fetching data: {e}")
        return None

def search(query):
    url = SEARCH_URL.format(GATEWAY_PORT=GATEWAY_PORT)
    payload = {"query": query}
    headers = get_authorization_header()
    try:
        response = requests.post(url, json=payload, headers=headers)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        st.error(f"Error searching: {e}")
        return None
    
def fetch_profile(entity_type, entity_id):
    url = f"http://localhost:{GATEWAY_PORT}/{entity_type}/{entity_id}"
    headers = get_authorization_header()
    try:
        response = requests.get(url, headers=headers)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        st.error(f"Error fetching profile: {e}")
        return None

def chunked(iterable, size):
    it = iter(iterable)
    return iter(lambda: list(islice(it, size)), [])

import validators
def get_image_url(image_url):
    if image_url == "" or not validators.url(image_url):
        return "https://via.placeholder.com/150"
    return image_url

def get_authorization_header():
    token = st.session_state.get("jwt_token", None)
    if token:
        return {"Authorization": f"Bearer {token}"}
    return {}

def login():
    st.subheader("Login")
    login = st.text_input("Username (email)", "")
    password = st.text_input("Password", "", type="password")
    
    if st.button("Login"):
        credentials = f"{login}:{password}"
        encoded_credentials = base64.b64encode(credentials.encode()).decode()
        headers = {"Authorization": f"Basic {encoded_credentials}"}
        
        try:
            response = requests.get(LOGIN_URL.format(GATEWAY_PORT=GATEWAY_PORT), headers=headers)
            response.raise_for_status()
            data = response.json()
            if "jwt" in data and "isAdmin" in data:
                st.session_state["jwt_token"] = data["jwt"]
                st.session_state["is_admin"] = data["isAdmin"]

                connect_to_notifications(data["jwt"])

                st.success("Logged in successfully!")
                st.rerun()

            else:
                st.error("Login failed: Invalid response format.")
        except requests.exceptions.RequestException as e:
            st.error(f"Login failed: {e}")

def register():
    st.subheader("Register")
    email = st.text_input("Email", "")
    nickname = st.text_input("Nickname", "")
    image_url = st.text_input("Image URL (optional)", "")
    password = st.text_input("Password", "", type="password")
    admin_pass = st.text_input("Admin Passphrase", "", type="password")

    is_admin = False
    if admin_pass == ADMIN_PASS:
        is_admin = True
    elif admin_pass:
        st.error("Invalid admin passphrase. You are not an admin.")

    if st.button("Register"):
        if email and nickname and password:
            registration_data = {
                "email": email,
                "nickname": nickname,
                "password": password,
                "isAdmin": is_admin,
                "imageURL": image_url if image_url else "-"
            }

            try:
                response = requests.post(REGISTRATION_URL.format(GATEWAY_PORT=GATEWAY_PORT), json=registration_data)
                response.raise_for_status()
                st.success("Registration successful! Please log in.")
            except requests.exceptions.RequestException as e:
                st.error(f"Registration failed: {e}")
        else:
            st.error("Please fill all required fields.")

def logout():
    token = st.session_state.get("jwt_token", None)
    if token:
        headers = {"Authorization": f"Bearer {token}"}
        try:
            response = requests.post(LOGOUT_URL.format(GATEWAY_PORT=GATEWAY_PORT), headers=headers)
            response.raise_for_status()
            st.session_state.pop("jwt_token", None)
            st.success("Logged out successfully!")
            st.rerun()
        except requests.exceptions.RequestException as e:
            st.error(f"Logout failed: {e}")

def delete_album(album_id):
    headers = get_authorization_header()
    
    try:
        response = requests.delete(DELETE_ALBUM.format(GATEWAY_PORT=GATEWAY_PORT, id=album_id), headers=headers)
        response.raise_for_status()
        st.session_state["reload"] = True
    except requests.exceptions.RequestException as e:
        st.error(f"Deletion failed: {e}")


def display_main_content():
    if "orders" not in st.session_state:
        unpaid_orders = fetch_orders()
        if unpaid_orders is not None:
            st.session_state["orders"] = unpaid_orders
        else:
            st.session_state["orders"] = []
    
    st.header("Search")
    search_query = st.text_input("Search for artists or albums", "")
    
    if search_query:
        search_results = search(search_query)

        if search_results:
            st.header("Search Results")

            if search_results.get("artists"):
                st.subheader("Artists")
                artist_columns = st.columns(min(len(search_results["artists"]), 5))
                for i, artist in enumerate(search_results["artists"]):
                    with artist_columns[i % 5]:
                        st.image(get_image_url(artist["imageURL"]), caption=artist["name"], use_container_width=True)
                        link = f"[View Artist Profile](?entity=artists&id={artist['id']})"
                        st.markdown(link, unsafe_allow_html=True)
                        
            if search_results.get("albums"):
                st.subheader("Albums")
                album_columns = st.columns(min(len(search_results["albums"]), 5))
                for i, album in enumerate(search_results["albums"]):
                    with album_columns[i % 5]:
                        st.image(get_image_url(album["imageURL"]), caption=album["name"], use_container_width=True)
                        link = f"[View Album Profile](?entity=albums&id={album['id']})"
                        st.markdown(link, unsafe_allow_html=True)

                        if "jwt_token" in st.session_state:
                            toggle_order_button(album['id'], 'search')
                        if "is_admin" in st.session_state and st.session_state["is_admin"] == True:
                            st.button("Delete", on_click=delete_album, args=(album['id'], ), key=f"delete_{album['id']}")
    if "data" not in st.session_state or st.session_state["data"] == None or ("reload" in st.session_state and st.session_state["reload"] == True):
        st.session_state["data"] = fetch_data()
        st.session_state["visited"] = True
        st.session_state["reload"] = False

    if st.session_state["data"]:
        data = st.session_state["data"]
        st.header("Artists")
        artist_columns = st.columns(len(data["artists"]))
        for artist, col in zip(data["artists"], artist_columns):
            with col:
                st.image(get_image_url(artist["imageURL"]), caption=artist["name"], use_container_width=True)
                link = f"[View Artist Profile](?entity=artists&id={artist['id']})"
                st.markdown(link, unsafe_allow_html=True)

        st.header("Albums")
        album_chunks = list(chunked(data["albums"], 5))
        for row in album_chunks:
            album_columns = st.columns(len(row))
            for album, col in zip(row, album_columns):
                with col:
                    st.image(get_image_url(album["imageURL"]), caption=album["name"], use_container_width=True)
                    link = f"[View Album Profile](?entity=albums&id={album['id']})"
                    st.markdown(link, unsafe_allow_html=True)
                    
                    if "jwt_token" in st.session_state:
                        toggle_order_button(album['id'])
                    if "is_admin" in st.session_state and st.session_state["is_admin"] == True:
                        st.button("Delete", on_click=delete_album, args=(album['id'], ), key=f"delete_{album['id']}")

def display_profile():
    query_params = st.query_params
    selected_entity = query_params.get("entity", [None])
    selected_id = query_params.get("id", [None])
    
    if selected_entity and selected_id:
        if selected_entity == "artists":
            artist_data = fetch_profile("artists", selected_id)
            if artist_data:
                artist = artist_data["artist"]
                albums = artist_data.get("albums", [])
                st.header(f"Artist: {artist['name']}")
                st.image(artist["imageURL"], width=150)
                st.write(f"**Genre:** {artist['genre']}")
                
                if albums:
                    st.subheader("Albums")
                    album_chunks = list(chunked(albums, 5))
                    for row in album_chunks:
                        album_columns = st.columns(len(row))
                        for album, col in zip(row, album_columns):
                            with col:
                                st.image(get_image_url(album["imageURL"]), caption=album["name"], width=150)
                                st.write(f"**Price:** ${album['price']:.2f}")
                                link = f"[View Album Profile](?entity=albums&id={album['id']})"
                                st.markdown(link, unsafe_allow_html=True)
                else:
                    st.write("This artist has no albums.")
        elif selected_entity == "albums":
            album_data = fetch_profile("albums", selected_id)
            if album_data:
                album = album_data["album"]
                st.header(f"Album: {album['name']}")
                st.image(album["imageURL"], width=150)
                st.write(f"**Author:** {album['author']['name']} ({album['author']['genre']})")
                st.image(album["author"]["imageURL"], caption="Author Image", width=150)
                st.write(f"**Price:** ${album['price']:.2f}")
                st.subheader("Tracks")
                for track in album["tracks"]:
                    st.write(f"{track['number']}. {track['name']}")
                
        else:
            st.error("Invalid entity or ID")


def sidebar_navigation():
    query_params = st.query_params
    if "entity" not in query_params:
        if "jwt_token" in st.session_state:
            if st.sidebar.button("Logout"):
                logout()
        else:
            auth_choice = st.sidebar.selectbox("Choose an option", ["Login", "Register"])

            if auth_choice == "Login":
                login()
            elif auth_choice == "Register":
                register()

query_params = st.query_params
sidebar_navigation()

if query_params.get("entity") and query_params.get("id"):
    display_profile()
else:
    display_main_content()
